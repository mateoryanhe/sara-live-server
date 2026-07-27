package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/upload"
)

// QueryAnchorList CMS分页查询主播列表
func QueryAnchorList(_ context.Context, req *accountdto.QueryAnchorListReq) (*httpserver.CMSQueryResp, error) {
	total, data := accountdao.GetAnchorList(req)
	for _, val := range data {
		if val == nil {
			continue
		}
		val.Avatar = upload.ResolveAvatarUrlForUser(val.ID, val.Avatar)
		guildId := val.GuildId
		if accountCache, _, _ := accountdao.FindAccountInCacheByID(val.ID); accountCache != nil {
			val.IP = accountCache.IP
		}
		if userInfoCache := userinfodao.GetUserInfoByUserId(val.ID); userInfoCache != nil {
			val.Nickname = userInfoCache.Nickname
			val.Phone = userInfoCache.Phone
			val.GuildId = userInfoCache.GuildId
			guildId = userInfoCache.GuildId
		}
		room := liveroomdao.GetRoomByAnchor(val.ID)
		if room == nil {
			room = liveroom.EnsureAnchorRoom(val.ID, guildId)
		}
		fillAnchorRoomFields(val, room)
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}

func fillAnchorRoomFields(item *accountdto.AnchorListItem, room *entity.LiveRoom) {
	if item == nil || room == nil {
		return
	}
	item.RoomTitle = room.Title
	item.RoomId = room.ID
	item.Category = room.Category
	item.PrivateInviteType = room.PrivateInviteType
	item.Ticket = room.Ticket
	item.Billing = room.Billing
	item.TotalIncome = room.TotalIncome
	item.TotalGiftIncome = room.TotalGiftIncome
	item.TotalPaidDanmakuIncome = room.TotalPaidDanmakuIncome
	item.TotalPrivateRoomTicketIncome = room.TotalPrivateRoomTicketIncome
	item.TotalPrivateRoomWatchIncome = room.TotalPrivateRoomWatchIncome
	item.TotalVideoCallIncome = room.TotalVideoCallIncome
	item.TotalVideoCallTicketIncome = room.TotalVideoCallTicketIncome
	item.TotalVideoCallBillingIncome = room.TotalVideoCallBillingIncome
	if room.LiveRecordId > 0 {
		item.LiveStatus = 1
	} else {
		item.LiveStatus = 0
	}
	item.BanApplyTime = room.BanApplyTime
	item.BanReason = room.BanReason
	item.Ban = liveroom.IsRoomBanned(room)
}
