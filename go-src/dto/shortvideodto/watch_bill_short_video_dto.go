package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// WatchBillShortVideoReq App端短视频观看扣费(建议每5秒调用一次)
type WatchBillShortVideoReq struct {
	g.Meta  `path:"/watchBillShortVideo" method:"post" summary:"短视频观看扣费" tags:"短视频"`
	VideoId uint64 `json:"videoId,string" v:"required#视频ID不能为空" dc:"短视频ID"`
}

type WatchBillShortVideoRes struct {
	Deducted          float64 `json:"deducted" dc:"本次扣除钻石数"`
	Diamond           float64 `json:"diamond" dc:"扣费后钻石余额"`
	BilledSeconds     uint32  `json:"billedSeconds" dc:"累计已计费观看秒数"`
	ChargeableSeconds uint32  `json:"chargeableSeconds" dc:"本次计费秒数(扣除免费时长后)"`
	CanContinue       bool    `json:"canContinue" dc:"余额是否足够继续观看"`
}
