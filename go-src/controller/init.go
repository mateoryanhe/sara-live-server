package controller

import (
	"xr-game-server/core/httpserver"
)

func Init() {
	initAuthApi()
	initAccountController()
	initGoldController()
	initGoldAppController()
	initDiamondController()
	initIndex()
	initWebSocket()
	initRank()
	initGlobalCfgController()
	initAppTokenController()
	initRoleController()                    // 添加角色控制器初始化
	initCMSUserController()                 // 添加CMS用户控制器初始化
	initUserInfoController()                // 用户基础信息(App)
	initGuildController()                   // 直播工会管理(CMS)
	initGuildAppController()                // 直播工会查询(App)
	initLiveRoomAppController()             // 直播间(App)
	initAgoraAppController()                // 声网(App)
	initAgoraCMSController()                // 声网(CMS)
	initAliyunTextModerationCMSController() // 阿里云文本审核(CMS)
	initGiftController()                    // 礼物配置(CMS)
	initGiftAppController()                 // 礼物列表(App)
	initBannerController()                  // 首页Banner(CMS)
	initBannerAppController()               // 首页Banner(App)
	initShortVideoController()              // 短视频(CMS)
	initShortVideoAppController()           // 短视频(App)
	initLiveFollowAppController()           // 关注主播(App)
	initRechargeCfgController()             // 充值配置管理(CMS)
	initRechargeCfgAppController()          // 充值配置查询(App)
	initRechargeOrderController()           // 充值订单(CMS:查询/手动充值)
	initCurrencyLogController()             // 货币流水(CMS)
	initLiveRevenueLogController()          // 直播收益流水(CMS)
	initLiveRecordCMSController()           // 直播记录(CMS)
	initRichRankAppController()             // 富豪榜(App)
	initAnchorRankAppController()           // 主播红人榜(App)
	initRechargeOrderAppController()        // 充值订单(App:发起/查询)
	initVipCfgController()                  // VIP配置(CMS)
	initGameCfgController()                 // 游戏配置(CMS)
	initGameCfgAppController()              // 游戏配置(App)
	initVipCfgAppController()               // VIP配置查询(App)
	initVipAppController()                  // VIP详情(App)
	initMessageAppController()              // 私信(App)
	initUploadController()                  // CMS文件上传
	initSysStatController()                 // 系统总数据/仪表盘(CMS)
	httpserver.InitWebsocket()
	go httpserver.InitHttpServer()
}
