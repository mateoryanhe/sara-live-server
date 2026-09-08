package golddto

import "github.com/gogf/gf/v2/frame/g"

// AppTransferGoldReq App端转赠金币(当前登录用户扣款,目标用户到账)
type AppTransferGoldReq struct {
	g.Meta       `path:"/transferGold" method:"post" summary:"转赠金币给指定用户" tags:"金币"`
	TargetUserId string  `json:"targetUserId" v:"required#目标用户ID不能为空" dc:"接收金币的用户ID"`
	Amount       float64 `json:"amount" v:"required|min:0.01#金币数量不能为空|金币数量至少0.01" dc:"转赠金币数量(最多2位小数)"`
}

type AppTransferGoldRes struct {
	Amount          float64 `json:"amount" dc:"实际转赠数量(已按2位小数规范化)"`
	SenderGold      float64 `json:"senderGold" dc:"转出后当前用户金币余额"`
	TargetUserId    string  `json:"targetUserId" dc:"接收用户ID"`
	TargetUserGold  float64 `json:"targetUserGold" dc:"接收后目标用户金币余额"`
}
