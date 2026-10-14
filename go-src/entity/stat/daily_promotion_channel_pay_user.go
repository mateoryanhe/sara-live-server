package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyPromotionChannelPayUser db.TbName = "daily_promotion_channel_pay_users"
)

const (
	DailyPromotionChannelPayUserUserId      db.TbCol = "user_id"
	DailyPromotionChannelPayUserStatDate      db.TbCol = "stat_date"
	DailyPromotionChannelPayUserChannelCode db.TbCol = "channel_code"
)

// DailyPromotionChannelPayUser 推广渠道日支付用户去重(主键=date_channelCode_userId)
type DailyPromotionChannelPayUser struct {
	ID          string     `gorm:"primaryKey;size:48;comment:复合ID" json:"id"`
	UserId      uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	StatDate    string     `gorm:"size:10;index;comment:统计日期" json:"statDate"`
	ChannelCode int        `gorm:"index;comment:渠道编码" json:"channelCode"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildDailyPromotionChannelPayUserId(date string, channelCode int, userId uint64) string {
	return fmt.Sprintf("%s_%d_%d", date, channelCode, userId)
}

func NewDailyPromotionChannelPayUser(date string, channelCode int, userId uint64) *DailyPromotionChannelPayUser {
	r := &DailyPromotionChannelPayUser{}
	r.ID = BuildDailyPromotionChannelPayUserId(date, channelCode, userId)
	r.SetUserId(userId)
	r.SetStatDate(date)
	r.SetChannelCode(channelCode)
	return r
}

func (r *DailyPromotionChannelPayUser) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyPromotionChannelPayUser) SetStatDate(v string) {
	r.StatDate = v
	syndb.AddData(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyPromotionChannelPayUser) SetChannelCode(v int) {
	r.ChannelCode = v
	syndb.AddData(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyPromotionChannelPayUser) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyPromotionChannelPayUser, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initDailyPromotionChannelPayUser() {
	syndb.RegLazy(TbDailyPromotionChannelPayUser, db.CreatedAtName)
	syndb.RegLazy(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserUserId)
	syndb.RegLazy(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserStatDate)
	syndb.RegLazy(TbDailyPromotionChannelPayUser, DailyPromotionChannelPayUserChannelCode)
	migrate.AutoMigrate(&DailyPromotionChannelPayUser{})
}
