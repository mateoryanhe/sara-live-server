package datasync

import (
	"context"
	"fmt"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
)

// SyncCfEmailCfg 从当前环境读取邮件 SMTP 配置(含密码)并推送到目标环境.
func SyncCfEmailCfg(_ context.Context, _ *datasyncdto.SyncCfEmailCfgReq) (*datasyncdto.SyncBatchRes, error) {
	row := cfgdao.LoadCfEmailCfgForDataSync()
	if row == nil || row.ID == 0 {
		return nil, errInvalidParam()
	}
	if strings.TrimSpace(row.SmtpHost) == "" || strings.TrimSpace(row.SmtpUsername) == "" || strings.TrimSpace(row.FromEmail) == "" {
		return nil, errInvalidParam()
	}
	if strings.TrimSpace(row.SmtpPassword) == "" {
		return nil, errInvalidParam()
	}
	payload := &datasyncdto.ReceiveCfEmailCfgReq{Row: row}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveCfEmailCfg", payload, &receiveRes); err != nil {
		return nil, err
	}
	return &datasyncdto.SyncBatchRes{
		Success:  receiveRes.Success,
		RowCount: receiveRes.RowCount,
		Message:  fmt.Sprintf("已同步邮件 SMTP 发信配置(%d 条)", receiveRes.RowCount),
	}, nil
}

// ReceiveCfEmailCfg 接收邮件 SMTP 配置同步,写入库并刷新内存缓存.
func ReceiveCfEmailCfg(_ context.Context, req *datasyncdto.ReceiveCfEmailCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil || req.Row == nil {
		return nil, errInvalidParam()
	}
	row := req.Row
	if strings.TrimSpace(row.SmtpHost) == "" || strings.TrimSpace(row.SmtpUsername) == "" || strings.TrimSpace(row.FromEmail) == "" {
		return nil, errInvalidParam()
	}
	if strings.TrimSpace(row.SmtpPassword) == "" {
		return nil, errInvalidParam()
	}
	if row.SmtpPort <= 0 {
		row.SmtpPort = 587
	}
	now := time.Now()
	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = now
	}
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if err := cfgdao.SaveCfEmailCfg(row); err != nil {
		return nil, fmt.Errorf("save cf email cfg id=%d: %w", row.ID, err)
	}
	cfgdao.ReloadCfEmailCfgCache()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: 1}, nil
}
