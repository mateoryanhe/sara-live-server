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
	if req.CurrencyType > 2 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := parseUint64Filter(req.UserId)
	langCode := currency.ParseLang(req.Lang)
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
			csvRows = append(csvRows, currencyLogToCSVRow(row, langCode))
		}
		return total, csvRows
	}, onProgress)
}

func currencyTypeText(currencyType uint8, langCode currency.Lang) string {
	isDiamond := currencyType == 2
	switch langCode {
	case currency.LangZHCN:
		if isDiamond {
			return "钻石"
		}
		return "金币"
	case currency.LangZHTW:
		if isDiamond {
			return "鑽石"
		}
		return "金幣"
	default:
		if isDiamond {
			return "Diamond"
		}
		return "Gold"
	}
}

func currencyLogToCSVRow(row *userentity.CurrencyLog, langCode currency.Lang) []string {
	if row == nil {
		return nil
	}
	actionText := "减少"
	if row.Action == 1 {
		actionText = "增加"
	}
	if langCode != currency.LangZHCN && langCode != currency.LangZHTW {
		actionText = "Decrease"
		if row.Action == 1 {
			actionText = "Increase"
		}
	} else if langCode == currency.LangZHTW {
		actionText = "減少"
		if row.Action == 1 {
			actionText = "增加"
		}
	}
	return []string{
		formatCSVUint(row.ID),
		currencyTypeText(row.Type, langCode),
		actionText,
		formatCSVFloat(row.Amount),
		formatCSVFloat(row.Before),
		formatCSVFloat(row.After),
		currency.Reason(row.Reason).Text(langCode),
		formatCSVTime(row.CreatedAt),
	}
}
