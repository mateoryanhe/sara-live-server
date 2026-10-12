package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	fiatentity "xr-game-server/entity/fiat"
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
}

type SyncAnchorSalarySocialShareCfgReq struct {
	g.Meta `path:"/syncAnchorSalarySocialShareCfg" method:"post" summary:"同步有底薪社交流水分佣配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的有底薪社交流水分佣配置" dc:"要同步的配置ID列表"`
}

type ReceiveAnchorSalarySocialShareCfgReq struct {
	g.Meta `path:"/receiveAnchorSalarySocialShareCfg" method:"post" summary:"接收有底薪社交流水分佣配置同步" tags:"数据同步"`
	Rows   []*entity.AnchorSalarySocialShareCfg `json:"rows"`
}

type SyncAnchorNoSalaryShareCfgReq struct {
	g.Meta `path:"/syncAnchorNoSalaryShareCfg" method:"post" summary:"同步无底薪社交流水分佣配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的无底薪社交流水分佣配置" dc:"要同步的配置ID列表"`
}

type ReceiveAnchorNoSalaryShareCfgReq struct {
	g.Meta `path:"/receiveAnchorNoSalaryShareCfg" method:"post" summary:"接收无底薪社交流水分佣配置同步" tags:"数据同步"`
	Rows   []*entity.AnchorNoSalaryShareCfg `json:"rows"`
}

type SyncAnchorGameShareCfgReq struct {
	g.Meta `path:"/syncAnchorGameShareCfg" method:"post" summary:"同步游戏流水档位分佣配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的游戏流水档位分佣配置" dc:"要同步的配置ID列表"`
}

type ReceiveAnchorGameShareCfgReq struct {
	g.Meta `path:"/receiveAnchorGameShareCfg" method:"post" summary:"接收游戏流水档位分佣配置同步" tags:"数据同步"`
	Rows   []*entity.AnchorGameShareCfg `json:"rows"`
}

type SyncFiatCurrencyReq struct {
	g.Meta `path:"/syncFiatCurrency" method:"post" summary:"同步法币配置到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的法币配置" dc:"要同步的法币配置ID列表"`
}

type ReceiveFiatCurrencyReq struct {
	g.Meta `path:"/receiveFiatCurrency" method:"post" summary:"接收法币配置同步" tags:"数据同步"`
	Rows   []*fiatentity.FiatCurrencyCfg `json:"rows"`
	Files  []*SyncFileItem               `json:"files"`
}
