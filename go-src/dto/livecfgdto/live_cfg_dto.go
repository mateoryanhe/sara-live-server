package livecfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetLiveCfgReq struct {
	g.Meta `path:"/getLiveCfg" method:"post" summary:"查询直播配置" tags:"直播配置"`
}

type LiveCfgItem struct {
	ID                              string  `json:"id"`
	PaidDanmakuPrice                float64 `json:"paidDanmakuPrice"`
	VideoCallTicketEnabled          bool    `json:"videoCallTicketEnabled"`
	AudienceListRefreshSeconds      uint32  `json:"audienceListRefreshSeconds"`
	OneToOneDailyFreeSeconds        uint32  `json:"oneToOneDailyFreeSeconds" dc:"1v1观众对单主播每日免费通话秒数(0=关闭)"`
	LiveRoomCallMinTotalRechargeUsd float64 `json:"liveRoomCallMinTotalRechargeUsd" dc:"历史配置项(当前不控制allowCallIcon)"`
	ShowcaseLiveRoomEnabled         bool    `json:"showcaseLiveRoomEnabled" dc:"App是否展示秀场直播间入口(category=1)"`
	AudienceGameListHidden          bool    `json:"audienceGameListHidden" dc:"开启后非主播用户App游戏列表与直播间推荐游戏为空"`
	AnchorGameListHidden            bool    `json:"anchorGameListHidden" dc:"开启后主播用户App游戏列表与直播间推荐游戏为空"`
	GameLiveRoomEnabled             bool    `json:"gameLiveRoomEnabled" dc:"App是否展示游戏直播间入口(category=2)"`
	VoiceChatLiveRoomEnabled        bool    `json:"voiceChatLiveRoomEnabled" dc:"App是否展示语聊房入口(category=5)"`
	OneToOneLiveRoomEnabled         bool    `json:"oneToOneLiveRoomEnabled" dc:"App是否展示1v1房间页面/入口"`
	MaxAnchorPhotos                 uint32  `json:"maxAnchorPhotos" dc:"主播相册照片上限,默认100"`
	MaxAnchorVideos                 uint32  `json:"maxAnchorVideos" dc:"主播相册视频上限,默认100"`
	MaxAnchorVideoFileSizeMB        uint32  `json:"maxAnchorVideoFileSizeMB" dc:"主播相册单个视频大小上限(MB),默认100,/cfg下发"`
	CreatedAt                       string  `json:"createdAt"`
	UpdatedAt                       string  `json:"updatedAt"`
}

type GetLiveCfgRes struct {
	Cfg *LiveCfgItem `json:"cfg"`
}

type SaveLiveCfgReq struct {
	g.Meta                          `path:"/saveLiveCfg" method:"post" summary:"保存直播配置" tags:"直播配置"`
	ID                              uint64  `json:"id" dc:"配置ID,首次保存可为0"`
	PaidDanmakuPrice                float64 `json:"paidDanmakuPrice" v:"required|min:0#付费弹幕价格不能为空|付费弹幕价格不能小于0"`
	VideoCallTicketEnabled          *bool   `json:"videoCallTicketEnabled" dc:"直播间视频通话接通是否扣门票;不传时保留原值,首次保存默认开启"`
	AudienceListRefreshSeconds      uint32  `json:"audienceListRefreshSeconds" v:"required|min:1#在线观众列表刷新时间不能为空|在线观众列表刷新时间不能小于1秒" dc:"在线观众列表刷新间隔(秒),默认300"`
	OneToOneDailyFreeSeconds        uint32  `json:"oneToOneDailyFreeSeconds" dc:"1v1观众对单主播每日免费通话秒数(0=关闭,默认30)"`
	LiveRoomCallMinTotalRechargeUsd float64 `json:"liveRoomCallMinTotalRechargeUsd" v:"min:0#累计充值门槛不能小于0" dc:"保留字段(当前不控制allowCallIcon),默认10"`
	ShowcaseLiveRoomEnabled         *bool   `json:"showcaseLiveRoomEnabled" dc:"App是否展示秀场入口;不传时保留原值,首次保存默认关闭"`
	AudienceGameListHidden          *bool   `json:"audienceGameListHidden" dc:"开启后非主播用户游戏列表为空;不传时保留原值,首次保存默认关闭"`
	AnchorGameListHidden            *bool   `json:"anchorGameListHidden" dc:"开启后主播用户游戏列表为空;不传时保留原值,首次保存默认关闭"`
	GameLiveRoomEnabled             *bool   `json:"gameLiveRoomEnabled" dc:"App是否展示游戏直播间入口;不传时保留原值,首次保存默认关闭"`
	VoiceChatLiveRoomEnabled        *bool   `json:"voiceChatLiveRoomEnabled" dc:"App是否展示语聊房入口;不传时保留原值,首次保存默认关闭"`
	OneToOneLiveRoomEnabled         *bool   `json:"oneToOneLiveRoomEnabled" dc:"App是否展示1v1页面/入口;不传时保留原值,首次保存默认关闭"`
	MaxAnchorPhotos                 uint32  `json:"maxAnchorPhotos" v:"required|min:1#相册照片上限不能为空|相册照片上限至少为1" dc:"主播相册照片上限,默认100"`
	MaxAnchorVideos                 uint32  `json:"maxAnchorVideos" v:"required|min:1#相册视频上限不能为空|相册视频上限至少为1" dc:"主播相册视频上限,默认100"`
	MaxAnchorVideoFileSizeMB        uint32  `json:"maxAnchorVideoFileSizeMB" v:"required|min:1#相册视频大小上限不能为空|相册视频大小上限至少为1MB" dc:"主播相册单个视频大小上限(MB),默认100"`
}

type SaveLiveCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
