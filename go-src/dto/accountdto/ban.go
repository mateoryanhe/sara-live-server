package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type BanReq struct {
	g.Meta       `path:"/ban" method:"post" summary:"封号" tags:"账号"`
	AccountId    uint64     `json:"accountId" dc:"账号id"`
	OpenId       string     `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel      uint       `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
	BanApplyTime *time.Time `json:"banApplyTime" dc:"封禁时间"`
}

// BanAnchorReq CMS封禁主播(含App推送,写入 live_rooms)
type BanAnchorReq struct {
	g.Meta       `path:"/banAnchor" method:"post" summary:"封禁主播" tags:"账号"`
	AccountId    uint64     `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
	OpenId       string     `json:"openId" dc:"登陆openId"`
	Channel      uint       `json:"channel" dc:"登陆渠道"`
	BanApplyTime *time.Time `json:"banApplyTime" v:"required#封禁截止时间不能为空" dc:"封禁截止时间"`
	BanReason    string     `json:"banReason" v:"required|length:1,512#封禁原因不能为空|封禁原因长度需在1到512之间" dc:"封禁原因"`
}

// UnBanAnchorReq CMS解封主播直播间
type UnBanAnchorReq struct {
	g.Meta    `path:"/unBanAnchor" method:"post" summary:"解封主播" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
	OpenId    string `json:"openId" dc:"登陆openId"`
	Channel   uint   `json:"channel" dc:"登陆渠道"`
}

type UnBanReq struct {
	g.Meta    `path:"/unBan" method:"post" summary:"解封" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
	OpenId    string `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel   uint   `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
}

type BanRes struct {
}

// BatchSetAnchorTimezoneReq CMS批量设置主播时区(仅限工会ID=0的主播)
type BatchSetAnchorTimezoneReq struct {
	g.Meta    `path:"/batchSetAnchorTimezone" method:"post" summary:"批量设置主播时区" tags:"账号"`
	AnchorIds []uint64 `json:"anchorIds" v:"required#主播ID列表不能为空" dc:"主播账号ID列表"`
	Timezone  int      `json:"timezone" v:"required#时区不能为空" dc:"时区值"`
}

type BatchSetAnchorTimezoneRes struct {
	SuccessCount int      `json:"successCount" dc:"成功数量"`
	FailCount    int      `json:"failCount" dc:"失败数量"`
	FailIds      []uint64 `json:"failIds" dc:"失败的主播ID"`
}

// ExitGuildReq CMS主播退出工会(将工会ID置为0)
type ExitGuildReq struct {
	g.Meta   `path:"/exitGuild" method:"post" summary:"退出工会" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播账号ID"`
}

type ExitGuildRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
