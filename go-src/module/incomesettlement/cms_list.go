package incomesettlement

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/entity/live"
)

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func fillCMSItemFromAnchor(row *entity.AnchorIncomeSettlementLog) *incomesettlementdto.CMSIncomeSettlementLogItem {
	if row == nil {
		return nil
	}
	item := &incomesettlementdto.CMSIncomeSettlementLogItem{
		Id:                           row.ID,
		RoomId:                       row.RoomId,
		TotalIncome:                  row.TotalIncome,
		TotalGiftIncome:              row.TotalGiftIncome,
		TotalPaidDanmakuIncome:       row.TotalPaidDanmakuIncome,
		TotalPrivateRoomTicketIncome: row.TotalPrivateRoomTicketIncome,
		TotalPrivateRoomWatchIncome:  row.TotalPrivateRoomWatchIncome,
		TotalVideoCallIncome:         row.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:   row.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome:  row.TotalVideoCallBillingIncome,
		TotalLiveDuration:            row.TotalLiveDuration,
		SettlementSalary:             row.SettlementSalary,
		SettlementShareAmount:        row.SettlementShareAmount,
		AnchorSharePercent:           row.AnchorSharePercent,
		CreatedAt:                    &row.CreatedAt,
	}
	return item
}

func fillCMSItemFromGuild(row *entity.GuildIncomeSettlementLog) *incomesettlementdto.CMSIncomeSettlementLogItem {
	if row == nil {
		return nil
	}
	item := &incomesettlementdto.CMSIncomeSettlementLogItem{
		Id:                           row.ID,
		GuildId:                      row.GuildId,
		TotalIncome:                  row.TotalIncome,
		TotalGiftIncome:              row.TotalGiftIncome,
		TotalPaidDanmakuIncome:       row.TotalPaidDanmakuIncome,
		TotalPrivateRoomTicketIncome: row.TotalPrivateRoomTicketIncome,
		TotalPrivateRoomWatchIncome:  row.TotalPrivateRoomWatchIncome,
		TotalVideoCallIncome:         row.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:   row.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome:  row.TotalVideoCallBillingIncome,
		TotalLiveDuration:            row.TotalLiveDuration,
		SettlementSalary:             row.SettlementSalary,
		SettlementShareAmount:        row.SettlementShareAmount,
		GuildSharePercent:            row.GuildSharePercent,
		CreatedAt:                    &row.CreatedAt,
	}
	return item
}

func collectAnchorRoomIds(rows []*entity.AnchorIncomeSettlementLog) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.RoomId > 0 {
			ids = append(ids, row.RoomId)
		}
	}
	return ids
}

func collectGuildIds(rows []*entity.GuildIncomeSettlementLog) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.GuildId > 0 {
			ids = append(ids, row.GuildId)
		}
	}
	return ids
}

// GetAnchorCMSList CMS分页查询主播结算流水
func GetAnchorCMSList(_ context.Context, req *incomesettlementdto.CMSAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	roomIds := liveroomdao.ParseLiveRecordAnchorIds(req.RoomId, "", "", req.AnchorIds)
	total, rows := liveroomdao.AnchorIncomeSettlementLogCMSList(&liveroomdao.AnchorIncomeSettlementLogCMSListFilter{
		RoomIds:   roomIds,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	nicknameMap := userinfodao.GetNicknameMapByUserIds(collectAnchorRoomIds(rows))
	list := make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		item := fillCMSItemFromAnchor(row)
		if item != nil && nicknameMap != nil {
			item.RoomNickname = nicknameMap[row.RoomId]
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// GetGuildCMSList CMS分页查询工会结算流水
func GetGuildCMSList(_ context.Context, req *incomesettlementdto.CMSGuildIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.GuildIncomeSettlementLogCMSList(&liveroomdao.GuildIncomeSettlementLogCMSListFilter{
		GuildId:   parseUint64Filter(req.GuildId),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	guildNameMap := guilddao.GetNameMapByIds(collectGuildIds(rows))
	list := make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		item := fillCMSItemFromGuild(row)
		if item != nil && guildNameMap != nil {
			item.GuildName = guildNameMap[row.GuildId]
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// GetAnchorCMSListByGuildIds CMS分页查询指定工会下主播结算流水
func GetAnchorCMSListByGuildIds(_ context.Context, guildIds []uint64, req *guilddto.CMSMyGuildAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	list := make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0)
	if len(guildIds) == 0 {
		return httpserver.NewCMSQueryResp(0, list), nil
	}
	total, rows := liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIds(&liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIdsFilter{
		GuildIds:  guildIds,
		RoomId:    parseUint64Filter(req.RoomId),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	roomIds := collectAnchorRoomIds(rows)
	nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
	guildIdMap := liveroomdao.GetGuildIdMapByRoomIds(roomIds)
	guildNameMap := guilddao.GetNameMapByIds(collectGuildIdsFromMap(guildIdMap))
	list = make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		item := fillCMSItemFromAnchor(row)
		if item == nil {
			continue
		}
		if nicknameMap != nil {
			item.RoomNickname = nicknameMap[row.RoomId]
		}
		if guildIdMap != nil {
			item.GuildId = guildIdMap[row.RoomId]
			if guildNameMap != nil {
				item.GuildName = guildNameMap[item.GuildId]
			}
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func collectGuildIdsFromMap(guildIdMap map[uint64]uint64) []uint64 {
	if len(guildIdMap) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(guildIdMap))
	seen := make(map[uint64]struct{}, len(guildIdMap))
	for _, guildId := range guildIdMap {
		if guildId == 0 {
			continue
		}
		if _, ok := seen[guildId]; ok {
			continue
		}
		seen[guildId] = struct{}{}
		ids = append(ids, guildId)
	}
	return ids
}
