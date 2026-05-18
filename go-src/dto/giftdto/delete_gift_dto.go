package giftdto

import "github.com/gogf/gf/v2/frame/g"

// DeleteGiftReq 删除礼物
type DeleteGiftReq struct {
	g.Meta `path:"/deleteGift" method:"post" summary:"删除礼物" tags:"礼物管理"`
	ID     uint64 `json:"id" v:"required#礼物ID不能为空" dc:"礼物ID"`
}

type DeleteGiftRes struct {
	Success bool `json:"success"`
}
