package liverevenue

import (
	"context"
	"strconv"
	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/giftdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverevenuedto"
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

func collectRevenueLogUserIds(rows []*entity.LiveRevenueLog) []uint64 {
	userIds := make([]uint64, 0, len(rows)*2)
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.SenderId > 0 {
			userIds = append(userIds, row.SenderId)
		}
		if row.ReceiverId > 0 {
			userIds = append(userIds, row.ReceiverId)
		}
	}
	return userIds
}

func toCMSItem(v *entity.LiveRevenueLog, nicknameMap map[uint64]string) *liverevenuedto.CMSLiveRevenueLogItem {
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
		ReceiverId:      v.ReceiverId,
		BizId:           v.BizId,
		Count:           v.Count,
		UnitPrice:       v.UnitPrice,
		TotalAmount:     v.TotalAmount,
		CreatedAt:       &v.CreatedAt,
	}
	if nicknameMap != nil {
		item.SenderNickname = nicknameMap[v.SenderId]
		item.ReceiverNickname = nicknameMap[v.ReceiverId]
	}
	if v.RevenueType == uint8(liverevenueconst.Gift) {
		if g := giftdao.GetGiftById(v.BizId); g != nil {
			item.BizName = g.Name
		}
	}
	return item
}

// GetCMSList CMS分页查询直播收益流水
func GetCMSList(_ context.Context, req *liverevenuedto.CMSLiveRevenueLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.RevenueLogCMSList(&liveroomdao.RevenueLogCMSListFilter{
		ReceiverId:  parseUint64Filter(req.ReceiverId),
		RevenueType: req.RevenueType,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		PageIndex:   req.PageIndex,
		PageSize:    req.PageSize,
	})
	nicknameMap := userinfodao.GetNicknameMapByUserIds(collectRevenueLogUserIds(rows))
	list := make([]*liverevenuedto.CMSLiveRevenueLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row, nicknameMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
