package staticdeploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var deployMu sync.Mutex

// DeployZip extracts a complete package into a sibling staging directory and
// replaces the current site only after extraction succeeds. The replacement
// removes files left by older releases while an invalid package leaves the
// current site untouched.
func DeployZip(zipPath, destRoot string) (fileCount, dirCount int, err error) {
	deployMu.Lock()
	defer deployMu.Unlock()

	destRoot = filepath.Clean(destRoot)
	parentDir := filepath.Dir(destRoot)
	baseName := filepath.Base(destRoot)
	if baseName == "." || baseName == string(os.PathSeparator) || destRoot == parentDir {
		return 0, 0, fmt.Errorf("invalid deploy dir: %s", destRoot)
	}
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("create deploy parent dir %s: %w", parentDir, err)
	}

	stageDir, err := os.MkdirTemp(parentDir, "."+baseName+"-stage-*")
	if err != nil {
		return 0, 0, fmt.Errorf("create deploy staging dir: %w", err)
	}
	defer os.RemoveAll(stageDir)
	if err := os.Chmod(stageDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("set deploy staging dir permission: %w", err)
	}

	fileCount, dirCount, err = extractZip(zipPath, stageDir)
	if err != nil {
		return fileCount, dirCount, err
	}
	if fileCount == 0 {
		return 0, dirCount, errors.New("zip package contains no deployable files")
	}
	if err := replaceDeployDir(stageDir, destRoot); err != nil {
		return fileCount, dirCount, err
	}
	return fileCount, dirCount, nil
}

func replaceDeployDir(stageDir, destRoot string) error {
	parentDir := filepath.Dir(destRoot)
	baseName := filepath.Base(destRoot)
	backupDir, err := os.MkdirTemp(parentDir, "."+baseName+"-backup-*")
	if err != nil {
		return fmt.Errorf("reserve deploy backup dir: %w", err)
	}
	if err := os.Remove(backupDir); err != nil {
		return fmt.Errorf("prepare deploy backup dir: %w", err)
	}

	hasCurrent := false
	if _, statErr := os.Stat(destRoot); statErr == nil {
		hasCurrent = true
		if err := os.Rename(destRoot, backupDir); err != nil {
			return fmt.Errorf("backup current deploy dir: %w", err)
		}
	} else if !os.IsNotExist(statErr) {
		return fmt.Errorf("inspect current deploy dir: %w", statErr)
	}

	if err := os.Rename(stageDir, destRoot); err != nil {
		if hasCurrent {
			if restoreErr := os.Rename(backupDir, destRoot); restoreErr != nil {
				return fmt.Errorf("activate new deploy dir: %v; restore previous deploy dir: %w", err, restoreErr)
			}
		}
		return fmt.Errorf("activate new deploy dir: %w", err)
	}
	if hasCurrent {
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf("new files deployed but remove previous deploy dir failed: %w", err)
		}
	}
	return nil
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

	mode := file.Mode().Perm()
	if mode == 0 {
		mode = 0644
	}
	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
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
	target := filepath.Clean(filepath.Join(base, filepath.FromSlash(name)))
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
