package guilddto

import "github.com/gogf/gf/v2/frame/g"

// 工会主播 CSV 导入失败原因
const (
	ImportAnchorFailUserNotFound       = 1 // 用户不存在
	ImportAnchorFailCancelCodeMismatch = 2 // 注销码错误
	ImportAnchorFailCancelCodeExpired  = 3 // 注销码已过期
	ImportAnchorFailAlreadyInGuild     = 4 // 已加入工会
	ImportAnchorFailCannotSetAnchor    = 5 // 无法设为主播(非普通用户)
)

// ImportGuildAnchorRow CSV 行
type ImportGuildAnchorRow struct {
	UserId     uint64 `json:"userId,string" dc:"用户ID"`
	CancelCode string `json:"cancelCode" dc:"注销码"`
}

// ImportGuildAnchorsReq CMS 按工会 CSV 导入主播
type ImportGuildAnchorsReq struct {
	g.Meta     `path:"/importGuildAnchors" method:"post" summary:"CSV导入工会主播" tags:"直播工会"`
	GuildId    uint64                  `json:"guildId" v:"required#工会ID不能为空" dc:"目标工会ID"`
	AnchorType uint8                   `json:"anchorType" v:"required|in:1,7#主播类型不能为空|仅支持普通主播或高级主播" dc:"主播类型(1普通主播,7高级主播)"`
	Rows       []*ImportGuildAnchorRow `json:"rows" v:"required|min-length:1#请至少导入一行数据" dc:"CSV解析后的行"`
}

type ImportGuildAnchorFailItem struct {
	UserId   string `json:"userId" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Reason   int    `json:"reason" dc:"失败原因(1用户不存在,2注销码错误,3注销码过期,4已加入工会,5无法设为主播)"`
}

type ImportGuildAnchorsRes struct {
	SuccessCount int                          `json:"successCount"`
	FailCount    int                          `json:"failCount"`
	Fails        []*ImportGuildAnchorFailItem `json:"fails"`
}
