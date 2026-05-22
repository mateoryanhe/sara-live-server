package golddto

import "github.com/gogf/gf/v2/frame/g"

type CMSAddGoldReq struct {
	g.Meta `path:"/add" method:"post" summary:"CMS增加用户金币" tags:"金币管理"`
	UserId uint64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Amount float64 `json:"amount" v:"required|min:0.0001#金币数量不能为空|金币数量必须大于0" dc:"增加数量(正数)"`
}

type CMSAddGoldRes struct {
	Gold float64 `json:"gold" dc:"变更后金币余额"`
}

type CMSSubGoldReq struct {
	g.Meta `path:"/sub" method:"post" summary:"CMS扣减用户金币" tags:"金币管理"`
	UserId uint64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Amount float64 `json:"amount" v:"required|min:0.0001#金币数量不能为空|金币数量必须大于0" dc:"扣减数量(正数)"`
}

type CMSSubGoldRes struct {
	Gold float64 `json:"gold" dc:"变更后金币余额"`
}

// GoldPushItem 金币余额推送载荷
type GoldPushItem struct {
	Gold string `json:"gold" dc:"最新金币余额"`
}
