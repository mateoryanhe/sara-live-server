package authdto

import "github.com/gogf/gf/v2/frame/g"

type PushAgainLoginReq struct {
	g.Meta `path:"/push/againLogin" method:"post" summary:"推送客户端挤号" description:"命令 8" tags:"权限"`
}
