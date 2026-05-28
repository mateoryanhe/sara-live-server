package apptokendto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type SaveAppTokenReq struct {
	g.Meta   `path:"/saveAppToken" method:"post" summary:"保存App Token" tags:"App Token"`
	Id       uint64     `json:"id" v:"required#用户ID不能为空" summary:"用户ID"`
	Token    string     `json:"token" summary:"Token值,新增时留空则自动生成"`
	ExpireAt *time.Time `json:"expireAt" summary:"过期时间,留空则默认7天"`
}
