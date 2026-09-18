package auth

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/firebasedto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func initFirebaseCfg() {
	cfgdao.InitFirebaseCfgDao()
	cfgdao.ReloadFirebaseCfgCache()
}

func GetFirebaseCfg(_ context.Context, _ *firebasedto.GetFirebaseCfgReq) (*firebasedto.GetFirebaseCfgRes, error) {
	cfg := cfgdao.GetFirebaseCfgCached()
	if cfg == nil {
		return &firebasedto.GetFirebaseCfgRes{Cfg: nil}, nil
	}
	return &firebasedto.GetFirebaseCfgRes{Cfg: toFirebaseCfgItem(cfg)}, nil
}

func SaveFirebaseCfg(ctx context.Context, req *firebasedto.SaveFirebaseCfgReq) (*firebasedto.SaveFirebaseCfgRes, error) {
	projectId := strings.TrimSpace(req.ProjectId)
	clientConfigJson := strings.TrimSpace(req.ClientConfigJson)
	serviceAccountJson := strings.TrimSpace(req.ServiceAccountJson)
	if _, err := parseFirebaseClientConfig(projectId, clientConfigJson); err != nil {
		g.Log().Warningf(ctx, "validate firebase client config failed: %v", err)
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := validateFirebaseConfig(ctx, projectId, serviceAccountJson); err != nil {
		g.Log().Warningf(ctx, "validate firebase config failed: %v", err)
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetFirebaseCfgCached()
	row := &sysentity.FirebaseCfg{
		ProjectId:          projectId,
		ClientConfigJson:   clientConfigJson,
		ServiceAccountJson: serviceAccountJson,
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
	if err := cfgdao.SaveFirebaseCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadFirebaseCfgCache()
	ResetFirebaseAuthClient()
	return &firebasedto.SaveFirebaseCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toFirebaseCfgItem(cfg *sysentity.FirebaseCfg) *firebasedto.FirebaseCfgItem {
	if cfg == nil {
		return nil
	}
	return &firebasedto.FirebaseCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		ProjectId:          cfg.ProjectId,
		ClientConfigJson:   cfg.ClientConfigJson,
		ServiceAccountJson: cfg.ServiceAccountJson,
		CreatedAt:          formatFirebaseTime(cfg.CreatedAt),
		UpdatedAt:          formatFirebaseTime(cfg.UpdatedAt),
	}
}

func formatFirebaseTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
