package activity

import (
	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/module/upload"
)

func privilegeRowsFromDTO(list []*activitydto.FirstRechargePrivilegeItem) []*activityentity.FirstRechargeActivityPrivilege {
	if len(list) == 0 {
		return nil
	}
	rows := make([]*activityentity.FirstRechargeActivityPrivilege, 0, len(list))
	for idx, item := range list {
		if item == nil {
			continue
		}
		icon := item.IconName
		if icon == "" {
			icon = item.Icon
		}
		sortVal := idx + 1
		rows = append(rows, &activityentity.FirstRechargeActivityPrivilege{
			Icon:   icon,
			DescEn: item.DescEn,
			DescEs: item.DescEs,
			DescPt: item.DescPt,
			DescHi: item.DescHi,
			DescId: item.DescId,
			Sort:   sortVal,
		})
	}
	return rows
}

func toPrivilegeDTOItems(list []*activityentity.FirstRechargeActivityPrivilege) []*activitydto.FirstRechargePrivilegeItem {
	if len(list) == 0 {
		return []*activitydto.FirstRechargePrivilegeItem{}
	}
	out := make([]*activitydto.FirstRechargePrivilegeItem, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		out = append(out, &activitydto.FirstRechargePrivilegeItem{
			IconName: item.Icon,
			Icon:     upload.GetUrlByName(item.Icon),
			DescEn:   item.DescEn,
			DescEs:   item.DescEs,
			DescPt:   item.DescPt,
			DescHi:   item.DescHi,
			DescId:   item.DescId,
		})
	}
	return out
}

func toPrivilegeAppItems(list []*activityentity.FirstRechargeActivityPrivilege) []*activitydto.AppFirstRechargePrivilegeItem {
	if len(list) == 0 {
		return []*activitydto.AppFirstRechargePrivilegeItem{}
	}
	out := make([]*activitydto.AppFirstRechargePrivilegeItem, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		out = append(out, &activitydto.AppFirstRechargePrivilegeItem{
			Icon:   upload.GetUrlByName(item.Icon),
			DescEn: item.DescEn,
			DescEs: item.DescEs,
			DescPt: item.DescPt,
			DescHi: item.DescHi,
			DescId: item.DescId,
		})
	}
	return out
}
