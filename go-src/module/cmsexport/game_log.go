package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/gamewindao"
	"xr-game-server/dto/cmsexportdto"
	gameentity "xr-game-server/entity/game"
)

func exportGameBetLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGameBetLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := gamebetdao.CMSList(&gamebetdao.CMSListFilter{
			UserId:       parseUint64Filter(req.UserId),
			GameCode:     req.GameCode,
			OrderId:      req.OrderId,
			PlatformType: req.PlatformType,
			PageIndex:    pageIndex,
			PageSize:     pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, gameBetLogToCSVRow(row))
		}
		return total, csvRows
	}, onProgress)
}

func exportGameWinLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportGameWinLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := gamewindao.CMSList(&gamewindao.CMSListFilter{
			UserId:       parseUint64Filter(req.UserId),
			GameCode:     req.GameCode,
			OrderId:      req.OrderId,
			PlatformType: req.PlatformType,
			PageIndex:    pageIndex,
			PageSize:     pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, gameWinLogToCSVRow(row))
		}
		return total, csvRows
	}, onProgress)
}

func gameBetLogToCSVRow(row *gameentity.GameBetLog) []string {
	if row == nil {
		return nil
	}
	return []string{
		formatCSVUint(row.ID),
		formatCSVUint(row.UserId),
		row.GameCode,
		row.NameEn,
		formatCSVFloat(row.Amount),
		row.PlatformType,
		formatCSVUint(row.LiveRoomId),
		formatCSVUint(row.LiveRecordId),
		row.OrderId,
		formatCSVTime(row.CreatedAt),
	}
}

func gameWinLogToCSVRow(row *gameentity.GameWinLog) []string {
	if row == nil {
		return nil
	}
	return []string{
		formatCSVUint(row.ID),
		formatCSVUint(row.UserId),
		row.GameCode,
		row.NameEn,
		formatCSVFloat(row.Amount),
		row.PlatformType,
		row.OrderId,
		formatCSVTime(row.CreatedAt),
	}
}
