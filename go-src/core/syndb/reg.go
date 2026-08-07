package syndb

import (
	"xr-game-server/constants/db"
)

var synCacheMap = make(map[string]*ColSynCache)

func cacheKey(tbName db.TbName, tbCol db.TbCol) string {
	return string(tbName) + ":" + string(tbCol)
}

// Reg 注册列缓冲(统一调度,quick/lazy 不再区分).
func Reg(tbName db.TbName, tbCol db.TbCol) {
	colKey := cacheKey(tbName, tbCol)
	if _, ok := synCacheMap[colKey]; ok {
		return
	}
	synCacheMap[colKey] = newColSynCache(string(tbName), string(tbCol))
}

// RegQuick 兼容旧接口,等价于 Reg.
func RegQuick(tbName db.TbName, tbCol db.TbCol) {
	Reg(tbName, tbCol)
}

// RegLazy 兼容旧接口,等价于 Reg.
func RegLazy(tbName db.TbName, tbCol db.TbCol) {
	Reg(tbName, tbCol)
}
