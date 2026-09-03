package shortvideodao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"sort"
	"strconv"
	"strings"
	"time"
	"xr-game-server/core/xrtime"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
)

var shortVideoCacheMgr = gmap.NewKVMap[uint64, *entity.ShortVideo](false)

func initShortVideoDao() {
	all := make([]*entity.ShortVideo, 0)
	g.Model(string(entity.TbShortVideo)).Scan(&all)
	for _, v := range all {
		shortVideoCacheMgr.Set(v.ID, v)
		GetStatByVideoId(v.ID)
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

// CountAuthorPublishedToday 作者当天已发布短视频数量(走内存缓存,按创建时间判断)
func CountAuthorPublishedToday(authorId uint64) int {
	if authorId == 0 {
		return 0
	}
	now := time.Now()
	count := 0
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil || video.AuthorId != authorId {
			continue
		}
		if xrtime.IsSameDay(video.CreatedAt, now) {
			count++
		}
	}
	return count
}

// HasAuthorPublishedToday 作者当天是否已发布短视频(走内存缓存,按创建时间判断)
func HasAuthorPublishedToday(authorId uint64) bool {
	return CountAuthorPublishedToday(authorId) > 0
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

// GetAuthorOnShelfShortVideos 查询指定作者已上架短视频(走内存缓存,不排序不分页)
func GetAuthorOnShelfShortVideos(authorId uint64) []*entity.ShortVideo {
	if authorId == 0 {
		return nil
	}
	filtered := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil || video.AuthorId != authorId {
			continue
		}
		if video.Status != entity.ShortVideoStatusOnShelf {
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
	authorKeyword := strings.ToLower(strings.TrimSpace(req.AuthorNickname))
	authorIdFilter := parseUint64Filter(req.AuthorId)
	filtered := make([]*entity.ShortVideo, 0)
	for _, video := range shortVideoCacheMgr.Values() {
		if video == nil {
			continue
		}
		if authorIdFilter > 0 && video.AuthorId != authorIdFilter {
			continue
		}
		if titleKeyword != "" && !strings.Contains(strings.ToLower(video.Title), titleKeyword) {
			continue
		}
		if authorKeyword != "" {
			if video.AuthorId == 0 {
				continue
			}
			user := userinfodao.GetUserInfoByUserId(video.AuthorId)
			if user == nil || !strings.Contains(strings.ToLower(user.Nickname), authorKeyword) {
				continue
			}
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

	sortShortVideoList(filtered, req.SortField)

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
	var likeCount, viewCount, watchCount uint64
	var totalDiamondIncome float64
	if stat := GetStatByVideoId(video.ID); stat != nil {
		likeCount = stat.LikeCount
		viewCount = stat.ViewCount
		watchCount = stat.WatchCount
		totalDiamondIncome = stat.TotalDiamondIncome
	}
	authorNickname := ""
	if user := userinfodao.GetUserInfoByUserId(video.AuthorId); user != nil {
		authorNickname = user.Nickname
	}
	return &shortvideodto.ShortVideoListRes{
		ID:                 strconv.FormatUint(video.ID, 10),
		Title:              video.Title,
		Video:              video.Video,
		Cover:              video.Cover,
		Sort:               video.Sort,
		Status:             video.Status,
		IsPaid:             video.IsPaid,
		PayDiamond:         video.PayDiamond,
		CategoryId:         video.CategoryId,
		Source:             video.Source,
		AuthorId:           strconv.FormatUint(video.AuthorId, 10),
		AuthorType:         video.AuthorType,
		AuthorNickname:     authorNickname,
		LikeCount:          likeCount,
		ViewCount:          viewCount,
		WatchCount:         watchCount,
		TotalDiamondIncome: totalDiamondIncome,
		Duration:           video.Duration,
		FreeWatchSeconds:   video.FreeWatchSeconds,
		CreatedAt:          formatShortVideoTime(video.CreatedAt),
		UpdatedAt:          formatShortVideoTime(video.UpdatedAt),
	}
}

// CountAll 短视频总记录数(走内存缓存)
func CountAll() int {
	if shortVideoCacheMgr == nil {
		return 0
	}
	return shortVideoCacheMgr.Size()
}

func sortShortVideoList(list []*entity.ShortVideo, sortField string) {
	switch strings.TrimSpace(sortField) {
	case "viewCount":
		sort.Slice(list, func(i, j int) bool {
			vi, vj := getVideoViewCount(list[i]), getVideoViewCount(list[j])
			if vi != vj {
				return vi < vj
			}
			return list[i].ID < list[j].ID
		})
	case "totalDiamondIncome":
		sort.Slice(list, func(i, j int) bool {
			vi, vj := getVideoTotalDiamondIncome(list[i]), getVideoTotalDiamondIncome(list[j])
			if vi != vj {
				return vi < vj
			}
			return list[i].ID < list[j].ID
		})
	default:
		sort.Slice(list, func(i, j int) bool {
			if list[i].CreatedAt.Equal(list[j].CreatedAt) {
				return list[i].ID > list[j].ID
			}
			return list[i].CreatedAt.After(list[j].CreatedAt)
		})
	}
}

func getVideoViewCount(video *entity.ShortVideo) uint64 {
	if video == nil {
		return 0
	}
	if stat := GetStatByVideoId(video.ID); stat != nil {
		return stat.ViewCount
	}
	return 0
}

func getVideoTotalDiamondIncome(video *entity.ShortVideo) float64 {
	if video == nil {
		return 0
	}
	if stat := GetStatByVideoId(video.ID); stat != nil {
		return stat.TotalDiamondIncome
	}
	return 0
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

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	n, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func formatShortVideoTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
