package cmsexport

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"xr-game-server/core/xrtime"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/fileexport"
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

func exportLiveDailyEffectiveLiveCSV(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveDailyEffectiveLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	guildIds, restrict, empty := cmsvis.VisibilityGuildFilterForUser(cmsUserId)
	if empty {
		return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
			return 0, nil
		}, onProgress)
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
			RoomIds:       liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
			GuildIds:      guildIds,
			FilterByGuild: restrict,
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Keyword:       req.Keyword,
			Settled:       req.Settled,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		roomIds := collectDailyAnchorRoomIds(rows)
		nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, liveDailyEffectiveLiveToCSVRow(row, nicknameMap, req.SettledYesText, req.SettledNoText))
		}
		return total, csvRows
	}, onProgress)
}

func parseGuildIdFilters(vals []string) []uint64 {
	if len(vals) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(vals))
	ids := make([]uint64, 0, len(vals))
	for _, val := range vals {
		id := parseUint64Filter(val)
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func resolveRequestedGuildIdsForExport(cmsUserId uint64, requested []uint64) ([]uint64, error) {
	if len(requested) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visible, restrict, empty := cmsvis.VisibilityGuildFilterForUser(cmsUserId)
	if empty {
		return nil, errExportEmpty
	}
	if !restrict {
		return requested, nil
	}
	visibleSet := make(map[uint64]struct{}, len(visible))
	for _, id := range visible {
		visibleSet[id] = struct{}{}
	}
	out := make([]uint64, 0, len(requested))
	for _, id := range requested {
		if _, ok := visibleSet[id]; ok {
			out = append(out, id)
		}
	}
	if len(out) == 0 {
		return nil, errExportEmpty
	}
	return out, nil
}

func sanitizeExportFilePart(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "guild"
	}
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_", "?", "_",
		"\"", "_", "<", "_", ">", "_", "|", "_",
	)
	name = replacer.Replace(name)
	name = strings.Trim(name, " .")
	if name == "" {
		return "guild"
	}
	return name
}

func writeCSVRowToBuilder(sb *strings.Builder, cells []string) {
	parts := make([]string, len(cells))
	for i, cell := range cells {
		parts[i] = escapeCSVCell(cell)
	}
	sb.WriteString(strings.Join(parts, ","))
	sb.WriteString("\r\n")
}

func exportLiveDailyEffectiveLiveByGuildZIP(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveDailyEffectiveLiveByGuildPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	if len(req.Headers) == 0 {
		return nil, errExportEmpty
	}
	guildIds, err := resolveRequestedGuildIdsForExport(cmsUserId, parseGuildIdFilters(req.GuildIds))
	if err != nil {
		return nil, err
	}

	nameMap := guilddao.GetNameMapByIds(guildIds)
	grandTotal := 0
	for _, guildId := range guildIds {
		total, _ := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
			GuildIds:      []uint64{guildId},
			FilterByGuild: true,
			LiveDateStart: req.LiveDateStart,
			LiveDateEnd:   req.LiveDateEnd,
			Keyword:       req.Keyword,
			Settled:       req.Settled,
			PageIndex:     1,
			PageSize:      1,
		})
		grandTotal += total
	}
	if grandTotal == 0 {
		return nil, errExportEmpty
	}

	record, filePath, err := createExportFile(".zip")
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		_ = fileexport.Delete(record.ExportID)
		return nil, err
	}

	abort := func(cause error) (*exportResult, error) {
		_ = file.Close()
		_ = fileexport.Delete(record.ExportID)
		if cause == nil {
			cause = errExportEmpty
		}
		return nil, cause
	}

	zw := zip.NewWriter(file)
	exportedRows := 0
	if onProgress != nil {
		onProgress(0, grandTotal)
	}

	for _, guildId := range guildIds {
		select {
		case <-ctx.Done():
			_ = zw.Close()
			return abort(ctx.Err())
		default:
		}

		guildName := ""
		if nameMap != nil {
			guildName = nameMap[guildId]
		}
		entryName := fmt.Sprintf("%d_%s.csv", guildId, sanitizeExportFilePart(guildName))
		entry, createErr := zw.Create(entryName)
		if createErr != nil {
			_ = zw.Close()
			return abort(createErr)
		}
		if _, writeErr := entry.Write([]byte("\uFEFF")); writeErr != nil {
			_ = zw.Close()
			return abort(writeErr)
		}
		headerLine := &strings.Builder{}
		writeCSVRowToBuilder(headerLine, req.Headers)
		if _, writeErr := entry.Write([]byte(headerLine.String())); writeErr != nil {
			_ = zw.Close()
			return abort(writeErr)
		}

		pageIndex := 1
		guildTotal := -1
		guildExported := 0
		for {
			select {
			case <-ctx.Done():
				_ = zw.Close()
				return abort(ctx.Err())
			default:
			}
			total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
				GuildIds:      []uint64{guildId},
				FilterByGuild: true,
				LiveDateStart: req.LiveDateStart,
				LiveDateEnd:   req.LiveDateEnd,
				Keyword:       req.Keyword,
				Settled:       req.Settled,
				PageIndex:     pageIndex,
				PageSize:      defaultExportPageSize,
			})
			if guildTotal < 0 {
				guildTotal = total
			}
			if len(rows) == 0 {
				break
			}
			roomIds := collectDailyAnchorRoomIds(rows)
			nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
			body := &strings.Builder{}
			for _, row := range rows {
				writeCSVRowToBuilder(body, liveDailyEffectiveLiveToCSVRow(row, nicknameMap, req.SettledYesText, req.SettledNoText))
			}
			if _, writeErr := entry.Write([]byte(body.String())); writeErr != nil {
				_ = zw.Close()
				return abort(writeErr)
			}
			guildExported += len(rows)
			exportedRows += len(rows)
			if onProgress != nil {
				onProgress(exportedRows, grandTotal)
			}
			if guildExported >= guildTotal {
				break
			}
			pageIndex++
		}
	}

	if err := zw.Close(); err != nil {
		return abort(err)
	}
	if err := file.Close(); err != nil {
		_ = fileexport.Delete(record.ExportID)
		return nil, err
	}
	if err := fileexport.Publish(record); err != nil {
		_ = fileexport.Delete(record.ExportID)
		return nil, err
	}
	return finalizeExportRecord(record, exportedRows), nil
}

