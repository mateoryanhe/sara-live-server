package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	liveentity "xr-game-server/entity/live"
)

// GetAnchorTransferInfo 按主播用户 ID 直查平台主播收款信息。
func GetAnchorTransferInfo(anchorId uint64) *liveentity.LiveAnchorTransferInfo {
	if anchorId == 0 {
		return nil
	}
	var row liveentity.LiveAnchorTransferInfo
	err := g.DB().Model(string(liveentity.TbLiveAnchorTransferInfo)).
		WherePri(anchorId).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetAnchorTransferInfoMapByIds 批量读取平台主播收款信息。
func GetAnchorTransferInfoMapByIds(anchorIds []uint64) map[uint64]*liveentity.LiveAnchorTransferInfo {
	ret := make(map[uint64]*liveentity.LiveAnchorTransferInfo)
	if len(anchorIds) == 0 {
		return ret
	}
	rows := make([]*liveentity.LiveAnchorTransferInfo, 0, len(anchorIds))
	err := g.DB().Model(string(liveentity.TbLiveAnchorTransferInfo)).
		Where("id IN (?)", anchorIds).
		Scan(&rows)
	if err != nil {
		return ret
	}
	for _, row := range rows {
		if row != nil && row.ID > 0 {
			ret[row.ID] = row
		}
	}
	return ret
}

// SaveAnchorTransferInfo 直写数据库保存平台主播收款信息。
func SaveAnchorTransferInfo(row *liveentity.LiveAnchorTransferInfo) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	now := time.Now()
	row.UpdatedAt = now
	existing := GetAnchorTransferInfo(row.ID)
	if existing == nil {
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		_, err := g.DB().Model(string(liveentity.TbLiveAnchorTransferInfo)).Data(row).Insert()
		return err
	}
	row.CreatedAt = existing.CreatedAt
	_, err := g.DB().Model(string(liveentity.TbLiveAnchorTransferInfo)).
		WherePri(row.ID).
		Data(row).
		Update()
	return err
}
