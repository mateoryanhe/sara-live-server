package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbRandomNickname db.TbName = "random_nicknames"
)

// RandomNickname CMS 导入的随机昵称库(按语言)
type RandomNickname struct {
	migrate.OneModel
	Nickname string `gorm:"size:64;not null;comment:昵称" json:"nickname"`
	Lang     uint8  `gorm:"not null;index:idx_random_nick_lang;comment:语言(1英2西3印4葡)" json:"lang"`
}

func initRandomNickname() {
	migrate.AutoMigrate(&RandomNickname{})
}
