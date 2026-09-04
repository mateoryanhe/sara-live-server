package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

func userExtLockKey(userId uint64) string {
	return fmt.Sprintf("user_ext:%d", userId)
}

const (
	TbUserExt db.TbName = "user_exts"
)

const (
	UserExtCanRank            db.TbCol = "can_rank"
	UserExtPrettyId           db.TbCol = "pretty_id"
	UserExtPackageName        db.TbCol = "package_name"
	UserExtAppVersion         db.TbCol = "app_version"
	UserExtFollowCount        db.TbCol = "follow_count"
	UserExtFollowerCount      db.TbCol = "follower_count"
	UserExtCancelCode         db.TbCol = "cancel_code"
	UserExtCancelCodeExpireAt db.TbCol = "cancel_code_expire_at"
	UserExtRechargeWhitelist         db.TbCol = "recharge_whitelist"
	UserExtFirstRecharge             db.TbCol = "first_recharge"
	UserExtShortVideoUnsettledIncome db.TbCol = "short_video_unsettled_income"
)

// UserExt 用户扩展信息(与用户一一对应,主键ID即用户ID)
type UserExt struct {
	migrate.OneModel
	CanRank            bool       `gorm:"default:1;comment:是否可上排行榜" json:"canRank"`
	PrettyId           uint64     `gorm:"default:0;comment:靓号(默认等于用户ID)" json:"prettyId"`
	PackageName        string     `gorm:"default:'';comment:注册包名" json:"packageName"`
	AppVersion         string     `gorm:"default:'';comment:注册版本号" json:"appVersion"`
	FollowCount        uint64     `gorm:"default:0;comment:当前关注数" json:"followCount"`
	FollowerCount      uint64     `gorm:"default:0;comment:当前粉丝数" json:"followerCount"`
	CancelCode         string     `gorm:"size:128;default:'';index;comment:注销码" json:"cancelCode"`
	CancelCodeExpireAt *time.Time `gorm:"comment:注销码过期时间" json:"cancelCodeExpireAt"`
	RechargeWhitelist         bool       `gorm:"default:0;comment:充值白名单(创建订单后直接到账)" json:"rechargeWhitelist"`
	FirstRecharge             bool       `gorm:"default:1;comment:是否首次充值(1=未首充,0=已首充)" json:"firstRecharge"`
	ShortVideoUnsettledIncome float64    `gorm:"type:decimal(16,4);default:0;comment:短视频未结算收益(非主播作者)" json:"shortVideoUnsettledIncome"`
}

func NewUserExt(userId uint64) *UserExt {
	ret := &UserExt{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetCanRank(true)
	ret.SetPrettyId(userId)
	ret.SetFirstRecharge(true)
	return ret
}

func (receiver *UserExt) SetCanRank(canRank bool) {
	receiver.CanRank = canRank
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtCanRank, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: canRank,
	})
}

func (receiver *UserExt) SetPrettyId(prettyId uint64) {
	receiver.PrettyId = prettyId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtPrettyId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: prettyId,
	})
}

func (receiver *UserExt) SetPackageName(packageName string) {
	receiver.PackageName = packageName
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtPackageName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: packageName,
	})
}

func (receiver *UserExt) SetAppVersion(appVersion string) {
	receiver.AppVersion = appVersion
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtAppVersion, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: appVersion,
	})
}

func (receiver *UserExt) SetCancelCode(cancelCode string) {
	receiver.CancelCode = cancelCode
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtCancelCode, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: cancelCode,
	})
}

func (receiver *UserExt) SetCancelCodeExpireAt(val *time.Time) {
	receiver.CancelCodeExpireAt = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtCancelCodeExpireAt, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserExt) SetRechargeWhitelist(v bool) {
	receiver.RechargeWhitelist = v
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtRechargeWhitelist, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: v,
	})
}

func (receiver *UserExt) SetFirstRecharge(v bool) {
	receiver.FirstRecharge = v
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtFirstRecharge, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: v,
	})
}

func (receiver *UserExt) AddShortVideoUnsettledIncome(amount float64) {
	if amount <= 0 {
		return
	}
	receiver.ShortVideoUnsettledIncome = math.AddFloat64(receiver.ShortVideoUnsettledIncome, amount)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtShortVideoUnsettledIncome, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.ShortVideoUnsettledIncome,
	})
}

// ClearShortVideoUnsettledIncome 清零短视频未结算收益并返回清零前的金额
func (receiver *UserExt) ClearShortVideoUnsettledIncome() float64 {
	prev := receiver.ShortVideoUnsettledIncome
	if prev <= 0 {
		return 0
	}
	receiver.ShortVideoUnsettledIncome = 0
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtShortVideoUnsettledIncome, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: float64(0),
	})
	return prev
}

func (receiver *UserExt) AddFollowCount(val uint64) {
	if val == 0 {
		return
	}
	receiver.FollowCount = math.Add(receiver.FollowCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserExt, UserExtFollowCount, &syndb.ColData{
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
	syndb.AddData(TbUserExt, UserExtFollowCount, &syndb.ColData{
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
	syndb.AddData(TbUserExt, UserExtFollowerCount, &syndb.ColData{
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
	syndb.AddData(TbUserExt, UserExtFollowerCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.FollowerCount,
	})
}

func (receiver *UserExt) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddData(TbUserExt, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserExt) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddData(TbUserExt, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserExt() {
	syndb.RegQuick(TbUserExt, db.CreatedAtName)
	syndb.RegQuick(TbUserExt, db.UpdatedAtName)
	syndb.RegQuick(TbUserExt, UserExtCanRank)
	syndb.RegQuick(TbUserExt, UserExtPrettyId)
	syndb.RegQuick(TbUserExt, UserExtPackageName)
	syndb.RegQuick(TbUserExt, UserExtAppVersion)
	syndb.RegQuick(TbUserExt, UserExtFollowCount)
	syndb.RegQuick(TbUserExt, UserExtFollowerCount)
	syndb.RegQuick(TbUserExt, UserExtCancelCode)
	syndb.RegQuick(TbUserExt, UserExtCancelCodeExpireAt)
	syndb.RegQuick(TbUserExt, UserExtRechargeWhitelist)
	syndb.RegQuick(TbUserExt, UserExtFirstRecharge)
	syndb.RegQuick(TbUserExt, UserExtShortVideoUnsettledIncome)

	migrate.AutoMigrate(&UserExt{})
}
