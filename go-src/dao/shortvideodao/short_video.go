package shortvideodao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
)

var shortVideoCacheMgr *cache.CacheMgr

func InitShortVideoDao() {
	shortVideoCacheMgr = cache.NewCacheMgr()
}

// GetShortVideoById 根据视频ID获取短视频(走缓存)
func GetShortVideoById(videoId uint64) *entity.ShortVideo {
	if videoId == 0 || shortVideoCacheMgr == nil {
		return nil
	}
	v := shortVideoCacheMgr.GetData(videoId, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.ShortVideo
		err = g.Model(string(entity.TbShortVideo)).Where(g.Map{
			string(db.IdName): videoId,
		}).Scan(&row)
		if row == nil || row.ID == 0 {
			return nil, nil
		}
		return row, err
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.ShortVideo)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

func GetById(id uint64) *entity.ShortVideo {
	return GetShortVideoById(id)
}

func GetByTitle(title string) *entity.ShortVideo {
	var row entity.ShortVideo
	if err := g.DB().Model(string(entity.TbShortVideo)).Where("title = ?", title).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.ShortVideo) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).Save(row)
	return err
}

func Update(row *entity.ShortVideo) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).Save(row)

	return err
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).WherePri(id).Delete()

	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetShortVideoList(req *shortvideodto.ShortVideoListReq) (int, []*shortvideodto.ShortVideoListRes) {
	sql := `select id, title, video, cover, sort, status, description, like_count, created_at, updated_at
            from short_videos
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*shortvideodto.ShortVideoListRes, 0)

	if req.Title != "" {
		sql += ` and title LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Title))
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.ShortVideoStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.ShortVideoStatusOnShelf)
	}

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
