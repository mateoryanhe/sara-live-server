package tracking

// Init 埋点事件订阅与异步消费队列(与 module/stat 相同模式)
func Init() {
	initTrackingQueue()
}
