package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	sysentity "xr-game-server/entity/sys"
)

type SyncCfEmailCfgReq struct {
	g.Meta `path:"/syncCfEmailCfg" method:"post" summary:"同步邮件SMTP发信配置到目标环境" tags:"数据同步"`
}

type ReceiveCfEmailCfgReq struct {
	g.Meta `path:"/receiveCfEmailCfg" method:"post" summary:"接收邮件SMTP发信配置同步" tags:"数据同步"`
	Row    *sysentity.CfEmailCfg `json:"row"`
}
