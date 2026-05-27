package accountdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// ListValidAppTokens 查询未过期的App Token(expire_at > 当前服务器时间)
func ListValidAppTokens() []*entity.AppToken {
	list := make([]*entity.AppToken, 0)
	_ = g.Model(string(entity.TbAppToken)).
		Where("expire_at > ?", time.Now()).
		Order("id desc").
		Scan(&list)
	return list
}
