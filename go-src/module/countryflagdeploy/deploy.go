package countryflagdeploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/countryflagdeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/module/upload"
)

var errImageRootNotConfigured = errors.New("image static root not configured")

// DeployZipFromRequest 上传 zip → 按头像同一套 upload 存储写入 country-flags/{version} → 入库 → 清理旧 version
func DeployZipFromRequest(r *ghttp.Request) (*countryflagdeploydto.DeployCountryFlagZipRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload file is empty")
	}
	root, err := getFlagsRoot()
	if err != nil {
		return nil, err
	}

	zipPath, err := saveUploadZip(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)

	version := time.Now().Format("20060102150405")
	fileCount, err := extractAndStoreFlagPNGs(zipPath, version)
	if err != nil {
		cleanupVersion(version)
		return nil, err
	}
	if fileCount == 0 {
		cleanupVersion(version)
		return nil, errors.New("zip has no png flag files")
	}

	if err := saveVersion(version); err != nil {
		cleanupVersion(version)
		return nil, err
	}
	ReloadCountryFlagCache()

	removed := cleanupOldVersions(root, version)
	return &countryflagdeploydto.DeployCountryFlagZipRes{
		Version:    version,
		FileCount:  fileCount,
		DeployPath: filepath.Join(root, version),
		UrlPrefix:  buildURLPrefix(version),
		Removed:    removed,
	}, nil
}

func saveVersion(version string) error {
	existing := cfgdao.LoadCountryFlagCfg()
	row := &entity.CountryFlagCfg{Version: version}
	now := time.Now()
	if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = now
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	return cfgdao.SaveCountryFlagCfg(row)
}

func cleanupOldVersions(root, keepVersion string) int {
	removed := 0
	if entries, err := os.ReadDir(root); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if name == keepVersion || name == "." || name == ".." {
				continue
			}
			path := filepath.Join(root, name)
			if err := os.RemoveAll(path); err == nil {
				removed++
			}
		}
	}
	removed += cleanupOldVersionsOnS3(keepVersion)
	return removed
}

func cleanupVersion(version string) {
	version = strings.TrimSpace(version)
	if version == "" {
		return
	}
	if root, err := getFlagsRoot(); err == nil {
		_ = os.RemoveAll(filepath.Join(root, version))
	}
	prefix := country.AssetDir + "/" + version
	objs, err := upload.ListStoredObjectsByPrefix(prefix)
	if err != nil {
		// 未开 S3 时 List 会报错,本地已在上面清理
		return
	}
	for _, o := range objs {
		upload.DeleteUploadedFile(o.StoredName)
	}
}

func cleanupOldVersionsOnS3(keepVersion string) int {
	if !upload.IsS3Enabled() {
		return 0
	}
	objs, err := upload.ListStoredObjectsByPrefix(country.AssetDir)
	if err != nil {
		return 0
	}
	removedKeys := map[string]struct{}{}
	removed := 0
	for _, o := range objs {
		name := strings.Trim(strings.ReplaceAll(o.StoredName, "\\", "/"), "/")
		parts := strings.Split(name, "/")
		// country-flags/{version}/xx.png
		if len(parts) < 3 || parts[0] != country.AssetDir {
			continue
		}
		ver := parts[1]
		if ver == keepVersion {
			continue
		}
		upload.DeleteUploadedFile(name)
		if _, ok := removedKeys[ver]; !ok {
			removedKeys[ver] = struct{}{}
			removed++
		}
	}
	return removed
}

func saveUploadZip(r *ghttp.Request) (string, error) {
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return "", mapUploadReadErr(err)
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", mapUploadReadErr(err)
		}
		if part.FormName() != "file" {
			part.Close()
			continue
		}
		ext := strings.ToLower(filepath.Ext(part.FileName()))
		if ext != ".zip" {
			part.Close()
			return "", fmt.Errorf("file ext not allowed: %s", ext)
		}
		tmpFile, err := os.CreateTemp("", "country-flag-deploy-*.zip")
		if err != nil {
			part.Close()
			return "", err
		}
		zipPath := tmpFile.Name()
		_, copyErr := io.Copy(tmpFile, part)
		part.Close()
		closeErr := tmpFile.Close()
		if copyErr != nil {
			os.Remove(zipPath)
			return "", mapUploadReadErr(copyErr)
		}
		if closeErr != nil {
			os.Remove(zipPath)
			return "", closeErr
		}
		return zipPath, nil
	}
	return "", errors.New("upload file is empty")
}

func extractAndStoreFlagPNGs(zipPath, version string) (fileCount int, err error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		name := strings.TrimSpace(file.Name)
		if name == "" || file.FileInfo().IsDir() || strings.HasSuffix(name, "/") {
			continue
		}
		name = filepath.ToSlash(name)
		if shouldSkipZipEntry(name) {
			continue
		}
		base := filepath.Base(name)
		if strings.ToLower(filepath.Ext(base)) != ".png" {
			continue
		}
		code := strings.ToLower(strings.TrimSuffix(base, filepath.Ext(base)))
		if len(code) != 2 {
			continue
		}
		rel := country.RelPath(code, version)
		if rel == "" {
			continue
		}
		src, openErr := file.Open()
		if openErr != nil {
			return fileCount, openErr
		}
		data, readErr := io.ReadAll(src)
		src.Close()
		if readErr != nil {
			return fileCount, readErr
		}
		// 与头像一致:开 S3 只写云,否则写统一 storagePath
		if err := upload.SaveUploadedFileBytes(rel, data); err != nil {
			return fileCount, err
		}
		fileCount++
	}
	return fileCount, nil
}

func shouldSkipZipEntry(name string) bool {
	base := filepath.Base(name)
	if base == ".DS_Store" || base == "Thumbs.db" {
		return true
	}
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "__macosx/") || strings.Contains(lower, "/__macosx/")
}

func mapUploadReadErr(err error) error {
	if err == nil {
		return nil
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return errors.New("upload request incomplete")
	}
	if strings.Contains(strings.ToLower(err.Error()), "unexpected eof") {
		return errors.New("upload request incomplete")
	}
	return err
}
