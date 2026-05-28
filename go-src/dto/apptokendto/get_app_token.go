package apptokendto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GetAppTokenReq struct {
	g.Meta `path:"/getAppToken" method:"post" summary:"查询App Token" tags:"App Token"`
	httpserver.CMSQueryReq
	UserId string `json:"userId" summary:"用户ID"`
}
