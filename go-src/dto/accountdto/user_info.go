package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

type QueryUserInfoReq struct {
	g.Meta `path:"/getUserInfo" method:"post" summary:"获取用户信息" tags:"账号"`
	httpserver.CMSQueryReq
	Key       string `json:"key" dc:"查询关键字"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
}

type UserInfoDto struct {
	ID           string     `json:"id"`
	CreatedAt    *time.Time `json:"createdAt"`
	OpenId       string     `json:"openId"`
	IP           string     `json:"ip"`
	Channel      uint       `json:"channel"`
	Ban          bool       `json:"ban"`
	BanTime      *time.Time `json:"banTime"`
	BanApplyTime *time.Time `json:"banApplyTime"`
	Cancel       bool       `json:"cancel"`
}
