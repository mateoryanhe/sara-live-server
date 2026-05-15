package globalcfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

type DelGlobalCfgReq struct {
	g.Meta `path:"/delGlobalCfg" method:"post" summary:"删除全局配置" tags:"全局配置"`
	*entity.GlobalCfg
}
