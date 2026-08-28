package liveroom

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/entity/live"
	"xr-game-server/module/upload"
)


func liveRecordToCMSItem(v *entity.LiveRecord) *liverecorddto.CMSLiveRecordItem {
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
		item.Avatar = upload.ResolveAvatarUrlForUser(v.AnchorId, u.Avatar)
	}
	return item
}

// GetCMSList CMS分页查询直播记录
func GetLiveRecordCMSList(_ context.Context, req *liverecorddto.CMSLiveRecordListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
		AnchorIds:    liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds),
		LiveRecordId: parseUint64Filter(req.LiveRecordId),
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	list := make([]*liverecorddto.CMSLiveRecordItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, liveRecordToCMSItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
