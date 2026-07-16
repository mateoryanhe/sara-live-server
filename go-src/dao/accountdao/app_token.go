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

// QueryAppTokens CMS分页查询App Token数据库记录
func QueryAppTokens(userId uint64, pageIndex, pageSize int) (int, []*entity.AppToken) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	list := make([]*entity.AppToken, 0)
	m := g.Model(string(entity.TbAppToken))
	if userId > 0 {
		m = m.Where(db.IdName, userId)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order(db.IdName + " desc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