func exportLiveWeeklyUnsettledLiveCSV(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveWeeklyUnsettledLivePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	guildIds, restrict, empty := cmsvis.VisibilityGuildFilterForUser(cmsUserId)
	if empty {
		return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
			return 0, nil
		}, onProgress)
	}
	weekStart, weekEnd := xrtime.WeekDateRange(time.Now())
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
			RoomIds:       liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
			GuildIds:      guildIds,
			FilterByGuild: restrict,
			LiveDateStart: weekStart,
			LiveDateEnd:   weekEnd,
			Keyword:       req.Keyword,
			Settled:       0,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
		})
		roomIds := collectDailyAnchorRoomIds(rows)
		nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
		weeklyIncomeMap := liveroomdao.ListWeeklyUnsettledIncomeByRoomIds(roomIds, weekStart, weekEnd)
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, liveWeeklyUnsettledLiveToCSVRow(row, nicknameMap, weeklyIncomeMap))
		}
		return total, csvRows
	}, onProgress)
}

func liveWeeklyUnsettledLiveToCSVRow(row *liveentity.DailyAnchorEffectiveLive, nicknameMap map[uint64]string, weeklyIncomeMap map[uint64]float64) []string {
	if row == nil {
		return nil
	}
	roomNickname := ""
	if nicknameMap != nil {
		roomNickname = nicknameMap[row.RoomId]
	}
	weeklyIncome := ""
	if weeklyIncomeMap != nil {
		weeklyIncome = formatCSVFloat(weeklyIncomeMap[row.RoomId])
	}
	return []string{
		row.LiveDate,
		formatCSVUint(row.RoomId),
		roomNickname,
		weeklyIncome,
		formatLiveDurationMinutes(row.LiveDuration),
		formatCSVFloat(row.TotalIncome),
	}
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
		formatLiveDurationMinutes(row.LiveDuration),
		formatLiveDurationMinutes(row.TotalLiveDuration),
		formatCSVFloat(row.TotalIncome),
		formatCSVFloat(row.TotalSocialIncome),
	)
	return cells
}

// liveDailyEffectiveLiveToCSVRow 列顺序须与 CMS buildLiveDailyEffectiveLiveListCsvColumns 表头一致:
// 日期,直播间ID,昵称,有效时长,心跳时长,9项收益,结算,创建,更新
func liveDailyEffectiveLiveToCSVRow(row *liveentity.DailyAnchorEffectiveLive, nicknameMap map[uint64]string, settledYesText, settledNoText string) []string {
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
		formatLiveDurationMinutes(row.LiveDuration),
		formatLiveDurationMinutes(row.TotalLiveDuration),
	}
	cells = append(cells, incomeAmountOnlyCSVCells(&row.LiveRoomIncomeAmounts)...)
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
