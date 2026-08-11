package datasync

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

// SyncActivityMessage 从当前环境读取指定活动消息及资源文件,推送到目标环境
func SyncActivityMessage(_ context.Context, req *datasyncdto.SyncActivityMessageReq) (*datasyncdto.SyncActivityMessageRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}

	rows := messagedao.GetActivityMessagesByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}

	fileNames := collectActivityMessageFileNames(rows)
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

	payload := &datasyncdto.ReceiveActivityMessageReq{
		Rows:  rows,
		Files: files,
	}
	var receiveRes datasyncdto.ReceiveActivityMessageRes
	if err := postSyncReceive("/dataSync/receiveActivityMessage", payload, &receiveRes); err != nil {
		return nil, err
	}

	return &datasyncdto.SyncActivityMessageRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条活动消息、%d 个资源文件", receiveRes.RowCount, receiveRes.FileCount),
	}, nil
}

// ReceiveActivityMessage 接收活动消息同步,按主键 upsert 并刷新缓存
func ReceiveActivityMessage(_ context.Context, req *datasyncdto.ReceiveActivityMessageReq) (*datasyncdto.ReceiveActivityMessageRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}

	fileCount := 0
	for _, item := range req.Files {
		if item == nil || strings.TrimSpace(item.Name) == "" || item.Data == "" {
			continue
		}
		content, err := base64.StdEncoding.DecodeString(item.Data)
		if err != nil {
			return nil, fmt.Errorf("decode file %s: %w", item.Name, err)
		}
		if err := upload.SaveUploadedFileBytes(item.Name, content); err != nil {
			return nil, fmt.Errorf("save file %s: %w", item.Name, err)
		}
		fileCount++
	}

	rowCount := 0
	for _, row := range req.Rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if _, err := g.DB().Model(string(entity.TbActivityMessage)).Save(row); err != nil {
			return nil, fmt.Errorf("save activity message id=%d: %w", row.ID, err)
		}
		rowCount++
	}

	messagedao.ReloadActivityMessageCaches()

	return &datasyncdto.ReceiveActivityMessageRes{
		Success:   true,
		RowCount:  rowCount,
		FileCount: fileCount,
	}, nil
}

func collectActivityMessageFileNames(rows []*entity.ActivityMessage) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		for _, name := range []string{
			row.IconEn, row.IconEs, row.IconPt, row.IconHi,
			row.BgEn, row.BgEs, row.BgPt, row.BgHi,
		} {
			name = strings.TrimSpace(name)
			if name == "" || strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
				continue
			}
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
	}
	return names
}
