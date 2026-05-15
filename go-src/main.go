package main

import (
	"xr-game-server/controller"
	"xr-game-server/core"
	"xr-game-server/core/shutdown"
	"xr-game-server/dao"
	"xr-game-server/entity"
	"xr-game-server/module"
)

func main() {
	//底层框架初始化
	core.Init()
	//数据库表结构初始化
	entity.Init()
	//数据库组件初始化
	dao.Init()
	//service模块初始化
	module.Init()
	//httpserver服务器模块启动
	controller.Init()
	//开始监听程序退出
	shutdown.ListenShutdown()
}
