package httpserver

// enableHotRestartGraceful 开启 GoFrame 热重启(fork 继承 fd),旧进程退出由业务自行处理.
func enableHotRestartGraceful() {
	httpServer.SetGraceful(true)
}
