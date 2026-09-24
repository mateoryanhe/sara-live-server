package guilddto

import "github.com/gogf/gf/v2/frame/g"

const MaxBatchImmediateSettleGuilds = 100

// BatchImmediateSettleGuildsReq CMS批量立即结算所选工会。
type BatchImmediateSettleGuildsReq struct {
	g.Meta   `path:"/batchImmediateSettleGuilds" method:"post" summary:"CMS批量立即结算工会" tags:"直播工会"`
	GuildIds []uint64 `json:"guildIds,string" v:"required|min-length:1#请至少选择一个工会" dc:"工会ID列表"`
}

// BatchImmediateSettleGuildsRes CMS批量立即结算结果。
type BatchImmediateSettleGuildsRes struct {
	SettledCount int      `json:"settledCount" dc:"成功生成工会结算单数量"`
	NoDataCount  int      `json:"noDataCount" dc:"没有待结算数据的工会数量"`
	FailCount    int      `json:"failCount" dc:"结算失败数量"`
	FailGuildIds []string `json:"failGuildIds" dc:"结算失败的工会ID"`
}
