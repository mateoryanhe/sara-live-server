package liverecord

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/entity/live"
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

func toCMSItem(v *entity.LiveRecord) *liverecorddto.CMSLiveRecordItem {
	if v == nil {
		return nil
	}
	item := &liverecorddto.CMSLiveRecordItem{
		Id:                           v.ID,
		AnchorId:                     v.AnchorId,
		StartTime:                    &v.StartTime,
		EndTime:                      v.EndTime,
		TotalAudience:                v.TotalAudience,
		TotalLiveDuration:            v.TotalLiveDuration,
		TotalIncome:                  v.TotalIncome,
		TotalGiftIncome:              v.TotalGiftIncome,
		TotalPaidDanmakuIncome:       v.TotalPaidDanmakuIncome,
		TotalPrivateRoomIncome:       v.TotalPrivateRoomIncome,
		TotalPrivateRoomTicketIncome: v.TotalPrivateRoomTicketIncome,
		TotalPrivateRoomWatchIncome:  v.TotalPrivateRoomWatchIncome,
		TotalVideoCallIncome:         v.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:   v.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome:  v.TotalVideoCallBillingIncome,
		TotalGameBet:                 v.TotalGameBet,
		TotalGiftSender:              v.TotalGiftSender,
		TotalNewFollower:             v.TotalNewFollower,
		CreatedAt:                    &v.CreatedAt,
	}
	if u := userinfodao.GetUserInfoByUserId(v.AnchorId); u != nil {
		item.Nickname = u.Nickname
		item.Avatar = u.Avatar
	}
	return item
}

// GetCMSList CMS分页查询直播记录
func GetCMSList(_ context.Context, req *liverecorddto.CMSLiveRecordListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
		AnchorIds: liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*liverecorddto.CMSLiveRecordItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
