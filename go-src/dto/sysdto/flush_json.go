package sysdto

import "github.com/gogf/gf/v2/frame/g"

type FlushJsonReq struct {
	g.Meta `path:"/flushJson" method:"get" summary:"刷新json文件" tags:"系统"`
}

type FlushJsonResp struct {
}
