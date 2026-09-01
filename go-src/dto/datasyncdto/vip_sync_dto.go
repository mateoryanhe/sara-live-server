package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

type SyncFileItem struct {
	Name string `json:"name"`
	Data string `json:"data"` // base64
}

type SyncVipCfgReq struct {
	g.Meta `path:"/syncVipCfg" method:"post" summary:"同步VIP配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的配置" dc:"要同步的VIP配置ID列表"`
}

type SyncVipCfgAssetsReq struct {
	g.Meta `path:"/syncVipCfgAssets" method:"post" summary:"同步VIP图标与动画资源到目标环境(更新资源文件、库中资源名并刷新缓存)" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的配置" dc:"要同步资源的VIP配置ID列表"`
}

type SyncVipCfgRes struct {
	Success   bool   `json:"success"`
	RowCount  int    `json:"rowCount"`
	FileCount int    `json:"fileCount"`
	Message   string `json:"message"`
}

type ReceiveVipCfgReq struct {
	g.Meta     `path:"/receiveVipCfg" method:"post" summary:"接收VIP配置同步" tags:"数据同步"`
	Rows       []*entity.VipCfg `json:"rows"`
	Files      []*SyncFileItem  `json:"files"`
	AssetsOnly bool             `json:"assetsOnly" dc:"仅更新图标/动画等资源名字段,不覆盖其它配置"`
}

type ReceiveVipCfgRes struct {
	Success   bool `json:"success"`
	RowCount  int  `json:"rowCount"`
	FileCount int  `json:"fileCount"`
}
