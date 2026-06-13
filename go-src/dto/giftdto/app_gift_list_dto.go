package giftdto

import "github.com/gogf/gf/v2/frame/g"

// AppGiftListReq App端查询礼物列表(仅返回已上架礼物,带缓存)
type AppGiftListReq struct {
	g.Meta `path:"/appGiftList" method:"post" summary:"App查询礼物列表(已上架)" tags:"礼物"`
}

// AppGiftItem App端礼物列表元素
type AppGiftItem struct {
	ID          uint64  `json:"id,string"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Animation   string  `json:"animation"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Sort        int     `json:"sort"`
	PublishedAt int64   `json:"publishedAt" dc:"发布时间(秒,0表示未设置)"`
	Description string  `json:"description"`
}

// AppGiftListRes App端礼物列表响应
type AppGiftListRes struct {
	List []*AppGiftItem `json:"list" dc:"礼物列表"`
}
