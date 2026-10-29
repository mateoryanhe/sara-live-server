package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

type SyncSofiePrivacyPolicyCfgReq struct {
	g.Meta `path:"/syncSofiePrivacyPolicyCfg" method:"post" summary:"同步 Sofie 隐私政策配置到目标环境" tags:"数据同步"`
}

type ReceiveSofiePrivacyPolicyCfgReq struct {
	g.Meta `path:"/receiveSofiePrivacyPolicyCfg" method:"post" summary:"接收 Sofie 隐私政策配置同步" tags:"数据同步"`
	Row    *entity.SofiePrivacyPolicyCfg `json:"row"`
}
