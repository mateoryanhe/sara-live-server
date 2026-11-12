package attrreport

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/attrreportdto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetAttrReportCfg(_ context.Context, _ *attrreportdto.GetAttrReportCfgReq) (*attrreportdto.GetAttrReportCfgRes, error) {
	cfg := cfgdao.ResolveAttrReportCfg()
	if cfg == nil {
		return &attrreportdto.GetAttrReportCfgRes{Cfg: nil}, nil
	}
	return &attrreportdto.GetAttrReportCfgRes{Cfg: toAttrReportCfgItem(cfg)}, nil
}

func SaveAttrReportCfg(_ context.Context, req *attrreportdto.SaveAttrReportCfgReq) (*attrreportdto.SaveAttrReportCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	baseApi := strings.TrimSpace(req.BaseApi)
	secretKey := strings.TrimSpace(req.SecretKey)
	if baseApi == "" || secretKey == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	cfgdao.EnsureAttrReportCfgDao()
	existing := cfgdao.ResolveAttrReportCfg()
	row := &sysentity.AttrReportCfg{
		BaseApi:   baseApi,
		SecretKey: secretKey,
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
	if err := cfgdao.SaveAttrReportCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadAttrReportCfgCache()
	return &attrreportdto.SaveAttrReportCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toAttrReportCfgItem(cfg *sysentity.AttrReportCfg) *attrreportdto.AttrReportCfgItem {
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	return &attrreportdto.AttrReportCfgItem{
		ID:        strconv.FormatUint(cfg.ID, 10),
		BaseApi:   cfg.BaseApi,
		SecretKey: cfg.SecretKey,
		CreatedAt: formatAttrReportTime(cfg.CreatedAt),
		UpdatedAt: formatAttrReportTime(cfg.UpdatedAt),
	}
}

func formatAttrReportTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
