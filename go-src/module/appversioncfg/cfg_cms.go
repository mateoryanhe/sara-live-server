package appversioncfg

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/appversioncfgdto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
)

func GetAppVersionCfg(_ context.Context, _ *appversioncfgdto.GetAppVersionCfgReq) (*appversioncfgdto.GetAppVersionCfgRes, error) {
	cfg := cfgdao.LoadAppVersionCfg()
	if cfg == nil {
		return &appversioncfgdto.GetAppVersionCfgRes{Cfg: &appversioncfgdto.AppVersionCfgItem{
			VersionQueryEnabled: false,
			Version:             "",
			BuildVersion:        "",
			DownloadUrl:         "",
			UpdateDetails:       []*appversioncfgdto.AppVersionUpdateDetailItem{},
		}}, nil
	}
	return &appversioncfgdto.GetAppVersionCfgRes{Cfg: toCfgItem(cfg, cfgdao.LoadAppVersionUpdateDetails())}, nil
}

func SaveAppVersionCfg(_ context.Context, req *appversioncfgdto.SaveAppVersionCfgReq) (*appversioncfgdto.SaveAppVersionCfgRes, error) {
	version := strings.TrimSpace(req.Version)
	buildVersion := strings.TrimSpace(req.BuildVersion)
	downloadUrl := strings.TrimSpace(req.DownloadUrl)

	existing := cfgdao.LoadAppVersionCfg()
	row := &entity.AppVersionCfg{
		VersionQueryEnabled: req.VersionQueryEnabled,
		Version:             version,
		BuildVersion:        buildVersion,
		DownloadUrl:         downloadUrl,
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
	if err := cfgdao.SaveAppVersionCfg(row); err != nil {
		return nil, err
	}
	if err := cfgdao.ReplaceAllAppVersionUpdateDetails(buildUpdateDetailRows(req.UpdateDetails)); err != nil {
		return nil, err
	}
	reloadCfgMemory()
	return &appversioncfgdto.SaveAppVersionCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func buildUpdateDetailRows(items []*appversioncfgdto.AppVersionUpdateDetailItem) []*entity.AppVersionUpdateDetail {
	if len(items) == 0 {
		return nil
	}
	rows := make([]*entity.AppVersionUpdateDetail, 0, len(items))
	for idx, item := range items {
		if item == nil {
			continue
		}
		content := strings.TrimSpace(item.Content)
		if content == "" {
			continue
		}
		sortVal := item.Sort
		if sortVal <= 0 {
			sortVal = idx + 1
		}
		rows = append(rows, &entity.AppVersionUpdateDetail{
			Content: content,
			Sort:    sortVal,
		})
	}
	return rows
}

func toCfgItem(cfg *entity.AppVersionCfg, details []*entity.AppVersionUpdateDetail) *appversioncfgdto.AppVersionCfgItem {
	if cfg == nil {
		return nil
	}
	return &appversioncfgdto.AppVersionCfgItem{
		ID:                  strconv.FormatUint(cfg.ID, 10),
		VersionQueryEnabled: cfg.VersionQueryEnabled,
		Version:             cfg.Version,
		BuildVersion:        cfg.BuildVersion,
		DownloadUrl:         cfg.DownloadUrl,
		UpdateDetails:       toUpdateDetailItems(details),
		CreatedAt:           formatTime(cfg.CreatedAt),
		UpdatedAt:           formatTime(cfg.UpdatedAt),
	}
}

func toUpdateDetailItems(details []*entity.AppVersionUpdateDetail) []*appversioncfgdto.AppVersionUpdateDetailItem {
	if len(details) == 0 {
		return []*appversioncfgdto.AppVersionUpdateDetailItem{}
	}
	items := make([]*appversioncfgdto.AppVersionUpdateDetailItem, 0, len(details))
	for _, detail := range details {
		if detail == nil || strings.TrimSpace(detail.Content) == "" {
			continue
		}
		items = append(items, &appversioncfgdto.AppVersionUpdateDetailItem{
			Content: detail.Content,
			Sort:    detail.Sort,
		})
	}
	return items
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
