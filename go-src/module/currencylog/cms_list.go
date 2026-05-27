package currencylog

import (
	"context"
	"strconv"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/currencylogdto"
	"xr-game-server/entity"
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

func toCMSItem(v *entity.CurrencyLog) *currencylogdto.CMSCurrencyLogItem {
	if v == nil {
		return nil
	}
	item := &currencylogdto.CMSCurrencyLogItem{
		Id:         v.ID,
		UserId:     v.UserId,
		Action:     v.Action,
		Amount:     v.Amount,
		Before:     v.Before,
		After:      v.After,
		Reason:     v.Reason,
		ReasonText: currency.Reason(v.Reason).Text(currency.LangZHCN),
		CreatedAt:  &v.CreatedAt,
	}
	if u := userinfodao.GetUserInfoByUserId(v.UserId); u != nil {
		item.Nickname = u.Nickname
	}
	return item
}

// GetCMSList CMS分页查询货币流水
func GetCMSList(_ context.Context, req *currencylogdto.CMSCurrencyLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := currencylogdao.CMSList(&currencylogdao.CMSListFilter{
		UserId:       parseUserIdFilter(req.UserId),
		CurrencyType: req.CurrencyType,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	list := make([]*currencylogdto.CMSCurrencyLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
