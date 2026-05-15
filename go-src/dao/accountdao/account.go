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

// GetAccountBy 根据玩家id拉取数据
func GetAccountBy(openId string, channel uint) *entity.Account {
	key := fmt.Sprintf("%v:%v", openId, channel)
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

func CancelAccount(accountId uint64) {
	account := GetAccountById(accountId)
	account.Cancel = true
	g.Model(string(entity.TbAccount)).Save(account)
	key := fmt.Sprintf("%v:%v", account.OpenId, account.Channel)
	accountCacheMgr.Cache.Remove(gctx.New(), key)
}

func UnCancelAccount(accountId uint64) {
	account := GetAccountById(accountId)
	account.Cancel = false
	g.Model(string(entity.TbAccount)).Save(account)
}

func InitAccountDao() {
	accountCacheMgr = cache.NewCacheMgr()
}

func GetUserInfo(req *accountdto.QueryUserInfoReq) (int, []*accountdto.UserInfoDto) {
	sql := `select  a.*
                    from accounts a
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
