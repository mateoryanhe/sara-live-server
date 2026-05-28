package apptokendto

import "github.com/gogf/gf/v2/frame/g"

type ReloadAppTokenReq struct {
	g.Meta `path:"/reloadAppToken" method:"post" summary:"重新加载App Token缓存" tags:"App Token"`
}
