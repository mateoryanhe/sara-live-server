package currencylog

import (
	"context"
	"strconv"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/currencylogdto"
	"xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

func parseUserIdFilter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func collectCurrencyLogUserIds(rows []*entity.CurrencyLog) []uint64 {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.UserId == 0 {
			continue
		}
		if _, ok := seen[row.UserId]; ok {
			continue
		}
		seen[row.UserId] = struct{}{}
		ids = append(ids, row.UserId)
	}
	return ids
}

func toCMSItem(v *entity.CurrencyLog, profileMap map[uint64]*entity.UserInfo) *currencylogdto.CMSCurrencyLogItem {
	if v == nil {
		return nil
	}
	item := &currencylogdto.CMSCurrencyLogItem{
		Id:               v.ID,
		UserId:           v.UserId,
		Action:           v.Action,
		Amount:           v.Amount,
		Before:           v.Before,
		After:            v.After,
		Reason:           v.Reason,
		ReasonText:       currency.Reason(v.Reason).Text(currency.LangZHCN),
		GameId:           v.GameId,
		GameName:         v.GameName,
		GameCategory:     v.GameCategory,
		BusinessType:     v.BusinessType,
		BusinessTypeText: currency.BusinessTypeText(v.BusinessType),
		CreatedAt:        &v.CreatedAt,
	}
	if profileMap != nil {
		if u := profileMap[v.UserId]; u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrlForUser(v.UserId, u.Avatar)
		}
	}
	return item
}

// GetCMSList CMS分页查询货币流水
func GetCMSList(_ context.Context, req *currencylogdto.CMSCurrencyLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := currencylogdao.CMSList(&currencylogdao.CMSListFilter{
		UserId:       parseUserIdFilter(req.UserId),
		CurrencyType: req.CurrencyType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	profileMap := userinfodao.GetUserProfileMapByUserIds(collectCurrencyLogUserIds(rows))
	list := make([]*currencylogdto.CMSCurrencyLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row, profileMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
