package livegiftlog

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/giftdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/livegiftlogdto"
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

func collectGiftLogUserIds(rows []*entity.LiveGiftLog) []uint64 {
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

func toCMSItem(v *entity.LiveGiftLog, nicknameMap map[uint64]string) *livegiftlogdto.CMSLiveGiftLogItem {
	if v == nil {
		return nil
	}
	item := &livegiftlogdto.CMSLiveGiftLogItem{
		Id:           v.ID,
		RoomId:       v.RoomId,
		LiveRecordId: v.LiveRecordId,
		SenderId:     v.SenderId,
		ReceiverId:   v.ReceiverId,
		GiftId:       v.GiftId,
		Count:        v.Count,
		UnitPrice:    v.UnitPrice,
		TotalCost:    v.TotalCost,
		CreatedAt:    &v.CreatedAt,
	}
	if nicknameMap != nil {
		item.SenderNickname = nicknameMap[v.SenderId]
		item.ReceiverNickname = nicknameMap[v.ReceiverId]
	}
	if g := giftdao.GetGiftById(v.GiftId); g != nil {
		item.GiftName = g.Name
	}
	return item
}

// GetCMSList CMS分页查询礼物流水
func GetCMSList(_ context.Context, req *livegiftlogdto.CMSLiveGiftLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.GiftLogCMSList(&liveroomdao.GiftLogCMSListFilter{
		ReceiverId: parseUint64Filter(req.ReceiverId),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
	})
	nicknameMap := userinfodao.GetNicknameMapByUserIds(collectGiftLogUserIds(rows))
	list := make([]*livegiftlogdto.CMSLiveGiftLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row, nicknameMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
