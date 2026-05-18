package userinfodto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetUserInfoReq 查询当前登录用户的基础信息
type GetUserInfoReq struct {
	g.Meta `path:"/get" method:"post" summary:"获取用户基础信息" tags:"用户基础信息"`
}

type GetUserInfoRes struct {
	UserId    uint64  `json:"userId"`
	Nickname  string  `json:"nickname"`
	Phone     string  `json:"phone"`
	Avatar    string  `json:"avatar"`
	Remark    string  `json:"remark"`
	Gold      float64 `json:"gold"`
	Diamond   float64 `json:"diamond"`
	ShareCode string  `json:"shareCode"`
}

// UpdateNicknameReq 修改昵称
type UpdateNicknameReq struct {
	g.Meta   `path:"/updateNickname" method:"post" summary:"修改用户昵称" tags:"用户基础信息"`
	Nickname string `json:"nickname" v:"required#昵称不能为空" dc:"用户昵称"`
}

type UpdateNicknameRes struct {
	Nickname string `json:"nickname"`
}

// GetCurrencyLogReq 查询用户货币流水
type GetCurrencyLogReq struct {
	g.Meta    `path:"/getCurrencyLog" method:"post" summary:"获取用户货币流水" tags:"用户基础信息"`
	UserId    uint64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	PageIndex int    `json:"pageIndex" dc:"页码,从1开始,默认1"`
	PageSize  int    `json:"pageSize" dc:"每页数量,默认20"`
}

type CurrencyLogItem struct {
	Id       uint64  `json:"id"`
	UserId   uint64  `json:"userId"`
	Type     uint8   `json:"type"`   // 1金币 2钻石
	Action   uint8   `json:"action"` // 1加 2减
	Amount   float64 `json:"amount"`
	Before   float64 `json:"before"`
	After    float64 `json:"after"`
	Reason   string  `json:"reason"`
	CreateAt int64   `json:"createAt"`
}

type GetCurrencyLogRes struct {
	Total int                `json:"total"`
	List  []*CurrencyLogItem `json:"list"`
}
