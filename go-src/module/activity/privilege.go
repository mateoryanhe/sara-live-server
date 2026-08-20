package activity

import (
	"encoding/json"

	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/module/upload"
)

func parsePrivileges(raw string) []activityentity.FirstRechargePrivilege {
	if raw == "" {
		return nil
	}
	var list []activityentity.FirstRechargePrivilege
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil
	}
	return list
}

func marshalPrivileges(list []activityentity.FirstRechargePrivilege) string {
	if len(list) == 0 {
		return ""
	}
	b, err := json.Marshal(list)
	if err != nil {
		return ""
	}
	return string(b)
}

func privilegesFromDTO(list []*activitydto.FirstRechargePrivilegeItem) []activityentity.FirstRechargePrivilege {
	if len(list) == 0 {
		return nil
	}
	out := make([]activityentity.FirstRechargePrivilege, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		icon := item.IconName
		if icon == "" {
			icon = item.Icon
		}
		out = append(out, activityentity.FirstRechargePrivilege{
			Icon:   icon,
			DescEn: item.DescEn,
			DescEs: item.DescEs,
			DescPt: item.DescPt,
			DescHi: item.DescHi,
			DescId: item.DescId,
		})
	}
	return out
}

func toPrivilegeDTOItems(list []activityentity.FirstRechargePrivilege) []*activitydto.FirstRechargePrivilegeItem {
	if len(list) == 0 {
		return []*activitydto.FirstRechargePrivilegeItem{}
	}
	out := make([]*activitydto.FirstRechargePrivilegeItem, 0, len(list))
	for _, item := range list {
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

func toPrivilegeAppItems(list []activityentity.FirstRechargePrivilege) []*activitydto.AppFirstRechargePrivilegeItem {
	if len(list) == 0 {
		return []*activitydto.AppFirstRechargePrivilegeItem{}
	}
	out := make([]*activitydto.AppFirstRechargePrivilegeItem, 0, len(list))
	for _, item := range list {
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
