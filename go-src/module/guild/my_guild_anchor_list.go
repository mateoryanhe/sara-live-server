package guild

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/guilddto"
	"xr-game-server/module/liveroom"
)

// GetMyGuildAnchorList CMS分页查询当前用户作为会长管理的指定工会名下主播(含未结算总收益)
func GetMyGuildAnchorList(ctx context.Context, req *guilddto.GetMyGuildAnchorListReq) (*httpserver.CMSQueryResp, error) {
	_, err := getGuildOwnedByCMSUser(ctx, req.GuildId)
	if err != nil {
		return nil, err
	}
	resp, err := liveroom.QueryAnchorList(ctx, &accountdto.QueryAnchorListReq{
		CMSQueryReq: req.CMSQueryReq,
		GuildId:     req.GuildId,
	})
	if err != nil {
		return nil, err
	}
	items, _ := resp.Data.([]*accountdto.AnchorListItem)
	list := make([]*guilddto.MyGuildAnchorListItem, 0, len(items))
	for _, item := range items {
		if row := toMyGuildAnchorListItem(item); row != nil {
			list = append(list, row)
		}
	}
	return httpserver.NewCMSQueryResp(resp.Total, list), nil
}

func toMyGuildAnchorListItem(item *accountdto.AnchorListItem) *guilddto.MyGuildAnchorListItem {
	if item == nil || item.ID == 0 {
		return nil
	}
	unsettledTotalIncome := 0.0
	if row := liveroomdao.GetLiveRoomIncomeUnsettledForCMS(item.ID); row != nil {
		unsettledTotalIncome = row.TotalIncome
	}
	return &guilddto.MyGuildAnchorListItem{
		ID:                   strconv.FormatUint(item.ID, 10),
		Nickname:             item.Nickname,
		Avatar:               item.Avatar,
		UnsettledTotalIncome: unsettledTotalIncome,
	}
}
