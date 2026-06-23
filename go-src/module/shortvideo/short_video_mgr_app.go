package shortvideo

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"sort"
	"strconv"
	"strings"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func PublishShortVideoApp(ctx context.Context, req *shortvideodto.AppPublishShortVideoReq) (*shortvideodto.AppPublishShortVideoRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := validateShortVideoAuthorId(authorId); err != nil {
		return nil, err
	}
	if shortvideodao.HasAuthorPublishedToday(authorId) {
		return nil, errercode.CreateCode(errercode.ShortVideoDailyUploadLimit)
	}
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, diamondPerMinute := entity.ShortVideoPaidNo, float64(0)

	if err := validateShortVideoCategoryId(req.CategoryId); err != nil {
		return nil, err
	}
	if err := validateShortVideoDuration(req.Duration); err != nil {
		return nil, err
	}
	videoName, err := uploadShortVideoFile(req.File)
	if err != nil {
		return nil, err
	}
	coverName, err := uploadShortVideoCoverFile(ctx, req.Cover)
	if err != nil {
		upload.DeleteUploadedFile(videoName)
		return nil, err
	}
	row := entity.NewShortVideo(
		snowflake.GetId(),
		req.Title,
		videoName,
		coverName,
		0,
		isPaid,
		diamondPerMinute,
		req.CategoryId,
		req.Source,
		authorId,
		req.Duration,
	)
	shortvideodao.AddShortVideoToCache(row)
	loadAppShortVideoListCache()
	res := &shortvideodto.AppPublishShortVideoRes{
		ID:    strconv.FormatUint(row.ID, 10),
		Video: upload.GetUrlByName(videoName),
	}
	if coverName != "" {
		res.Cover = upload.GetUrlByName(coverName)
	}
	return res, nil
}

func GetAppShortVideoUploadRecordList(ctx context.Context, req *shortvideodto.AppShortVideoUploadRecordListReq) (*shortvideodto.AppShortVideoUploadRecordListRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	rows := shortvideodao.GetAuthorShortVideos(authorId)
	sort.Slice(rows, func(i, j int) bool {
		return compareShortVideoByCreatedAt(rows[i], rows[j])
	})
	total := len(rows)
	start, end := appListPageRange(total, page, pageSize)
	pageRows := rows[start:end]
	list := make([]*shortvideodto.AppShortVideoUploadRecordItem, 0, len(pageRows))
	for _, row := range pageRows {
		list = append(list, toAppShortVideoUploadRecordItem(row))
	}
	return &shortvideodto.AppShortVideoUploadRecordListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// compareShortVideoByCreatedAt 按创建时间降序,相同时间按 ID 降序
func compareShortVideoByCreatedAt(a, b *entity.ShortVideo) bool {
	if a == nil || b == nil {
		return a != nil
	}
	if a.CreatedAt.Equal(b.CreatedAt) {
		return a.ID > b.ID
	}
	return a.CreatedAt.After(b.CreatedAt)
}

func toAppShortVideoUploadRecordItem(row *entity.ShortVideo) *shortvideodto.AppShortVideoUploadRecordItem {
	var likeCount, viewCount uint64
	var totalDiamondIncome float64
	if stat := shortvideodao.GetStatByVideoId(row.ID); stat != nil {
		likeCount = stat.LikeCount
		viewCount = stat.ViewCount
		totalDiamondIncome = stat.TotalDiamondIncome
	}
	return &shortvideodto.AppShortVideoUploadRecordItem{
		ID:                 strconv.FormatUint(row.ID, 10),
		Title:              row.Title,
		Video:              upload.GetUrlByName(row.Video),
		Cover:              upload.GetUrlByName(row.Cover),
		Status:             row.Status,
		CategoryId:         row.CategoryId,
		Source:             row.Source,
		Duration:           row.Duration,
		LikeCount:          likeCount,
		ViewCount:          viewCount,
		TotalDiamondIncome: totalDiamondIncome,
		CreatedAt:          formatShortVideoUploadTime(row.CreatedAt),
		UpdatedAt:          formatShortVideoUploadTime(row.UpdatedAt),
	}
}

func formatShortVideoUploadTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func validateShortVideoAuthorId(authorId uint64) error {
	if authorId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(authorId) == nil {
		return errercode.CreateCode(errercode.SysError)
	}
	return nil
}

func validateShortVideoDuration(duration uint32) error {
	if duration == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	maxDuration := getShortVideoMaxDuration()
	if maxDuration > 0 && duration > maxDuration {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func uploadShortVideoFile(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	maxSize := getShortVideoMaxFileSize()
	if file.Size > int64(maxSize) {
		return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	name, err := upload.UploadShortVideoFile(file, maxSize)
	if err != nil {
		if strings.Contains(err.Error(), "file too large") {
			return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
		}
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return name, nil
}

func uploadShortVideoCoverFile(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", nil
	}
	maxSize := int64(getShortVideoMaxCoverFileSize()) * 1024 * 1024
	if file.Size > maxSize {
		return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	name, err := upload.UploadImageForApp(ctx, file)
	if err != nil {
		return "", err
	}
	return name, nil
}
