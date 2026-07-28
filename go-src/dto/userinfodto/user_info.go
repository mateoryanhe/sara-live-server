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
	UserId        uint64  `json:"userId,string"`
	Nickname      string  `json:"nickname"`
	Phone         string  `json:"phone"`
	Avatar        string  `json:"avatar"`
	Remark        string  `json:"remark"`
	Gold          float64 `json:"gold"`
	Diamond       float64 `json:"diamond"`
	ShareCode     string  `json:"shareCode"`
	VipLevel      uint32  `json:"vipLevel"`
	IsAnchor      bool    `json:"isAnchor" dc:"是否主播"`
	UserType      uint8   `json:"userType" dc:"用户类型(6测试型主播仅App端使用,CMS展示同普通主播)"`
	HasLiveRoom   bool    `json:"hasLiveRoom" dc:"是否已创建直播间"`
	Gender        uint8   `json:"gender" dc:"性别(0未知,1男,2女)"`
	Birthday      string  `json:"birthday" dc:"出生日期(YYYY-MM-DD,空表示未设置)"`
	FollowCount   int     `json:"followCount" dc:"用户关注数"`
	FollowerCount int     `json:"followerCount" dc:"用户粉丝数"`
	FollowStatus  uint8   `json:"followStatus" dc:"关注状态(0未关注,1已关注,2互为好友)"`
	TotalIncome   float64 `json:"totalIncome" dc:"最近30天收益(钻石,来自主播红人榜,未上榜为0)"`
	Age           int64   `json:"age"`
}

// GetUserExtReq 查询用户扩展信息
type GetUserExtReq struct {
	g.Meta `path:"/getUserExt" method:"post" summary:"获取用户扩展信息" tags:"用户基础信息"`
	UserId uint64 `json:"userId" dc:"用户ID,不传则查当前登录用户"`
}

type GetUserExtRes struct {
	UserId        uint64 `json:"userId,string"`
	CanRank       bool   `json:"canRank" dc:"是否可上排行榜"`
	PackageName   string `json:"packageName" dc:"注册包名"`
	AppVersion    string `json:"appVersion" dc:"注册版本号"`
	FollowCount   uint64 `json:"followCount" dc:"当前关注数"`
	FollowerCount uint64 `json:"followerCount" dc:"当前粉丝数"`
	CancelCode    string `json:"cancelCode" dc:"注销码"`
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

// CancelAccountReq App端销户(需登录)
type CancelAccountReq struct {
	g.Meta `path:"/cancelAccount" method:"post" summary:"App端销户(需登录)" tags:"用户基础信息"`
}

type CancelAccountRes struct {
	Success bool `json:"success"`
}

// CancelAccountByCodeReq 通过注销码注销账号(官网)
type CancelAccountByCodeReq struct {
	g.Meta     `path:"/cancelAccountByCode" method:"post" summary:"通过注销码注销账号" tags:"用户基础信息"`
	CancelCode string `json:"cancelCode" v:"required#注销码不能为空" dc:"注销码"`
}

type CancelAccountByCodeRes struct {
	Success bool `json:"success"`
}

// AppFeedbackReq App端用户反馈(不解析body,占位接口)
type AppFeedbackReq struct {
	g.Meta `path:"/feedback" method:"post" summary:"用户反馈" tags:"用户反馈"`
}

type AppFeedbackRes struct {
	Success bool `json:"success"`
}

// AppReportReq App端用户举报(不解析body,占位接口)
type AppReportReq struct {
	g.Meta `path:"/report" method:"post" summary:"用户举报" tags:"用户举报"`
}

type AppReportRes struct {
	Success bool `json:"success"`
}
