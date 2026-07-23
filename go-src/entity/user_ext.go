package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserExt db.TbName = "user_exts"
)

const (
	UserExtCanRank       db.TbCol = "can_rank"
	UserExtPackageName   db.TbCol = "package_name"
	UserExtAppVersion    db.TbCol = "app_version"
	UserExtFollowCount   db.TbCol = "follow_count"
	UserExtFollowerCount db.TbCol = "follower_count"
	UserExtCancelCode    db.TbCol = "cancel_code"
)

// UserExt 用户扩展信息(与用户一一对应,主键ID即用户ID)
type UserExt struct {
	migrate.OneModel
	CanRank       bool   `gorm:"default:1;comment:是否可上排行榜" json:"canRank"`
	PackageName   string `gorm:"default:'';comment:注册包名" json:"packageName"`
	AppVersion    string `gorm:"default:'';comment:注册版本号" json:"appVersion"`
	FollowCount   uint64 `gorm:"default:0;comment:当前关注数" json:"followCount"`
	FollowerCount uint64 `gorm:"default:0;comment:当前粉丝数" json:"followerCount"`
	CancelCode    string `gorm:"size:128;default:'';comment:注销码" json:"cancelCode"`
}

func NewUserExt(userId uint64) *UserExt {
	ret := &UserExt{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetCanRank(true)
	return ret
}

func (receiver *UserExt) SetCanRank(canRank bool) {
	receiver.CanRank = canRank
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtCanRank, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: canRank,
	})
}

func (receiver *UserExt) SetPackageName(packageName string) {
	receiver.PackageName = packageName
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtPackageName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: packageName,
	})
}

func (receiver *UserExt) SetAppVersion(appVersion string) {
	receiver.AppVersion = appVersion
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtAppVersion, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: appVersion,
	})
}

func (receiver *UserExt) SetCancelCode(cancelCode string) {
	receiver.CancelCode = cancelCode
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtCancelCode, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: cancelCode,
	})
}

func (receiver *UserExt) AddFollowCount(val uint64) {
	if val == 0 {
		return
	}
	receiver.FollowCount = math.Add(receiver.FollowCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtFollowCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.FollowCount,
	})
}

func (receiver *UserExt) SubFollowCount(val uint64) {
	if val == 0 {
		return
	}
	receiver.FollowCount = math.Sub(receiver.FollowCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtFollowCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.FollowCount,
	})
}

func (receiver *UserExt) AddFollowerCount(val uint64) {
	if val == 0 {
		return
	}
	receiver.FollowerCount = math.Add(receiver.FollowerCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtFollowerCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.FollowerCount,
	})
}

func (receiver *UserExt) SubFollowerCount(val uint64) {
	if val == 0 {
		return
	}
	receiver.FollowerCount = math.Sub(receiver.FollowerCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtFollowerCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.FollowerCount,
	})
}

func (receiver *UserExt) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbUserExt, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserExt) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbUserExt, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserExt() {
	syndb.RegQuick(TbUserExt, db.CreatedAtName)
	syndb.RegQuick(TbUserExt, db.UpdatedAtName)
	syndb.RegQuick(TbUserExt, UserExtCanRank)
	syndb.RegQuick(TbUserExt, UserExtPackageName)
	syndb.RegQuick(TbUserExt, UserExtAppVersion)
	syndb.RegQuick(TbUserExt, UserExtFollowCount)
	syndb.RegQuick(TbUserExt, UserExtFollowerCount)
	syndb.RegQuick(TbUserExt, UserExtCancelCode)

	migrate.AutoMigrate(&UserExt{})
}
