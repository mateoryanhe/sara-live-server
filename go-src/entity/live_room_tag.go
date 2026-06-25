package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveRoomTag db.TbName = "live_room_tags"
)

// LiveRoomTag 直播间标签(CMS 管理)
type LiveRoomTag struct {
	migrate.OneModel
	Name string `gorm:"size:64;comment:标签名称" json:"name"`
	Sort int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
}

func initLiveRoomTag() {
	migrate.AutoMigrate(&LiveRoomTag{})
}
