package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	activityentity "xr-game-server/entity/activity"
)

type SyncFirstRechargeActivityCfgReq struct {
	g.Meta `path:"/syncFirstRechargeActivityCfg" method:"post" summary:"同步首充活动配置到目标环境" tags:"数据同步"`
}

type SyncFirstRechargeActivityCfgRes struct {
	Success   bool   `json:"success"`
	RowCount  int    `json:"rowCount"`
	FileCount int    `json:"fileCount"`
	Message   string `json:"message"`
}

type ReceiveFirstRechargeActivityCfgReq struct {
	g.Meta     `path:"/receiveFirstRechargeActivityCfg" method:"post" summary:"接收首充活动配置同步" tags:"数据同步"`
	Row        *activityentity.FirstRechargeActivityCfg           `json:"row"`
	Privileges []*activityentity.FirstRechargeActivityPrivilege `json:"privileges"`
	Files      []*SyncFileItem                                    `json:"files"`
}

type ReceiveFirstRechargeActivityCfgRes struct {
	Success   bool `json:"success"`
	RowCount  int  `json:"rowCount"`
	FileCount int  `json:"fileCount"`
}
