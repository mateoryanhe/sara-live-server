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
	// LiveRoomBanned 直播间已被封禁
	LiveRoomBanned XRCode = 53
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
	// LiveRoomChatMuted 您已被禁言
	LiveRoomChatMuted XRCode = 108
	// LiveRoomCannotMuteSelf 不能禁言自己
	LiveRoomCannotMuteSelf XRCode = 109
	// LiveRoomAudienceNotOnline 观众不在该直播间
	LiveRoomAudienceNotOnline XRCode = 110
	// LiveRoomCannotKickSelf 不能踢出自己
	LiveRoomCannotKickSelf XRCode = 117
	// LiveRoomKickBanned 被踢出后暂不可进入直播间
	LiveRoomKickBanned XRCode = 118
	// TextSensitiveWord 文本包含铭感词
	TextSensitiveWord XRCode = 111
	// AliyunTextModerationCfgInvalid 阿里云文本审核配置无效
	AliyunTextModerationCfgInvalid XRCode = 112
	// AliyunTextModerationFailed 阿里云文本审核调用失败
	AliyunTextModerationFailed XRCode = 113
	// ImageSensitiveContent 图片内容不合规
	ImageSensitiveContent XRCode = 114
	// ImageModerationCfgInvalid 图片审核配置无效
	ImageModerationCfgInvalid XRCode = 115
	// ImageModerationFailed 图片审核服务调用失败
	ImageModerationFailed XRCode = 116
	// TicketExist 门票名称已存在
	TicketExist XRCode = 119
	// TicketNonExist 门票不存在
	TicketNonExist XRCode = 121
	// PrivateRoomBillingNonExist 私密直播间计费配置不存在
	PrivateRoomBillingNonExist XRCode = 123
	// LiveRoomPrivateAudienceFull 私密直播间观众已满
	LiveRoomPrivateAudienceFull XRCode = 124
	// LiveRoomNotLive 直播间未在直播中(未开播或已关播)
	LiveRoomNotLive XRCode = 125
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
