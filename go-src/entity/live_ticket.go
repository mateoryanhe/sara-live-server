package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveTicket db.TbName = "live_tickets"
)

const (
	LiveTicketStatusOffShelf uint8 = 0
	LiveTicketStatusOnShelf  uint8 = 1
)

// LiveTicket 直播门票配置(CMS 管理)
type LiveTicket struct {
	migrate.OneModel
	Price  float64 `gorm:"type:decimal(10,4);default:0;comment:钻石价格" json:"price"`
	Sort   int     `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status uint8   `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
}

func initLiveTicket() {
	migrate.AutoMigrate(&LiveTicket{})
}
