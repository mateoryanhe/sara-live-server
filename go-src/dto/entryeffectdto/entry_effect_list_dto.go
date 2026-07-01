package entryeffectdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type EntryEffectListReq struct {
	g.Meta `path:"/entryEffectList" method:"post" summary:"获取进场特效列表" tags:"进场特效"`
	httpserver.CMSQueryReq
	Name         string `json:"name" dc:"名称(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
}

type EntryEffectListRes struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	LevelStart    int    `json:"levelStart"`
	LevelEnd      int    `json:"levelEnd"`
	Animation     string `json:"animation" dc:"动画完整URL(列表展示)"`
	AnimationName string `json:"animationName" dc:"动画资源文件名(编辑保存用)"`
	Status        uint8  `json:"status"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
