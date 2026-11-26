package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveRoomShareCfg db.TbName = "live_room_share_cfgs"
)

// LiveRoomShareCfg 直播间分享文案与 AppsFlyer 链接(CMS 单条配置).
type LiveRoomShareCfg struct {
	migrate.OneModel
	ShareTextZhCN string `gorm:"type:text;comment:分享文案(简体中文)" json:"shareTextZhCN"`
	ShareTextZhTW string `gorm:"type:text;comment:分享文案(繁体中文)" json:"shareTextZhTW"`
	ShareTextEn   string `gorm:"type:text;comment:分享文案(英文)" json:"shareTextEn"`
	ShareTextEs   string `gorm:"type:text;comment:分享文案(西班牙语)" json:"shareTextEs"`
	ShareTextPt   string `gorm:"type:text;comment:分享文案(葡萄牙语)" json:"shareTextPt"`
	ShareTextHi   string `gorm:"type:text;comment:分享文案(印地语)" json:"shareTextHi"`
	ShareTextId   string `gorm:"type:text;comment:分享文案(印尼语)" json:"shareTextId"`
	AfLink        string `gorm:"size:1024;default:'';comment:AppsFlyer 深度链接(不含 inviteCode/roomId)" json:"afLink"`
}

func initLiveRoomShareCfg() {
	migrate.AutoMigrate(&LiveRoomShareCfg{})
}