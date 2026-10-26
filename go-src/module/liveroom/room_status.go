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

// BatchOffShelfShowcaseLiveRooms CMS一键批量下架已上架的秀场直播间(category=1)
func BatchOffShelfShowcaseLiveRooms(ctx context.Context, _ *accountdto.BatchOffShelfShowcaseLiveRoomsReq) (*accountdto.BatchOffShelfShowcaseLiveRoomsRes, error) {
	rooms := getRoomListCache()
	anchorIds := make([]uint64, 0)
	for _, room := range rooms {
		if room == nil || room.ID == 0 {
			continue
		}
		if room.Status != entity.LiveRoomStatusOnShelf {
			continue
		}
		if liveRoomCategoryFromCache(room) != entity.LiveRoomCategoryHot {
			continue
		}
		anchorIds = append(anchorIds, room.ID)
	}
	res := &accountdto.BatchOffShelfShowcaseLiveRoomsRes{
		TotalMatched: len(anchorIds),
		FailIds:      make([]uint64, 0),
	}
	if len(anchorIds) == 0 {
		return res, nil
	}
	changed := false
	for _, anchorId := range anchorIds {
		if _, err := doOffShelfLiveRoom(ctx, anchorId); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, anchorId)
			continue
		}
		res.SuccessCount++
		changed = true
	}
	if changed {
		RefreshRoomListCache(ctx)
	}
	return res, nil
}

// BatchOnShelfShowcaseLiveRooms CMS一键批量上架已下架的秀场直播间(category=1,直查DB)
func BatchOnShelfShowcaseLiveRooms(ctx context.Context, _ *accountdto.BatchOnShelfShowcaseLiveRoomsReq) (*accountdto.BatchOnShelfShowcaseLiveRoomsRes, error) {
	anchorIds := liveroomdao.ListOffShelfRoomIdsByCategory(entity.LiveRoomCategoryHot)
	res := &accountdto.BatchOnShelfShowcaseLiveRoomsRes{
		TotalMatched: len(anchorIds),
		FailIds:      make([]uint64, 0),
	}
	if len(anchorIds) == 0 {
		return res, nil
	}
	changed := false
	for _, anchorId := range anchorIds {
		if _, err := doOnShelfLiveRoom(ctx, anchorId); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, anchorId)
			continue
		}
		res.SuccessCount++
		changed = true
	}
	if changed {
		RefreshRoomListCache(ctx)
	}
	return res, nil
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

// doOffShelfLiveRoom 下架主播直播间:只改 status=下架并移出列表缓存.
// 若正在直播则先停播(写 end_time/有效时长等走既有 stopLive);不归档/不清未结算收益,不清日表有效开播时长.
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
	room.SetStatus(entity.LiveRoomStatusOffShelf)
	liveroomdao.FlushRoomCache(room)
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
