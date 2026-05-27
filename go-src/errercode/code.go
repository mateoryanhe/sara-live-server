package errercode

type XRCode int

const (
	Success = 0

	EmptyUserId         = 1
	EmptyToken          = 2
	Token               = 3
	TestEnvClose        = 4
	SysError            = 5
	LoginFail    XRCode = 6
	Ban          XRCode = 7

	CMSLoginFail    XRCode = 16
	ServerClose     XRCode = 24
	GuildExist      XRCode = 27
	GuildNonExist   XRCode = 29
	NoPermission    XRCode = 30
	GuildKickSelf   XRCode = 31
	GuildApplyExist XRCode = 34
	// DiamondAmountInvalid 钻石数量必须为正数
	DiamondAmountInvalid XRCode = 40
	// DiamondNotEnough 钻石余额不足
	DiamondNotEnough XRCode = 41
	// GoldAmountInvalid 金币数量必须为正数
	GoldAmountInvalid XRCode = 42
	// GoldNotEnough 金币余额不足
	GoldNotEnough XRCode = 43
	// DiamondOverflow 钻石数量溢出
	DiamondOverflow XRCode = 44
	// GoldOverflow 金币数量溢出
	GoldOverflow XRCode = 45
	// LiveRoomExist 主播已存在直播间
	LiveRoomExist XRCode = 50
	// LiveRoomNotAnchor 用户尚未成为主播(未加入工会)
	LiveRoomNotAnchor XRCode = 51
	// LiveRoomNotExist 直播间不存在
	LiveRoomNotExist XRCode = 52
	// GiftExist 礼物名称已存在
	GiftExist XRCode = 60
	// GiftNonExist 礼物不存在
	GiftNonExist XRCode = 61
	// GiftCountInvalid 送礼数量非法
	GiftCountInvalid XRCode = 62
	// GiftOffShelf 礼物已下架
	GiftOffShelf XRCode = 63
	// RechargeCfgExist 充值配置名称已存在
	RechargeCfgExist XRCode = 70
	// RechargeCfgNonExist 充值配置不存在
	RechargeCfgNonExist XRCode = 71
	// RechargeOrderNonExist 充值订单不存在
	RechargeOrderNonExist XRCode = 80
	// RechargeAmountInvalid 充值金额非法(必须为正)
	RechargeAmountInvalid XRCode = 81
	// RechargeGoldInvalid 充值发放金币数非法(必须为正)
	RechargeGoldInvalid XRCode = 82
	// RechargeOrderStateInvalid 订单状态不允许此操作(如重复完成)
	RechargeOrderStateInvalid XRCode = 83
	// RechargeCfgOffShelf 引用的充值配置已下架
	RechargeCfgOffShelf XRCode = 84
	// InvalidParam 参数无效
	InvalidParam XRCode = 90
	// RequestTooFrequent 请求过于频繁
	RequestTooFrequent XRCode = 91
	// DailyLimitExceeded 超过每日发送限制
	DailyLimitExceeded XRCode = 92
	// VerifyCodeInvalid 验证码无效
	VerifyCodeInvalid XRCode = 93
	// VerifyCodeExpired 验证码已过期
	VerifyCodeExpired XRCode = 94
	// AccountAlreadyExists 账号已存在
	AccountAlreadyExists XRCode = 95
	// BannerExist Banner标题已存在
	BannerExist XRCode = 96
	// BannerNonExist Banner不存在
	BannerNonExist XRCode = 97
	// UserAlreadyAnchor 用户已是主播
	UserAlreadyAnchor XRCode = 98
	// ShortVideoExist 短视频标题已存在
	ShortVideoExist XRCode = 99
	// ShortVideoNonExist 短视频不存在
	ShortVideoNonExist XRCode = 100
	// ShortVideoAlreadyLiked 短视频已点赞
	ShortVideoAlreadyLiked XRCode = 101
	// VipCfgExist VIP等级已存在
	VipCfgExist XRCode = 102
	// VipCfgNonExist VIP配置不存在
	VipCfgNonExist XRCode = 103
	// ShortVideoFileTooLarge 短视频文件超过大小限制
	ShortVideoFileTooLarge XRCode = 104
	// GameCfgExist 游戏编码已存在
	GameCfgExist XRCode = 105
	// GameCfgNonExist 游戏配置不存在
	GameCfgNonExist XRCode = 106
	// AgoraCfgInvalid 声网配置无效
	AgoraCfgInvalid XRCode = 107
)

type XError struct {
	error error // Wrapped error.
	xcode XRCode
	msg   string
	Param any
}

func (e *XError) Error() string {
	return e.msg
}

func (e *XError) Code() XRCode {
	return e.xcode
}

func CreateCode(code XRCode) error {
	return &XError{
		xcode: code,
	}
}

func CreateCodeAndParam(code XRCode, param any) error {
	return &XError{
		xcode: code,
		Param: param,
	}
}
