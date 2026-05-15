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
	initRoleController()    // 添加角色控制器初始化
	initCMSUserController() // 添加CMS用户控制器初始化
	httpserver.InitWebsocket()
	go httpserver.InitHttpServer()
}
