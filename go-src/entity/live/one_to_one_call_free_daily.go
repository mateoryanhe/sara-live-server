package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbOneToOneCallFreeDailyUse db.TbName = "one_to_one_call_free_daily_uses"
)

// OneToOneCallFreeDailyUse 观众对某主播在某一自然日已消耗 1v1 免费通话额度(整段作废)。
type OneToOneCallFreeDailyUse struct {
	migrate.OneModel
	PayerId  uint64 `gorm:"uniqueIndex:uk_o2o_free_daily,priority:1;comment:付费用户ID" json:"payerId"`
	AnchorId uint64 `gorm:"uniqueIndex:uk_o2o_free_daily,priority:2;comment:主播用户ID" json:"anchorId"`
	StatDate string `gorm:"size:10;uniqueIndex:uk_o2o_free_daily,priority:3;comment:服务器自然日YYYY-MM-DD" json:"statDate"`
}

func initOneToOneCallFreeDailyUse() {
	migrate.AutoMigrate(&OneToOneCallFreeDailyUse{})
}
