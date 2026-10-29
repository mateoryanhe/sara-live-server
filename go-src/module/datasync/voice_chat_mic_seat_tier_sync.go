package datasync

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/module/voicechatmicseattier"
)

// SyncVoiceChatMicSeatTier 从当前环境读取全部语聊麦位挡位并推送到目标环境.
func SyncVoiceChatMicSeatTier(_ context.Context, _ *datasyncdto.SyncVoiceChatMicSeatTierReq) (*datasyncdto.SyncBatchRes, error) {
	rows := cfgdao.GetAllVoiceChatMicSeatTiersForSync()
	payload := &datasyncdto.ReceiveVoiceChatMicSeatTierReq{Rows: rows}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveVoiceChatMicSeatTier", payload, &receiveRes); err != nil {
		return nil, err
	}
	return &datasyncdto.SyncBatchRes{
		Success:  receiveRes.Success,
		RowCount: receiveRes.RowCount,
		Message:  fmt.Sprintf("已同步语聊麦位挡位(%d 条)", receiveRes.RowCount),
	}, nil
}

// ReceiveVoiceChatMicSeatTier 接收语聊麦位挡位同步,按主键 upsert 并删除目标环境多余记录,再刷新缓存.
func ReceiveVoiceChatMicSeatTier(_ context.Context, req *datasyncdto.ReceiveVoiceChatMicSeatTierReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}
	rowCount, err := saveVoiceChatMicSeatTiers(req.Rows)
	if err != nil {
		return nil, err
	}
	voicechatmicseattier.ReloadTierCache()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: rowCount}, nil
}

func saveVoiceChatMicSeatTiers(rows []*liveentity.VoiceChatMicSeatTier) (int, error) {
	now := time.Now()
	keepIDs := make([]uint64, 0, len(rows))
	saved := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		if row.UpdatedAt.IsZero() {
			row.UpdatedAt = now
		}
		if _, err := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier)).Save(row); err != nil {
			return saved, fmt.Errorf("save voice chat mic seat tier id=%d: %w", row.ID, err)
		}
		keepIDs = append(keepIDs, row.ID)
		saved++
	}
	model := g.DB().Model(string(liveentity.TbVoiceChatMicSeatTier))
	var err error
	if len(keepIDs) == 0 {
		_, err = model.Where("1=1").Delete()
	} else {
		_, err = model.WhereNotIn("id", keepIDs).Delete()
	}
	if err != nil {
		return saved, err
	}
	return saved, nil
}
