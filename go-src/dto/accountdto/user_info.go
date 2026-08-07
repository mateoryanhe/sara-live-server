package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

type QueryUserInfoReq struct {
	g.Meta `path:"/getUserInfo" method:"post" summary:"获取用户信息" tags:"账号"`
	httpserver.CMSQueryReq
	Key       string `json:"key" dc:"查询关键字(用户ID模糊/openId精确)"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
}

type UserInfoDto struct {
	ID              uint64     `json:"id,string"`
	CreatedAt       *time.Time `json:"createdAt"`
	OpenId          string     `json:"openId"`
	IP              string     `json:"ip"`
	RegisterIp      string     `json:"registerIp"`
	RegisterCountry string     `json:"registerCountry"`
	LoginCountry    string     `json:"loginCountry"`
	Channel         uint       `json:"channel"`
	Ban             bool       `json:"ban"`
	BanTime         *time.Time `json:"banTime"`
	BanApplyTime    *time.Time `json:"banApplyTime"`
	Cancel          bool       `json:"cancel"`
	PhoneAreaCode   string     `json:"phoneAreaCode"`
	// 以下字段来自 user_infos 表(LEFT JOIN,可能为空)
	Nickname      string     `json:"nickname"`
	Phone         string     `json:"phone"`
	Avatar        string     `json:"avatar"`
	Remark        string     `json:"remark"`
	Gold          float64    `json:"gold"`
	Diamond       float64    `json:"diamond"`
	ShareCode     string     `json:"shareCode"`
	GuildId       uint64     `json:"guildId"`
	UserType      uint8      `json:"userType"`
	IsAnchor      bool       `json:"isAnchor"`
	VipLevel      uint32     `json:"vipLevel"`
	LastLoginTime *time.Time `json:"lastLoginTime"`
	DeviceType    string     `json:"deviceType"`
	PackageName   string     `json:"packageName"`
	AppVersion    string     `json:"appVersion"`
	CanRank       bool       `json:"canRank"`
	CancelCode    string     `json:"cancelCode"`
}

type SetAnchorReq struct {
	g.Meta    `path:"/setAnchor" method:"post" summary:"设为用户主播(不可回退)" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#用户ID不能为空" dc:"用户ID"`
}

type SetAnchorRes struct {
	Success bool `json:"success"`
}

type SetSeniorAnchorReq struct {
	g.Meta    `path:"/setSeniorAnchor" method:"post" summary:"设为高级主播(不可回退)" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#用户ID不能为空" dc:"用户ID"`
}

type SetSeniorAnchorRes struct {
	Success bool `json:"success"`
}

type SetUserTypeReq struct {
	g.Meta    `path:"/setUserType" method:"post" summary:"修改用户类型" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#用户ID不能为空" dc:"用户ID"`
	UserType  uint8  `json:"userType" v:"required|in:0,4#用户类型不能为空|仅允许普通用户或测试人员" dc:"用户类型(0普通用户,4测试人员)"`
}

type SetUserTypeRes struct {
	Success bool `json:"success"`
}

type SetCanRankReq struct {
	g.Meta    `path:"/setCanRank" method:"post" summary:"设置用户是否可上排行榜" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#用户ID不能为空" dc:"用户ID"`
	CanRank   bool   `json:"canRank" dc:"是否可上排行榜"`
}

type SetCanRankRes struct {
	Success bool `json:"success"`
}
