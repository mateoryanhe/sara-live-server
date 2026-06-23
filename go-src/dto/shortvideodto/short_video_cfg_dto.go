package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type GetShortVideoCfgReq struct {
	g.Meta `path:"/getShortVideoCfg" method:"post" summary:"查询短视频配置" tags:"短视频配置"`
}

type ShortVideoCfgItem struct {
	ID               string `json:"id"`
	MaxFileSize      uint64 `json:"maxFileSize"`
	MaxCoverFileSize uint32 `json:"maxCoverFileSize"`
	MaxDuration      uint32 `json:"maxDuration"`
	FreeWatchSeconds uint32 `json:"freeWatchSeconds"`
	EntryEnabled     uint8  `json:"entryEnabled"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetShortVideoCfgRes struct {
	Cfg *ShortVideoCfgItem `json:"cfg"`
}

type SaveShortVideoCfgReq struct {
	g.Meta           `path:"/saveShortVideoCfg" method:"post" summary:"保存短视频配置" tags:"短视频配置"`
	ID               uint64 `json:"id,string" dc:"配置ID,首次保存传0"`
	MaxFileSize      uint64 `json:"maxFileSize" v:"required|min:1#最大文件大小不能为空|最大文件大小必须大于0" dc:"最大文件大小(字节)"`
	MaxCoverFileSize uint32 `json:"maxCoverFileSize" v:"required|min:1#封面图片大小不能为空|封面图片大小必须大于0" dc:"封面图片最大大小(M)"`
	MaxDuration      uint32 `json:"maxDuration" v:"required|min:1#最大时长不能为空|最大时长必须大于0" dc:"最大时长(秒)"`
	FreeWatchSeconds uint32 `json:"freeWatchSeconds" v:"required|min:0#免费观看时长不能为空|免费观看时长不能小于0" dc:"免费观看时长(秒)"`
	EntryEnabled     uint8  `json:"entryEnabled" v:"in:0,1#入口开关取值无效" dc:"入口开关(0关闭,1开启)"`
}

type SaveShortVideoCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type AppShortVideoCfgReq struct {
	g.Meta `path:"/appShortVideoCfg" method:"post" summary:"App查询短视频配置" tags:"短视频"`
}

type AppShortVideoCfgRes struct {
	MaxFileSize      uint64 `json:"maxFileSize"`
	MaxCoverFileSize uint32 `json:"maxCoverFileSize"`
	MaxDuration      uint32 `json:"maxDuration"`
	FreeWatchSeconds uint32 `json:"freeWatchSeconds"`
	EntryEnabled     uint8  `json:"entryEnabled"`
}
