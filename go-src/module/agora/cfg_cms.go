package agora

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetAgoraCfg(_ context.Context, _ *agoradto.GetAgoraCfgReq) (*agoradto.GetAgoraCfgRes, error) {
	cfg := cfgdao.LoadAgoraCfg()
	if cfg == nil {
		return &agoradto.GetAgoraCfgRes{Cfg: nil}, nil
	}
	return &agoradto.GetAgoraCfgRes{Cfg: toAgoraCfgItem(cfg)}, nil
}

func SaveAgoraCfg(_ context.Context, req *agoradto.SaveAgoraCfgReq) (*agoradto.SaveAgoraCfgRes, error) {
	if !isValidAgoraTokenCfg(req.TokenExpireSeconds, req.TokenRefreshSeconds) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadAgoraCfg()
	row := &entity.AgoraCfg{
		AppId:               req.AppId,
		AppCertificate:      req.AppCertificate,
		RestCustomerId:      req.RestCustomerId,
		RestCustomerSecret:  req.RestCustomerSecret,
		CloudPlayerRegion:   normalizeCloudPlayerRegion(req.CloudPlayerRegion),
		TokenExpireSeconds:  req.TokenExpireSeconds,
		TokenRefreshSeconds: req.TokenRefreshSeconds,
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
	if err := cfgdao.SaveAgoraCfg(row); err != nil {
		return nil, err
	}
	if isAgoraTokenCfgChanged(existing, row) {
		clearSubscriberTokenCache()
	}
	reloadAgoraCfgMemory()
	return &agoradto.SaveAgoraCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toAgoraCfgItem(cfg *entity.AgoraCfg) *agoradto.AgoraCfgItem {
	if cfg == nil {
		return nil
	}
	expireSeconds, refreshSeconds := normalizeAgoraTokenCfg(cfg.TokenExpireSeconds, cfg.TokenRefreshSeconds)
	return &agoradto.AgoraCfgItem{
		ID:                  strconv.FormatUint(cfg.ID, 10),
		AppId:               cfg.AppId,
		AppCertificate:      cfg.AppCertificate,
		RestCustomerId:      cfg.RestCustomerId,
		RestCustomerSecret:  cfg.RestCustomerSecret,
		CloudPlayerRegion:   normalizeCloudPlayerRegion(cfg.CloudPlayerRegion),
		TokenExpireSeconds:  expireSeconds,
		TokenRefreshSeconds: refreshSeconds,
		CreatedAt:           formatAgoraCfgTime(cfg.CreatedAt),
		UpdatedAt:           formatAgoraCfgTime(cfg.UpdatedAt),
	}
}

func formatAgoraCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func isAgoraTokenCfgChanged(before, after *entity.AgoraCfg) bool {
	if before == nil || after == nil {
		return false
	}
	return before.AppId != after.AppId ||
		before.AppCertificate != after.AppCertificate ||
		before.TokenExpireSeconds != after.TokenExpireSeconds ||
		before.TokenRefreshSeconds != after.TokenRefreshSeconds
}
