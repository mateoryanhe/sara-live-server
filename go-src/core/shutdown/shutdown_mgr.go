package shutdown

import "github.com/gogf/gf/v2/os/gproc"

// ListenShutdown 阻塞 main,等待进程退出(GF 信号处理).
func ListenShutdown() {
	gproc.Listen()
}
