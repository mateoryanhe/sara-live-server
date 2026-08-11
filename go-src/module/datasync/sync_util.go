package datasync

import (
	"encoding/base64"
	"fmt"
	"strings"

	"xr-game-server/dto/datasyncdto"
	"xr-game-server/module/upload"
)

func isSyncableFileName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	return !strings.HasPrefix(name, "http://") && !strings.HasPrefix(name, "https://")
}

func appendUniqueFileName(seen map[string]struct{}, names []string, name string) []string {
	name = strings.TrimSpace(name)
	if !isSyncableFileName(name) {
		return names
	}
	if _, ok := seen[name]; ok {
		return names
	}
	seen[name] = struct{}{}
	return append(names, name)
}

func buildSyncFiles(fileNames []string) ([]*datasyncdto.SyncFileItem, error) {
	files := make([]*datasyncdto.SyncFileItem, 0, len(fileNames))
	for _, name := range fileNames {
		content, err := upload.ReadUploadedFileBytes(name)
		if err != nil {
			return nil, fmt.Errorf("read resource file %s: %w", name, err)
		}
		files = append(files, &datasyncdto.SyncFileItem{
			Name: name,
			Data: base64.StdEncoding.EncodeToString(content),
		})
	}
	return files, nil
}

func saveSyncFiles(files []*datasyncdto.SyncFileItem) (int, error) {
	fileCount := 0
	for _, item := range files {
		if item == nil || strings.TrimSpace(item.Name) == "" || item.Data == "" {
			continue
		}
		content, err := base64.StdEncoding.DecodeString(item.Data)
		if err != nil {
			return fileCount, fmt.Errorf("decode file %s: %w", item.Name, err)
		}
		if err := upload.SaveUploadedFileBytes(item.Name, content); err != nil {
			return fileCount, fmt.Errorf("save file %s: %w", item.Name, err)
		}
		fileCount++
	}
	return fileCount, nil
}

func newSyncBatchRes(receiveRes *datasyncdto.ReceiveBatchRes, label string) *datasyncdto.SyncBatchRes {
	return &datasyncdto.SyncBatchRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条%s、%d 个资源文件", receiveRes.RowCount, label, receiveRes.FileCount),
	}
}
