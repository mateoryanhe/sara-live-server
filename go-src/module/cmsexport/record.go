package cmsexport

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	appcfg "xr-game-server/core/cfg"
	"xr-game-server/core/xrtimer"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

type exportRecord struct {
	exportID  string
	fileName  string
	absPath   string
	fileURL   string
	createdAt time.Time
}

type exportResult struct {
	ExportID string
	FileName string
	FileUrl  string
	Total    int
}

var (
	errExportUnavailable = errors.New("CMS文件导出未配置")
	errExportEmpty       = errors.New("没有可导出的数据")
)

var (
	exportRecords  sync.Map
	exportInitOnce sync.Once
)

func initExportCleanup() {
	exportInitOnce.Do(func() {
		if appcfg.GetCMSFileExportRoot() == "" {
			return
		}
		xrtimer.AddSingleton(gctx.New(), time.Minute, cleanupExports)
	})
}

func ensureExportReady() error {
	if appcfg.GetCMSFileExportRoot() == "" {
		return errExportUnavailable
	}
	initExportCleanup()
	return nil
}

func createExportFile(suffix string) (*exportRecord, string, error) {
	if err := ensureExportReady(); err != nil {
		return nil, "", err
	}
	exportID := guid.S()
	fileName := exportID + suffix
	absPath := appcfg.JoinCMSFileExportPath(fileName)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return nil, "", err
	}
	record := &exportRecord{
		exportID:  exportID,
		fileName:  fileName,
		absPath:   absPath,
		fileURL:   appcfg.BuildCMSFileExportURL(fileName),
		createdAt: time.Now(),
	}
	exportRecords.Store(exportID, record)
	return record, absPath, nil
}

func finalizeExportRecord(record *exportRecord, total int) *exportResult {
	if record == nil {
		return nil
	}
	return &exportResult{
		ExportID: record.exportID,
		FileName: record.fileName,
		FileUrl:  record.fileURL,
		Total:    total,
	}
}

func deleteExport(exportID string) error {
	if exportID == "" {
		return errExportUnavailable
	}
	if value, ok := exportRecords.Load(exportID); ok {
		if record, ok := value.(*exportRecord); ok {
			removeFile(record.absPath)
		}
		exportRecords.Delete(exportID)
	}
	removeFile(appcfg.JoinCMSFileExportPath(exportID + ".csv"))
	return nil
}

func cleanupExports(_ context.Context) {
	if appcfg.GetCMSFileExportRoot() == "" {
		return
	}
	expireBefore := time.Now().Add(-time.Duration(appcfg.GetCMSFileExportTTLMinutes()) * time.Minute)
	exportRecords.Range(func(key, value any) bool {
		record, ok := value.(*exportRecord)
		if !ok || record.createdAt.After(expireBefore) {
			return true
		}
		removeFile(record.absPath)
		exportRecords.Delete(key)
		return true
	})
	cleanExportDir(appcfg.ResolveCMSFileExportDir(), expireBefore)
}

func removeFile(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}

func cleanExportDir(dir string, expireBefore time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil || info.ModTime().After(expireBefore) {
			continue
		}
		removeFile(filepath.Join(dir, entry.Name()))
	}
}
