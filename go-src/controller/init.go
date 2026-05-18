package controller

import (
	"xr-game-server/core/httpserver"
)

func Init() {
	initAuthApi()
	initAccountController()
	initIndex()
	initWebSocket()
	initRank()
	initGlobalCfgController()
	initRoleController()        // 添加角色控制器初始化
	initCMSUserController()     // 添加CMS用户控制器初始化
	initUserInfoController()    // 用户基础信息(App)
	initGuildController()       // 直播工会管理(CMS)
	initGuildAppController()    // 直播工会查询(App)
	initLiveRoomAppController() // 直播间(App)
	httpserver.InitWebsocket()
	go httpserver.InitHttpServer()
}
