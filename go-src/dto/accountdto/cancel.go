package accountdto

import "github.com/gogf/gf/v2/frame/g"

type CancelReq struct {
	g.Meta    `path:"/cancel" method:"post" summary:"注销" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
}

type UnCancelReq struct {
	g.Meta    `path:"/unCancel" method:"post" summary:"取消注销" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
}
