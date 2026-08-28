package guild

import (
	"context"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

// QueryGuildDetail CMS获取工会详情收益(缓存优先,否则直查DB;基本信息由前端传入)
func QueryGuildDetail(_ context.Context, req *guilddto.GetGuildDetailReq) (*guilddto.GetGuildDetailRes, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &guilddto.GetGuildDetailRes{
		IncomeUnsettled: toIncomeUnsettledItem(liveroomdao.GetGuildIncomeUnsettledForCMS(req.GuildId)),
		IncomeSettled:   toIncomeSettledItem(liveroomdao.GetGuildIncomeSettledForCMS(req.GuildId)),
		IncomeTotal:     toIncomeTotalItem(liveroomdao.GetGuildIncomeTotalForCMS(req.GuildId)),
	}, nil
}

// QueryGuildIncomeArchives CMS获取工会下架归档(直查DB)
func QueryGuildIncomeArchives(_ context.Context, req *guilddto.GetGuildIncomeArchivesReq) (*guilddto.GetGuildIncomeArchivesRes, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &guilddto.GetGuildIncomeArchivesRes{
		List: toGuildIncomeArchiveItems(liveroomdao.ListGuildIncomeUnsettledArchives(req.GuildId, 50)),
	}, nil
}

func toIncomeAmountsItem(a *liveentity.LiveRoomIncomeAmounts) accountdto.LiveRoomIncomeAmountsItem {
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

func toIncomeUnsettledItem(row *liveentity.GuildIncomeUnsettled) *accountdto.LiveRoomIncomeUnsettledItem {
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

func toIncomeSettledItem(row *liveentity.GuildIncomeSettled) *accountdto.LiveRoomIncomeSettledItem {
	if row == nil {
		return nil
	}
	item := &accountdto.LiveRoomIncomeSettledItem{
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
		SettlementSalary:          row.SettlementSalary,
		SettlementShareAmount:     row.SettlementShareAmount,
		SettlementShareAmountUsd:  row.SettlementShareAmountUsd,
		SettlementReceivableUsd:   row.SettlementReceivableUsd,
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toIncomeTotalItem(row *liveentity.GuildIncomeTotal) *accountdto.LiveRoomIncomeTotalItem {
	if row == nil {
		return nil
	}
	item := &accountdto.LiveRoomIncomeTotalItem{
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
		SettlementSalary:          row.SettlementSalary,
		SettlementShareAmount:     row.SettlementShareAmount,
		SettlementShareAmountUsd:  row.SettlementShareAmountUsd,
		SettlementReceivableUsd:   row.SettlementReceivableUsd,
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toGuildIncomeArchiveItems(rows []*liveentity.GuildIncomeUnsettledArchive) []*guilddto.GuildIncomeArchiveItem {
	list := make([]*guilddto.GuildIncomeArchiveItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := &guilddto.GuildIncomeArchiveItem{
			ID:                        row.ID,
			GuildId:                   row.GuildId,
			LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
		}
		if !row.CreatedAt.IsZero() {
			createdAt := row.CreatedAt
			item.CreatedAt = &createdAt
		}
		list = append(list, item)
	}
	return list
}
