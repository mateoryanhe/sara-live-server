package datasync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
	"xr-game-server/module/vip"
)

const syncHTTPTimeout = 120 * time.Second

type syncAPIResponse struct {
	Code int                           `json:"code"`
	Data *datasyncdto.ReceiveVipCfgRes `json:"data"`
}

// SyncVipCfg 从当前环境读取指定 VIP 配置及资源文件,推送到目标环境
func SyncVipCfg(_ context.Context, req *datasyncdto.SyncVipCfgReq) (*datasyncdto.SyncVipCfgRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}
	cfg := cfgdao.GetDataSyncCfg()
	if cfg == nil {
		return nil, errInvalidParam()
	}
	targetBase := strings.TrimRight(strings.TrimSpace(cfg.TargetApiBase), "/")
	syncToken := strings.TrimSpace(cfg.Token)
	if targetBase == "" || syncToken == "" {
		return nil, errInvalidParam()
	}

	rows := cfgdao.GetVipCfgsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	fileNames := collectVipCfgFileNames(rows)
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

	payload := &datasyncdto.ReceiveVipCfgReq{
		Rows:  rows,
		Files: files,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqCtx, cancel := context.WithTimeout(context.Background(), syncHTTPTimeout)
	defer cancel()

	client := gclient.New().SetTimeout(syncHTTPTimeout)
	client.SetHeader(HeaderDataSyncToken, syncToken)
	client.SetHeader("Content-Type", "application/json")

	url := targetBase + "/dataSync/receiveVipCfg"
	resp, err := client.Post(reqCtx, url, body)
	if err != nil {
		return nil, fmt.Errorf("sync request failed: %w", err)
	}
	defer resp.Close()

	raw := resp.ReadAll()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sync http status=%d body=%s", resp.StatusCode, string(raw))
	}

	var envelope syncAPIResponse
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("parse sync response: %w", err)
	}
	if envelope.Code != errercode.Success {
		return nil, fmt.Errorf("sync rejected code=%d", envelope.Code)
	}
	if envelope.Data == nil {
		return nil, fmt.Errorf("sync response data is empty")
	}

	return &datasyncdto.SyncVipCfgRes{
		Success:   envelope.Data.Success,
		RowCount:  envelope.Data.RowCount,
		FileCount: envelope.Data.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条配置、%d 个资源文件", envelope.Data.RowCount, envelope.Data.FileCount),
	}, nil
}

// ReceiveVipCfg 接收 VIP 配置同步,按主键 upsert 并写入资源文件
func ReceiveVipCfg(_ context.Context, req *datasyncdto.ReceiveVipCfgReq) (*datasyncdto.ReceiveVipCfgRes, error) {
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
		if err := cfgdao.CreateVipCfg(row); err != nil {
			return nil, fmt.Errorf("save vip cfg id=%d: %w", row.ID, err)
		}
		rowCount++
	}

	vip.ReloadVipCfgMemory()

	return &datasyncdto.ReceiveVipCfgRes{
		Success:   true,
		RowCount:  rowCount,
		FileCount: fileCount,
	}, nil
}

func collectVipCfgFileNames(rows []*entity.VipCfg) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		for _, name := range []string{
			row.LevelIcon,
			row.WithdrawIcon,
			row.Animation,
			row.AnimationIcon,
			row.CommentEffect,
			row.CommentEffectIcon,
			row.CustomerServiceIcon,
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
