package controller

import (
	"xr-game-server/core/httpserver"
)

func Init() {
	initAuthApi()
	initAccountController()
	initGoldController()
	initDiamondController()
	initIndex()
	initWebSocket()
	initRank()
	initGlobalCfgController()
	initRoleController()             // 添加角色控制器初始化
	initCMSUserController()          // 添加CMS用户控制器初始化
	initUserInfoController()         // 用户基础信息(App)
	initGuildController()            // 直播工会管理(CMS)
	initGuildAppController()         // 直播工会查询(App)
	initLiveRoomAppController()      // 直播间(App)
	initGiftController()             // 礼物配置(CMS)
	initGiftAppController()          // 礼物列表(App)
	initBannerController()           // 首页Banner(CMS)
	initBannerAppController()        // 首页Banner(App)
	initShortVideoController()       // 短视频(CMS)
	initShortVideoAppController()    // 短视频(App)
	initLiveFollowAppController()    // 关注主播(App)
	initRechargeCfgController()      // 充值配置管理(CMS)
	initRechargeCfgAppController()   // 充值配置查询(App)
	initRechargeOrderController()    // 充值订单(CMS:查询/手动充值)
	initRechargeOrderAppController() // 充值订单(App:发起/查询)
	initUploadController()           // CMS文件上传
	httpserver.InitWebsocket()
	go httpserver.InitHttpServer()
}
