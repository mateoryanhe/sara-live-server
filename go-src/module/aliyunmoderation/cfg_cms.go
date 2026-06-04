package aliyunmoderation

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/aliyuntextmoderationcfgdao"
	"xr-game-server/dto/aliyuntextmoderationdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetCfg(_ context.Context, _ *aliyuntextmoderationdto.GetCfgReq) (*aliyuntextmoderationdto.GetCfgRes, error) {
	cfg := aliyuntextmoderationcfgdao.Load()
	if cfg == nil {
		return &aliyuntextmoderationdto.GetCfgRes{Cfg: nil}, nil
	}
	return &aliyuntextmoderationdto.GetCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SaveCfg(_ context.Context, req *aliyuntextmoderationdto.SaveCfgReq) (*aliyuntextmoderationdto.SaveCfgRes, error) {
	existing := aliyuntextmoderationcfgdao.Load()
	row := &entity.AliyunTextModerationCfg{
		Enabled:         req.Enabled,
		AccessKeyId:     strings.TrimSpace(req.AccessKeyId),
		AccessKeySecret: strings.TrimSpace(req.AccessKeySecret),
		RegionId:        strings.TrimSpace(req.RegionId),
		Endpoint:        strings.TrimSpace(req.Endpoint),
		ChatService:     strings.TrimSpace(req.ChatService),
		NicknameService: strings.TrimSpace(req.NicknameService),
		CommentService:  strings.TrimSpace(req.CommentService),
	}
	if row.Enabled {
		if row.AccessKeyId == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing == nil && row.AccessKeySecret == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing != nil && row.AccessKeySecret == "" {
			row.AccessKeySecret = existing.AccessKeySecret
		}
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
	if err := aliyuntextmoderationcfgdao.Save(row); err != nil {
		return nil, err
	}
	invalidateGreenClient()
	reloadCfgMemory()
	return &aliyuntextmoderationdto.SaveCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *entity.AliyunTextModerationCfg) *aliyuntextmoderationdto.CfgItem {
	if cfg == nil {
		return nil
	}
	return &aliyuntextmoderationdto.CfgItem{
		ID:              strconv.FormatUint(cfg.ID, 10),
		Enabled:         cfg.Enabled,
		AccessKeyId:     cfg.AccessKeyId,
		AccessKeySecret: maskSecret(cfg.AccessKeySecret),
		RegionId:        cfg.RegionId,
		Endpoint:        cfg.Endpoint,
		ChatService:     cfg.ChatService,
		NicknameService: cfg.NicknameService,
		CommentService:  cfg.CommentService,
		CreatedAt:       formatTime(cfg.CreatedAt),
		UpdatedAt:       formatTime(cfg.UpdatedAt),
	}
}

func maskSecret(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
