package entity

import "xr-game-server/core/migrate"

const (
	TbGlobalCfg = "global_cfgs"
)

type GlobalCfgModule string

type GlobalCfg struct {
	migrate.OneModel
	Module     string `json:"module" gorm:"column:module;comment:模块编码"`
	ModuleName string `json:"moduleName" gorm:"column:module_name;default:'';comment:模块名称"`
	Key        string `json:"key" gorm:"column:key"`
	Value      string `json:"value" gorm:"column:value"`
	Desc       string `json:"desc" gorm:"column:desc"`
}

func initGlobalCfg() {
	migrate.AutoMigrate(&GlobalCfg{})
}
