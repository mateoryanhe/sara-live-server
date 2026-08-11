package guilddao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// SaveGuildAnchor 保存工会主播关联(新增或更新)
func SaveGuildAnchor(row *entity.LiveGuildAnchor) error {
	if row == nil {
		return nil
	}
	now := time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	_, err := g.DB().Model(string(entity.TbLiveGuildAnchor)).Save(row)
	return err
}

// GetGuildAnchorById 按主键查询
func GetGuildAnchorById(id uint64) *entity.LiveGuildAnchor {
	if id == 0 {
		return nil
	}
	var row entity.LiveGuildAnchor
	if err := g.DB().Model(string(entity.TbLiveGuildAnchor)).Where("id = ?", id).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetGuildAnchorByGuildAndRoom 按工会与直播间查询
func GetGuildAnchorByGuildAndRoom(guildId, roomId uint64) *entity.LiveGuildAnchor {
	if guildId == 0 || roomId == 0 {
		return nil
	}
	var row entity.LiveGuildAnchor
	err := g.DB().Model(string(entity.TbLiveGuildAnchor)).
		Where(string(entity.GuildAnchorGuildId), guildId).
		Where(string(entity.GuildAnchorRoomId), roomId).
		Limit(1).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListGuildAnchorsByGuildId 查询工会名下全部主播(直播间)
func ListGuildAnchorsByGuildId(guildId uint64) []*entity.LiveGuildAnchor {
	if guildId == 0 {
		return nil
	}
	rows := make([]*entity.LiveGuildAnchor, 0)
	_ = g.DB().Model(string(entity.TbLiveGuildAnchor)).
		Where(string(entity.GuildAnchorGuildId), guildId).
		OrderDesc("created_at").
		Scan(&rows)
	return rows
}

// ListGuildAnchorRoomIdsByGuildId 查询工会名下直播间 ID 列表
func ListGuildAnchorRoomIdsByGuildId(guildId uint64) []uint64 {
	rows := ListGuildAnchorsByGuildId(guildId)
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.RoomId == 0 {
			continue
		}
		ids = append(ids, row.RoomId)
	}
	return ids
}

// DeleteGuildAnchorById 按主键删除
func DeleteGuildAnchorById(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveGuildAnchor)).WherePri(id).Delete()
	return err
}

// DeleteGuildAnchorByGuildAndRoom 按工会与直播间删除
func DeleteGuildAnchorByGuildAndRoom(guildId, roomId uint64) error {
	if guildId == 0 || roomId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveGuildAnchor)).
		Where(string(entity.GuildAnchorGuildId), guildId).
		Where(string(entity.GuildAnchorRoomId), roomId).
		Delete()
	return err
}

// DeleteGuildAnchorsByGuildId 删除工会名下全部主播关联
func DeleteGuildAnchorsByGuildId(guildId uint64) error {
	if guildId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveGuildAnchor)).
		Where(string(entity.GuildAnchorGuildId), guildId).
		Delete()
	return err
}
