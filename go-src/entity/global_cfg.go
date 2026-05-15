package entity

import "xr-game-server/core/migrate"

const (
	TbGlobalCfg = "global_cfgs"
)

type GlobalCfgModule string

type GlobalCfg struct {
	migrate.OneModel
	Module string `json:"module" gorm:"column:module"`
	Key    string `json:"key" gorm:"column:key"`
	Value  string `json:"value" gorm:"column:value"`
	Desc   string `json:"desc" gorm:"column:desc"`
}

func initGlobalCfg() {
	migrate.AutoMigrate(&GlobalCfg{})
}
