package cmsexport

import (
	"context"
	"encoding/json"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/entity/live"
	"xr-game-server/module/cmsvis"
)

func exportLiveRecordCSV(ctx context.Context, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveRecordPayload
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
		total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
			AnchorIds:     liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
			GuildIds:      guildIds,
			FilterByGuild: restrict,
			LiveRecordId:  parseUint64Filter(req.LiveRecordId),
			Keyword:       req.Keyword,
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			PageIndex:     pageIndex,
			PageSize:      pageSize,
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
