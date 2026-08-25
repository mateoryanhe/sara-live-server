package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/guild"
)

func exportAnchorDailyEffectiveLiveCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportAnchorDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	if req.AnchorId == 0 || liveroomdao.ResolveRoom(req.AnchorId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSList(&liveroomdao.DailyAnchorEffectiveLiveCMSListFilter{
			RoomId:        req.AnchorId,
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, anchorDailyEffectiveLiveToCSVRow(row, req.SettledYesText, req.SettledNoText))
		}
		return total, csvRows
	}, onProgress)
}

func exportGuildDailyEffectiveLiveCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGuildDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	if req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyGuildEffectiveLiveCMSList(&liveroomdao.DailyGuildEffectiveLiveCMSListFilter{
			GuildId:       req.GuildId,
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, guildDailyEffectiveLiveToCSVRow(row, req.SettledYesText, req.SettledNoText))
		}
		return total, csvRows
	}, onProgress)
}

func exportLiveDailyEffectiveLiveCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
			RoomIds:       liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		roomIds := collectDailyAnchorRoomIds(rows)
		nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
		unsettledIncomeMap := liveroomdao.ListLiveRoomIncomeUnsettledTotalForCMS(roomIds)
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, liveDailyEffectiveLiveToCSVRow(row, nicknameMap, unsettledIncomeMap, req.SettledYesText, req.SettledNoText))
		}
		return total, csvRows
	}, onProgress)
}

func exportGuildAnchorDailyEffectiveLiveCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGuildAnchorDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	guildId := parseUint64Filter(req.GuildId)
	if guildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSListByGuild(&liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildFilter{
			GuildId:       guildId,
			RoomId:        parseUint64Filter(req.RoomId),
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		roomIds := collectDailyAnchorRoomIds(rows)
		nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
		unsettledIncomeMap := liveroomdao.ListLiveRoomIncomeUnsettledTotalForCMS(roomIds)
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, guildAnchorDailyEffectiveLiveToCSVRow(row, nicknameMap, unsettledIncomeMap))
		}
		return total, csvRows
	}, onProgress)
}

func exportMyGuildAnchorDailyEffectiveLiveCSV(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportMyGuildAnchorDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	guildIds := guild.CollectOwnedGuildIdsForCMSUser(cmsUserId)
	if len(guildIds) == 0 {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	roomIds := guild.ResolveOwnedGuildAnchorRoomIds(req.RoomId, req.RoomIds)
	if len(roomIds) > 0 {
		if err := guild.ValidateMyOwnedGuildAnchorsAccessForUser(cmsUserId, roomIds); err != nil {
			return nil, err
		}
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		filter := &liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIdsFilter{
			GuildIds:      guildIds,
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		}
		if len(roomIds) == 1 {
			filter.RoomId = roomIds[0]
		} else if len(roomIds) > 1 {
			filter.RoomIds = roomIds
		}
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIds(filter)
		anchorRoomIds := collectDailyAnchorRoomIds(rows)
		nicknameMap := userinfodao.GetNicknameMapByUserIds(anchorRoomIds)
		unsettledIncomeMap := liveroomdao.ListLiveRoomIncomeUnsettledTotalForCMS(anchorRoomIds)
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, guildAnchorDailyEffectiveLiveToCSVRow(row, nicknameMap, unsettledIncomeMap))
		}
		return total, csvRows
	}, onProgress)
}

func anchorDailyEffectiveLiveToCSVRow(row *liveentity.DailyAnchorEffectiveLive, settledYesText, settledNoText string) []string {
	if row == nil {
		return nil
	}
	cells := []string{
		row.ID,
		row.LiveDate,
		formatLiveDurationMinutes(row.LiveDuration),
		formatLiveDurationMinutes(row.TotalLiveDuration),
	}
	cells = append(cells, incomeAmountCSVCells(&row.LiveRoomIncomeAmounts)...)
	cells = append(cells,
		formatSettledText(row.Settled, settledYesText, settledNoText),
		formatCSVTime(row.CreatedAt),
		formatCSVTime(row.UpdatedAt),
	)
	return cells
}

func guildDailyEffectiveLiveToCSVRow(row *liveentity.DailyGuildEffectiveLive, settledYesText, settledNoText string) []string {
	if row == nil {
		return nil
	}
	cells := []string{
		row.ID,
		row.LiveDate,
		formatLiveDurationMinutes(row.LiveDuration),
		formatLiveDurationMinutes(row.TotalLiveDuration),
	}
	cells = append(cells, incomeAmountCSVCells(&row.LiveRoomIncomeAmounts)...)
	cells = append(cells,
		formatSettledText(row.Settled, settledYesText, settledNoText),
		formatCSVTime(row.CreatedAt),
		formatCSVTime(row.UpdatedAt),
	)
	return cells
}

func guildAnchorDailyEffectiveLiveToCSVRow(row *liveentity.DailyAnchorEffectiveLive, nicknameMap map[uint64]string, unsettledIncomeMap map[uint64]float64) []string {
	if row == nil {
		return nil
	}
	roomNickname := ""
	if nicknameMap != nil {
		roomNickname = nicknameMap[row.RoomId]
	}
	cells := []string{
		formatCSVUint(row.RoomId),
		roomNickname,
		row.LiveDate,
	}
	if unsettledIncomeMap != nil {
		cells = append(cells, formatCSVFloat(unsettledIncomeMap[row.RoomId]))
	} else {
		cells = append(cells, "")
	}
	cells = append(cells,
		formatCSVFloat(row.TotalIncome),
		formatLiveDurationMinutes(row.LiveDuration),
	)
	return cells
}

func liveDailyEffectiveLiveToCSVRow(row *liveentity.DailyAnchorEffectiveLive, nicknameMap map[uint64]string, unsettledIncomeMap map[uint64]float64, settledYesText, settledNoText string) []string {
	if row == nil {
		return nil
	}
	roomNickname := ""
	if nicknameMap != nil {
		roomNickname = nicknameMap[row.RoomId]
	}
	cells := []string{
		row.LiveDate,
		formatCSVUint(row.RoomId),
		roomNickname,
	}
	if unsettledIncomeMap != nil {
		cells = append(cells, formatCSVFloat(unsettledIncomeMap[row.RoomId]))
	} else {
		cells = append(cells, "")
	}
	cells = append(cells,
		formatLiveDurationMinutes(row.LiveDuration),
		formatLiveDurationMinutes(row.TotalLiveDuration),
	)
	cells = append(cells, incomeAmountCSVCells(&row.LiveRoomIncomeAmounts)...)
	cells = append(cells,
		formatSettledText(row.Settled, settledYesText, settledNoText),
		formatCSVTime(row.CreatedAt),
		formatCSVTime(row.UpdatedAt),
	)
	return cells
}

func collectDailyAnchorRoomIds(rows []*liveentity.DailyAnchorEffectiveLive) []uint64 {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.RoomId > 0 {
			ids = append(ids, row.RoomId)
		}
	}
	return ids
}
