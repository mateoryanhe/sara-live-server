package shortvideo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
)

// GetCMSAuthorSettlementLogList CMS 分页查询短视频作者结算日志
func GetCMSAuthorSettlementLogList(_ context.Context, req *shortvideodto.CMSAuthorSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := shortvideodao.AuthorSettlementLogCMSList(&shortvideodao.AuthorSettlementLogCMSListFilter{
		UserId:    parseUint64Filter(req.UserId),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*shortvideodto.CMSAuthorSettlementLogItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toCMSAuthorSettlementLogItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toCMSAuthorSettlementLogItem(row *entity.ShortVideoAuthorSettlementLog) *shortvideodto.CMSAuthorSettlementLogItem {
	if row == nil {
		return nil
	}
	createdAt := row.CreatedAt
	return &shortvideodto.CMSAuthorSettlementLogItem{
		Id:                 row.ID,
		UserId:             row.UserId,
		UnsettledIncome:    row.UnsettledIncome,
		SettlementDiamond:  row.SettlementDiamond,
		AnchorSharePercent: row.AnchorSharePercent,
		CreatedAt:          &createdAt,
	}
}
