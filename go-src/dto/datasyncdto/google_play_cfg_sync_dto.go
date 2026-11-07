package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/recharge"
)

type SyncGooglePlayCfgReq struct {
	g.Meta `path:"/syncGooglePlayCfg" method:"post" summary:"同步 Google Play 配置到目标环境" tags:"数据同步"`
}

type ReceiveGooglePlayCfgReq struct {
	g.Meta `path:"/receiveGooglePlayCfg" method:"post" summary:"接收 Google Play 配置同步" tags:"数据同步"`
	Row    *entity.GooglePlayCfg `json:"row"`
}
