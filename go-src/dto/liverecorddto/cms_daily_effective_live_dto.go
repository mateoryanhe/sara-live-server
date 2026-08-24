package liverecorddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/guilddto"
)

// CMSDailyEffectiveLiveListReq CMS分页查询每日流水
type CMSDailyEffectiveLiveListReq struct {
	g.Meta `path:"/cmsDailyEffectiveLiveList" method:"post" summary:"CMS查询每日流水" tags:"直播记录"`
	httpserver.CMSQueryReq
	AnchorId         string   `json:"anchorId"  dc:"主播ID(可选,留空查全部)"`
	PlatformAnchorId string   `json:"platformAnchorId" dc:"平台主播ID(可选,兼容旧参数)"`
	GuildAnchorId    string   `json:"guildAnchorId" dc:"工会主播ID(可选,兼容旧参数)"`
	AnchorIds        []string `json:"anchorIds" dc:"主播ID列表(可选,多选)"`
	LiveDateStart    string   `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd      string   `json:"liveDateEnd"   dc:"日期止(YYYY-MM-DD,可选)"`
	Settled          int8     `json:"settled"       dc:"结算状态(-1全部,0未结算,1已结算)"`
}

// CMSDailyEffectiveLiveItem CMS每日流水列表项
type CMSDailyEffectiveLiveItem = guilddto.GuildAnchorDailyEffectiveLiveItem
