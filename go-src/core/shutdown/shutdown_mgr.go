package shutdown

import "github.com/gogf/gf/v2/os/gproc"

// ListenShutdown 阻塞 main,等待进程退出(GF 负责 SIGUSR1 热重启与 SIGTERM 关 HTTP).
func ListenShutdown() {
	gproc.Listen()
}
