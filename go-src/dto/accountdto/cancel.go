package accountdto

import "github.com/gogf/gf/v2/frame/g"

type CancelReq struct {
	g.Meta    `path:"/cancel" method:"post" summary:"注销" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
	OpenId    string `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel   uint   `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
}

type UnCancelReq struct {
	g.Meta    `path:"/unCancel" method:"post" summary:"取消注销" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
	OpenId    string `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel   uint   `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
}
