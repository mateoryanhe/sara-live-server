package cmsexport

import (
	"context"
	"encoding/json"
	"strconv"

	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/dao/calldao"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/cmsexportdto"
	callentity "xr-game-server/entity/call"
	liveentity "xr-game-server/entity/live"
)

func exportLiveRevenueLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportLiveRevenueLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	receiverIds := liveroomdao.ParseRevenueLogReceiverIds(req.ReceiverId, req.PlatformAnchorId, req.GuildAnchorId, req.ReceiverIds)
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		total, rows := liveroomdao.RevenueLogCMSList(&liveroomdao.RevenueLogCMSListFilter{
			ReceiverIds:  receiverIds,
			LiveRecordId: parseUint64Filter(req.LiveRecordId),
			Keyword:      req.Keyword,
			RevenueType:  req.RevenueType,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			PageIndex:   pageIndex,
			PageSize:    pageSize,
		})
		nicknameMap := userinfodao.GetNicknameMapByUserIds(collectRevenueLogUserIds(rows))
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, liveRevenueLogToCSVRow(row, nicknameMap))
		}
		return total, csvRows
	}, onProgress)
}

func exportVideoCallLogCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportVideoCallLogPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		receiverIds := liveroomdao.ParseRevenueLogReceiverIds(req.ReceiverId, req.PlatformAnchorId, req.GuildAnchorId, req.ReceiverIds)
		total, rows := calldao.CallOrderCMSList(&calldao.CallOrderCMSListFilter{
			CallerId:    parseUint64Filter(req.CallerId),
			ReceiverIds: receiverIds,
			Source:     req.Source,
			Status:     req.Status,
			CallType:   callentity.CallOrderTypeVideo,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
			PageIndex:  pageIndex,
			PageSize:   pageSize,
		})
		nicknameMap := userinfodao.GetNicknameMapByUserIds(collectCallOrderUserIds(rows))
		csvRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			csvRows = append(csvRows, videoCallLogToCSVRow(row, nicknameMap))
		}
		return total, csvRows
	}, onProgress)
}

func collectRevenueLogUserIds(rows []*liveentity.LiveRevenueLog) []uint64 {
	userIds := make([]uint64, 0, len(rows)*2)
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.SenderId > 0 {
			userIds = append(userIds, row.SenderId)
		}
		if row.RoomId > 0 {
			userIds = append(userIds, row.RoomId)
		}
	}
	return userIds
}

func collectCallOrderUserIds(rows []*callentity.CallOrder) []uint64 {
	userIds := make([]uint64, 0, len(rows)*2)
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.CallerId > 0 {
			userIds = append(userIds, row.CallerId)
		}
		if row.ReceiverId > 0 {
			userIds = append(userIds, row.ReceiverId)
		}
	}
	return userIds
}

func liveRevenueLogToCSVRow(v *liveentity.LiveRevenueLog, nicknameMap map[uint64]string) []string {
	if v == nil {
		return nil
	}
	bizName := ""
	if v.RevenueType == uint8(liverevenueconst.Gift) {
		if g := cfgdao.GetGiftById(v.BizId); g != nil {
			bizName = g.Name
		}
	}
	senderNickname := ""
	receiverNickname := ""
	if nicknameMap != nil {
		senderNickname = nicknameMap[v.SenderId]
		receiverNickname = nicknameMap[v.RoomId]
	}
	return []string{
		formatCSVUint(v.ID),
		liverevenueconst.Text(liverevenueconst.Type(v.RevenueType)),
		formatCSVUint(v.RoomId),
		formatCSVUint(v.LiveRecordId),
		formatCSVUint(v.SenderId),
		senderNickname,
		formatCSVUint(v.RoomId),
		receiverNickname,
		formatCSVUint(v.BizId),
		bizName,
		strconv.Itoa(v.Count),
		formatCSVFloat(v.UnitPrice),
		formatCSVFloat(v.TotalAmount),
		revenueLogStatusText(v.Status),
		formatCSVTime(v.CreatedAt),
	}
}

func revenueLogStatusText(status uint8) string {
	switch status {
	case liveentity.LiveRevenueLogStatusRefunded:
		return "已退款"
	default:
		return "正常"
	}
}

func videoCallLogToCSVRow(v *callentity.CallOrder, nicknameMap map[uint64]string) []string {
	if v == nil {
		return nil
	}
	callerNickname := ""
	receiverNickname := ""
	if nicknameMap != nil {
		callerNickname = nicknameMap[v.CallerId]
		receiverNickname = nicknameMap[v.ReceiverId]
	}
	return []string{
		formatCSVUint(v.ID),
		v.StatusText(),
		callSourceText(v.Source),
		formatCSVUint(v.CallerId),
		callerNickname,
		formatCSVUint(v.ReceiverId),
		receiverNickname,
		formatCSVTimePtr(&v.CallStartTime),
		formatCSVTimePtr(v.AnswerTime),
		formatCSVTimePtr(v.CallerHeartTime),
		formatCSVTimePtr(v.ReceiverHeartTime),
		formatCSVTimePtr(v.OrderEndTime),
		strconv.FormatUint(uint64(v.CallDuration), 10),
		formatCSVFloat(v.TicketPrice),
		formatCSVFloat(v.PricePerMinute),
		strconv.FormatUint(uint64(v.BillingDuration), 10),
		formatCSVFloat(v.TotalCost),
		formatCSVTimePtr(v.ChargeTime),
		formatCSVTime(v.CreatedAt),
	}
}

func callSourceText(v uint8) string {
	switch v {
	case callentity.CallOrderSourceLiveRoom:
		return "直播间"
	case callentity.CallOrderSourcePrivateMessage:
		return "私信"
	default:
		return "未知"
	}
}
