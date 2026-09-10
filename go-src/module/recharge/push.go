package recharge

import (
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity/recharge"
)

// pushRechargeSuccessToApp 充值到账成功推送(每次到账都会推)
func pushRechargeSuccessToApp(userId uint64, order *entity.RechargeOrder, goldBalance float64) {
	if userId == 0 || order == nil {
		return
	}
	push.Data(userId, cmd.RechargeSuccessPush, &rechargeorderdto.RechargeSuccessPushItem{
		OrderId:     strconv.FormatUint(order.ID, 10),
		Gold:        uint64(order.Gold),
		GoldBalance: uint64(goldBalance),
		Price:       order.Price,
		PayChannel:  order.PayChannel,
	})
}
