package agoradto

import "github.com/gogf/gf/v2/frame/g"

// GetAppIdReq App端获取声网AppId
type GetAppIdReq struct {
	g.Meta `path:"/appId" method:"post" summary:"获取声网AppId" tags:"声网"`
}

// GetAppIdRes App端声网AppId
type GetAppIdRes struct {
	AppId string `json:"appId" dc:"声网AppId"`
}
