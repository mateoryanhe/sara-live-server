package diamonddto

import "github.com/gogf/gf/v2/frame/g"

type CMSAddDiamondReq struct {
	g.Meta `path:"/add" method:"post" summary:"CMS增加用户钻石" tags:"钻石管理"`
	UserId uint64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Amount float64 `json:"amount" v:"required|min:0.0001#钻石数量不能为空|钻石数量必须大于0" dc:"增加数量(正数)"`
}

type CMSAddDiamondRes struct {
	Diamond float64 `json:"diamond" dc:"变更后钻石余额"`
}

type CMSSubDiamondReq struct {
	g.Meta `path:"/sub" method:"post" summary:"CMS扣减用户钻石" tags:"钻石管理"`
	UserId uint64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Amount float64 `json:"amount" v:"required|min:0.0001#钻石数量不能为空|钻石数量必须大于0" dc:"扣减数量(正数)"`
}

type CMSSubDiamondRes struct {
	Diamond float64 `json:"diamond" dc:"变更后钻石余额"`
}

// DiamondPushItem 钻石余额推送载荷
type DiamondPushItem struct {
	Diamond float64 `json:"diamond" dc:"最新钻石余额"`
}
