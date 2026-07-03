package sys

import (
	"context"
	"time"
	"xr-game-server/dto/sysdto"
	"xr-game-server/module/livecfg"
)

func GetSysCfg(ctx context.Context, req *sysdto.SysCfgReq) (*sysdto.SysCfgResp, error) {
	return &sysdto.SysCfgResp{
		SysTime:                     time.Now().UnixMilli(),
		PaidDanmakuPrice:            livecfg.GetPaidDanmakuPrice(),
		PrivateRoomFreeWatchSeconds: livecfg.GetPrivateRoomFreeWatchSeconds(),
	}, nil
}
