package giftdto

import "github.com/gogf/gf/v2/frame/g"

// ReceivedGiftIdsReq App端查询主播收到过的礼物列表
type ReceivedGiftIdsReq struct {
	g.Meta   `path:"/receivedGiftIds" method:"post" summary:"查询主播收到过的礼物ID与数量" tags:"礼物"`
	AnchorId uint64 `json:"anchorId,string" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

// ReceivedGiftItem 单个礼物累计
type ReceivedGiftItem struct {
	GiftId uint64 `json:"giftId,string" dc:"礼物ID"`
	Count  uint64 `json:"count" dc:"累计收到数量"`
}

// ReceivedGiftIdsRes 主播收到过的礼物(按礼物去重,最多2000,含数量)
type ReceivedGiftIdsRes struct {
	List []*ReceivedGiftItem `json:"list" dc:"收礼列表"`
}
