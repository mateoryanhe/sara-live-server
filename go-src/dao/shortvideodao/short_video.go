package shortvideodao

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
)

var shortVideoCacheMgr = gmap.NewKVMap[uint64, *entity.ShortVideo](false)

func initShortVideoDao() {
	all := make([]*entity.ShortVideo, 0)
	g.Model(string(entity.TbShortVideo)).Scan(&all)
	for _, v := range all {
		shortVideoCacheMgr.Set(v.ID, v)
	}
}

// GetShortVideoById 根据视频ID获取短视频(走缓存)
func GetShortVideoById(videoId uint64) *entity.ShortVideo {
	if videoId == 0 || shortVideoCacheMgr == nil {
		return nil
	}
	v := shortVideoCacheMgr.Get(videoId)

	return v
}

func FlushShortVideo(data *entity.ShortVideo) {
	shortVideoCacheMgr.Set(data.ID, data)
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
	FlushShortVideo(row)
	return err
}

func Update(row *entity.ShortVideo) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).Save(row)
	FlushShortVideo(row)
	return err
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).WherePri(id).Delete()
	shortVideoCacheMgr.Remove(id)
	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

// GetOnShelfShortVideos 获取全部已上架短视频(按点赞数降序,再按ID降序)
func GetOnShelfShortVideos() []*entity.ShortVideo {
	ret := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video.Status == entity.ShortVideoStatusOnShelf {
			ret = append(ret, video)
		}
	}
	return ret
}

func GetShortVideoList(req *shortvideodto.ShortVideoListReq) (int, []*shortvideodto.ShortVideoListRes) {
	sql := `select v.id, v.title, v.video, v.cover, v.sort, v.status, v.is_paid, v.diamond_per_second, v.description,
                   coalesce(s.like_count, 0) as like_count, v.created_at, v.updated_at
            from short_videos v
            left join short_video_stats s on s.id = v.id
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*shortvideodto.ShortVideoListRes, 0)

	if req.Title != "" {
		sql += ` and v.title LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Title))
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and v.status = ?`
		param = append(param, entity.ShortVideoStatusOffShelf)
	case 2:
		sql += ` and v.status = ?`
		param = append(param, entity.ShortVideoStatusOnShelf)
	}

	sql += ` order by v.sort desc, v.created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
