package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrpool"

	"github.com/gogf/gf/v2/os/gctx"
)

func Init() {
	initPushController()
	initAuthApi()
	initAuthAppController()
	initAccountController()
	initBotAnchorController()
	initGoldController()
	initGoldAppController()
	initDiamondController()
	initIndex()
	initSysController()
	initAppTokenController()
	initRoleController()                    // 添加角色控制器初始化
	initCMSUserController()                 // 添加CMS用户控制器初始化
	initUserInfoController()                // 用户基础信息(App)
	initUserInfoPublicController()          // 用户公开接口(官网销户等)
	initGuildController()                   // 直播工会管理(CMS)
	initGuildAppController()                // 直播工会查询(App)
	initLiveRoomAppController()             // 直播间(App)
	initLiveRoomTagController()             // 直播间标签(CMS)
	initLiveRoomTagAppController()          // 直播间标签(App)
	initAgoraAppController()                // 声网(App)
	initCallAppController()                 // 通话(App)
	initCallCMSController()                 // 通话(CMS)
	initAgoraCMSController()                // 声网(CMS)
	initLiveCfgCMSController()              // 直播配置(CMS)
	initAliyunTextModerationCMSController() // 阿里云文本审核(CMS)
	initPrivacyPolicyCMSController()        // 隐私政策配置(CMS)
	initCustomerServiceCMSController()      // 客服联系配置(CMS)
	initCustomerServiceAppController()      // 客服联系配置(App)
	initAccountCfgCMSController()           // 账号配置(CMS)
	initGiftController()                    // 礼物配置(CMS)
	initGiftAppController()                 // 礼物列表(App)
	initBannerController()                  // 首页Banner(CMS)
	initActivityMessageController()         // 活动消息(CMS)
	initBannerAppController()               // 首页Banner(App)
	initTicketController()                  // 门票(CMS)
	initTicketAppController()               // 门票(App)
	initPrivateRoomBillingController()      // 私密直播间计费(CMS)
	initPrivateRoomBillingAppController()   // 私密直播间计费(App)
	initShortVideoController()              // 短视频(CMS)
	initShortVideoAppController()           // 短视频(App)
	initLiveFollowAppController()           // 关注主播(App)
	initRechargeCfgController()             // 充值配置管理(CMS)
	initRechargeCfgAppController()          // 充值配置查询(App)
	initRechargeOrderController()           // 充值订单(CMS:查询/手动充值)
	initCurrencyLogController()             // 货币流水(CMS)
	initLiveRevenueLogController()          // 直播收益流水(CMS)
	initLiveRecordCMSController()           // 直播记录(CMS)
	initLiveRecordAppController()           // 直播记录(App)
	initRichRankAppController()             // 富豪榜(App)
	initAnchorRankAppController()           // 主播红人榜(App)
	initRechargeOrderAppController()        // 充值订单(App:发起/查询)
	initVipCfgController()                  // VIP配置(CMS)
	initAppPkgController()                  // App包管理(CMS)
	initGameCfgController()                 // 游戏配置(CMS)
	initGameCfgAppController()              // 游戏配置(App)
	initVipCfgAppController()               // VIP配置查询(App)
	initVipAppController()                  // VIP详情(App)
	initMessageAppController()              // 私信(App)
	initUploadController()                  // CMS文件上传
	initRandomNicknameController()          // 随机昵称库(CMS)
	initSysStatController()                 // 系统总数据/仪表盘(CMS)
	initResourceMetricController()          // 系统资源监控(CMS)
	initLogQueryController()                // 日志查询(CMS)
	httpserver.InitWebsocket()
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		httpserver.InitHttpServer()
	})
}
