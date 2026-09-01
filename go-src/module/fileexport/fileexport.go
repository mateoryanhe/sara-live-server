// Package fileexport 统一 CMS 业务导出与日志查询导出的落盘/TTL 清理。
// 前端可主动 deleteExport；后端对每个文件挂 AddOnce timer 兜底。
// 不做目录扫描按 mtime 删除，避免误伤同目录上传资源。
package fileexport

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
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/guid"
)

var ErrUnavailable = errors.New("CMS文件导出未配置")

// 前端主动删除时，若内存记录已丢(重启)，按常见后缀兜底删一次。
var knownSuffixes = []string{
	".csv",
	".log",
	".trace.log",
	".stats.tsv",
	".trend.tsv",
}

type Record struct {
	ExportID  string
	FileName  string
	AbsPath   string
	FileURL   string
	CreatedAt time.Time
}

type entry struct {
	record *Record
	timer  *gtimer.Entry
}

var records sync.Map // exportID -> *entry

func Ready() bool {
	return appcfg.GetCMSFileExportRoot() != ""
}

// Create 分配 exportID 与落盘路径并注册 TTL timer。
// 写出失败时请调用 Delete 取消 timer 并清理文件。
func Create(suffix string) (*Record, error) {
	if !Ready() {
		return nil, ErrUnavailable
	}
	exportID := guid.S()
	return Register(exportID, exportID+suffix)
}

// Register 登记已生成的导出文件并启动 TTL 删除 timer（幂等：同 ID 先取消旧 timer）。
func Register(exportID, fileName string) (*Record, error) {
	if !Ready() {
		return nil, ErrUnavailable
	}
	if exportID == "" || fileName == "" {
		return nil, ErrUnavailable
	}
	absPath := appcfg.JoinCMSFileExportPath(fileName)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return nil, err
	}
	rec := &Record{
		ExportID:  exportID,
		FileName:  fileName,
		AbsPath:   absPath,
		FileURL:   appcfg.BuildCMSFileExportURL(fileName),
		CreatedAt: time.Now(),
	}
	schedule(rec)
	return rec, nil
}

// Delete 取消 timer 并删除该导出文件（前端主动删除 / 导出失败回滚 / timer 到期）。
func Delete(exportID string) error {
	if exportID == "" {
		return ErrUnavailable
	}
	cancelTimer(exportID)

	if v, ok := records.LoadAndDelete(exportID); ok {
		if e, ok := v.(*entry); ok && e != nil && e.record != nil {
			removeFile(e.record.AbsPath)
		}
	}
	// 重启后 map 为空，或文件名后缀与登记不一致时兜底
	for _, suffix := range knownSuffixes {
		removeFile(appcfg.JoinCMSFileExportPath(exportID + suffix))
	}
	return nil
}

func schedule(rec *Record) {
	if rec == nil || rec.ExportID == "" {
		return
	}
	cancelTimer(rec.ExportID)

	ttlMin := appcfg.GetCMSFileExportTTLMinutes()
	if ttlMin <= 0 {
		ttlMin = 30
	}
	delay := time.Duration(ttlMin) * time.Minute
	exportID := rec.ExportID
	timer := xrtimer.AddOnce(gctx.New(), delay, func(ctx context.Context) {
		_ = Delete(exportID)
	})
	records.Store(exportID, &entry{record: rec, timer: timer})
}

func cancelTimer(exportID string) {
	if v, ok := records.Load(exportID); ok {
		if e, ok := v.(*entry); ok && e != nil && e.timer != nil {
			e.timer.Close()
			e.timer = nil
		}
	}
}

func removeFile(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}
