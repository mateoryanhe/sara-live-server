package usermaxid

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

var userMaxId *entity.UserMaxId

const userMaxIdRecordID uint64 = 1

// 参考 MAX(user_infos.id) 时排除压测/脏数据用的超大 id(须大于业务号段且低于 JS Number 安全上限)
const userInfoMaxIdSanityMax uint64 = 10_000_000_000

// Init 加载/初始化用户最大ID(启动时调用)
func Init() {
	row := &entity.UserMaxId{}
	err := g.Model(string(entity.TbUserMaxId)).Order("id asc").Limit(1).Scan(row)
	if err != nil || row.ID == 0 {
		g.Log().Warningf(gctx.New(), "user_max_id load empty, create start=%d err=%v", entity.UserMaxIdDefaultID, err)
		userMaxId = entity.NewUserMaxId(userMaxIdRecordID)
		// 初始最大ID=1000000, 第一次 NextUserId 得到 1000001
		userMaxId.SetMaxId(entity.UserMaxIdDefaultID)
		return
	}
	userMaxId = row
	// 兼容错误初始化(曾写成0导致发号从1开始): 拉回到默认起点
	if userMaxId.Add(0) < entity.UserMaxIdDefaultID {
		g.Log().Warningf(gctx.New(), "user_max_id too small=%d, reset to %d", userMaxId.Add(0), entity.UserMaxIdDefaultID)
		userMaxId.SetMaxId(entity.UserMaxIdDefaultID)
	}
}

// NextUserId 原子分配下一个用户ID
func NextUserId() uint64 {
	if userMaxId == nil {
		Init()
	}
	return userMaxId.Add(1)
}

// CurrentUserMaxId 当前已分配的最大用户ID(不自增)
func CurrentUserMaxId() uint64 {
	if userMaxId == nil {
		return 0
	}
	return userMaxId.Add(0)
}

// MaxUserInfoIdInDB 查询 user_infos 合理号段内的最大 id(过滤 UINT64_MAX 等测试脏行)
func MaxUserInfoIdInDB() uint64 {
	v, err := g.Model(string(entity.TbUserInfo)).
		Fields("IFNULL(MAX(id), 0)").
		WhereBetween(string(db.IdName), entity.UserMaxIdDefaultID, userInfoMaxIdSanityMax).
		Value()
	if err != nil || v == nil || v.IsNil() {
		return 0
	}
	return v.Uint64()
}

// SetMaxIdFromCMS 将发号器最大 ID 调整到 target(下一注册为 target+1); 不得低于当前发号器与 user_infos 最大 id
func SetMaxIdFromCMS(target uint64) error {
	if userMaxId == nil {
		Init()
	}
	if target < entity.UserMaxIdDefaultID {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	current := userMaxId.Add(0)
	dbMax := MaxUserInfoIdInDB()
	minAllowed := current
	if dbMax > minAllowed {
		minAllowed = dbMax
	}
	if target < minAllowed {
		g.Log().Warningf(gctx.New(), "user_max_id cms rejected target=%d min=%d counter=%d user_infos.max=%d", target, minAllowed, current, dbMax)
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if target == current {
		return nil
	}
	userMaxId.SetMaxId(target)
	g.Log().Infof(gctx.New(), "user_max_id cms set max_id=%d prev=%d user_infos.max=%d", target, current, dbMax)
	return nil
}