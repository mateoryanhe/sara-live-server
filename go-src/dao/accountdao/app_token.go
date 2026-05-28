package accountdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
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

// GetAppTokenByUserId 根据用户ID查询App Token
func GetAppTokenByUserId(userId uint64) *entity.AppToken {
	ret := &entity.AppToken{}
	err := g.Model(string(entity.TbAppToken)).Where(db.IdName, userId).Scan(ret)
	if err != nil || ret.ID == 0 {
		return nil
	}
	return ret
}
