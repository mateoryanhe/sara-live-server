package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveCfg db.TbName = "live_cfgs"
)

// LiveCfg 直播配置(CMS管理,通常仅一条记录)
type LiveCfg struct {
	migrate.OneModel
	PaidDanmakuPrice           float64 `gorm:"type:decimal(10,4);default:0;comment:付费弹幕价格(钻石)" json:"paidDanmakuPrice"`
	VideoCallTicketEnabled     bool    `gorm:"default:1;comment:直播间视频通话接通是否扣门票" json:"videoCallTicketEnabled"`
	AudienceListRefreshSeconds uint32  `gorm:"default:300;comment:在线观众列表刷新间隔(秒)" json:"audienceListRefreshSeconds"`
	OneToOneDailyFreeSeconds   uint32  `gorm:"default:30;comment:1v1观众对单主播每日免费通话秒数(0=关闭)" json:"oneToOneDailyFreeSeconds"`
}

func initLiveCfg() {
	migrate.AutoMigrate(&LiveCfg{})
}
