package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// PlatformAnchorListForVisibilityReq 可见性管理页拉取全部平台主播。
type PlatformAnchorListForVisibilityReq struct {
	g.Meta `path:"/platformAnchorListForVisibility" method:"post" summary:"可见性管理拉取全部平台主播" tags:"账号"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"查询关键字(用户ID/昵称/手机号/分享码)"`
}

// PlatformAnchorVisibilityItem 平台主播可见性记录。
type PlatformAnchorVisibilityItem struct {
	Id           string `json:"id" dc:"可见性记录ID"`
	AnchorId     string `json:"anchorId" dc:"平台主播ID"`
	AnchorName   string `json:"anchorName" dc:"平台主播昵称"`
	AnchorAvatar string `json:"anchorAvatar" dc:"平台主播头像"`
	CmsUserId    string `json:"cmsUserId" dc:"CMS用户ID"`
	CmsUserName  string `json:"cmsUserName" dc:"CMS用户名"`
	CreatedAt    string `json:"createdAt" dc:"授权时间"`
}

// PlatformAnchorVisibilityListReq 查询平台主播已授权的 CMS 用户。
type PlatformAnchorVisibilityListReq struct {
	g.Meta   `path:"/platformAnchorVisibilityList" method:"post" summary:"查询平台主播可见性授权列表" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#平台主播ID不能为空" dc:"平台主播ID"`
}

type PlatformAnchorVisibilityListRes struct {
	List []*PlatformAnchorVisibilityItem `json:"list"`
}

// PlatformAnchorVisibilityByUserListReq 查询指定 CMS 用户已授权的平台主播。
type PlatformAnchorVisibilityByUserListReq struct {
	g.Meta    `path:"/platformAnchorVisibilityByUserList" method:"post" summary:"查询CMS用户已授权平台主播列表" tags:"账号"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

type PlatformAnchorVisibilityByUserListRes struct {
	List []*PlatformAnchorVisibilityItem `json:"list"`
}

type GrantPlatformAnchorVisibilityReq struct {
	g.Meta    `path:"/grantPlatformAnchorVisibility" method:"post" summary:"授权CMS用户平台主播可见性" tags:"账号"`
	AnchorId  uint64 `json:"anchorId" v:"required#平台主播ID不能为空" dc:"平台主播ID"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

type GrantPlatformAnchorVisibilityRes struct {
	Success bool `json:"success"`
}

type BatchGrantPlatformAnchorVisibilityReq struct {
	g.Meta    `path:"/batchGrantPlatformAnchorVisibility" method:"post" summary:"批量授权CMS用户平台主播可见性" tags:"账号"`
	CmsUserId uint64   `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
	AnchorIds []uint64 `json:"anchorIds" v:"required#平台主播ID列表不能为空" dc:"平台主播ID列表"`
}

type BatchGrantPlatformAnchorVisibilityRes struct {
	Success      bool `json:"success"`
	GrantedCount int  `json:"grantedCount" dc:"新授权条数"`
}

type RevokePlatformAnchorVisibilityReq struct {
	g.Meta    `path:"/revokePlatformAnchorVisibility" method:"post" summary:"撤销CMS用户平台主播可见性" tags:"账号"`
	AnchorId  uint64 `json:"anchorId" v:"required#平台主播ID不能为空" dc:"平台主播ID"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

type RevokePlatformAnchorVisibilityRes struct {
	Success bool `json:"success"`
}

type BatchRevokePlatformAnchorVisibilityReq struct {
	g.Meta    `path:"/batchRevokePlatformAnchorVisibility" method:"post" summary:"批量撤销CMS用户平台主播可见性" tags:"账号"`
	CmsUserId uint64   `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
	AnchorIds []uint64 `json:"anchorIds" v:"required#平台主播ID列表不能为空" dc:"平台主播ID列表"`
}

type BatchRevokePlatformAnchorVisibilityRes struct {
	Success      bool `json:"success"`
	RevokedCount int  `json:"revokedCount" dc:"撤销条数"`
}
