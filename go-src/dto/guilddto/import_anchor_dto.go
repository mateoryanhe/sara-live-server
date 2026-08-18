package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ImportGuildAnchorsReq struct {
	g.Meta     `path:"/importGuildAnchors" method:"post" summary:"CSV导入工会主播" tags:"直播工会"`
	GuildId    uint64                  `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	AnchorType uint8                   `json:"anchorType" v:"required#主播类型不能为空" dc:"主播类型(1=普通主播,7=高级主播)"`
	Rows       []*ImportGuildAnchorRow `json:"rows" v:"required#导入数据不能为空" dc:"导入数据行"`
}

type ImportGuildAnchorRow struct {
	UserId uint64 `json:"userId" dc:"用户ID"`
}

type ImportGuildAnchorsRes struct {
	SuccessCount int                          `json:"successCount" dc:"成功数量"`
	FailCount    int                          `json:"failCount" dc:"失败数量"`
	Fails        []*ImportGuildAnchorFailItem `json:"fails" dc:"失败列表"`
}

type ImportGuildAnchorFailItem struct {
	UserId   string `json:"userId" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"用户昵称"`
	Reason   int    `json:"reason" dc:"失败原因代码"`
}

const (
	ImportAnchorFailUserNotFound       = 1 // 用户不存在
	ImportAnchorFailCancelCodeMismatch = 2 // 注销码不匹配(已废弃)
	ImportAnchorFailCancelCodeExpired  = 3 // 注销码已过期(已废弃)
	ImportAnchorFailAlreadyInGuild     = 4 // 已在其他工会
	ImportAnchorFailCannotSetAnchor    = 5 // 无法设置为主播
	ImportAnchorFailAlreadyHasLiveRoom = 6 // 主播间缓存已存在
)

func GetImportAnchorFailReasonText(reason int) string {
	switch reason {
	case ImportAnchorFailUserNotFound:
		return "User not found"
	case ImportAnchorFailCancelCodeMismatch:
		return "Cancel code mismatch"
	case ImportAnchorFailCancelCodeExpired:
		return "Cancel code expired"
	case ImportAnchorFailAlreadyInGuild:
		return "Already in another guild"
	case ImportAnchorFailCannotSetAnchor:
		return "Cannot set as anchor"
	case ImportAnchorFailAlreadyHasLiveRoom:
		return "Live room already exists"
	default:
		return "Unknown error"
	}
}

type SetAnchorGuildReq struct {
	g.Meta     `path:"/joinGuildAnchor" method:"post" summary:"CMS加入工会主播" tags:"直播工会"`
	GuildId    uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	UserId     uint64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	AnchorType uint8  `json:"anchorType" v:"required#主播类型不能为空" dc:"主播类型(1=普通主播,7=高级主播)"`
}

type SetAnchorGuildRes struct {
	Success  bool   `json:"success" dc:"是否成功"`
	Reason   int    `json:"reason" dc:"失败原因代码(0=成功)"`
	Nickname string `json:"nickname" dc:"用户昵称"`
}

// SetGuildAnchorTypeReq CMS设置工会主播类型
type SetGuildAnchorTypeReq struct {
	g.Meta     `path:"/setGuildAnchorType" method:"post" summary:"CMS设置工会主播类型" tags:"直播工会"`
	GuildId    uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	UserId     uint64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	AnchorType uint8  `json:"anchorType" v:"required#主播类型不能为空" dc:"主播类型(1=普通主播,7=高级主播)"`
}

type SetGuildAnchorTypeRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

type BatchExitGuildReq struct {
	GuildId   uint64   `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	AnchorIds []uint64 `json:"anchorIds" v:"required#主播ID列表不能为空" dc:"主播ID列表"`
}

type BatchExitGuildRes struct {
	SuccessCount int      `json:"successCount" dc:"成功数量"`
	FailCount    int      `json:"failCount" dc:"失败数量"`
	FailIds      []uint64 `json:"failIds" dc:"失败的主播ID"`
}
