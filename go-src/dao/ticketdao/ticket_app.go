package ticketdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

// GetAll 返回全部门票配置，供进程缓存初始化使用。
func GetAll() []*entity.LiveTicket {
	ret := make([]*entity.LiveTicket, 0)
	_ = g.DB().Model(string(entity.TbLiveTicket)).
		Order("sort desc, created_at desc").
		Scan(&ret)
	return ret
}
