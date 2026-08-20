package activity

import (
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dto/activitydto"
	"xr-game-server/entity/recharge"
)

func pushFirstRechargeSuccessToApp(userId uint64, order *entity.RechargeOrder) {
	if userId == 0 || order == nil {
		return
	}
	push.Data(userId, cmd.FirstRechargeSuccessPush, &activitydto.FirstRechargeSuccessPushItem{
		FirstRecharge: false,
		Gold:          uint64(order.Gold),
		OrderId:       strconv.FormatUint(order.ID, 10),
	})
}
