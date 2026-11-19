package anchorphotodao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	entity "xr-game-server/entity/live"
)

const photoListFromDBMax = 2000

// PreloadAnchorPhotoToCache 启动时批量预热全部有效相册照片(按主播分组列表 + 行缓存).
func PreloadAnchorPhotoToCache() (anchorCount int, rowCount int) {
	if rowCacheMgr == nil || listCacheMgr == nil {
		return 0, 0
	}
	rows := loadAllActiveAnchorPhotosFromDB()
	if len(rows) == 0 {
		return 0, 0
	}
	ctx := gctx.New()
	grouped := make(map[uint64][]*entity.AnchorPhoto)
	perAnchor := make(map[uint64]int)
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.AnchorId == 0 {
			continue
		}
		rowCacheMgr.PublishRow(ctx, row.ID, row)
		rowCount++
		n := perAnchor[row.AnchorId]
		if n >= photoListFromDBMax {
			continue
		}
		grouped[row.AnchorId] = append(grouped[row.AnchorId], row)
		perAnchor[row.AnchorId] = n + 1
	}
	for anchorId, list := range grouped {
		listCacheMgr.PublishList(ctx, listCacheKey(anchorId), list)
		anchorCount++
	}
	return anchorCount, rowCount
}

func loadAllActiveAnchorPhotosFromDB() []*entity.AnchorPhoto {
	out := make([]*entity.AnchorPhoto, 0)
	_ = g.Model(string(entity.TbAnchorPhoto)).
		Where("is_deleted = 0").
		Order("anchor_id asc, created_at desc, id desc").
		Scan(&out)
	return out
}