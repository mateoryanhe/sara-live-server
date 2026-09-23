package liveroomdao

import (
	"github.com/gogf/gf/v2/frame/g"
	liveentity "xr-game-server/entity/live"
)

// AddPlatformAnchorVisibility 写入平台主播可见性，重复记录会被忽略。
func AddPlatformAnchorVisibility(anchorId, cmsUserId uint64) error {
	if anchorId == 0 || cmsUserId == 0 {
		return nil
	}
	if HasPlatformAnchorVisibility(anchorId, cmsUserId) {
		return nil
	}
	row := liveentity.NewLivePlatformAnchorVisibility(anchorId, cmsUserId)
	_, err := g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).Save(row)
	return err
}

// HasPlatformAnchorVisibility 判断可见性记录是否存在。
func HasPlatformAnchorVisibility(anchorId, cmsUserId uint64) bool {
	if anchorId == 0 || cmsUserId == 0 {
		return false
	}
	count, err := g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Where(string(liveentity.LivePlatformAnchorVisibilityAnchorId), anchorId).
		Where(string(liveentity.LivePlatformAnchorVisibilityCmsUserId), cmsUserId).
		Count()
	return err == nil && count > 0
}

// ListVisiblePlatformAnchorIds 查询 CMS 用户可见的平台主播 ID。
func ListVisiblePlatformAnchorIds(cmsUserId uint64) []uint64 {
	if cmsUserId == 0 {
		return nil
	}
	type row struct {
		AnchorId uint64 `json:"anchorId"`
	}
	rows := make([]row, 0)
	_ = g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Fields(string(liveentity.LivePlatformAnchorVisibilityAnchorId)).
		Where(string(liveentity.LivePlatformAnchorVisibilityCmsUserId), cmsUserId).
		Scan(&rows)
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, item := range rows {
		if item.AnchorId == 0 {
			continue
		}
		if _, ok := seen[item.AnchorId]; ok {
			continue
		}
		seen[item.AnchorId] = struct{}{}
		ids = append(ids, item.AnchorId)
	}
	return ids
}

// ListPlatformAnchorVisibilitiesByCmsUserId 按 CMS 用户查询可见性记录。
func ListPlatformAnchorVisibilitiesByCmsUserId(cmsUserId uint64) []*liveentity.LivePlatformAnchorVisibility {
	if cmsUserId == 0 {
		return nil
	}
	rows := make([]*liveentity.LivePlatformAnchorVisibility, 0)
	_ = g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Where(string(liveentity.LivePlatformAnchorVisibilityCmsUserId), cmsUserId).
		Order("created_at asc").
		Scan(&rows)
	return rows
}

// ListPlatformAnchorVisibilitiesByAnchorId 按平台主播查询可见性记录。
func ListPlatformAnchorVisibilitiesByAnchorId(anchorId uint64) []*liveentity.LivePlatformAnchorVisibility {
	if anchorId == 0 {
		return nil
	}
	rows := make([]*liveentity.LivePlatformAnchorVisibility, 0)
	_ = g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Where(string(liveentity.LivePlatformAnchorVisibilityAnchorId), anchorId).
		Order("created_at asc").
		Scan(&rows)
	return rows
}

// RemovePlatformAnchorVisibility 删除单条平台主播可见性。
func RemovePlatformAnchorVisibility(anchorId, cmsUserId uint64) error {
	if anchorId == 0 || cmsUserId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Where(string(liveentity.LivePlatformAnchorVisibilityAnchorId), anchorId).
		Where(string(liveentity.LivePlatformAnchorVisibilityCmsUserId), cmsUserId).
		Delete()
	return err
}

// DeletePlatformAnchorVisibilities 删除某平台主播的全部可见性记录。
func DeletePlatformAnchorVisibilities(anchorId uint64) error {
	if anchorId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(liveentity.TbLivePlatformAnchorVisibility)).
		Where(string(liveentity.LivePlatformAnchorVisibilityAnchorId), anchorId).
		Delete()
	return err
}
