package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbMetaPixelCfg db.TbName = "meta_pixel_cfgs"
)

// MetaPixelCfg Meta Pixel Conversions API(CMS 管理,通常仅一条).
type MetaPixelCfg struct {
	migrate.OneModel
	Enabled        uint8  `gorm:"default:0;comment:是否启用(0否1是)" json:"enabled"`
	PixelId        string `gorm:"size:32;default:'';comment:Meta Pixel ID" json:"pixelId"`
	AccessToken    string `gorm:"type:text;comment:Conversions API Access Token" json:"accessToken"`
	TestEventCode  string `gorm:"size:64;default:'';comment:测试事件代码(可选)" json:"testEventCode"`
}

func (MetaPixelCfg) TableName() string {
	return string(TbMetaPixelCfg)
}

func (c *MetaPixelCfg) IsActive() bool {
	return c != nil && c.ID > 0 && c.Enabled == 1 &&
		c.PixelId != "" && c.AccessToken != ""
}

func initMetaPixelCfg() {
	migrate.AutoMigrate(&MetaPixelCfg{})
}
