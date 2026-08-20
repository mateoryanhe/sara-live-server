package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/entity/live"
)

func exportLiveRecordCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveRecordPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	anchorId := parseUint64Filter(req.AnchorId)
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
			AnchorId:  anchorId,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		})
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, liveRecordToCSVRow(row))
		}
		return total, csvRows
	}, onProgress)
}

func liveRecordToCSVRow(v *entity.LiveRecord) []string {
	if v == nil {
		return nil
	}
	nickname := ""
	if u := userinfodao.GetUserInfoByUserId(v.AnchorId); u != nil {
		nickname = u.Nickname
	}
	return []string{
		formatCSVUint(v.ID),
		formatCSVUint(v.AnchorId),
		nickname,
		formatCSVTime(v.StartTime),
		formatCSVTimePtr(v.EndTime),
		formatCSVUint(v.TotalAudience),
		formatLiveDurationMinutes(v.TotalLiveDuration),
		formatCSVFloat(v.TotalIncome),
		formatCSVFloat(v.TotalGiftIncome),
		formatCSVFloat(v.TotalPaidDanmakuIncome),
		formatCSVFloat(v.TotalVideoCallTicketIncome),
		formatCSVFloat(v.TotalVideoCallBillingIncome),
		formatCSVFloat(v.TotalVideoCallIncome),
		formatCSVUint(v.TotalGiftSender),
		formatCSVUint(v.TotalNewFollower),
		formatCSVFloat(v.TotalGameBet),
		formatCSVTime(v.CreatedAt),
	}
}
