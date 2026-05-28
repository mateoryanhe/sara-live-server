package liverecord

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverecorddto"
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

func collectAnchorIds(rows []*entity.LiveRecord) []uint64 {
	userIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.AnchorId == 0 {
			continue
		}
		userIds = append(userIds, row.AnchorId)
	}
	return userIds
}

func toCMSItem(v *entity.LiveRecord, nicknameMap map[uint64]string) *liverecorddto.CMSLiveRecordItem {
	if v == nil {
		return nil
	}
	item := &liverecorddto.CMSLiveRecordItem{
		Id:                v.ID,
		AnchorId:          v.AnchorId,
		StartTime:         &v.StartTime,
		EndTime:           v.EndTime,
		TotalAudience:     v.TotalAudience,
		TotalLiveDuration: v.TotalLiveDuration,
		TotalIncome:       v.TotalIncome,
		TotalGameBet:      v.TotalGameBet,
		CreatedAt:         &v.CreatedAt,
	}
	if nicknameMap != nil {
		item.Nickname = nicknameMap[v.AnchorId]
	}
	return item
}

// GetCMSList CMS分页查询直播记录
func GetCMSList(_ context.Context, req *liverecorddto.CMSLiveRecordListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
		AnchorId:  parseUint64Filter(req.AnchorId),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	nicknameMap := userinfodao.GetNicknameMapByUserIds(collectAnchorIds(rows))
	list := make([]*liverecorddto.CMSLiveRecordItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSItem(row, nicknameMap))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
