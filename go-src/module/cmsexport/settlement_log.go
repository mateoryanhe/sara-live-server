package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func exportAnchorIncomeSettlementLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportAnchorIncomeSettlementLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		roomIds := liveroomdao.ParseLiveRecordAnchorIds(req.RoomId, "", "", req.AnchorIds)
		total, rows := liveroomdao.AnchorIncomeSettlementLogCMSList(&liveroomdao.AnchorIncomeSettlementLogCMSListFilter{
			RoomIds:   roomIds,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		nicknameMap := userinfodao.GetNicknameMapByUserIds(collectAnchorSettlementRoomIds(rows))
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, anchorSettlementLogToCSVRow(row, nicknameMap))
		}
		return total, csvRows
	}, onProgress)
}

func exportGuildIncomeSettlementLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGuildIncomeSettlementLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.GuildIncomeSettlementLogCMSList(&liveroomdao.GuildIncomeSettlementLogCMSListFilter{
			GuildId:   parseUint64Filter(req.GuildId),
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		guildNameMap := guilddao.GetNameMapByIds(collectGuildSettlementGuildIds(rows))
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, guildSettlementLogToCSVRow(row, guildNameMap))
		}
		return total, csvRows
	}, onProgress)
}

func exportGuildAnchorIncomeSettlementLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGuildAnchorIncomeSettlementLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	guildId := parseUint64Filter(req.GuildId)
	if guildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIds(&liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIdsFilter{
			GuildIds:  []uint64{guildId},
			RoomId:    parseUint64Filter(req.RoomId),
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		csvRows := guildAnchorSettlementRowsToCSV(rows)
		return total, csvRows
	}, onProgress)
}

func exportMyGuildAnchorIncomeSettlementLogCSV(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportMyGuildAnchorIncomeSettlementLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	if cmsUserId == 0 {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	guildIds := collectOwnedGuildIds(guilddao.ListGuildsByLeaderId(cmsUserId))
	if len(guildIds) == 0 {
		return nil, errExportEmpty
	}
	if filterGuildId := parseUint64Filter(req.GuildId); filterGuildId > 0 {
		if !containsGuildId(guildIds, filterGuildId) {
			return nil, errercode.CreateCode(errercode.NoPermission)
		}
		guildIds = []uint64{filterGuildId}
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIds(&liveroomdao.AnchorIncomeSettlementLogCMSListByGuildIdsFilter{
			GuildIds:  guildIds,
			RoomId:    parseUint64Filter(req.RoomId),
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		return total, guildAnchorSettlementRowsToCSV(rows)
	}, onProgress)
}

func guildAnchorSettlementRowsToCSV(rows []*liveentity.AnchorIncomeSettlementLog) [][]string {
	roomIds := collectAnchorSettlementRoomIds(rows)
	nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
	guildIdMap := liveroomdao.GetGuildIdMapByRoomIds(roomIds)
	guildNameMap := guilddao.GetNameMapByIds(collectGuildIdsFromMap(guildIdMap))
	csvRows := make([][]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		guildName := ""
		if guildIdMap != nil && guildNameMap != nil {
			guildName = guildNameMap[guildIdMap[row.RoomId]]
		}
		roomNickname := ""
		if nicknameMap != nil {
			roomNickname = nicknameMap[row.RoomId]
		}
		csvRows = append(csvRows, guildAnchorSettlementLogToCSVRow(row, guildName, roomNickname))
	}
	return csvRows
}

func anchorSettlementLogToCSVRow(row *liveentity.AnchorIncomeSettlementLog, nicknameMap map[uint64]string) []string {
	if row == nil {
		return nil
	}
	roomNickname := ""
	if nicknameMap != nil {
		roomNickname = nicknameMap[row.RoomId]
	}
	cells := []string{
		formatCSVUint(row.ID),
		formatCSVUint(row.RoomId),
		roomNickname,
	}
	cells = append(cells, anchorSettlementAmountCells(row)...)
	cells = append(cells, formatCSVTime(row.CreatedAt))
	return cells
}

func guildSettlementLogToCSVRow(row *liveentity.GuildIncomeSettlementLog, guildNameMap map[uint64]string) []string {
	if row == nil {
		return nil
	}
	guildName := ""
	if guildNameMap != nil {
		guildName = guildNameMap[row.GuildId]
	}
	cells := []string{
		formatCSVUint(row.ID),
		formatCSVUint(row.GuildId),
		guildName,
	}
	cells = append(cells, guildSettlementAmountCells(row)...)
	cells = append(cells, formatCSVTime(row.CreatedAt))
	return cells
}

func guildAnchorSettlementLogToCSVRow(row *liveentity.AnchorIncomeSettlementLog, guildName, roomNickname string) []string {
	if row == nil {
		return nil
	}
	cells := []string{
		formatCSVUint(row.ID),
		guildName,
		formatCSVUint(row.RoomId),
		roomNickname,
	}
	cells = append(cells, anchorSettlementAmountCells(row)...)
	cells = append(cells, formatCSVTime(row.CreatedAt))
	return cells
}

func anchorSettlementAmountCells(row *liveentity.AnchorIncomeSettlementLog) []string {
	return append(incomeAmountCSVCells(&row.LiveRoomIncomeAmounts),
		formatCSVFloat(row.SettlementSalary),
		formatCSVFloat(row.AnchorSharePercent),
		formatCSVFloat(row.SettlementShareAmount),
	)
}

func guildSettlementAmountCells(row *liveentity.GuildIncomeSettlementLog) []string {
	return append(incomeAmountCSVCells(&row.LiveRoomIncomeAmounts),
		formatCSVFloat(row.GuildSharePercent),
		formatCSVFloat(row.SettlementShareAmount),
	)
}

func collectAnchorSettlementRoomIds(rows []*liveentity.AnchorIncomeSettlementLog) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.RoomId > 0 {
			ids = append(ids, row.RoomId)
		}
	}
	return ids
}

func collectGuildSettlementGuildIds(rows []*liveentity.GuildIncomeSettlementLog) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.GuildId > 0 {
			ids = append(ids, row.GuildId)
		}
	}
	return ids
}

func collectOwnedGuildIds(rows []*liveentity.LiveGuild) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.ID > 0 {
			ids = append(ids, row.ID)
		}
	}
	return ids
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

func containsGuildId(ids []uint64, guildId uint64) bool {
	for _, id := range ids {
		if id == guildId {
			return true
		}
	}
	return false
}
