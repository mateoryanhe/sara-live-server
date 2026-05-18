package giftdto

import "github.com/gogf/gf/v2/frame/g"

// OnShelfGiftReq 上架礼物
type OnShelfGiftReq struct {
	g.Meta `path:"/onShelfGift" method:"post" summary:"上架礼物" tags:"礼物管理"`
	ID     uint64 `json:"id" v:"required#礼物ID不能为空" dc:"礼物ID"`
}

type OnShelfGiftRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status" dc:"状态(0-下架,1-上架)"`
}

// OffShelfGiftReq 下架礼物
type OffShelfGiftReq struct {
	g.Meta `path:"/offShelfGift" method:"post" summary:"下架礼物" tags:"礼物管理"`
	ID     uint64 `json:"id" v:"required#礼物ID不能为空" dc:"礼物ID"`
}

type OffShelfGiftRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status" dc:"状态(0-下架,1-上架)"`
}
