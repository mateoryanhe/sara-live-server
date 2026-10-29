package datasync

import (
	"context"
	"fmt"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/module/privacypolicy"
)

// SyncSofiePrivacyPolicyCfg 从当前环境读取 Sofie 隐私政策配置并推送到目标环境.
func SyncSofiePrivacyPolicyCfg(_ context.Context, _ *datasyncdto.SyncSofiePrivacyPolicyCfgReq) (*datasyncdto.SyncBatchRes, error) {
	row := cfgdao.LoadSofiePrivacyPolicyCfg()
	if row == nil || row.ID == 0 {
		return nil, errInvalidParam()
	}
	payload := &datasyncdto.ReceiveSofiePrivacyPolicyCfgReq{Row: row}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveSofiePrivacyPolicyCfg", payload, &receiveRes); err != nil {
		return nil, err
	}
	return &datasyncdto.SyncBatchRes{
		Success:  receiveRes.Success,
		RowCount: receiveRes.RowCount,
		Message:  fmt.Sprintf("已同步 Sofie 隐私政策配置(%d 条)", receiveRes.RowCount),
	}, nil
}

// ReceiveSofiePrivacyPolicyCfg 接收 Sofie 隐私政策配置同步,写入库并刷新内存.
func ReceiveSofiePrivacyPolicyCfg(_ context.Context, req *datasyncdto.ReceiveSofiePrivacyPolicyCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil || req.Row == nil || req.Row.ID == 0 {
		return nil, errInvalidParam()
	}
	row := req.Row
	now := time.Now()
	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = now
	}
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if err := cfgdao.SaveSofiePrivacyPolicyCfg(row); err != nil {
		return nil, fmt.Errorf("save sofie privacy policy cfg id=%d: %w", row.ID, err)
	}
	privacypolicy.ReloadSofieCfgMemory()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: 1}, nil
}
