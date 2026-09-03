package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/entity/shortvideo"
)

func exportShortVideoAuthorSettlementLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportShortVideoAuthorSettlementLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := shortvideodao.AuthorSettlementLogCMSList(&shortvideodao.AuthorSettlementLogCMSListFilter{
			UserId:    parseUint64Filter(req.UserId),
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, shortVideoAuthorSettlementLogToCSVRow(row))
		}
		return total, csvRows
	}, onProgress)
}

func shortVideoAuthorSettlementLogToCSVRow(row *entity.ShortVideoAuthorSettlementLog) []string {
	if row == nil {
		return nil
	}
	return []string{
		formatCSVUint(row.ID),
		formatCSVFloat(row.UnsettledIncome),
		formatCSVFloat(row.SettlementDiamond),
		formatCSVFloat(row.AnchorSharePercent),
		formatCSVTime(row.CreatedAt),
	}
}
