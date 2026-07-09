package call

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
)

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func collectCallOrderUserIds(rows []*entity.CallOrder) []uint64 {
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

func callTypeText(v uint8) string {
	switch v {
	case entity.CallOrderTypeVoice:
		return "语音"
	case entity.CallOrderTypeVideo:
		return "视频"
	default:
		return "未知"
	}
}

func callSourceText(v uint8) string {
	switch v {
	case entity.CallOrderSourceLiveRoom:
		return "直播间"
	case entity.CallOrderSourcePrivateMessage:
		return "私信"
	default:
		return "未知"
	}
}

func callOrderStatusText(o *entity.CallOrder) string {
	if o == nil {
		return ""
	}
	if o.HasEnded() {
		return "已结束"
	}
	if o.IsCallStarted() {
		return "通话中"
	}
	if o.IsAccepted() {
		return "已接听"
	}
	if o.IsCalling() {
		return "呼叫中"
	}
	return "未知"
}

func toCMSVideoCallItem(v *entity.CallOrder, nicknameMap map[uint64]string) *calldto.CMSVideoCallLogItem {
	if v == nil {
		return nil
	}
	item := &calldto.CMSVideoCallLogItem{
		Id:              v.ID,
		CallerId:        v.CallerId,
		ReceiverId:      v.ReceiverId,
		CallType:        v.CallType,
		CallTypeText:    callTypeText(v.CallType),
		Source:          v.Source,
		SourceText:      callSourceText(v.Source),
		StatusText:      callOrderStatusText(v),
		CallStartTime:   &v.CallStartTime,
		AnswerTime:      v.AnswerTime,
		OrderEndTime:    v.OrderEndTime,
		CallDuration:    v.CallDuration,
		TicketPrice:     v.TicketPrice,
		PricePerMinute:  v.PricePerMinute,
		TotalCost:       v.TotalCost,
		BillingDuration: v.BillingDuration,
		ChargeTime:      v.ChargeTime,
		CreatedAt:       &v.CreatedAt,
	}
	if nicknameMap != nil {
		item.CallerNickname = nicknameMap[v.CallerId]
		item.ReceiverNickname = nicknameMap[v.ReceiverId]
	}
	return item
}

// GetCMSVideoCallLogList CMS分页查询视频通话日志
func GetCMSVideoCallLogList(_ context.Context, req *calldto.CMSVideoCallLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := calldao.CallOrderCMSList(&calldao.CallOrderCMSListFilter{
		CallerId:   parseUint64Filter(req.CallerId),
		ReceiverId: parseUint64Filter(req.ReceiverId),
		Source:     req.Source,
		CallType:   entity.CallOrderTypeVideo,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
	})
	nicknameMap := userinfodao.GetNicknameMapByUserIds(collectCallOrderUserIds(rows))
	list := make([]*calldto.CMSVideoCallLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSVideoCallItem(row, nicknameMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
