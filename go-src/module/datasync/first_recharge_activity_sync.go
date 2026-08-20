package datasync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/module/activity"
	"xr-game-server/module/upload"
)

// SyncFirstRechargeActivityCfg 从当前环境读取首充活动配置及资源文件,推送到目标环境
func SyncFirstRechargeActivityCfg(_ context.Context, _ *datasyncdto.SyncFirstRechargeActivityCfgReq) (*datasyncdto.SyncFirstRechargeActivityCfgRes, error) {
	row := cfgdao.LoadFirstRechargeActivityCfg()
	if row == nil || row.ID == 0 {
		return nil, errInvalidParam()
	}

	fileNames := collectFirstRechargeActivityFileNames(row)
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

	payload := &datasyncdto.ReceiveFirstRechargeActivityCfgReq{
		Row:   row,
		Files: files,
	}
	var receiveRes datasyncdto.ReceiveFirstRechargeActivityCfgRes
	if err := postSyncReceive("/dataSync/receiveFirstRechargeActivityCfg", payload, &receiveRes); err != nil {
		return nil, err
	}

	return &datasyncdto.SyncFirstRechargeActivityCfgRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步首充活动配置、%d 个资源文件", receiveRes.FileCount),
	}, nil
}

// ReceiveFirstRechargeActivityCfg 接收首充活动配置同步,upsert 并刷新缓存
func ReceiveFirstRechargeActivityCfg(_ context.Context, req *datasyncdto.ReceiveFirstRechargeActivityCfgReq) (*datasyncdto.ReceiveFirstRechargeActivityCfgRes, error) {
	if req == nil || req.Row == nil || req.Row.ID == 0 {
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

	if err := cfgdao.SaveFirstRechargeActivityCfg(req.Row); err != nil {
		return nil, fmt.Errorf("save first recharge activity cfg id=%d: %w", req.Row.ID, err)
	}
	activity.ReloadFirstRechargeActivityCache()

	return &datasyncdto.ReceiveFirstRechargeActivityCfgRes{
		Success:   true,
		RowCount:  1,
		FileCount: fileCount,
	}, nil
}

func collectFirstRechargeActivityFileNames(row *activityentity.FirstRechargeActivityCfg) []string {
	if row == nil {
		return nil
	}
	seen := make(map[string]struct{})
	names := make([]string, 0)
	appendFileName := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" || strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
			return
		}
		if _, ok := seen[name]; ok {
			return
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	appendFileName(row.Icon)

	if row.Privileges != "" {
		var privileges []activityentity.FirstRechargePrivilege
		if err := json.Unmarshal([]byte(row.Privileges), &privileges); err == nil {
			for _, item := range privileges {
				appendFileName(item.Icon)
			}
		}
	}
	return names
}
