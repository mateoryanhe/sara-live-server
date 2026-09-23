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
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/upload"
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
		Id:                          row.ID,
		RoomId:                      row.RoomId,
		TotalIncome:                 row.TotalIncome,
		TotalSocialIncome:           row.TotalSocialIncome,
		TotalGiftIncome:             row.TotalGiftIncome,
		TotalPaidDanmakuIncome:      row.TotalPaidDanmakuIncome,
		TotalVideoCallIncome:        row.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:  row.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome: row.TotalVideoCallBillingIncome,
		TotalShortVideoIncome:       row.TotalShortVideoIncome,
		TotalGameIncome:             row.TotalGameIncome,
		TotalLiveDuration:           row.TotalLiveDuration,
		SettlementSalary:            row.SettlementSalary,
		SettlementShareAmount:       row.SettlementShareAmount,
		SettlementShareAmountUsd:    row.SettlementShareAmountUsd,
		AnchorSharePercent:          row.AnchorSharePercent,
		HasSalary:                   row.HasSalary,
		AnchorSocialSharePercent:    row.AnchorSocialSharePercent,
		GuildSocialSharePercent:     row.GuildSocialSharePercent,
		AnchorGameSharePercent:      row.AnchorGameSharePercent,
		GuildGameSharePercent:       row.GuildGameSharePercent,
		AnchorSocialShareAmount:     row.AnchorSocialShareAmount,
		GuildSocialShareAmount:      row.GuildSocialShareAmount,
		AnchorGameShareAmountGold:   row.AnchorGameShareAmountGold,
		GuildGameShareAmountGold:    row.GuildGameShareAmountGold,
		SettlementRuleType:          row.SettlementRuleType,
		DirectPayout:                row.DirectPayout,
		SettlementReceivableUsd:     row.SettlementReceivableUsd,
		GoldToDiamondRate:           row.GoldToDiamondRate,
		UsdToGoldRate:               row.UsdToGoldRate,
		GameShareAmountDiamond:      row.GameShareAmountDiamond,
		TotalSettlementDiamond:      row.TotalSettlementDiamond,
		Status:                      row.Status,
		TransferAt:                  row.TransferAt,
		TransferOrderId:             row.TransferOrderId,
		TransferPlatformNo:          row.TransferPlatformNo,
		TransferLocalAmount:         row.TransferLocalAmount,
		TransferFailMsg:             row.TransferFailMsg,
		TransferCurrency:            row.TransferCurrency,
		CreatedAt:                   &row.CreatedAt,
		UpdatedAt:                   &row.UpdatedAt,
	}
	return item
}

func fillCMSItemFromGuild(row *entity.GuildIncomeSettlementLog) *incomesettlementdto.CMSIncomeSettlementLogItem {
	if row == nil {
		return nil
	}
	item := &incomesettlementdto.CMSIncomeSettlementLogItem{
		Id:                          row.ID,
		GuildId:                     row.GuildId,
		TotalIncome:                 row.TotalIncome,
		TotalSocialIncome:           row.TotalSocialIncome,
		TotalGiftIncome:             row.TotalGiftIncome,
		TotalPaidDanmakuIncome:      row.TotalPaidDanmakuIncome,
		TotalVideoCallIncome:        row.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:  row.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome: row.TotalVideoCallBillingIncome,
		TotalShortVideoIncome:       row.TotalShortVideoIncome,
		TotalGameIncome:             row.TotalGameIncome,
		TotalLiveDuration:           row.TotalLiveDuration,
		SettlementSalary:            row.SettlementSalary,
		SettlementShareAmount:       row.SettlementShareAmount,
		SettlementShareAmountUsd:    row.SettlementShareAmountUsd,
		SettlementReceivableUsd:     row.SettlementReceivableUsd,
		GuildSharePercent:           row.GuildSharePercent,
		SettlementRuleType:          row.SettlementRuleType,
		AnchorSocialShareAmount:     row.AnchorSocialShareAmount,
		GuildSocialShareAmount:      row.GuildSocialShareAmount,
		AnchorGameShareAmountGold:   row.AnchorGameShareAmountGold,
		GuildGameShareAmountGold:    row.GuildGameShareAmountGold,
		GoldToDiamondRate:           row.GoldToDiamondRate,
		UsdToGoldRate:               row.UsdToGoldRate,
		GameShareAmountDiamond:      row.GameShareAmountDiamond,
		TotalSettlementDiamond:      row.TotalSettlementDiamond,
		Status:                      row.Status,
		TransferAt:                  row.TransferAt,
		TransferOrderId:             row.TransferOrderId,
		TransferPlatformNo:          row.TransferPlatformNo,
		TransferLocalAmount:         row.TransferLocalAmount,
		TransferFailMsg:             row.TransferFailMsg,
		TransferCurrency:            row.TransferCurrency,
		CreatedAt:                   &row.CreatedAt,
		UpdatedAt:                   &row.UpdatedAt,
	}
	return item
}

