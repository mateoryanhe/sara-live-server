package guilddao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	liveentity "xr-game-server/entity/live"
)

// GetGuildTransferInfo 按工会ID直查收款信息(主键=工会ID)
func GetGuildTransferInfo(guildId uint64) *liveentity.LiveGuildTransferInfo {
	if guildId == 0 {
		return nil
	}
	var row liveentity.LiveGuildTransferInfo
	err := g.DB().Model(string(liveentity.TbLiveGuildTransferInfo)).
		WherePri(guildId).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// SaveGuildTransferInfo 直写数据库保存工会收款信息(存在则更新)
func SaveGuildTransferInfo(row *liveentity.LiveGuildTransferInfo) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	now := time.Now()
	row.UpdatedAt = now
	existing := GetGuildTransferInfo(row.ID)
	if existing == nil {
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		_, err := g.DB().Model(string(liveentity.TbLiveGuildTransferInfo)).Data(row).Insert()
		return err
	}
	row.CreatedAt = existing.CreatedAt
	_, err := g.DB().Model(string(liveentity.TbLiveGuildTransferInfo)).
		WherePri(row.ID).
		Data(row).
		Update()
	return err
}
