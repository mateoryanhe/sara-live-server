package name

import (
	"context"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/cache"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/namedao"
	"xr-game-server/entity"
)

// 初始化布隆过滤器（预计元素数量500万，误判率0.01）
var filter = bloom.New(5000000, 5) // 参数：位数组长度、哈希函数数量
var nameCacheMgr = cache.NewCacheMgr()

func Init() {
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		initName(entity.RoleNameType)
	})
}

func initName(nameType entity.DbNameType) {
	names := namedao.GetValBy(nameType)
	for _, val := range names {
		filter.Add([]byte(fmt.Sprintf("%v:%s", nameType, val.String())))
	}
}

// addNameToFilter 添加名称到过滤器
func addNameToFilter(name string, nameType entity.DbNameType) {
	lockStr := "name_add"
	gmlock.Lock(lockStr)
	defer gmlock.Unlock(lockStr)
	filter.Add([]byte(fmt.Sprintf("%v:%s", nameType, name)))
}

func ChkName(target string, nameType entity.DbNameType) bool {
	lockName := fmt.Sprintf("name_chk_%v", nameType)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)
	flag := filter.Test([]byte(fmt.Sprintf("%v:%s", nameType, target)))
	if flag == false {
		AddName(target, nameType)
		return true
	}
	//查询缓存
	cacheRet, _ := nameCacheMgr.Cache.Get(gctx.New(), fmt.Sprintf("%v:%s", nameType, target))
	if cacheRet == nil {
		//查询数据库
		//查询数据库
		dbRet := namedao.GetNameBy(nameType, target)
		//检查是否存在数据库
		if dbRet == nil {
			AddName(target, nameType)
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func AddName(target string, nameType entity.DbNameType) {
	addNameToFilter(target, nameType)
	dbRet := entity.NewName(target, nameType)
	nameCacheMgr.FlushCache(fmt.Sprintf("%v:%s", nameType, target), dbRet)
}
