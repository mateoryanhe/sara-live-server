package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// GetUserDetailReq CMS获取用户详情
type GetUserDetailReq struct {
	g.Meta `path:"/getUserDetail" method:"post" summary:"CMS获取用户详情" tags:"账号"`
	UserId uint64 `json:"userId,string" v:"required#用户ID不能为空" dc:"用户ID"`
}

// UserAccountDetailItem 账号信息(accounts)
type UserAccountDetailItem struct {
	ID              uint64     `json:"id,string"`
	OpenId          string     `json:"openId"`
	IP              string     `json:"ip"`
	RegisterIp      string     `json:"registerIp"`
	RegisterCountry string     `json:"registerCountry"`
	LoginCountry    string     `json:"loginCountry"`
	Channel         uint       `json:"channel"`
	PhoneAreaCode   string     `json:"phoneAreaCode"`
	Ban             bool       `json:"ban"`
	BanTime         *time.Time `json:"banTime"`
	BanApplyTime    *time.Time `json:"banApplyTime"`
	Cancel          bool       `json:"cancel"`
	CreatedAt       *time.Time `json:"createdAt"`
}

// UserProfileDetailItem 用户基础信息(user_infos,不含钱包)
type UserProfileDetailItem struct {
	Nickname        string     `json:"nickname"`
	Phone           string     `json:"phone"`
	Avatar          string     `json:"avatar"`
	Remark          string     `json:"remark"`
	ShareCode       string     `json:"shareCode"`
	UserType        uint8      `json:"userType"`
	IsAnchor        bool       `json:"isAnchor"`
	InviterId       uint64     `json:"inviterId,string"`
	VipLevel        uint32     `json:"vipLevel"`
	LastLoginTime   *time.Time `json:"lastLoginTime"`
	LiveRoomId      uint64     `json:"liveRoomId,string"`
	LiveRoomVer     uint64     `json:"liveRoomVer,string"`
	Gender          uint8      `json:"gender"`
	Birthday        *time.Time `json:"birthday"`
	BotAnchorStatus uint8      `json:"botAnchorStatus"`
	GuildId         uint64     `json:"guildId,string"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

// UserWalletDetailItem 钱包(user_infos.gold/diamond)
type UserWalletDetailItem struct {
	Gold    float64 `json:"gold"`
	Diamond float64 `json:"diamond"`
}

// UserExtDetailItem 用户扩展信息(user_exts)
type UserExtDetailItem struct {
	CanRank            bool       `json:"canRank"`
	PrettyId           uint64     `json:"prettyId,string"`
	PackageName        string     `json:"packageName"`
	AppVersion         string     `json:"appVersion"`
	FollowCount        uint64     `json:"followCount"`
	FollowerCount      uint64     `json:"followerCount"`
	CancelCode         string     `json:"cancelCode"`
	CancelCodeExpireAt *time.Time `json:"cancelCodeExpireAt"`
	RechargeWhitelist  bool       `json:"rechargeWhitelist"`
	UpdatedAt          *time.Time `json:"updatedAt"`
}

// UserCumulativeStatDetailItem 累计统计(user_cumulative_stats)
type UserCumulativeStatDetailItem struct {
	TotalRecharge       float64    `json:"totalRecharge"`
	TotalWithdraw       float64    `json:"totalWithdraw"`
	TotalPayCount       uint64     `json:"totalPayCount"`
	TotalDiamondConsume float64    `json:"totalDiamondConsume"`
	TotalGoldConsume    float64    `json:"totalGoldConsume"`
	TotalLiveDuration   float64    `json:"totalLiveDuration"`
	UpdatedAt           *time.Time `json:"updatedAt"`
}

// UserLoginDeviceDetailItem 登录设备(user_login_devices)
type UserLoginDeviceDetailItem struct {
	DeviceType  string     `json:"deviceType"`
	DeviceModel string     `json:"deviceModel"`
	CpuModel    string     `json:"cpuModel"`
	OsVersion   string     `json:"osVersion"`
	AppVersion  string     `json:"appVersion"`
	DeviceId    string     `json:"deviceId"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

// GetUserDetailRes CMS用户详情
type GetUserDetailRes struct {
	Account        *UserAccountDetailItem        `json:"account"`
	Profile        *UserProfileDetailItem        `json:"profile"`
	Wallet         *UserWalletDetailItem         `json:"wallet"`
	UserExt        *UserExtDetailItem            `json:"userExt"`
	CumulativeStat *UserCumulativeStatDetailItem `json:"cumulativeStat"`
	LoginDevice    *UserLoginDeviceDetailItem    `json:"loginDevice"`
}
