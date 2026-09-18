package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/firebasedto"
)

func GetFirebaseClientCfgForApp(ctx context.Context, _ *firebasedto.GetFirebaseClientCfgForAppReq) (*firebasedto.GetFirebaseClientCfgForAppRes, error) {
	cfg := cfgdao.GetFirebaseCfgCached()
	if cfg == nil || strings.TrimSpace(cfg.ClientConfigJson) == "" {
		return &firebasedto.GetFirebaseClientCfgForAppRes{}, nil
	}
	clientCfg, err := parseFirebaseClientConfig(cfg.ProjectId, cfg.ClientConfigJson)
	if err != nil {
		g.Log().Errorf(ctx, "cached firebase client config is invalid: %v", err)
		return &firebasedto.GetFirebaseClientCfgForAppRes{}, nil
	}
	return &firebasedto.GetFirebaseClientCfgForAppRes{Cfg: clientCfg}, nil
}

func parseFirebaseClientConfig(projectId, raw string) (*firebasedto.FirebaseClientCfgForApp, error) {
	projectId = strings.TrimSpace(projectId)
	raw = strings.TrimSpace(raw)
	if projectId == "" || raw == "" {
		return nil, fmt.Errorf("firebase project id and client config are required")
	}
	var cfg firebasedto.FirebaseClientCfgForApp
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, fmt.Errorf("parse firebase client config json: %w", err)
	}
	cfg.ApiKey = strings.TrimSpace(cfg.ApiKey)
	cfg.AuthDomain = strings.TrimSpace(cfg.AuthDomain)
	cfg.ProjectId = strings.TrimSpace(cfg.ProjectId)
	cfg.StorageBucket = strings.TrimSpace(cfg.StorageBucket)
	cfg.MessagingSenderId = strings.TrimSpace(cfg.MessagingSenderId)
	cfg.AppId = strings.TrimSpace(cfg.AppId)
	cfg.MeasurementId = strings.TrimSpace(cfg.MeasurementId)
	cfg.DatabaseUrl = strings.TrimSpace(cfg.DatabaseUrl)
	if cfg.ApiKey == "" || cfg.ProjectId == "" || cfg.AppId == "" {
		return nil, fmt.Errorf("firebase client config requires apiKey, projectId and appId")
	}
	if cfg.ProjectId != projectId {
		return nil, fmt.Errorf("firebase client project id does not match server project id")
	}
	return &cfg, nil
}
