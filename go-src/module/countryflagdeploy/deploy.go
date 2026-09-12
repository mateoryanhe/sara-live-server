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
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/countryflagdeploydto"
	"xr-game-server/entity/sys"
)

var errImageRootNotConfigured = errors.New("image static root not configured")

// DeployZipFromRequest 上传 zip → 新 version 目录 → 入库 → 删除旧 version 目录
func DeployZipFromRequest(r *ghttp.Request) (*countryflagdeploydto.DeployCountryFlagZipRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload file is empty")
	}
	root, err := getFlagsRoot()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create flags root %s: %w", root, err)
	}

	zipPath, err := saveUploadZip(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)

	version := time.Now().Format("20060102150405")
	versionDir := filepath.Join(root, version)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return nil, fmt.Errorf("create version dir %s: %w", versionDir, err)
	}

	fileCount, err := extractFlagPNGs(zipPath, versionDir)
	if err != nil {
		os.RemoveAll(versionDir)
		return nil, err
	}
	if fileCount == 0 {
		os.RemoveAll(versionDir)
		return nil, errors.New("zip has no png flag files")
	}

	if err := saveVersion(version); err != nil {
		os.RemoveAll(versionDir)
		return nil, err
	}
	ReloadCountryFlagCache()

	removed := cleanupOldVersions(root, version)
	return &countryflagdeploydto.DeployCountryFlagZipRes{
		Version:    version,
		FileCount:  fileCount,
		DeployPath: versionDir,
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
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0
	}
	removed := 0
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

func extractFlagPNGs(zipPath, destDir string) (fileCount int, err error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	destDir = filepath.Clean(destDir)
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
		targetPath := filepath.Join(destDir, code+".png")
		if err := extractZipFile(file, targetPath); err != nil {
			return fileCount, err
		}
		fileCount++
	}
	return fileCount, nil
}

func extractZipFile(file *zip.File, targetPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func joinFlagsRoot(imageRoot string) string {
	return filepath.Join(imageRoot, countryflagdeploydto.CountryFlagStaticDir)
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
