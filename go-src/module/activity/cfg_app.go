package activity

import (
	"context"

	"xr-game-server/dto/activitydto"
)

func GetAppFirstRechargeActivityCfg(_ context.Context, _ *activitydto.AppFirstRechargeActivityCfgReq) (*activitydto.AppFirstRechargeActivityCfgRes, error) {
	return toAppRes(getCfgCache()), nil
}
