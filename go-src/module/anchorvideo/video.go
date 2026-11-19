package anchorvideo

import (
	"context"
	"strings"
	"unicode/utf8"

	"xr-game-server/constants/common"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorvideodao"
	"xr-game-server/dao/anchorvideopurchasedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorvideodto"
	entity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	maxAnchorVideoDescriptionRunes = 200
	defaultVideoListPageSize       = 20
	maxVideoListPageSize           = 100
)

// GetVideoList 读主播相册列表缓存(最多1000)再分页;新增/删除会刷新缓存;paid:价格0/主播本人/已购为true.
func GetVideoList(ctx context.Context, req *anchorvideodto.VideoListReq) (*anchorvideodto.VideoListRes, error) {
	if req == nil || req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = defaultVideoListPageSize
	}
	if pageSize > maxVideoListPageSize {
		pageSize = maxVideoListPageSize
	}

	rows := anchorvideodao.GetCachedVideoList(req.AnchorId)
	total := len(rows)
	offset := (pageIndex - 1) * pageSize
	if offset > total {
		offset = total
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	pageRows := rows[offset:end]

	viewerId := httpserver.GetAuthId(ctx)
	purchased := map[uint64]struct{}{}
	if viewerId > 0 && viewerId != req.AnchorId {
		purchased = anchorvideopurchasedao.GetPurchasedVideoIds(viewerId, req.AnchorId)
	}

	list := make([]*anchorvideodto.VideoItem, 0, len(pageRows))
	for _, row := range pageRows {
		if row == nil || row.ID == 0 || row.IsDeleted {
			continue
		}
		item := toVideoItem(row)
		item.Paid = isVideoPaidForViewer(viewerId, row, purchased)
		list = append(list, item)
	}
	return &anchorvideodto.VideoListRes{
		List:      list,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, nil
}

func isVideoPaidForViewer(viewerId uint64, row *entity.AnchorVideo, purchased map[uint64]struct{}) bool {
	if row == nil {
		return false
	}
	// 价格为 0: 免费,默认已付费
	if row.Price <= common.Zero {
		return true
	}
	if viewerId == 0 {
		return false
	}
	// 主播本人
	if viewerId == row.AnchorId {
		return true
	}
	_, ok := purchased[row.ID]
	return ok
}

// AddVideoFromStoredFile 写入单个相册视频;price=0 免费,>0 为自定义钻石价;description 可空.
// cover/fileType 原样入库; duration 为视频秒数.
func AddVideoFromStoredFile(ctx context.Context, storedName string, price float64, description, cover, fileType string, duration uint32) (*anchorvideodto.UploadVideoRes, error) {
	storedName = strings.TrimSpace(storedName)
	if storedName == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if price < common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	description = strings.TrimSpace(description)
	if utf8.RuneCountInString(description) > maxAnchorVideoDescriptionRunes {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil || !info.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	if anchorvideodao.CountVideos(userId) >= anchorvideodao.MaxVideos() {
		return nil, errercode.CreateCode(errercode.AnchorVideoLimit)
	}
	row := anchorvideodao.TryAddVideo(userId, storedName, price, description, cover, fileType, duration)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoLimit)
	}
	item := toVideoItem(row)
	return &anchorvideodto.UploadVideoRes{
		Id:          item.Id,
		Video:       item.Video,
		Cover:       item.Cover,
		FileType:    item.FileType,
		Duration:    item.Duration,
		Price:       item.Price,
		Description: item.Description,
	}, nil
}

// UpdateVideo 修改自己相册视频的价格/描述,并刷新列表缓存.
func UpdateVideo(ctx context.Context, req *anchorvideodto.UpdateVideoReq) (*anchorvideodto.UpdateVideoRes, error) {
	if req == nil || req.VideoId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.Price < common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	description := strings.TrimSpace(req.Description)
	if utf8.RuneCountInString(description) > maxAnchorVideoDescriptionRunes {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil || !info.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	row := anchorvideodao.UpdateVideoMeta(userId, req.VideoId, req.Price, description)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoNonExist)
	}
	item := toVideoItem(row)
	return &anchorvideodto.UpdateVideoRes{
		Id:          item.Id,
		Video:       item.Video,
		Cover:       item.Cover,
		FileType:    item.FileType,
		Duration:    item.Duration,
		Price:       item.Price,
		Description: item.Description,
	}, nil
}

// DeleteVideo 批量删除当前登录主播自己的相册视频
func DeleteVideo(ctx context.Context, req *anchorvideodto.DeleteVideoReq) (*anchorvideodto.DeleteVideoRes, error) {
	if req == nil || len(req.VideoIds) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorvideodao.SoftDeleteVideos(userId, req.VideoIds)
	if len(rows) == 0 {
		return nil, errercode.CreateCode(errercode.AnchorVideoNonExist)
	}
	deletedIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		upload.DeleteUploadedFile(row.Video)
		deletedIds = append(deletedIds, row.ID)
	}
	return &anchorvideodto.DeleteVideoRes{
		Success:    true,
		DeletedIds: deletedIds,
	}, nil
}

func toVideoItem(row *entity.AnchorVideo) *anchorvideodto.VideoItem {
	if row == nil {
		return &anchorvideodto.VideoItem{}
	}
	return &anchorvideodto.VideoItem{
		Id:          row.ID,
		Video:       upload.GetUrlByName(row.Video),
		Cover:       row.Cover,
		FileType:    row.FileType,
		Duration:    row.Duration,
		Price:       row.Price,
		Description: row.Description,
	}
}
