package liverevenue

import (
	"context"
	"strconv"
	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverevenuedto"
	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
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

func collectRevenueLogUserIds(rows []*entity.LiveRevenueLog) []uint64 {
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

func revenueLogStatusText(status uint8) string {
	switch status {
	case entity.LiveRevenueLogStatusRefunded:
		return "已退款"
	default:
		return "正常"
	}
}

func toCMSItem(v *entity.LiveRevenueLog, profileMap map[uint64]*userentity.UserInfo) *liverevenuedto.CMSLiveRevenueLogItem {
	if v == nil {
		return nil
	}
	item := &liverevenuedto.CMSLiveRevenueLogItem{
		Id:              v.ID,
		RevenueType:     v.RevenueType,
		RevenueTypeText: liverevenueconst.Text(liverevenueconst.Type(v.RevenueType)),
		RoomId:          v.RoomId,
		LiveRecordId:    v.LiveRecordId,
		SenderId:        v.SenderId,
		ReceiverId:      v.RoomId,
		BizId:           v.BizId,
		Count:           v.Count,
		UnitPrice:       v.UnitPrice,
		TotalAmount:     v.TotalAmount,
		Status:          v.Status,
		StatusText:      revenueLogStatusText(v.Status),
		CreatedAt:       &v.CreatedAt,
	}
	if profileMap != nil {
		if sender := profileMap[v.SenderId]; sender != nil {
			item.SenderNickname = sender.Nickname
			item.SenderAvatar = upload.ResolveAvatarUrlForUser(v.SenderId, sender.Avatar)
		}
		if receiver := profileMap[v.RoomId]; receiver != nil {
			item.ReceiverNickname = receiver.Nickname
			item.ReceiverAvatar = upload.ResolveAvatarUrlForUser(v.RoomId, receiver.Avatar)
		}
	}
	if v.RevenueType == uint8(liverevenueconst.Gift) {
		if g := cfgdao.GetGiftById(v.BizId); g != nil {
			item.BizName = g.Name
		}
	}
	return item
}

// GetCMSList CMS分页查询直播收益流水
func GetCMSList(_ context.Context, req *liverevenuedto.CMSLiveRevenueLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.RevenueLogCMSList(&liveroomdao.RevenueLogCMSListFilter{
		ReceiverIds:  liveroomdao.ParseRevenueLogReceiverIds(req.ReceiverId, req.PlatformAnchorId, req.GuildAnchorId, req.ReceiverIds),
		LiveRecordId: parseUint64Filter(req.LiveRecordId),
		Keyword:      req.Keyword,
		RevenueType:  req.RevenueType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		PageIndex:    req.PageIndex,
		PageSize:     req.PageSize,
	})
	profileMap := userinfodao.GetUserProfileMapByUserIds(collectRevenueLogUserIds(rows))
	list := make([]*liverevenuedto.CMSLiveRevenueLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row, profileMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
