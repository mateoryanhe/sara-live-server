package thirdpaydeploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/cfg"
	"xr-game-server/dto/thirdpaydeploydto"
)

func DeployZipFromRequest(r *ghttp.Request) (*thirdpaydeploydto.DeployThirdPayZipRes, error) {
	if r == nil || r.Request == nil {
		return nil, errors.New("upload file is empty")
	}
	deployDir, err := getDeployDir()
	if err != nil {
		return nil, err
	}

	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, mapUploadReadErr(err)
	}

	var zipPath string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, mapUploadReadErr(err)
		}
		if part.FormName() != "file" {
			part.Close()
			continue
		}
		ext := strings.ToLower(filepath.Ext(part.FileName()))
		if ext != ".zip" {
			part.Close()
			return nil, fmt.Errorf("file ext not allowed: %s", ext)
		}
		tmpFile, err := os.CreateTemp("", "third-pay-deploy-*.zip")
		if err != nil {
			part.Close()
			return nil, err
		}
		zipPath = tmpFile.Name()
		_, copyErr := io.Copy(tmpFile, part)
		part.Close()
		closeErr := tmpFile.Close()
		if copyErr != nil {
			os.Remove(zipPath)
			return nil, mapUploadReadErr(copyErr)
		}
		if closeErr != nil {
			os.Remove(zipPath)
			return nil, closeErr
		}
		break
	}
	if zipPath == "" {
		return nil, errors.New("upload file is empty")
	}
	defer os.Remove(zipPath)

	fileCount, dirCount, err := extractZip(zipPath, deployDir)
	if err != nil {
		return nil, err
	}
	return &thirdpaydeploydto.DeployThirdPayZipRes{
		FileCount:  fileCount,
		DirCount:   dirCount,
		DeployPath: deployDir,
		UrlPrefix:  thirdpaydeploydto.ThirdPayStaticPrefix,
	}, nil
}

func getDeployDir() (string, error) {
	root := strings.TrimSpace(cfg.GetStaticPathRoot(thirdpaydeploydto.ThirdPayStaticPrefix))
	if root == "" {
		return "", fmt.Errorf("static path not configured for prefix %s", thirdpaydeploydto.ThirdPayStaticPrefix)
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", fmt.Errorf("create deploy dir %s: %w", root, err)
	}
	return root, nil
}

func extractZip(zipPath, destRoot string) (fileCount, dirCount int, err error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close()

	destRoot = filepath.Clean(destRoot)
	for _, file := range reader.File {
		name := strings.TrimSpace(file.Name)
		if name == "" {
			continue
		}
		name = filepath.ToSlash(name)
		if shouldSkipZipEntry(name) {
			continue
		}
		targetPath, err := safeJoin(destRoot, name)
		if err != nil {
			return fileCount, dirCount, err
		}
		if file.FileInfo().IsDir() || strings.HasSuffix(name, "/") {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fileCount, dirCount, err
			}
			dirCount++
			continue
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fileCount, dirCount, err
		}
		if err := extractZipFile(file, targetPath); err != nil {
			return fileCount, dirCount, err
		}
		fileCount++
	}
	return fileCount, dirCount, nil
}

func extractZipFile(file *zip.File, targetPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode().Perm())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func safeJoin(base, name string) (string, error) {
	name = strings.TrimPrefix(filepath.ToSlash(name), "/")
	if name == "" || strings.Contains(name, "..") {
		return "", fmt.Errorf("invalid zip entry path: %s", name)
	}
	target := filepath.Join(base, filepath.FromSlash(name))
	target = filepath.Clean(target)
	baseClean := filepath.Clean(base)
	if target != baseClean && !strings.HasPrefix(target, baseClean+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid zip entry path: %s", name)
	}
	return target, nil
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
