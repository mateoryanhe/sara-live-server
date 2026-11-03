package channeldownloadlandingdeploy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/channeldownloadlandingdeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetChannelDownloadLandingDeployInfo(ctx context.Context, _ *channeldownloadlandingdeploydto.GetChannelDownloadLandingDeployInfoReq) (*channeldownloadlandingdeploydto.GetChannelDownloadLandingDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &channeldownloadlandingdeploydto.GetChannelDownloadLandingDeployInfoRes{
		Info: &channeldownloadlandingdeploydto.ChannelDownloadLandingDeployInfoItem{
			ID:           formatUintID(snap.ID),
			Domain:       cfg.GetStaticSiteDomains(channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix),
			UrlPrefix:    channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			UpdatedAt:    snap.UpdatedAt,
			LastUploadAt: domainsite.GetLastUploadAt(ctx, channeldownloadlandingdeploydto.ChannelDownloadLandingSiteKey),
		},
	}, nil
}

func SaveChannelDownloadLandingDeployCfg(ctx context.Context, req *channeldownloadlandingdeploydto.SaveChannelDownloadLandingDeployCfgReq) (*channeldownloadlandingdeploydto.SaveChannelDownloadLandingDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadChannelDownloadLandingDeployCfg()
	row := &entity.ChannelDownloadLandingDeployCfg{}
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
		channeldownloadlandingdeploydto.ChannelDownloadLandingSiteKey,
		channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	if err := cfgdao.SaveChannelDownloadLandingDeployCfg(row); err != nil {
		return nil, err
	}
	ReloadChannelDownloadLandingDeployCache()
	return &channeldownloadlandingdeploydto.SaveChannelDownloadLandingDeployCfgRes{
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
