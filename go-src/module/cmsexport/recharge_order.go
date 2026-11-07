package cmsexport

import (
	"context"
	"encoding/json"
	"strconv"

	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	entity "xr-game-server/entity/recharge"
	"xr-game-server/module/ipgeo"
)

func exportRechargeOrderCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportRechargeOrderPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := rechargeorderdao.CMSList(&rechargeorderdao.CMSListFilter{
			UserId:    parseUint64Filter(req.UserId),
			OrderId:   parseUint64Filter(req.OrderId),
			StatusVal: resolveRechargeOrderStatusFilter(req.StatusFilter),
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, rechargeOrderToCSVRow(row, req.StatusLabels, req.SourceLabels))
		}
		return total, csvRows
	}, onProgress)
}

// resolveRechargeOrderStatusFilter 与 module/recharge.resolveStatusFilter 一致:
// 0=全部,1=待支付,2=已完成,3=超时取消,4=失败 → DAO statusVal(<0=不过滤)
func resolveRechargeOrderStatusFilter(f int) int {
	switch f {
	case 1:
		return int(entity.RechargeOrderStatusPending)
	case 2:
		return int(entity.RechargeOrderStatusCompleted)
	case 3:
		return int(entity.RechargeOrderStatusCancelled)
	case 4:
		return int(entity.RechargeOrderStatusFailed)
	default:
		return -1
	}
}

func rechargeOrderToCSVRow(row *entity.RechargeOrder, statusLabels, sourceLabels map[string]string) []string {
	if row == nil {
		return nil
	}
	nickname := ""
	if u := userinfodao.GetUserInfoByUserId(row.UserId); u != nil {
		nickname = u.Nickname
	}
	cfgId := ""
	if row.CfgId > 0 {
		cfgId = formatCSVUint(row.CfgId)
	}
	payAmount := ""
	if row.PayAmount != 0 {
		payAmount = formatCSVFloat(row.PayAmount)
	}
	currency := row.Currency
	if currency == "" {
		currency = "-"
	}
	return []string{
		formatCSVUint(row.ID),
		ipgeo.FormatCountryDisplay(row.UserCountry),
		formatCSVUint(row.UserId),
		nickname,
		cfgId,
		formatCSVFloat(row.Price),
		payAmount,
		currency,
		formatCSVFloat(row.Gold),
		lookupExportLabel(statusLabels, row.Status),
		lookupExportLabel(sourceLabels, row.Source),
		row.ThirdOrderId,
		row.Remark,
		formatCSVTime(row.CreatedAt),
		formatCSVTime(row.PaidAt),
	}
}

func lookupExportLabel(labels map[string]string, code uint8) string {
	key := strconv.FormatUint(uint64(code), 10)
	if labels != nil {
		if text, ok := labels[key]; ok && text != "" {
			return text
		}
	}
	return key
}
