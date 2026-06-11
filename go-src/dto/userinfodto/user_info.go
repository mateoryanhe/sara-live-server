package userinfodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// GetUserInfoReq 查询当前登录用户的基础信息
type GetUserInfoReq struct {
	g.Meta `path:"/get" method:"post" summary:"获取用户基础信息" tags:"用户基础信息"`
	UserId uint64 `json:"userId"`
}

type GetUserInfoRes struct {
	UserId        uint64  `json:"userId"`
	Nickname      string  `json:"nickname"`
	Phone         string  `json:"phone"`
	Avatar        string  `json:"avatar"`
	Remark        string  `json:"remark"`
	Gold          float64 `json:"gold"`
	Diamond       float64 `json:"diamond"`
	ShareCode     string  `json:"shareCode"`
	VipLevel      uint32  `json:"vipLevel"`
	IsAnchor      bool    `json:"isAnchor" dc:"是否主播"`
	HasLiveRoom   bool    `json:"hasLiveRoom" dc:"是否已创建直播间"`
	Gender        uint8   `json:"gender" dc:"性别(0未知,1男,2女)"`
	Birthday      string  `json:"birthday" dc:"出生日期(YYYY-MM-DD,空表示未设置)"`
	FollowCount   int     `json:"followCount" dc:"用户关注数"`
	FollowerCount int     `json:"followerCount" dc:"用户粉丝数"`
	TotalIncome   float64 `json:"totalIncome" dc:"主播总收益"`
}

// UpdateGenderReq 修改性别
type UpdateGenderReq struct {
	g.Meta `path:"/updateGender" method:"post" summary:"修改用户性别" tags:"用户基础信息"`
	Gender uint8 `json:"gender" v:"required#性别不能为空" dc:"性别(0未知,1男,2女)"`
}

type UpdateGenderRes struct {
	Gender uint8 `json:"gender"`
}

// UpdateBirthdayReq 修改出生日期
type UpdateBirthdayReq struct {
	g.Meta   `path:"/updateBirthday" method:"post" summary:"修改用户出生日期" tags:"用户基础信息"`
	Birthday string `json:"birthday" v:"required#出生日期不能为空" dc:"出生日期,格式YYYY-MM-DD"`
}

type UpdateBirthdayRes struct {
	Birthday string `json:"birthday"`
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
	Reason   uint8   `json:"reason"` // 货币变动原因枚举,参见 constants/currency.Reason
	CreateAt int64   `json:"createAt"`
}

type GetCurrencyLogRes struct {
	Total int                `json:"total"`
	List  []*CurrencyLogItem `json:"list"`
}

// UploadAvatarReq 上传头像
type UploadAvatarReq struct {
	g.Meta `path:"/uploadAvatar" method:"post" mime:"multipart/form-data" summary:"上传用户头像" tags:"用户基础信息"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择头像图片" dc:"头像图片文件"`
}

type UploadAvatarRes struct {
	Avatar string `json:"avatar" dc:"头像文件名"`
}

// CancelAccountReq App端销户
type CancelAccountReq struct {
	g.Meta `path:"/cancelAccount" method:"post" summary:"销户" tags:"用户基础信息"`
}

type CancelAccountRes struct {
	Success bool `json:"success"`
}
