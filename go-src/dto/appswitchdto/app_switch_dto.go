package appswitchdto

import "github.com/gogf/gf/v2/frame/g"

// GetAppSwitchesReq App 登录后下发的功能开关(与 /sysInfo/cfg 免登录大包分离,便于按需拉取与扩展).
type GetAppSwitchesReq struct {
	g.Meta `path:"/getAppSwitches" method:"post" summary:"App功能开关" tags:"App开关"`
}

type GetAppSwitchesRes struct {
	// PersonalAnchorAlbumEnabled 个人主播相册入口是否展示
	PersonalAnchorAlbumEnabled bool `json:"personalAnchorAlbumEnabled" dc:"个人主播相册开关"`
}