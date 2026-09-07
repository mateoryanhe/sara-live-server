package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GuildVisibilityListReq 查询工会可见的 CMS 用户列表
type GuildVisibilityListReq struct {
	g.Meta  `path:"/guildVisibilityList" method:"post" summary:"查询工会可见性授权列表" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// GuildVisibilityItem 工会可见性条目
type GuildVisibilityItem struct {
	Id          string `json:"id" dc:"可见性记录ID"`
	GuildId     string `json:"guildId" dc:"工会ID"`
	GuildName   string `json:"guildName" dc:"工会名称"`
	CmsUserId   string `json:"cmsUserId" dc:"CMS用户ID"`
	CmsUserName string `json:"cmsUserName" dc:"CMS用户名"`
	CreatedAt   string `json:"createdAt" dc:"授权时间"`
}

// GuildVisibilityListRes 工会可见性列表
type GuildVisibilityListRes struct {
	List []*GuildVisibilityItem `json:"list"`
}

// GuildVisibilityByUserListReq 查询指定 CMS 用户已授权的工会
type GuildVisibilityByUserListReq struct {
	g.Meta    `path:"/guildVisibilityByUserList" method:"post" summary:"查询CMS用户已授权工会列表" tags:"直播工会"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

// GuildVisibilityByUserListRes 指定用户已授权工会
type GuildVisibilityByUserListRes struct {
	List []*GuildVisibilityItem `json:"list"`
}

// GrantGuildVisibilityReq 授权 CMS 用户可见某工会
type GrantGuildVisibilityReq struct {
	g.Meta    `path:"/grantGuildVisibility" method:"post" summary:"授权CMS用户工会可见性" tags:"直播工会"`
	GuildId   uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

// GrantGuildVisibilityRes 授权结果
type GrantGuildVisibilityRes struct {
	Success bool `json:"success"`
}

// BatchGrantGuildVisibilityReq 批量授权 CMS 用户可见多个工会
type BatchGrantGuildVisibilityReq struct {
	g.Meta    `path:"/batchGrantGuildVisibility" method:"post" summary:"批量授权CMS用户工会可见性" tags:"直播工会"`
	CmsUserId uint64   `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
	GuildIds  []uint64 `json:"guildIds" v:"required#工会ID列表不能为空" dc:"工会ID列表"`
}

// BatchGrantGuildVisibilityRes 批量授权结果
type BatchGrantGuildVisibilityRes struct {
	Success      bool `json:"success"`
	GrantedCount int  `json:"grantedCount" dc:"新授权条数"`
}

// RevokeGuildVisibilityReq 撤销 CMS 用户对某工会的可见性
type RevokeGuildVisibilityReq struct {
	g.Meta    `path:"/revokeGuildVisibility" method:"post" summary:"撤销CMS用户工会可见性" tags:"直播工会"`
	GuildId   uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	CmsUserId uint64 `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

// RevokeGuildVisibilityRes 撤销结果
type RevokeGuildVisibilityRes struct {
	Success bool `json:"success"`
}

// BatchRevokeGuildVisibilityReq 批量撤销 CMS 用户对多个工会的可见性
type BatchRevokeGuildVisibilityReq struct {
	g.Meta    `path:"/batchRevokeGuildVisibility" method:"post" summary:"批量撤销CMS用户工会可见性" tags:"直播工会"`
	CmsUserId uint64   `json:"cmsUserId" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
	GuildIds  []uint64 `json:"guildIds" v:"required#工会ID列表不能为空" dc:"工会ID列表"`
}

// BatchRevokeGuildVisibilityRes 批量撤销结果
type BatchRevokeGuildVisibilityRes struct {
	Success      bool `json:"success"`
	RevokedCount int  `json:"revokedCount" dc:"撤销条数"`
}
