package activity

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/errercode"
)

func GetFirstRechargeActivityCfg(_ context.Context, _ *activitydto.GetFirstRechargeActivityCfgReq) (*activitydto.GetFirstRechargeActivityCfgRes, error) {
	return &activitydto.GetFirstRechargeActivityCfgRes{Cfg: toCfgItemFromSnapshot(getCfgCache())}, nil
}

func SaveFirstRechargeActivityCfg(_ context.Context, req *activitydto.SaveFirstRechargeActivityCfgReq) (*activitydto.SaveFirstRechargeActivityCfgRes, error) {
	existing := cfgdao.LoadFirstRechargeActivityCfg()
	row := &activityentity.FirstRechargeActivityCfg{
		Enabled:           req.Enabled,
		Icon:              req.Icon,
		TitleEn:           req.TitleEn,
		TitleEs:           req.TitleEs,
		TitlePt:           req.TitlePt,
		TitleHi:           req.TitleHi,
		TitleId:           req.TitleId,
		RechargeBtnTextEn: req.RechargeBtnTextEn,
		RechargeBtnTextEs: req.RechargeBtnTextEs,
		RechargeBtnTextPt: req.RechargeBtnTextPt,
		RechargeBtnTextHi: req.RechargeBtnTextHi,
		RechargeBtnTextId: req.RechargeBtnTextId,
		FirstRechargeRatio: normalizeFirstRechargeRatio(req.FirstRechargeRatio),
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
	if err := cfgdao.SaveFirstRechargeActivityCfg(row); err != nil {
		return nil, err
	}
	if err := cfgdao.ReplaceAllFirstRechargeActivityPrivileges(privilegeRowsFromDTO(req.Privileges)); err != nil {
		return nil, err
	}
	reloadCfgMemory()
	return &activitydto.SaveFirstRechargeActivityCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func defaultCfgItem() *activitydto.FirstRechargeActivityCfgItem {
	return &activitydto.FirstRechargeActivityCfgItem{
		FirstRechargeRatio: defaultFirstRechargeRatio,
		Privileges:        []*activitydto.FirstRechargePrivilegeItem{},
	}
}

func formatUintID(id uint64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatUint(id, 10)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
