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

func GetSofieShareAppDownloadCfg(_ context.Context, _ *shareappdownloaddto.GetSofieShareAppDownloadCfgReq) (*shareappdownloaddto.GetSofieShareAppDownloadCfgRes, error) {
	row := cfgdao.LoadSofieShareAppDownloadCfg()
	item := toSofieShareCfgItem(row)
	if item == nil {
		item = &shareappdownloaddto.SofieShareAppDownloadCfgItem{ID: "0"}
	}
	return &shareappdownloaddto.GetSofieShareAppDownloadCfgRes{Cfg: item}, nil
}

func SaveSofieShareAppDownloadCfg(_ context.Context, req *shareappdownloaddto.SaveSofieShareAppDownloadCfgReq) (*shareappdownloaddto.SaveSofieShareAppDownloadCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	packageName := strings.TrimSpace(req.PackageName)
	googlePlayUrl := strings.TrimSpace(req.GooglePlayUrl)
	shareText := strings.TrimSpace(req.ShareText)
	if googlePlayUrl == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.LoadSofieShareAppDownloadCfg()
	row := &entity.SofieShareAppDownloadCfg{
		PackageName:   packageName,
		GooglePlayUrl: googlePlayUrl,
		ShareText:     shareText,
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
	if err := cfgdao.SaveSofieShareAppDownloadCfg(row); err != nil {
		return nil, err
	}
	reloadSofieShareMemory()
	return &shareappdownloaddto.SaveSofieShareAppDownloadCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toSofieShareCfgItem(cfg *entity.SofieShareAppDownloadCfg) *shareappdownloaddto.SofieShareAppDownloadCfgItem {
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	pkg := strings.TrimSpace(cfg.PackageName)
	return &shareappdownloaddto.SofieShareAppDownloadCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		PackageName:      pkg,
		GooglePlayUrl:    strings.TrimSpace(cfg.GooglePlayUrl),
		ShortUrl:         BuildSofieShortUrl(pkg),
		ShareText:        strings.TrimSpace(cfg.ShareText),
		PreviewShareCode: BuildSofiePreviewShareCode(pkg),
		CreatedAt:        formatTime(cfg.CreatedAt),
		UpdatedAt:        formatTime(cfg.UpdatedAt),
	}
}