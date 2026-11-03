package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyH5AppPromotionChannelPayUser db.TbName = "daily_h5_app_promotion_channel_pay_users"
)

const (
	DailyH5AppPromotionChannelPayUserUserId       db.TbCol = "user_id"
	DailyH5AppPromotionChannelPayUserStatDate     db.TbCol = "stat_date"
	DailyH5AppPromotionChannelPayUserChannelCode db.TbCol = "channel_code"
)

// DailyH5AppPromotionChannelPayUser H5下载App推广日支付用户去重
type DailyH5AppPromotionChannelPayUser struct {
	ID          string     `gorm:"primaryKey;size:48;comment:复合ID" json:"id"`
	UserId      uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	StatDate    string     `gorm:"size:10;index;comment:统计日期" json:"statDate"`
	ChannelCode int        `gorm:"index;comment:推广URL渠道编码" json:"channelCode"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildDailyH5AppPromotionChannelPayUserId(date string, channelCode int, userId uint64) string {
	return fmt.Sprintf("%s_%d_%d", date, channelCode, userId)
}

func NewDailyH5AppPromotionChannelPayUser(date string, channelCode int, userId uint64) *DailyH5AppPromotionChannelPayUser {
	r := &DailyH5AppPromotionChannelPayUser{}
	r.ID = BuildDailyH5AppPromotionChannelPayUserId(date, channelCode, userId)
	r.SetUserId(userId)
	r.SetStatDate(date)
	r.SetChannelCode(channelCode)
	return r
}

func (r *DailyH5AppPromotionChannelPayUser) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyH5AppPromotionChannelPayUser) SetStatDate(v string) {
	r.StatDate = v
	syndb.AddData(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyH5AppPromotionChannelPayUser) SetChannelCode(v int) {
	r.ChannelCode = v
	syndb.AddData(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyH5AppPromotionChannelPayUser) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyH5AppPromotionChannelPayUser, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initDailyH5AppPromotionChannelPayUser() {
	syndb.RegLazy(TbDailyH5AppPromotionChannelPayUser, db.CreatedAtName)
	syndb.RegLazy(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserUserId)
	syndb.RegLazy(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserStatDate)
	syndb.RegLazy(TbDailyH5AppPromotionChannelPayUser, DailyH5AppPromotionChannelPayUserChannelCode)
	migrate.AutoMigrate(&DailyH5AppPromotionChannelPayUser{})
}
