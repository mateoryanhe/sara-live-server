package liveroom

import (
	"context"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// QueryAnchorDetail CMS获取主播详情(含直播间与各收益表)
func QueryAnchorDetail(_ context.Context, req *accountdto.GetAnchorDetailReq) (*accountdto.GetAnchorDetailRes, error) {
	room := liveroomdao.ResolveRoom(req.AnchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	anchor := buildAnchorListItem(room)
	if anchor == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &accountdto.GetAnchorDetailRes{
		Anchor:          anchor,
		LiveRoom:        buildAnchorLiveRoomDetail(room),
		IncomeUnsettled: toIncomeUnsettledItem(liveroomdao.GetLiveRoomIncomeUnsettledForCMS(room.ID)),
		IncomeSettled:   toIncomeSettledItem(liveroomdao.GetLiveRoomIncomeSettledForCMS(room.ID)),
		IncomeTotal:     toIncomeTotalItem(liveroomdao.GetLiveRoomIncomeTotalForCMS(room.ID)),
		IncomeArchives:  toIncomeArchiveItems(liveroomdao.ListLiveRoomIncomeUnsettledArchives(room.ID, 50)),
	}, nil
}

func buildAnchorLiveRoomDetail(room *entity.LiveRoom) *accountdto.AnchorLiveRoomDetailItem {
	if room == nil {
		return nil
	}
	item := &accountdto.AnchorLiveRoomDetailItem{
		ID:           room.ID,
		GuildId:      room.GuildId,
		Title:        room.Title,
		Cover:        upload.GetUrlByName(room.Cover),
		Notice:       room.Notice,
		LiveRecordId: room.LiveRecordId,
		HeartTime:    room.HeartTime,
		Ban:          IsRoomBanned(room),
		BanApplyTime: room.BanApplyTime,
		BanReason:    room.BanReason,
		Status:       room.Status,
	}
	if room.LiveRecordId > 0 {
		item.LiveStatus = 1
	}
	if cfg := liveroomdao.GetLiveRoomCfgFromCache(room.ID); cfg != nil {
		item.Category = cfg.Category
		item.PrivateInviteType = entity.NormalizePrivateInviteType(cfg.PrivateInviteType, cfg.Category)
		item.Ticket = cfg.Ticket
		item.Billing = cfg.Billing
	} else if cfg := liveroomdao.GetLiveRoomCfgFromDB(room.ID); cfg != nil {
		item.Category = cfg.Category
		item.PrivateInviteType = entity.NormalizePrivateInviteType(cfg.PrivateInviteType, cfg.Category)
		item.Ticket = cfg.Ticket
		item.Billing = cfg.Billing
	}
	if !room.CreatedAt.IsZero() {
		createdAt := room.CreatedAt
		item.CreatedAt = &createdAt
	}
	if !room.UpdatedAt.IsZero() {
		updatedAt := room.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toIncomeAmountsItem(a *entity.LiveRoomIncomeAmounts) accountdto.LiveRoomIncomeAmountsItem {
	if a == nil {
		return accountdto.LiveRoomIncomeAmountsItem{}
	}
	return accountdto.LiveRoomIncomeAmountsItem{
		TotalIncome:                  a.TotalIncome,
		TotalGiftIncome:              a.TotalGiftIncome,
		TotalPaidDanmakuIncome:       a.TotalPaidDanmakuIncome,
		TotalPrivateRoomTicketIncome: a.TotalPrivateRoomTicketIncome,
		TotalPrivateRoomWatchIncome:  a.TotalPrivateRoomWatchIncome,
		TotalVideoCallIncome:         a.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:   a.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome:  a.TotalVideoCallBillingIncome,
		TotalShortVideoIncome:        a.TotalShortVideoIncome,
		TotalGameIncome:              a.TotalGameIncome,
		TotalLiveDuration:            a.TotalLiveDuration,
	}
}

func toIncomeUnsettledItem(row *entity.LiveRoomIncomeUnsettled) *accountdto.LiveRoomIncomeUnsettledItem {
	if row == nil {
		return nil
	}
	item := &accountdto.LiveRoomIncomeUnsettledItem{
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toIncomeSettledItem(row *entity.LiveRoomIncomeSettled) *accountdto.LiveRoomIncomeSettledItem {
	if row == nil {
		return nil
	}
	item := &accountdto.LiveRoomIncomeSettledItem{
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
		SettlementSalary:          row.SettlementSalary,
		SettlementShareAmount:     row.SettlementShareAmount,
		SettlementShareAmountUsd:  row.SettlementShareAmountUsd,
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toIncomeTotalItem(row *entity.LiveRoomIncomeTotal) *accountdto.LiveRoomIncomeTotalItem {
	if row == nil {
		return nil
	}
	item := &accountdto.LiveRoomIncomeTotalItem{
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
		SettlementSalary:          row.SettlementSalary,
		SettlementShareAmount:     row.SettlementShareAmount,
		SettlementShareAmountUsd:  row.SettlementShareAmountUsd,
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toIncomeArchiveItems(rows []*entity.LiveRoomIncomeUnsettledArchive) []*accountdto.LiveRoomIncomeArchiveItem {
	list := make([]*accountdto.LiveRoomIncomeArchiveItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := &accountdto.LiveRoomIncomeArchiveItem{
			ID:                        row.ID,
			RoomId:                    row.RoomId,
			GuildId:                   row.GuildId,
			LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
			SettlementSalary:          row.SettlementSalary,
		}
		if !row.CreatedAt.IsZero() {
			createdAt := row.CreatedAt
			item.CreatedAt = &createdAt
		}
		list = append(list, item)
	}
	return list
}
