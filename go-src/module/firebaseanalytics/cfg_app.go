package firebaseanalytics

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/firebaseanalyticsdto"
)

func GetClientCfgForApp(ctx context.Context, _ *firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppReq) (*firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppRes, error) {
	cfg := cfgdao.ResolveFirebaseAnalyticsCfg()
	if !isRuntimeCfgActive(cfg) {
		return &firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppRes{Enabled: 0}, nil
	}
	parsed, err := parseClientConfig(cfg.ProjectId, cfg.ClientConfigJson)
	if err != nil {
		g.Log().Errorf(ctx, "firebase analytics client config invalid in cache: %v", err)
		return &firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppRes{Enabled: 0}, nil
	}
	return &firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppRes{
		Enabled: 1,
		Cfg:     toAppClientCfg(parsed),
	}, nil
}

func toAppClientCfg(cfg *ClientConfig) *firebaseanalyticsdto.FirebaseAnalyticsClientCfgForApp {
	if cfg == nil {
		return nil
	}
	return &firebaseanalyticsdto.FirebaseAnalyticsClientCfgForApp{
		ApiKey:            cfg.ApiKey,
		AuthDomain:        cfg.AuthDomain,
		ProjectId:         cfg.ProjectId,
		StorageBucket:     cfg.StorageBucket,
		MessagingSenderId: cfg.MessagingSenderId,
		AppId:             cfg.AppId,
		MeasurementId:     cfg.MeasurementId,
		DatabaseUrl:       cfg.DatabaseUrl,
	}
}
