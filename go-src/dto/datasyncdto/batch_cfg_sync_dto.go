package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
	rechargeentity "xr-game-server/entity/recharge"
)

type SyncBatchRes struct {
	Success   bool   `json:"success"`
	RowCount  int    `json:"rowCount"`
	FileCount int    `json:"fileCount"`
	Message   string `json:"message"`
}

type ReceiveBatchRes struct {
	Success   bool `json:"success"`
	RowCount  int  `json:"rowCount"`
	FileCount int  `json:"fileCount"`
}

type SyncBannerReq struct {
	g.Meta `path:"/syncBanner" method:"post" summary:"同步Banner到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的Banner" dc:"要同步的Banner ID列表"`
}

type ReceiveBannerReq struct {
	g.Meta `path:"/receiveBanner" method:"post" summary:"接收Banner同步" tags:"数据同步"`
	Rows   []*entity.HomeBanner `json:"rows"`
	Files  []*SyncFileItem      `json:"files"`
}

type SyncGiftReq struct {
	g.Meta `path:"/syncGift" method:"post" summary:"同步礼物到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的礼物" dc:"要同步的礼物ID列表"`
}

type SyncGiftAssetsReq struct {
	g.Meta `path:"/syncGiftAssets" method:"post" summary:"同步礼物图标与动画资源到目标环境(更新资源文件、库中资源名并刷新缓存)" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的礼物" dc:"要同步资源的礼物ID列表"`
}

type ReceiveGiftReq struct {
	g.Meta     `path:"/receiveGift" method:"post" summary:"接收礼物同步" tags:"数据同步"`
	Rows       []*entity.LiveGift `json:"rows"`
	Files      []*SyncFileItem    `json:"files"`
	AssetsOnly bool               `json:"assetsOnly" dc:"仅更新 icon/animation 资源名,不覆盖其它配置字段"`
}

type SyncRechargeCfgReq struct {
	g.Meta `path:"/syncRechargeCfg" method:"post" summary:"同步充值配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的充值配置" dc:"要同步的充值配置ID列表"`
}

type ReceiveRechargeCfgReq struct {
	g.Meta `path:"/receiveRechargeCfg" method:"post" summary:"接收充值配置同步" tags:"数据同步"`
	Rows   []*rechargeentity.RechargeCfg `json:"rows"`
	Files  []*SyncFileItem               `json:"files"`
}
