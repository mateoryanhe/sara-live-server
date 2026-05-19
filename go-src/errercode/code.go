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
