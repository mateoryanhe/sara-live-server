package accountdao

import (
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/user"
)

var accountCacheMgr *cache.ListCache[*entity.Account]

func GetAccountById(accountId uint64) *entity.Account {
	var account *entity.Account
	g.Model(string(entity.TbAccount)).Unscoped().Where(db.IdName, accountId).Scan(&account)
	if account != nil {
		return account
	}
	return nil
}

func InitAccountDao() {
	accountCacheMgr = cache.NewListCache[*entity.Account]()
}

func GetUserInfo(req *accountdto.QueryUserInfoReq) (int, []*accountdto.UserInfoDto) {
	sql := `select  a.*,
                    u.nickname, u.phone, u.avatar, u.remark,
                    u.gold, u.diamond, u.share_code, IFNULL(r.guild_id, 0) as guild_id, u.user_type, u.vip_level, u.last_login_time,
                    d.device_type, e.package_name, e.app_version,
                    IFNULL(e.can_rank, 1) as can_rank,
                    IFNULL(e.recharge_whitelist, 0) as recharge_whitelist
                    from accounts a
                    left join user_infos u on u.id = a.id
                    left join live_rooms r on r.id = a.id
                    left join user_login_devices d on d.id = a.id
                    left join user_exts e on e.id = a.id
                    where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*accountdto.UserInfoDto, 0)
	if req.Key != "" {
		sql += ` and (CAST(a.id AS CHAR) LIKE ? or a.open_id = ?)`
		param = append(param, "%"+req.Key+"%", req.Key)
	}
	if req.StartTime != "" {
		sql += ` and (a.created_at between ? and ?)`
		startTime, _ := time.Parse("2006-01-02", req.StartTime)
		endTime, _ := time.Parse("2006-01-02", req.EndTime)
		param = append(param, startTime, endTime)
	}
	if req.RechargeWhitelist != nil {
		if *req.RechargeWhitelist == 1 {
			sql += ` and IFNULL(e.recharge_whitelist, 0) = 1`
		} else {
			sql += ` and IFNULL(e.recharge_whitelist, 0) = 0`
		}
	}
	if req.IsAnchor != nil {
		if *req.IsAnchor == 1 {
			sql += ` and u.user_type in (?, ?, ?)`
			param = append(param, entity.UserTypeAnchor, entity.UserTypeBotAnchor, entity.UserTypeSeniorAnchor)
		} else {
			sql += ` and IFNULL(u.user_type, 0) not in (?, ?, ?)`
			param = append(param, entity.UserTypeAnchor, entity.UserTypeBotAnchor, entity.UserTypeSeniorAnchor)
		}
	}
	sql += ` order by a.id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
