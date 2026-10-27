package shareappdownload

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetShareAppDownloadCfg(_ context.Context, _ *shareappdownloaddto.GetShareAppDownloadCfgReq) (*shareappdownloaddto.GetShareAppDownloadCfgRes, error) {
	cfg := cfgdao.LoadShareAppDownloadCfg()
	item := toCfgItem(cfg)
	if item == nil {
		item = &shareappdownloaddto.ShareAppDownloadCfgItem{ID: "0"}
	}
	// 下载链接 / 短链 / 最终预览只读
	item.DownloadUrl = GetDownloadUrl()
	item.ShortUrl = GetApkShortUrl()
	item.PreviewShareCode = BuildPreviewShareCode()
	return &shareappdownloaddto.GetShareAppDownloadCfgRes{Cfg: item}, nil
}

func SaveShareAppDownloadCfg(_ context.Context, req *shareappdownloaddto.SaveShareAppDownloadCfgReq) (*shareappdownloaddto.SaveShareAppDownloadCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	shareText := strings.TrimSpace(req.ShareText)

	existing := cfgdao.LoadShareAppDownloadCfg()
	row := &entity.ShareAppDownloadCfg{
		ShareText: shareText,
	}
	if existing != nil {
		// 保留历史 DownloadUrl 列值(不再由本页维护)
		row.DownloadUrl = existing.DownloadUrl
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
	if err := cfgdao.SaveShareAppDownloadCfg(row); err != nil {
		return nil, err
	}
	reloadCfgMemory()
	return &shareappdownloaddto.SaveShareAppDownloadCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *entity.ShareAppDownloadCfg) *shareappdownloaddto.ShareAppDownloadCfgItem {
	if cfg == nil {
		return nil
	}
	return &shareappdownloaddto.ShareAppDownloadCfgItem{
		ID:        strconv.FormatUint(cfg.ID, 10),
		ShareText: strings.TrimSpace(cfg.ShareText),
		CreatedAt: formatTime(cfg.CreatedAt),
		UpdatedAt: formatTime(cfg.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
