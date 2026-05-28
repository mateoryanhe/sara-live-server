package statdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/snowflake"
	"xr-game-server/entity"
)

var sysSata *entity.SystemTotalStat

func initSysSata() {

	g.Model(string(entity.TbSystemTotalStat)).Limit(1).Scan(&sysSata)
	if sysSata == nil {
		sysSata = entity.NewSystemTotalStat(snowflake.GetId())
	}
}

func GetSysStat() *entity.SystemTotalStat {
	return sysSata
}
