package guilddto

import "github.com/gogf/gf/v2/frame/g"

type BatchUpdateGuildTimezoneReq struct {
	g.Meta   `path:"/batchUpdateGuildTimezone" method:"post" summary:"批量更新工会时区" tags:"直播工会"`
	GuildIds []uint64 `json:"guildIds" v:"required#请选择要更新的工会" dc:"工会ID列表"`
	Timezone int      `json:"timezone" v:"required#请选择时区" dc:"时区偏移量"`
}

type BatchUpdateGuildTimezoneRes struct {
	Success bool `json:"success"`
}
