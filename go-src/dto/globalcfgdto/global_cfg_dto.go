package globalcfgdto

import (
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/entity"
)

type GlobalCfgDto struct {
	Id         string `json:"id"`
	Module     string `json:"module" gorm:"column:module"`
	ModuleName string `json:"moduleName" gorm:"column:module_name"`
	Key        string `json:"key" gorm:"column:key"`
	Value      string `json:"value" gorm:"column:value"`
	Desc       string `json:"desc" gorm:"column:desc"`
}

func NewGlobalCfgDto(cfg *entity.GlobalCfg) *GlobalCfgDto {
	ret := &GlobalCfgDto{}
	gconv.Struct(cfg, ret)
	return ret
}
