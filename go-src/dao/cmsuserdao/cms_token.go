package cmsuserdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// ListValidCmsTokens 查询未过期的CMS Token(expire_at > 当前服务器时间)
func ListValidCmsTokens() []*entity.CmsToken {
	list := make([]*entity.CmsToken, 0)
	_ = g.Model(string(entity.TbCmsToken)).
		Where("expire_at > ?", time.Now()).
		Order("id desc").
		Scan(&list)
	return list
}
