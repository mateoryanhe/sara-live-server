package sys

import (
	"context"
	"time"
	"xr-game-server/dto/sysdto"
)

func GetSysCfg(ctx context.Context, req *sysdto.SysCfgReq) (*sysdto.SysCfgResp, error) {
	return &sysdto.SysCfgResp{
		SysTime: time.Now().UnixMilli(),
	}, nil
}
