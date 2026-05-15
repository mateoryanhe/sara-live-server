package rankdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/constants/common"
	"xr-game-server/entity"
)

func GetRankDataBy(typeId uint32, roleId uint64) *entity.PlayerRank {
	rets := make([]*entity.PlayerRank, common.Zero)
	g.Model(string(entity.TbPlayerRank)).Unscoped().Where(g.Map{
		string(entity.PlayerRankRoleId): roleId,
		string(entity.PlayerRankTypeId): typeId,
	}).Scan(&rets)
	if len(rets) == common.Zero {
		return nil
	}
	return rets[common.Zero]
}

// GetServerRankData 根据服务器id获取
func GetServerRankData(serverId uint64, typeId uint32, asc bool, max int) []*entity.PlayerRank {
	rets := make([]*entity.PlayerRank, common.Zero)
	sql := "select pr.* from player_ranks pr join player_roles r on r.id=pr.role_id where pr.val>0 and  r.server_id=? and type_id=? "
	if asc {
		sql += " order by pr.val asc limit 0,?"
	} else {
		sql += " order by pr.val desc limit 0,?"
	}
	lst, _ := g.DB().GetAll(gctx.New(), sql, serverId, typeId, max)
	if lst != nil {
		ret := new(entity.PlayerRank)
		gconv.Struct(lst[common.Zero], ret)
		rets = append(rets, ret)
	}
	return rets
}

func GetZoneData(zoneId uint64, typeId uint32, asc bool, max int) []*entity.PlayerRank {
	rets := make([]*entity.PlayerRank, common.Zero)
	sql := "select t.* from " +
		"(select pr.*, r.server_id  from player_ranks pr  join player_roles r on r.id = pr.role_id) as t" +
		"  join " +
		" (select s.id ,z.id zone_id  from servers s left join zones z on z.id = s.zone_id) as k   " +
		"  on k.id = t.server_id " +
		"where  t.val>0 and  k.zone_id = ? and t.type_id=? "
	if asc {
		sql += " order by t.val asc limit 0,?"
	} else {
		sql += " order by t.val desc limit 0,?"
	}
	lst, _ := g.DB().GetAll(gctx.New(), sql, zoneId, typeId, max)
	if lst != nil {
		for _, record := range lst {
			ret := new(entity.PlayerRank)
			gconv.Struct(record, ret)
			rets = append(rets, ret)
		}
	}
	return rets
}

// UpdateRankByTypeId 重置数值
func UpdateRankByTypeId(typeId uint32) {
	g.Model(string(entity.TbPlayerRank)).Unscoped().Data(string(entity.PlayerRankVal), common.Zero).Where(string(entity.PlayerRankTypeId), typeId).Update()
}