func fillCMSGuildListItem(row *entity.GuildIncomeSettlementLog) *incomesettlementdto.CMSIncomeSettlementLogItem {
	if row == nil {
		return nil
	}
	return &incomesettlementdto.CMSIncomeSettlementLogItem{
		Id:                      row.ID,
		GuildId:                 row.GuildId,
		SettlementReceivableUsd: row.SettlementReceivableUsd,
		Status:                  row.Status,
		CreatedAt:               &row.CreatedAt,
		UpdatedAt:               &row.UpdatedAt,
	}
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
func GetAnchorCMSList(ctx context.Context, req *incomesettlementdto.CMSAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	roomIds := liveroomdao.ParseLiveRecordAnchorIds(req.RoomId, "", "", req.AnchorIds)
	if req.DirectPayout != nil && *req.DirectPayout {
		visibleIds, restrict, empty := cmsvis.PlatformAnchorVisibilityFilter(ctx)
		if empty {
			return httpserver.NewCMSQueryResp(0, []*incomesettlementdto.CMSIncomeSettlementLogItem{}), nil
		}
		if restrict {
			roomIds = intersectUint64Ids(roomIds, visibleIds)
			if len(roomIds) == 0 {
				return httpserver.NewCMSQueryResp(0, []*incomesettlementdto.CMSIncomeSettlementLogItem{}), nil
			}
		}
	}
	total, rows := liveroomdao.AnchorIncomeSettlementLogCMSList(&liveroomdao.AnchorIncomeSettlementLogCMSListFilter{
		RoomIds:                  roomIds,
		StartTime:                req.StartTime,
		EndTime:                  req.EndTime,
		Status:                   req.Status,
		DirectPayout:             req.DirectPayout,
		OrderByReceivableUsdDesc: req.OrderByReceivableUsdDesc,
		PageIndex:                req.PageIndex,
		PageSize:                 req.PageSize,
	})
	profileRoomIds := collectAnchorRoomIds(rows)
	nicknameMap := userinfodao.GetNicknameMapByUserIds(profileRoomIds)
	profileMap := userinfodao.GetUserProfileMapByUserIds(profileRoomIds)
	var transferMap map[uint64]*entity.LiveAnchorTransferInfo
	if req.IncludeTransferInfo {
		transferMap = liveroomdao.GetAnchorTransferInfoMapByIds(profileRoomIds)
	}
	list := make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		item := fillCMSItemFromAnchor(row)
		if item != nil {
			if nicknameMap != nil {
				item.RoomNickname = nicknameMap[row.RoomId]
			}
			if profileMap != nil {
				if profile := profileMap[row.RoomId]; profile != nil {
					item.RoomAvatar = upload.ResolveAvatarUrlForUser(row.RoomId, profile.Avatar)
				}
			}
			if transferMap != nil {
				if info := transferMap[row.RoomId]; info != nil {
					if item.TransferCurrency == "" {
						item.TransferCurrency = info.Currency
					}
					item.TransferPayeeName = info.PayeeName
					item.TransferBankName = info.BankName
					item.TransferAccountNo = info.AccountNo
					item.TransferBankCode = info.BankCode
				}
			}
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func intersectUint64Ids(selected, allowed []uint64) []uint64 {
	if len(allowed) == 0 {
		return nil
	}
	if len(selected) == 0 {
		return append([]uint64(nil), allowed...)
	}
	set := make(map[uint64]struct{}, len(allowed))
	for _, id := range allowed {
		set[id] = struct{}{}
	}
	out := make([]uint64, 0, len(selected))
	for _, id := range selected {
		if _, ok := set[id]; ok {
			out = append(out, id)
		}
	}
	return out
}

// GetGuildCMSList CMS分页查询工会结算流水
func GetGuildCMSList(ctx context.Context, req *incomesettlementdto.CMSGuildIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	guildIds, restrict, empty := cmsvis.VisibilityGuildFilter(ctx)
	if empty {
		return httpserver.NewCMSQueryResp(0, []*incomesettlementdto.CMSIncomeSettlementLogItem{}), nil
	}
	total, rows := liveroomdao.GuildIncomeSettlementLogCMSList(&liveroomdao.GuildIncomeSettlementLogCMSListFilter{
		GuildId:                  parseUint64Filter(req.GuildId),
		GuildIds:                 guildIds,
		FilterByGuild:            restrict,
		GuildType:                req.GuildType,
		StartTime:                req.StartTime,
		EndTime:                  req.EndTime,
		TransferStartTime:        req.TransferStartTime,
		TransferEndTime:          req.TransferEndTime,
		PayoutOnly:               req.PayoutOnly,
		Status:                   req.Status,
		OrderByReceivableUsdDesc: req.OrderByReceivableUsdDesc,
		IncludeDetail:            req.IncludeDetail,
		IncludeTransfer:          req.IncludeTransfer,
		PageIndex:                req.PageIndex,
		PageSize:                 req.PageSize,
	})
	rowGuildIds := collectGuildIds(rows)
	guildNameMap := guilddao.GetNameMapByIds(rowGuildIds)
	var transferMap map[uint64]*entity.LiveGuildTransferInfo
	if req.IncludeTransferInfo {
		transferMap = guilddao.GetGuildTransferInfoMapByIds(rowGuildIds)
	}
	list := make([]*incomesettlementdto.CMSIncomeSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		item := fillCMSGuildListItem(row)
		if req.IncludeDetail || req.IncludeTransfer {
			item = fillCMSItemFromGuild(row)
		}
		if item == nil {
			continue
		}
		if guildNameMap != nil {
			item.GuildName = guildNameMap[row.GuildId]
		}
		if transferMap != nil {
			if info := transferMap[row.GuildId]; info != nil {
				if item.TransferCurrency == "" {
					item.TransferCurrency = info.Currency
				}
				item.TransferPayeeName = info.PayeeName
				item.TransferBankName = info.BankName
				item.TransferAccountNo = info.AccountNo
				item.TransferBankCode = info.BankCode
			}
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// GetGuildCMSDetail 查询一条工会结算主单，并装配结算快照、代付过程和当前收款信息。
func GetGuildCMSDetail(ctx context.Context, req *incomesettlementdto.CMSGuildIncomeSettlementLogDetailReq) (*incomesettlementdto.CMSGuildIncomeSettlementLogDetailRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil || id == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := guildVisibleSet(ctx)
	if empty {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	row := liveroomdao.GetGuildIncomeSettlementLogById(id)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if restrict {
		if _, ok := visibleSet[row.GuildId]; !ok {
			return nil, errercode.CreateCode(errercode.NoPermission)
		}
	}
	item := fillCMSItemFromGuild(row)
	if names := guilddao.GetNameMapByIds([]uint64{row.GuildId}); names != nil {
		item.GuildName = names[row.GuildId]
	}
	if info := guilddao.GetGuildTransferInfo(row.GuildId); info != nil {
		if item.TransferCurrency == "" {
			item.TransferCurrency = info.Currency
		}
		item.TransferPayeeName = info.PayeeName
		item.TransferBankName = info.BankName
		item.TransferAccountNo = info.AccountNo
		item.TransferBankCode = info.BankCode
	}
	return &incomesettlementdto.CMSGuildIncomeSettlementLogDetailRes{Item: item}, nil
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
	profileMap := userinfodao.GetUserProfileMapByUserIds(roomIds)
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
		if profileMap != nil {
			if profile := profileMap[row.RoomId]; profile != nil {
				item.RoomAvatar = upload.ResolveAvatarUrlForUser(row.RoomId, profile.Avatar)
			}
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
