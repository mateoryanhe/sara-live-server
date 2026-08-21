package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/constants/currency"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dto/cmsexportdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func exportCurrencyLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportCurrencyLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	if req.CurrencyType == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := parseUint64Filter(req.UserId)
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := currencylogdao.CMSList(&currencylogdao.CMSListFilter{
			UserId:       userId,
			CurrencyType: req.CurrencyType,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
			PageIndex:    pageIndex,
			PageSize:     pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, currencyLogToCSVRow(row))
		}
		return total, csvRows
	}, onProgress)
}

func currencyLogToCSVRow(row *userentity.CurrencyLog) []string {
	if row == nil {
		return nil
	}
	actionText := "减少"
	if row.Action == 1 {
		actionText = "增加"
	}
	return []string{
		formatCSVUint(row.ID),
		actionText,
		formatCSVFloat(row.Amount),
		formatCSVFloat(row.Before),
		formatCSVFloat(row.After),
		currency.Reason(row.Reason).Text(currency.LangZHCN),
		row.GameName,
		row.GameCategory,
		currency.BusinessTypeText(row.BusinessType),
		formatCSVTime(row.CreatedAt),
	}
}
