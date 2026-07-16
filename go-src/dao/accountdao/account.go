package accountdao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
)

var accountCacheMgr *cache.CacheMgr

const accountMissCacheTime = 10 * time.Minute

// GetAccountBy 根据玩家id拉取数据，不存在时创建内存账号并异步落库（仅用于注册等写场景）
func GetAccountBy(openId string, channel uint) *entity.Account {
	key := accountCacheKey(openId, channel)
	//命中不了缓存，从数据库拉取数据
	cacheData := accountCacheMgr.GetData(key, func(ctx context.Context) (value interface{}, err error) {
		//从数据库拉取数据
		var account *entity.Account
		err = g.Model(string(entity.TbAccount)).Unscoped().Where(g.Map{
			string(entity.AccountOpenId):  openId,
			string(entity.AccountChannel): channel,
		}).Scan(&account)
		if account != nil {
			return account, nil
		} else {
			return entity.NewAccount(openId, channel), nil
		}

	})
	return cacheData.(*entity.Account)
}

// FindAccountBy 只读查询账号，不存在时返回 nil，不会创建新记录
func FindAccountBy(openId string, channel uint) *entity.Account {
	key := accountCacheKey(openId, channel)
	ctx := gctx.New()

	if ok, _ := accountCacheMgr.Cache.Contains(ctx, accountMissCacheKey(key)); ok {
		return nil
	}

	if cached := accountCacheMgr.GetFromCache(key); cached != nil {
		return cached.(*entity.Account)
	}

	var account *entity.Account
	_ = g.Model(string(entity.TbAccount)).Unscoped().Where(g.Map{
		string(entity.AccountOpenId):  openId,
		string(entity.AccountChannel): channel,
	}).Scan(&account)
	if account == nil {
		_ = accountCacheMgr.Cache.Set(ctx, accountMissCacheKey(key), 1, accountMissCacheTime)
		return nil
	}

	accountCacheMgr.FlushCache(key, account)
	return account
}

func accountCacheKey(openId string, channel uint) string {
	return fmt.Sprintf("%v:%v", openId, channel)
}

func accountMissCacheKey(key string) string {
	return "miss:" + key
}

func GetAccountById(accountId uint64) *entity.Account {
	//从数据库拉取数据
	var account *entity.Account
	g.Model(string(entity.TbAccount)).Unscoped().Where(db.IdName, accountId).Scan(&account)
	if account != nil {
		return account
	} else {
		return nil
	}
}

func InitAccountDao() {
	accountCacheMgr = cache.NewCacheMgr()
}

func GetUserInfo(req *accountdto.QueryUserInfoReq) (int, []*accountdto.UserInfoDto) {
	sql := `select  a.*,
                    u.nickname, u.phone, u.avatar, u.remark,
                    u.gold, u.diamond, u.share_code, u.guild_id, u.user_type, u.vip_level,
                    d.device_type
                    from accounts a
                    left join user_infos u on u.id = a.id
                    left join user_login_devices d on d.id = a.id
                    where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*accountdto.UserInfoDto, 0)
	if req.Key != "" {
		sql += ` and (a.id =? or a.open_id=? )`
		param = append(param, req.Key, req.Key)
	}
	if req.StartTime != "" {
		sql += ` and (a.created_at between ? and ?)`
		startTime, _ := time.Parse("2006-01-02", req.StartTime)
		endTime, _ := time.Parse("2006-01-02", req.EndTime)
		param = append(param, startTime, endTime)
	}
	sql += ` order by a.id desc`
	//获取总数
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
