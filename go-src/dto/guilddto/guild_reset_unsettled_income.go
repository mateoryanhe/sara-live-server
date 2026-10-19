package guilddto

import "github.com/gogf/gf/v2/frame/g"

// BatchResetGuildUnsettledIncomeReq CMS批量清零工会及名下主播未结算收益(不写结算/归档日志)
type BatchResetGuildUnsettledIncomeReq struct {
	g.Meta   `path:"/batchResetGuildUnsettledIncome" method:"post" summary:"CMS批量清零工会未结算收益" tags:"直播工会"`
	GuildIds []uint64 `json:"guildIds,string" v:"required|min-length:1#请至少选择一个工会" dc:"工会ID列表"`
}

// BatchResetGuildUnsettledIncomeRes CMS批量清零工会未结算收益结果
type BatchResetGuildUnsettledIncomeRes struct {
	ResetCount   int      `json:"resetCount" dc:"成功清零的工会数量"`
	NoDataCount  int      `json:"noDataCount" dc:"工会及名下主播未结算已为0的数量"`
	FailCount    int      `json:"failCount" dc:"失败数量"`
	FailGuildIds []string `json:"failGuildIds" dc:"失败的工会ID"`
}
