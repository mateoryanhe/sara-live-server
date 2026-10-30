package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyApkPromotionChannelPayUser db.TbName = "daily_apk_promotion_channel_pay_users"
)

const (
	DailyApkPromotionChannelPayUserUserId       db.TbCol = "user_id"
	DailyApkPromotionChannelPayUserStatDate     db.TbCol = "stat_date"
	DailyApkPromotionChannelPayUserChannelCode db.TbCol = "channel_code"
)

// DailyApkPromotionChannelPayUser APK 推广渠道日支付用户去重(主键=date_channelCode_userId)
type DailyApkPromotionChannelPayUser struct {
	ID          string     `gorm:"primaryKey;size:128;comment:复合ID" json:"id"`
	UserId      uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	StatDate    string     `gorm:"size:10;index;comment:统计日期" json:"statDate"`
	ChannelCode string     `gorm:"size:64;index;comment:渠道编码" json:"channelCode"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildDailyApkPromotionChannelPayUserId(date, channelCode string, userId uint64) string {
	return fmt.Sprintf("%s_%s_%d", date, channelCode, userId)
}

func NewDailyApkPromotionChannelPayUser(date, channelCode string, userId uint64) *DailyApkPromotionChannelPayUser {
	r := &DailyApkPromotionChannelPayUser{}
	r.ID = BuildDailyApkPromotionChannelPayUserId(date, channelCode, userId)
	r.SetUserId(userId)
	r.SetStatDate(date)
	r.SetChannelCode(channelCode)
	return r
}

func (r *DailyApkPromotionChannelPayUser) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyApkPromotionChannelPayUser) SetStatDate(v string) {
	r.StatDate = v
	syndb.AddData(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyApkPromotionChannelPayUser) SetChannelCode(v string) {
	r.ChannelCode = v
	syndb.AddData(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyApkPromotionChannelPayUser) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyApkPromotionChannelPayUser, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initDailyApkPromotionChannelPayUser() {
	syndb.RegLazy(TbDailyApkPromotionChannelPayUser, db.CreatedAtName)
	syndb.RegLazy(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserUserId)
	syndb.RegLazy(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserStatDate)
	syndb.RegLazy(TbDailyApkPromotionChannelPayUser, DailyApkPromotionChannelPayUserChannelCode)
	migrate.AutoMigrate(&DailyApkPromotionChannelPayUser{})
}
