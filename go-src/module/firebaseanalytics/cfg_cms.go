package firebaseanalytics

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/firebaseanalyticsdto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetFirebaseAnalyticsCfg(_ context.Context, _ *firebaseanalyticsdto.GetFirebaseAnalyticsCfgReq) (*firebaseanalyticsdto.GetFirebaseAnalyticsCfgRes, error) {
	cfg := cfgdao.ResolveFirebaseAnalyticsCfg()
	if cfg == nil {
		return &firebaseanalyticsdto.GetFirebaseAnalyticsCfgRes{Cfg: nil}, nil
	}
	return &firebaseanalyticsdto.GetFirebaseAnalyticsCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SaveFirebaseAnalyticsCfg(ctx context.Context, req *firebaseanalyticsdto.SaveFirebaseAnalyticsCfgReq) (*firebaseanalyticsdto.SaveFirebaseAnalyticsCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	projectId := strings.TrimSpace(req.ProjectId)
	clientConfigJson := strings.TrimSpace(req.ClientConfigJson)
	serviceAccountJson := strings.TrimSpace(req.ServiceAccountJson)
	measurementApiSecret := strings.TrimSpace(req.MeasurementApiSecret)
	if req.Enabled == 1 {
		if _, err := parseClientConfig(projectId, clientConfigJson); err != nil {
			g.Log().Warningf(ctx, "validate firebase analytics client config failed: %v", err)
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if err := validateServiceAccountJSON(ctx, projectId, serviceAccountJson); err != nil {
			g.Log().Warningf(ctx, "validate firebase analytics service account failed: %v", err)
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if measurementApiSecret == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
	}
	cfgdao.EnsureFirebaseAnalyticsCfgDao()
	existing := cfgdao.ResolveFirebaseAnalyticsCfg()
	row := &sysentity.FirebaseAnalyticsCfg{
		Enabled:            req.Enabled,
		ProjectId:          projectId,
		ClientConfigJson:   clientConfigJson,
		ServiceAccountJson:   serviceAccountJson,
		MeasurementApiSecret: measurementApiSecret,
	}
	if existing != nil && existing.ID > 0 {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	} else if req.ID > 0 {
		row.ID = req.ID
	}
	now := time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	if err := cfgdao.SaveFirebaseAnalyticsCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadFirebaseAnalyticsCfgCache()
	OnCfgSaved(row)
	return &firebaseanalyticsdto.SaveFirebaseAnalyticsCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *sysentity.FirebaseAnalyticsCfg) *firebaseanalyticsdto.FirebaseAnalyticsCfgItem {
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	return &firebaseanalyticsdto.FirebaseAnalyticsCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		Enabled:            cfg.Enabled,
		ProjectId:          cfg.ProjectId,
		ClientConfigJson:   cfg.ClientConfigJson,
		ServiceAccountJson:   cfg.ServiceAccountJson,
		MeasurementApiSecret: cfg.MeasurementApiSecret,
		CreatedAt:            formatCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
