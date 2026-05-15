package globalcfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

type SaveGlobalCfgReq struct {
	g.Meta `path:"/saveGlobalCfg" method:"post" summary:"保存全局配置" tags:"全局配置"`
	*entity.GlobalCfg
}
