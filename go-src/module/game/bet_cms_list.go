package game

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/entity"
)

func parseBetCMSUserIdFilter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func toCMSBetLogItem(v *entity.GameBetLog) *gamebetdto.CMSGameBetLogItem {
	if v == nil {
		return nil
	}
	item := &gamebetdto.CMSGameBetLogItem{
		Id:           v.ID,
		UserId:       v.UserId,
		GameCode:     v.GameCode,
		NameEn:       v.NameEn,
		Cover:        BuildGameCoverUrl(v.Cover),
		Amount:       v.Amount,
		PlatformType: v.PlatformType,
		OrderId:      v.OrderId,
		LiveRoomId:   v.LiveRoomId,
		LiveRecordId: v.LiveRecordId,
		CreatedAt:    &v.CreatedAt,
	}
	if u := userinfodao.GetUserInfoByUserId(v.UserId); u != nil {
		item.Nickname = u.Nickname
	}
	fillCMSBetLogLiveRoomFields(item, v.LiveRoomId)
	return item
}

func fillCMSBetLogLiveRoomFields(item *gamebetdto.CMSGameBetLogItem, liveRoomId uint64) {
	if item == nil || liveRoomId == 0 {
		return
	}
	if room := liveroomdao.GetRoomById(liveRoomId); room != nil {
		item.LiveRoomTitle = room.Title
	}
	if anchor := userinfodao.GetUserInfoByUserId(liveRoomId); anchor != nil {
		item.AnchorNickname = anchor.Nickname
	}
}

// GetBetCMSList CMS 分页查询游戏下注记录
func GetBetCMSList(_ context.Context, req *gamebetdto.CMSGameBetLogListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return httpserver.NewCMSQueryResp(0, []*gamebetdto.CMSGameBetLogItem{}), nil
	}
	total, rows := gamebetdao.CMSList(&gamebetdao.CMSListFilter{
		UserId:       parseBetCMSUserIdFilter(req.UserId),
		GameCode:     req.GameCode,
		OrderId:      req.OrderId,
		PlatformType: req.PlatformType,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	list := make([]*gamebetdto.CMSGameBetLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSBetLogItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
