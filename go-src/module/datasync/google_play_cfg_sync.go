package datasync

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/module/recharge"
)

// SyncGooglePlayCfg 从当前环境读取 Google Play 配置并推送到目标环境.
func SyncGooglePlayCfg(_ context.Context, _ *datasyncdto.SyncGooglePlayCfgReq) (*datasyncdto.SyncBatchRes, error) {
	row := cfgdao.LoadGooglePlayCfg()
	if row == nil || row.ID == 0 {
		return nil, errInvalidParam()
	}
	if strings.TrimSpace(row.ServiceAccountJson) == "" || !json.Valid([]byte(row.ServiceAccountJson)) {
		return nil, errInvalidParam()
	}
	payload := &datasyncdto.ReceiveGooglePlayCfgReq{Row: row}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveGooglePlayCfg", payload, &receiveRes); err != nil {
		return nil, err
	}
	return &datasyncdto.SyncBatchRes{
		Success:  receiveRes.Success,
		RowCount: receiveRes.RowCount,
		Message:  fmt.Sprintf("已同步 Google Play 配置(%d 条)", receiveRes.RowCount),
	}, nil
}

// ReceiveGooglePlayCfg 接收 Google Play 配置同步,写入库并刷新内存/客户端.
func ReceiveGooglePlayCfg(_ context.Context, req *datasyncdto.ReceiveGooglePlayCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil || req.Row == nil || req.Row.ID == 0 {
		return nil, errInvalidParam()
	}
	row := req.Row
	if strings.TrimSpace(row.ServiceAccountJson) == "" || !json.Valid([]byte(row.ServiceAccountJson)) {
		return nil, errInvalidParam()
	}
	now := time.Now()
	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = now
	}
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if err := cfgdao.SaveGooglePlayCfg(row); err != nil {
		return nil, fmt.Errorf("save google play cfg id=%d: %w", row.ID, err)
	}
	cfgdao.ReloadGooglePlayCfgCache()
	recharge.ResetGooglePlayClient()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: 1}, nil
}
