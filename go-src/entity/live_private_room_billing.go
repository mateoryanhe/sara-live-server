package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLivePrivateRoomBilling db.TbName = "live_private_room_billings"
)

const (
	LivePrivateRoomBillingStatusOffShelf uint8 = 0
	LivePrivateRoomBillingStatusOnShelf  uint8 = 1
)

// LivePrivateRoomBilling 私密直播间按分钟计费配置(CMS 管理)
type LivePrivateRoomBilling struct {
	migrate.OneModel
	PricePerMinute float64 `gorm:"type:decimal(10,4);default:0;comment:每分钟钻石价格" json:"pricePerMinute"`
	Sort           int     `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status         uint8   `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
}

func initLivePrivateRoomBilling() {
	migrate.AutoMigrate(&LivePrivateRoomBilling{})
}
