package accountdto

import "github.com/gogf/gf/v2/frame/g"

const MaxBatchImmediateSettlePlatformAnchors = 100

// BatchImmediateSettlePlatformAnchorsReq CMS批量立即结算平台主播。
type BatchImmediateSettlePlatformAnchorsReq struct {
	g.Meta    `path:"/batchImmediateSettlePlatformAnchors" method:"post" summary:"CMS批量立即结算平台主播" tags:"账号"`
	AnchorIds []uint64 `json:"anchorIds,string" v:"required|min-length:1#请至少选择一个平台主播" dc:"平台主播ID列表"`
}

// BatchImmediateSettlePlatformAnchorsRes CMS批量立即结算平台主播结果。
type BatchImmediateSettlePlatformAnchorsRes struct {
	SettledCount  int      `json:"settledCount" dc:"成功生成主播结算单数量"`
	NoDataCount   int      `json:"noDataCount" dc:"没有待结算数据的平台主播数量"`
	FailCount     int      `json:"failCount" dc:"结算失败数量"`
	FailAnchorIds []string `json:"failAnchorIds" dc:"结算失败的平台主播ID"`
}
