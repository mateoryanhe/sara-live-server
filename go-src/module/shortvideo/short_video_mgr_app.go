package shortvideo

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"sort"
	"strconv"
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

func PublishShortVideoAppFromRequest(ctx context.Context, r *ghttp.Request) (*shortvideodto.AppPublishShortVideoRes, error) {
	input, err := parseAppPublishShortVideoMultipart(ctx, r)
	if err != nil {
		return nil, err
	}
	return publishShortVideoApp(ctx, input)
}

func publishShortVideoApp(ctx context.Context, input *appPublishShortVideoInput) (*shortvideodto.AppPublishShortVideoRes, error) {
	if input == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := validateShortVideoAuthorId(authorId); err != nil {
		return nil, err
	}
	user := userinfodao.GetUserInfoByUserId(authorId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	dailyLimit := getAuthorDailyUploadLimit(user.UserType)
	if uint32(shortvideodao.CountAuthorPublishedToday(authorId)) >= dailyLimit {
		return nil, errercode.CreateCode(errercode.ShortVideoDailyUploadLimit)
	}
	if existing := shortvideodao.GetByTitle(input.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, payDiamond, err := normalizeShortVideoPaid(input.IsPaid, input.PayDiamond)
	if err != nil {
		return nil, err
	}
	freeWatchSeconds := normalizeShortVideoFreeWatchSeconds(isPaid, input.FreeWatchSeconds)
	if err := validateShortVideoCategoryId(input.CategoryId); err != nil {
		return nil, err
	}
	if err := validateShortVideoDuration(input.Duration); err != nil {
		return nil, err
	}
	row := entity.NewShortVideo(
		snowflake.GetId(),
		input.Title,
		input.VideoName,
		input.CoverName,
		0,
		isPaid,
		payDiamond,
		input.CategoryId,
		input.Source,
		authorId,
		entity.ShortVideoAuthorTypeApp,
		input.Duration,
		freeWatchSeconds,
	)
	shortvideodao.AddShortVideoToCache(row)
	loadAppShortVideoListCache()
	res := &shortvideodto.AppPublishShortVideoRes{
		ID:    strconv.FormatUint(row.ID, 10),
		Video: upload.GetUrlByName(input.VideoName),
	}
	if input.CoverName != "" {
		res.Cover = upload.GetUrlByName(input.CoverName)
	}
	return res, nil
}

func GetAppShortVideoUploadRecordList(ctx context.Context, req *shortvideodto.AppShortVideoUploadRecordListReq) (*shortvideodto.AppShortVideoUploadRecordListRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	return paginateAppShortVideoUploadRecordList(shortvideodao.GetAuthorShortVideos(authorId), req.Page, req.PageSize, compareShortVideoByCreatedAt), nil
}

// GetAppShortVideoPendingReviewList App端分页查询本人发布的全部短视频(审核中优先,同状态按创建时间降序)
func GetAppShortVideoPendingReviewList(ctx context.Context, req *shortvideodto.AppShortVideoPendingReviewListReq) (*shortvideodto.AppShortVideoUploadRecordListRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	return paginateAppShortVideoUploadRecordList(shortvideodao.GetAuthorShortVideos(authorId), req.Page, req.PageSize, compareShortVideoByStatusThenCreatedAt), nil
}

// DeleteShortVideoApp App端删除本人未审核通过的短视频(status=0)
func DeleteShortVideoApp(ctx context.Context, req *shortvideodto.AppDeleteShortVideoReq) (*shortvideodto.AppDeleteShortVideoRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.AuthorId != authorId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if row.Status == entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoOnShelfCannotDelete)
	}
	if err := removeShortVideoRow(row); err != nil {
		return nil, err
	}
	return &shortvideodto.AppDeleteShortVideoRes{Success: true}, nil
}

func paginateAppShortVideoUploadRecordList(rows []*entity.ShortVideo, page, pageSize int, less func(a, b *entity.ShortVideo) bool) *shortvideodto.AppShortVideoUploadRecordListRes {
	page, pageSize = normalizeAppListPage(page, pageSize)
	sort.Slice(rows, func(i, j int) bool {
		return less(rows[i], rows[j])
	})
	total := len(rows)
	start, end := appListPageRange(total, page, pageSize)
	list := make([]*shortvideodto.AppShortVideoUploadRecordItem, 0, end-start)
	for _, row := range rows[start:end] {
		list = append(list, toAppShortVideoUploadRecordItem(row))
	}
	return &shortvideodto.AppShortVideoUploadRecordListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}
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

// compareShortVideoByStatusThenCreatedAt 审核中(status=0)优先,同状态按创建时间降序
func compareShortVideoByStatusThenCreatedAt(a, b *entity.ShortVideo) bool {
	if a == nil || b == nil {
		return a != nil
	}
	if a.Status != b.Status {
		return a.Status < b.Status
	}
	return compareShortVideoByCreatedAt(a, b)
}

func toAppShortVideoUploadRecordItem(row *entity.ShortVideo) *shortvideodto.AppShortVideoUploadRecordItem {
	var likeCount, viewCount, watchCount uint64
	var totalDiamondIncome float64
	if stat := shortvideodao.GetStatByVideoId(row.ID); stat != nil {
		likeCount = stat.LikeCount
		viewCount = stat.ViewCount
		watchCount = stat.WatchCount
		totalDiamondIncome = stat.TotalDiamondIncome
	}
	author := getShortVideoAuthorInfo(row.AuthorId)
	return &shortvideodto.AppShortVideoUploadRecordItem{
		ID:                 strconv.FormatUint(row.ID, 10),
		Title:              row.Title,
		Video:              upload.GetUrlByName(row.Video),
		Cover:              upload.GetUrlByName(row.Cover),
		Status:             row.Status,
		CategoryId:         row.CategoryId,
		Source:             row.Source,
		AuthorId:           author.AuthorId,
		AuthorNickname:     author.AuthorNickname,
		AuthorAvatar:       author.AuthorAvatar,
		Duration:           row.Duration,
		LikeCount:          likeCount,
		ViewCount:          viewCount,
		WatchCount:         watchCount,
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
	if err := ensureShortVideoUploadDiskSpace(); err != nil {
		return "", err
	}
	name, err := upload.UploadShortVideoFile(file)
	if err != nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return name, nil
}

func uploadShortVideoCoverFile(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", nil
	}
	name, err := upload.UploadImageForApp(ctx, file)
	if err != nil {
		return "", err
	}
	return name, nil
}
