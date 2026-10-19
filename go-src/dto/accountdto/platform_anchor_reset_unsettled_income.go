package accountdto

import "github.com/gogf/gf/v2/frame/g"

// BatchResetPlatformAnchorUnsettledIncomeReq CMS批量清零平台主播未结算收益(不写结算/归档日志)
type BatchResetPlatformAnchorUnsettledIncomeReq struct {
	g.Meta    `path:"/batchResetPlatformAnchorUnsettledIncome" method:"post" summary:"CMS批量清零平台主播未结算收益" tags:"账号"`
	AnchorIds []uint64 `json:"anchorIds,string" v:"required|min-length:1#请至少选择一个平台主播" dc:"平台主播ID列表"`
}

// BatchResetPlatformAnchorUnsettledIncomeRes CMS批量清零平台主播未结算收益结果
type BatchResetPlatformAnchorUnsettledIncomeRes struct {
	ResetCount    int      `json:"resetCount" dc:"成功清零数量"`
	NoDataCount   int      `json:"noDataCount" dc:"未结算已为0的数量"`
	FailCount     int      `json:"failCount" dc:"失败数量"`
	FailAnchorIds []string `json:"failAnchorIds" dc:"失败的平台主播ID"`
}
