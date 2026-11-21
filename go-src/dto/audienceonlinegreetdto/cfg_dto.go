package audienceonlinegreetdto

import "github.com/gogf/gf/v2/frame/g"

type GetAudienceOnlineAnchorGreetCfgReq struct {
	g.Meta `path:"/getAudienceOnlineAnchorGreetCfg" method:"post" summary:"查询观众上线主播招呼配置" tags:"观众上线主播招呼"`
}

type AudienceOnlineAnchorGreetCfgItem struct {
	ID             string   `json:"id"`
	Enabled        bool     `json:"enabled"`
	AnchorCount    uint8    `json:"anchorCount" dc:"每次最多几名在线主播打招呼"`
	IntervalMinSec uint32   `json:"intervalMinSec" dc:"排队招呼间隔下限(秒)"`
	IntervalMaxSec uint32   `json:"intervalMaxSec" dc:"排队招呼间隔上限(秒)"`
	ContentsEn     []string `json:"contentsEn" dc:"英文招呼列表"`
	ContentsEs     []string `json:"contentsEs" dc:"西班牙语招呼列表"`
	ContentsPt     []string `json:"contentsPt" dc:"葡萄牙语招呼列表"`
	ContentsHi     []string `json:"contentsHi" dc:"印地语招呼列表"`
	ContentsId     []string `json:"contentsId" dc:"印尼语招呼列表"`
	// UsingBuiltinContents 各语种 contents* 均为空时为 true(运行时走 Go 内置)
	UsingBuiltinContents bool     `json:"usingBuiltinContents"`
	BuiltinContentsEn    []string `json:"builtinContentsEn" dc:"Go 内置英文(只读参考)"`
	BuiltinContentsEs    []string `json:"builtinContentsEs"`
	BuiltinContentsPt    []string `json:"builtinContentsPt"`
	BuiltinContentsHi    []string `json:"builtinContentsHi"`
	BuiltinContentsId    []string `json:"builtinContentsId"`
	CreatedAt            string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type GetAudienceOnlineAnchorGreetCfgRes struct {
	Cfg *AudienceOnlineAnchorGreetCfgItem `json:"cfg"`
}

type SaveAudienceOnlineAnchorGreetCfgReq struct {
	g.Meta         `path:"/saveAudienceOnlineAnchorGreetCfg" method:"post" summary:"保存观众上线主播招呼配置" tags:"观众上线主播招呼"`
	ID             uint64   `json:"id" dc:"配置ID,首次保存可为0"`
	Enabled        bool     `json:"enabled" dc:"是否开启"`
	AnchorCount    uint8    `json:"anchorCount" v:"required|between:1,20#主播人数须在1-20之间" dc:"每次最多几名在线主播(按1v1通话次数选取)"`
	IntervalMinSec uint32 `json:"intervalMinSec" v:"between:0,3600#间隔下限须在0-3600秒" dc:"排队招呼间隔下限(秒)"`
	IntervalMaxSec uint32 `json:"intervalMaxSec" v:"between:0,3600#间隔上限须在0-3600秒" dc:"排队招呼间隔上限(秒)"`
}

type SaveAudienceOnlineAnchorGreetCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}