package livecfg

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/livecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetLiveCfg(_ context.Context, _ *livecfgdto.GetLiveCfgReq) (*livecfgdto.GetLiveCfgRes, error) {
	meta := cfgdao.LoadLiveCfg()
	if meta == nil {
		return &livecfgdto.GetLiveCfgRes{Cfg: nil}, nil
	}
	return &livecfgdto.GetLiveCfgRes{Cfg: liveCfgItemFromMemory(meta)}, nil
}

func SaveLiveCfg(_ context.Context, req *livecfgdto.SaveLiveCfgReq) (*livecfgdto.SaveLiveCfgRes, error) {
	existing := cfgdao.LoadLiveCfg()
	videoCallTicketEnabled := true
	showcaseLiveRoomEnabled := false
	gameLiveRoomEnabled := false
	voiceChatLiveRoomEnabled := false
	oneToOneLiveRoomEnabled := false
	if existing != nil {
		videoCallTicketEnabled = existing.VideoCallTicketEnabled
		showcaseLiveRoomEnabled = existing.ShowcaseLiveRoomEnabled
		gameLiveRoomEnabled = existing.GameLiveRoomEnabled
		voiceChatLiveRoomEnabled = existing.VoiceChatLiveRoomEnabled
		oneToOneLiveRoomEnabled = existing.OneToOneLiveRoomEnabled
	}
	if req.VideoCallTicketEnabled != nil {
		videoCallTicketEnabled = *req.VideoCallTicketEnabled
	}
	if req.ShowcaseLiveRoomEnabled != nil {
		showcaseLiveRoomEnabled = *req.ShowcaseLiveRoomEnabled
	}
	if req.GameLiveRoomEnabled != nil {
		gameLiveRoomEnabled = *req.GameLiveRoomEnabled
	}
	if req.VoiceChatLiveRoomEnabled != nil {
		voiceChatLiveRoomEnabled = *req.VoiceChatLiveRoomEnabled
	}
	if req.OneToOneLiveRoomEnabled != nil {
		oneToOneLiveRoomEnabled = *req.OneToOneLiveRoomEnabled
	}
	row := &entity.LiveCfg{
		PaidDanmakuPrice:                req.PaidDanmakuPrice,
		VideoCallTicketEnabled:          videoCallTicketEnabled,
		ShowcaseLiveRoomEnabled:         showcaseLiveRoomEnabled,
		GameLiveRoomEnabled:             gameLiveRoomEnabled,
		VoiceChatLiveRoomEnabled:        voiceChatLiveRoomEnabled,
		OneToOneLiveRoomEnabled:         oneToOneLiveRoomEnabled,
		AudienceListRefreshSeconds:      NormalizeAudienceListRefreshSeconds(req.AudienceListRefreshSeconds),
		OneToOneDailyFreeSeconds:        NormalizeOneToOneDailyFreeSeconds(req.OneToOneDailyFreeSeconds),
		LiveRoomCallMinTotalRechargeUsd: NormalizeLiveRoomCallMinTotalRechargeUsd(req.LiveRoomCallMinTotalRechargeUsd),
	}
	if existing != nil {
		row.LuckyGiftRtp = existing.LuckyGiftRtp
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveLiveCfg(row); err != nil {
		return nil, err
	}
	// 入库后刷 live_cfgs 内存快照，App /cfg 与业务侧只读内存
	reloadLiveCfgMemory()
	return &livecfgdto.SaveLiveCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toLiveCfgItem(cfg *entity.LiveCfg) *livecfgdto.LiveCfgItem {
	if cfg == nil {
		return nil
	}
	return liveCfgItemFromMemory(cfg)
}

// liveCfgItemFromMemory CMS 展示与 App /cfg 一致：业务字段读内存，id/时间戳用库表元数据。
func liveCfgItemFromMemory(meta *entity.LiveCfg) *livecfgdto.LiveCfgItem {
	if meta == nil {
		return nil
	}
	return &livecfgdto.LiveCfgItem{
		ID:                              strconv.FormatUint(meta.ID, 10),
		PaidDanmakuPrice:                GetPaidDanmakuPrice(),
		VideoCallTicketEnabled:          IsVideoCallTicketEnabled(),
		AudienceListRefreshSeconds:      GetAudienceListRefreshSeconds(),
		OneToOneDailyFreeSeconds:        GetOneToOneDailyFreeSeconds(),
		LiveRoomCallMinTotalRechargeUsd: GetLiveRoomCallMinTotalRechargeUSD(),
		ShowcaseLiveRoomEnabled:         IsShowcaseLiveRoomEnabled(),
		GameLiveRoomEnabled:             IsGameLiveRoomEnabled(),
		VoiceChatLiveRoomEnabled:        IsVoiceChatLiveRoomEnabled(),
		OneToOneLiveRoomEnabled:         IsOneToOneLiveRoomEnabled(),
		CreatedAt:                       formatLiveCfgTime(meta.CreatedAt),
		UpdatedAt:                       formatLiveCfgTime(meta.UpdatedAt),
	}
}

func formatLiveCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
