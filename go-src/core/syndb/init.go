package syndb

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/shutdown"
	"xr-game-server/core/xrtimer"
)

// InitSynCache 初始化同步数据库缓存
func InitSynCache() {
	//先注册好通道,固定10毫秒处理一下缓冲数据
	xrtimer.ModuleLoop(10*time.Millisecond, consume)
	shutdown.RegCommonShutDownHandler(SysExit)
}

// RegWithSelf 指定同步频率,缓冲大小
func RegWithSelf(tbName db.TbName, tbCol db.TbCol, synTime time.Duration, bufferSize int) {
	colKey := string(tbName) + ":" + string(tbCol)
	lazyMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, bufferSize),
		TempData:  make([]*ColData, 0),
		Period:    synTime,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}
