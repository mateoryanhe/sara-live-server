package indexdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type IndexReq struct {
	g.Meta `path:"/getAll" method:"post" summary:"首页基础信息" tags:"首页"`
}
