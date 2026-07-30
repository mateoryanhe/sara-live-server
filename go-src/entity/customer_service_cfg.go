package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCustomerServiceCfg db.TbName = "customer_service_cfgs"
)

// CustomerServiceCfg 客服联系配置(CMS 管理,通常仅一条)
type CustomerServiceCfg struct {
	migrate.OneModel
	TelegramUrl string `gorm:"size:512;default:'';comment:Telegram联系方式" json:"telegramUrl"`
	FacebookUrl string `gorm:"size:512;default:'';comment:Facebook联系方式" json:"facebookUrl"`
	WhatsappUrl string `gorm:"size:512;default:'';comment:WhatsApp联系方式" json:"whatsappUrl"`
}

func initCustomerServiceCfg() {
	migrate.AutoMigrate(&CustomerServiceCfg{})
}
