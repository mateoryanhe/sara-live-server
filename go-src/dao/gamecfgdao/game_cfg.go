package gamecfgdao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/gamecfgdto"
	"xr-game-server/entity"
)

const gameCfgCacheKey = "all"

var gameCfgCacheMgr *cache.CacheMgr

func InitGameCfgDao() {
	gameCfgCacheMgr = cache.NewCacheMgr()
}

func GetById(id uint64) *entity.GameCfg {
	var row entity.GameCfg
	if err := g.DB().Model(string(entity.TbGameCfg)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByCode(code string) *entity.GameCfg {
	if code == "" {
		return nil
	}
	var row entity.GameCfg
	if err := g.DB().Model(string(entity.TbGameCfg)).Where("code = ?", code).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.GameCfg) error {
	_, err := g.DB().Model(string(entity.TbGameCfg)).Save(row)
	return err
}

func Update(row *entity.GameCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbGameCfg)).WherePri(id).Delete()
	return err
}

func loadAllFromDB() []*entity.GameCfg {
	rows := make([]*entity.GameCfg, 0)
	_ = g.DB().Model(string(entity.TbGameCfg)).Order("sort desc, id desc").Scan(&rows)
	return rows
}

// ReloadCache 配置变更后清除缓存并从数据库重新加载
func ReloadCache() {
	if gameCfgCacheMgr == nil {
		return
	}
	gameCfgCacheMgr.FlushCache(gameCfgCacheKey, loadAllFromDB())
	//设置永不过期
	gameCfgCacheMgr.Cache.UpdateExpire(gctx.New(), gameCfgCacheKey, time.Hour*24*365*100)
}

// GetAllCached 获取全部游戏配置(优先读缓存,未命中再查库)
func GetAllCached() []*entity.GameCfg {
	if gameCfgCacheMgr == nil {
		return loadAllFromDB()
	}
	v := gameCfgCacheMgr.GetData(gameCfgCacheKey, func(ctx context.Context) (value interface{}, err error) {
		return loadAllFromDB(), nil
	})
	//设置永不过期
	gameCfgCacheMgr.Cache.UpdateExpire(gctx.New(), gameCfgCacheKey, time.Hour*24*365*100)
	if v == nil {
		return make([]*entity.GameCfg, 0)
	}
	list, _ := v.([]*entity.GameCfg)
	if list == nil {
		return make([]*entity.GameCfg, 0)
	}
	return list
}

func GetList(req *gamecfgdto.GameCfgListReq) (int, []*gamecfgdto.GameCfgListRes) {
	sql := `select id, name, code, live_cover, link, sort, status, created_at, updated_at
            from game_cfgs
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*gamecfgdto.GameCfgListRes, 0)

	if req.Name != "" {
		sql += ` and name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.Code != "" {
		sql += ` and code LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Code))
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.GameCfgStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.GameCfgStatusOnShelf)
	}

	sql += ` order by sort desc, id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa((req.PageIndex-1)*req.PageSize)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
