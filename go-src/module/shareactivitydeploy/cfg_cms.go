package shareactivitydeploy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/shareactivitydeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetShareActivityDeployInfo(ctx context.Context, _ *shareactivitydeploydto.GetShareActivityDeployInfoReq) (*shareactivitydeploydto.GetShareActivityDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &shareactivitydeploydto.GetShareActivityDeployInfoRes{
		Info: &shareactivitydeploydto.ShareActivityDeployInfoItem{
			ID:           formatUintID(snap.ID),
			Domain:       cfg.GetStaticSiteDomains(shareactivitydeploydto.ShareActivityStaticPrefix),
			UrlPrefix:    shareactivitydeploydto.ShareActivityStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			UpdatedAt:    snap.UpdatedAt,
			LastUploadAt: domainsite.GetLastUploadAt(ctx, shareactivitydeploydto.ShareActivitySiteKey),
		},
	}, nil
}

func SaveShareActivityDeployCfg(ctx context.Context, req *shareactivitydeploydto.SaveShareActivityDeployCfgReq) (*shareactivitydeploydto.SaveShareActivityDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadShareActivityDeployCfg()
	row := &entity.ShareActivityDeployCfg{}
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
		shareactivitydeploydto.ShareActivitySiteKey,
		shareactivitydeploydto.ShareActivityStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	if err := cfgdao.SaveShareActivityDeployCfg(row); err != nil {
		return nil, err
	}
	ReloadShareActivityDeployCache()
	return &shareactivitydeploydto.SaveShareActivityDeployCfgRes{
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
