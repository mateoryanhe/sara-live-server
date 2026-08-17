package liveroom

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// IsRoomOffShelf 直播间是否已下架(不在缓存或 status!=上架)
func IsRoomOffShelf(room *entity.LiveRoom) bool {
	if room == nil {
		return true
	}
	return room.Status != entity.LiveRoomStatusOnShelf
}

// SetLiveRoomStatus CMS 上架/下架主播直播间
func SetLiveRoomStatus(ctx context.Context, req *accountdto.SetLiveRoomStatusReq) (*accountdto.SetLiveRoomStatusRes, error) {
	if req.Status != entity.LiveRoomStatusOffShelf && req.Status != entity.LiveRoomStatusOnShelf {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	if req.Status == entity.LiveRoomStatusOffShelf {
		return offShelfLiveRoom(ctx, req.AnchorId)
	}
	return onShelfLiveRoom(ctx, req.AnchorId)
}

// OffShelfGuildLiveRooms 下架指定工会下全部主播间(停播+清缓存+写库)
func OffShelfGuildLiveRooms(ctx context.Context, guildId uint64) {
	if guildId == 0 {
		return
	}
	changed := false
	for _, room := range liveroomdao.ListRoomsByGuild(guildId) {
		if room == nil || room.ID == 0 {
			continue
		}
		if _, err := doOffShelfLiveRoom(ctx, room.ID); err == nil {
			changed = true
		}
	}
	if changed {
		RefreshRoomListCache(ctx)
	}
}

// OnShelfGuildLiveRooms 上架指定工会下全部已下架主播间
func OnShelfGuildLiveRooms(ctx context.Context, guildId uint64) {
	if guildId == 0 {
		return
	}
	changed := false
	for _, room := range liveroomdao.ListRoomsByGuildFromDB(guildId) {
		if room == nil || room.ID == 0 {
			continue
		}
		if room.Status == entity.LiveRoomStatusOnShelf {
			continue
		}
		if _, err := doOnShelfLiveRoom(ctx, room.ID); err == nil {
			changed = true
		}
	}
	if changed {
		RefreshRoomListCache(ctx)
	}
}

func offShelfLiveRoom(ctx context.Context, anchorId uint64) (*accountdto.SetLiveRoomStatusRes, error) {
	res, err := doOffShelfLiveRoom(ctx, anchorId)
	if err == nil {
		RefreshRoomListCache(ctx)
	}
	return res, err
}

func doOffShelfLiveRoom(ctx context.Context, anchorId uint64) (*accountdto.SetLiveRoomStatusRes, error) {
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.Status == entity.LiveRoomStatusOffShelf {
		liveroomdao.RemoveRoomFromCache(anchorId)
		return &accountdto.SetLiveRoomStatusRes{Success: true, Status: entity.LiveRoomStatusOffShelf}, nil
	}

	if room.LiveRecordId > 0 {
		stopLive(anchorId)
	} else {
		clearRoomAudienceCaches(anchorId)
	}
	// 未结算收益新建归档记录后清零,避免下架丢数据
	liveroomdao.ArchiveAndClearUnsettledIncome(anchorId, room.GuildId)
	// 最近未结算日有效直播次数清零(直查DB)
	liveroomdao.ClearRecentUnsettledDailyEffectiveLiveCount(anchorId)
	room.SetStatus(entity.LiveRoomStatusOffShelf)
	liveroomdao.RemoveRoomFromCache(anchorId)
	return &accountdto.SetLiveRoomStatusRes{Success: true, Status: entity.LiveRoomStatusOffShelf}, nil
}

func onShelfLiveRoom(ctx context.Context, anchorId uint64) (*accountdto.SetLiveRoomStatusRes, error) {
	res, err := doOnShelfLiveRoom(ctx, anchorId)
	if err == nil {
		RefreshRoomListCache(ctx)
	}
	return res, err
}

func doOnShelfLiveRoom(ctx context.Context, anchorId uint64) (*accountdto.SetLiveRoomStatusRes, error) {
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		room = liveroomdao.GetRoomFromDB(anchorId)
	}
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.Status != entity.LiveRoomStatusOnShelf {
		room.SetStatus(entity.LiveRoomStatusOnShelf)
	}
	liveroomdao.AddRoomToCache(room)
	userinfodao.PreloadUserCumulativeStatToCache([]uint64{anchorId})
	return &accountdto.SetLiveRoomStatusRes{Success: true, Status: entity.LiveRoomStatusOnShelf}, nil
}

// QueryOffShelfLiveRoomList CMS 回收站:直查 DB 已下架直播间
func QueryOffShelfLiveRoomList(_ context.Context, req *accountdto.QueryOffShelfLiveRoomListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.ListOffShelfRooms(req.PageIndex, req.PageSize, req.Key)
	roomIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			roomIds = append(roomIds, row.ID)
		}
	}
	incomeMap := liveroomdao.ListLiveRoomIncomeTotalFromDB(roomIds)
	list := make([]*accountdto.OffShelfLiveRoomItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &accountdto.OffShelfLiveRoomItem{
			ID:           row.ID,
			Nickname:     row.Nickname,
			Phone:        row.Phone,
			Avatar:       upload.ResolveAvatarUrlForUser(row.ID, row.Avatar),
			UserType:     row.UserType,
			GuildId:      row.GuildId,
			RoomTitle:    row.Title,
			RoomId:       row.ID,
			Category:     row.Category,
			Ban:          row.Ban,
			BanApplyTime: row.BanApplyTime,
			BanReason:    row.BanReason,
			Status:       entity.LiveRoomStatusOffShelf,
			UpdatedAt:    row.UpdatedAt,
			CreatedAt:    row.CreatedAt,
		}
		if income := incomeMap[row.ID]; income != nil {
			item.TotalIncome = income.TotalIncome
			item.TotalGiftIncome = income.TotalGiftIncome
			item.TotalPaidDanmakuIncome = income.TotalPaidDanmakuIncome
			item.TotalPrivateRoomTicketIncome = income.TotalPrivateRoomTicketIncome
			item.TotalPrivateRoomWatchIncome = income.TotalPrivateRoomWatchIncome
			item.TotalVideoCallIncome = income.TotalVideoCallIncome
			item.TotalVideoCallTicketIncome = income.TotalVideoCallTicketIncome
			item.TotalVideoCallBillingIncome = income.TotalVideoCallBillingIncome
		}
		list = append(list, item)
	}
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  list,
	}, nil
}
