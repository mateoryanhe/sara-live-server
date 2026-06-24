package shortvideodao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"sort"
	"strconv"
	"strings"
	"time"
	"xr-game-server/core/xrtime"
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
	return shortVideoCacheMgr.Get(videoId)
}

func FlushShortVideo(data *entity.ShortVideo) {
	if data == nil {
		return
	}
	shortVideoCacheMgr.Set(data.ID, data)
}

func GetById(id uint64) *entity.ShortVideo {
	return GetShortVideoById(id)
}

func GetByTitle(title string) *entity.ShortVideo {
	if title == "" {
		return nil
	}
	for _, row := range shortVideoCacheMgr.Values() {
		if row != nil && row.Title == title {
			return row
		}
	}
	return nil
}

// HasAuthorPublishedToday 作者当天是否已发布短视频(走内存缓存,按创建时间判断)
func HasAuthorPublishedToday(authorId uint64) bool {
	if authorId == 0 {
		return false
	}
	now := time.Now()
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil || video.AuthorId != authorId {
			continue
		}
		if xrtime.IsSameDay(video.CreatedAt, now) {
			return true
		}
	}
	return false
}

// GetAuthorShortVideos 查询指定作者的全部短视频(走内存缓存,不排序不分页)
func GetAuthorShortVideos(authorId uint64) []*entity.ShortVideo {
	if authorId == 0 {
		return nil
	}
	filtered := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil || video.AuthorId != authorId {
			continue
		}
		filtered = append(filtered, video)
	}
	return filtered
}

func AddShortVideoToCache(row *entity.ShortVideo) {
	FlushShortVideo(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbShortVideo)).WherePri(id).Delete()
	shortVideoCacheMgr.Remove(id)
	return err
}

func UpdateStatus(id uint64, status uint8) {
	row := GetShortVideoById(id)
	if row == nil {
		return
	}
	row.SetStatus(status)
	FlushShortVideo(row)
}

// GetOnShelfShortVideos 获取全部已上架短视频
func GetOnShelfShortVideos() []*entity.ShortVideo {
	ret := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video != nil && video.Status == entity.ShortVideoStatusOnShelf {
			ret = append(ret, video)
		}
	}
	return ret
}

func GetShortVideoList(req *shortvideodto.ShortVideoListReq) (int, []*shortvideodto.ShortVideoListRes) {
	titleKeyword := strings.ToLower(strings.TrimSpace(req.Title))
	filtered := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil {
			continue
		}
		if titleKeyword != "" && !strings.Contains(strings.ToLower(video.Title), titleKeyword) {
			continue
		}
		switch req.StatusFilter {
		case 1:
			if video.Status != entity.ShortVideoStatusOffShelf {
				continue
			}
		case 2:
			if video.Status != entity.ShortVideoStatusOnShelf {
				continue
			}
		}
		filtered = append(filtered, video)
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Sort != filtered[j].Sort {
			return filtered[i].Sort > filtered[j].Sort
		}
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	pageIndex, pageSize := normalizeShortVideoListPage(req.PageIndex, req.PageSize)
	start, end := shortVideoListPageRange(total, pageIndex, pageSize)

	ret := make([]*shortvideodto.ShortVideoListRes, 0, end-start)
	for _, video := range filtered[start:end] {
		ret = append(ret, toShortVideoListRes(video))
	}
	return total, ret
}

func toShortVideoListRes(video *entity.ShortVideo) *shortvideodto.ShortVideoListRes {
	var likeCount uint64
	var totalDiamondIncome float64
	if stat := GetStatByVideoId(video.ID); stat != nil {
		likeCount = stat.LikeCount
		totalDiamondIncome = stat.TotalDiamondIncome
	}
	return &shortvideodto.ShortVideoListRes{
		ID:                 strconv.FormatUint(video.ID, 10),
		Title:              video.Title,
		Video:              video.Video,
		Cover:              video.Cover,
		Sort:               video.Sort,
		Status:             video.Status,
		IsPaid:             video.IsPaid,
		DiamondPerMinute:   video.DiamondPerMinute,
		CategoryId:         video.CategoryId,
		Source:             video.Source,
		AuthorId:           strconv.FormatUint(video.AuthorId, 10),
		LikeCount:          likeCount,
		TotalDiamondIncome: totalDiamondIncome,
		Duration:           video.Duration,
		FreeWatchSeconds:   video.FreeWatchSeconds,
		CreatedAt:          formatShortVideoTime(video.CreatedAt),
		UpdatedAt:          formatShortVideoTime(video.UpdatedAt),
	}
}

func normalizeShortVideoListPage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return pageIndex, pageSize
}

func shortVideoListPageRange(total, pageIndex, pageSize int) (int, int) {
	start := (pageIndex - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}

func formatShortVideoTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
