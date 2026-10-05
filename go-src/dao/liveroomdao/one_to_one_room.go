package liveroomdao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

var oneToOneRoomCache = gmap.NewKVMap[uint64, *entity.OneToOneRoom](false)

func initOneToOneRoomDao() {
	rows := make([]*entity.OneToOneRoom, 0)
	_ = g.Model(string(entity.TbOneToOneRoom)).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			oneToOneRoomCache.Set(row.ID, row)
		}
	}
}

func GetOneToOneRoom(anchorId uint64) *entity.OneToOneRoom {
	if anchorId == 0 || !oneToOneRoomCache.Contains(anchorId) {
		return nil
	}
	return oneToOneRoomCache.Get(anchorId)
}

func IsOneToOneOnShelf(anchorId uint64) bool {
	row := GetOneToOneRoom(anchorId)
	return row != nil && row.IsOnShelf()
}

func GetAllOneToOneRooms() []*entity.OneToOneRoom {
	return oneToOneRoomCache.Values()
}

func AddOneToOneRoomToCache(row *entity.OneToOneRoom) {
	if row == nil || row.ID == 0 {
		return
	}
	oneToOneRoomCache.Set(row.ID, row)
}

func CreateOneToOneRoom(anchorId uint64, billing float64) *entity.OneToOneRoom {
	if anchorId == 0 {
		return nil
	}
	if existing := GetOneToOneRoom(anchorId); existing != nil {
		return existing
	}
	row := entity.NewOneToOneRoom(anchorId)
	row.SetBilling(billing)
	oneToOneRoomCache.Set(anchorId, row)
	return row
}
