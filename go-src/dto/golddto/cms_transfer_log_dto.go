package golddto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSCoinMerchantTransferLogListReq CMS 分页查询币商金币转账流水。
type CMSCoinMerchantTransferLogListReq struct {
	g.Meta `path:"/cmsCoinMerchantTransferLogList" method:"post" summary:"CMS查询币商金币转账流水" tags:"金币"`
	httpserver.CMSQueryReq
	CoinMerchantUserId string `json:"coinMerchantUserId" dc:"币商用户ID(可选,留空查全部)"`
	StartTime          int64  `json:"startTime" dc:"开始时间(Unix秒,可选,包含)"`
	EndTime            int64  `json:"endTime" dc:"结束时间(Unix秒,可选,包含)"`
}

// CMSCoinMerchantTransferLogItem CMS 币商金币转账流水列表项。
type CMSCoinMerchantTransferLogItem struct {
	Id                   uint64     `json:"id,string"`
	CoinMerchantUserId   uint64     `json:"coinMerchantUserId,string"`
	CoinMerchantNickname string     `json:"coinMerchantNickname"`
	CoinMerchantAvatar   string     `json:"coinMerchantAvatar"`
	TargetUserId         uint64     `json:"targetUserId,string"`
	TargetNickname       string     `json:"targetNickname"`
	TargetAvatar         string     `json:"targetAvatar"`
	Amount               float64    `json:"amount"`
	SenderGoldBefore     float64    `json:"senderGoldBefore"`
	SenderGoldAfter      float64    `json:"senderGoldAfter"`
	TargetGoldBefore     float64    `json:"targetGoldBefore"`
	TargetGoldAfter      float64    `json:"targetGoldAfter"`
	CreatedAt            *time.Time `json:"createdAt"`
}
