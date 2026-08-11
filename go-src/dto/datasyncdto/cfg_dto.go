package datasyncdto

import "github.com/gogf/gf/v2/frame/g"

type GetDataSyncCfgReq struct {
	g.Meta `path:"/getDataSyncCfg" method:"post" summary:"查询数据同步配置" tags:"数据同步配置"`
}

type DataSyncCfgItem struct {
	ID            string `json:"id"`
	TargetApiBase string `json:"targetApiBase"`
	Token         string `json:"token"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type GetDataSyncCfgRes struct {
	Cfg *DataSyncCfgItem `json:"cfg"`
}

type SaveDataSyncCfgReq struct {
	g.Meta        `path:"/saveDataSyncCfg" method:"post" summary:"保存数据同步配置" tags:"数据同步配置"`
	ID            uint64 `json:"id" dc:"配置ID,新建传0"`
	TargetApiBase string `json:"targetApiBase" v:"required#目标API根地址不能为空" dc:"目标API根地址"`
	Token         string `json:"token" v:"required#同步Token不能为空" dc:"同步Token"`
}

type SaveDataSyncCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
