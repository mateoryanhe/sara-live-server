package game

import (
	"context"
	"strconv"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/gamewindao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gamewindto"
)

func parseWinCMSUserIdFilter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func toCMSWinLogItem(v *gamewindao.CMSListItem) *gamewindto.CMSGameWinLogItem {
	if v == nil {
		return nil
	}
	item := &gamewindto.CMSGameWinLogItem{
		Id:           v.ID,
		UserId:       v.UserId,
		GameCode:     v.GameCode,
		NameEn:       v.NameEn,
		Cover:        BuildGameCoverUrl(v.Cover),
		Amount:       v.Amount,
		BetAmount:    v.BetAmount,
		PlatformType: v.PlatformType,
		OrderId:      v.OrderId,
		CreatedAt:    &v.CreatedAt,
	}
	if u := userinfodao.GetUserInfoByUserId(v.UserId); u != nil {
		item.Nickname = u.Nickname
	}
	return item
}

// GetWinCMSList CMS 分页查询游戏派彩记录
func GetWinCMSList(_ context.Context, req *gamewindto.CMSGameWinLogListReq) (*gamewindto.CMSGameWinLogListRes, error) {
	if req == nil {
		return &gamewindto.CMSGameWinLogListRes{
			Data:      []*gamewindto.CMSGameWinLogItem{},
			Platforms: cfgdao.ListDistinctVendorGamePlatformsUpper(),
		}, nil
	}
	total, rows := gamewindao.CMSList(&gamewindao.CMSListFilter{
		UserId:       parseWinCMSUserIdFilter(req.UserId),
		GameCode:     req.GameCode,
		GameCodes:    req.GameCodes,
		OrderId:      req.OrderId,
		PlatformType: req.PlatformType,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	list := make([]*gamewindto.CMSGameWinLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSWinLogItem(row))
	}
	return &gamewindto.CMSGameWinLogListRes{
		Total:     total,
		Data:      list,
		Platforms: cfgdao.ListDistinctVendorGamePlatformsUpper(),
	}, nil
}
