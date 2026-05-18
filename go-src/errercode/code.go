package errercode

type XRCode int

const (
	Success = 0

	EmptyUserId          = 1
	EmptyToken           = 2
	Token                = 3
	TestEnvClose         = 4
	SysError             = 5
	LoginFail     XRCode = 6
	Ban           XRCode = 7
	RoleNameEmpty XRCode = 8
	RoleNameOut   XRCode = 9
	NameRepeat    XRCode = 10
	RoleRepeat    XRCode = 11
	RoleError     XRCode = 12
	// ItemNon 道具不存在
	ItemNon XRCode = 13
	// PropNon 背包处理器不存在
	PropNon XRCode = 14
	// ItemNumber 物品数量不足
	ItemNumber       XRCode = 15
	CMSLoginFail     XRCode = 16
	ServerNotExist   XRCode = 17
	ServerNotOpen    XRCode = 18
	EmptyString      XRCode = 19
	EmptyZoneId      XRCode = 20
	TaskNon          XRCode = 21
	TaskNotFinished  XRCode = 22
	TaskGot          XRCode = 23
	ServerClose      XRCode = 24
	RankNon          XRCode = 25
	ChatCD           XRCode = 26
	GuildExist       XRCode = 27
	GuildApplyBlack  XRCode = 28
	GuildNonExist    XRCode = 29
	NoPermission     XRCode = 30
	GuildKickSelf    XRCode = 31
	GuildFull        XRCode = 32
	GuildMaster      XRCode = 33
	GuildApplyExist  XRCode = 34
	GuildApplyNumMax XRCode = 35
	MailIsRead       XRCode = 36
	MailIsReceived   XRCode = 37
	SysBuyLimit      XRCode = 38
	GuildMasterRole  XRCode = 39
	// DiamondAmountInvalid 钻石数量必须为正数
	DiamondAmountInvalid XRCode = 40
	// DiamondNotEnough 钻石余额不足
	DiamondNotEnough XRCode = 41
	// GoldAmountInvalid 金币数量必须为正数
	GoldAmountInvalid XRCode = 42
	// GoldNotEnough 金币余额不足
	GoldNotEnough XRCode = 43
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
