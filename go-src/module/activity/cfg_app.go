package activity

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/activitydto"
)

func GetAppFirstRechargeActivityCfg(ctx context.Context, _ *activitydto.AppFirstRechargeActivityCfgReq) (*activitydto.AppFirstRechargeActivityCfgRes, error) {
	res := toAppRes(getCfgCache())
	userId := httpserver.GetAuthId(ctx)
	if userId > 0 {
		res.FirstRecharge = userinfodao.GetUserExtByUserId(userId).FirstRecharge
	}
	return res, nil
}
