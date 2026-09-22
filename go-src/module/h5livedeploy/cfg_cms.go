package h5livedeploy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/h5livedeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetH5LiveDeployInfo(ctx context.Context, _ *h5livedeploydto.GetH5LiveDeployInfoReq) (*h5livedeploydto.GetH5LiveDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &h5livedeploydto.GetH5LiveDeployInfoRes{
		Info: &h5livedeploydto.H5LiveDeployInfoItem{
			ID:           formatUintID(snap.ID),
			Domain:       cfg.GetStaticSiteDomains(h5livedeploydto.H5LiveStaticPrefix),
			UrlPrefix:    h5livedeploydto.H5LiveStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			DeploySecret: snap.DeploySecret,
			UpdatedAt:    snap.UpdatedAt,
			LastUploadAt: domainsite.GetLastUploadAt(ctx, h5livedeploydto.H5LiveSiteKey),
		},
	}, nil
}

func SaveH5LiveDeployCfg(ctx context.Context, req *h5livedeploydto.SaveH5LiveDeployCfgReq) (*h5livedeploydto.SaveH5LiveDeployCfgRes, error) {
	secret := strings.TrimSpace(req.DeploySecret)
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if secret == "" || domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadH5LiveDeployCfg()
	row := &entity.H5LiveDeployCfg{
		DeploySecret: secret,
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
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		h5livedeploydto.H5LiveSiteKey,
		h5livedeploydto.H5LiveStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	if err := cfgdao.SaveH5LiveDeployCfg(row); err != nil {
		return nil, err
	}
	ReloadH5LiveDeployCache()
	return &h5livedeploydto.SaveH5LiveDeployCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func formatUintID(id uint64) string {
	if id == 0 {
		return "0"
	}
	return strconv.FormatUint(id, 10)
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
