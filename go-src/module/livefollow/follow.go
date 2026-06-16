package livefollow

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/userstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/livefollowdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Follow 关注主播(同一对粉丝/主播仅一行,通过 Status 切换)
//  1. 不允许关注自己
//  2. 主播必须存在(以 user_info 是否存在为准)
//  3. 已关注则幂等返回当前状态
func Follow(ctx context.Context, req *livefollowdto.FollowReq) (*livefollowdto.FollowRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if req.AnchorId == 0 || req.AnchorId == userId {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	if anchor := userinfodao.GetUserInfoByUserId(req.AnchorId); anchor == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}

	existing := livefollowdao.GetByUserAnchor(userId, req.AnchorId)
	if existing == nil {
		f := entity.NewLiveFollow(userId, req.AnchorId)
		livefollowdao.AddFollowToCache(f)
	} else if existing.Status != entity.LiveFollowStatusFollow {
		existing.SetStatus(entity.LiveFollowStatusFollow)
		livefollowdao.AddFollowToCache(existing)
	}

	return &livefollowdto.FollowRes{
		Following:     true,
		FollowerCount: len(livefollowdao.GetFollowersByAnchor(req.AnchorId)),
	}, nil
}

// Unfollow 取消关注;不存在或已是未关注状态都按幂等处理
func Unfollow(ctx context.Context, req *livefollowdto.UnfollowReq) (*livefollowdto.UnfollowRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	existing := livefollowdao.GetByUserAnchor(userId, req.AnchorId)
	if existing != nil && existing.Status != entity.LiveFollowStatusUnfollow {
		existing.SetStatus(entity.LiveFollowStatusUnfollow)
		livefollowdao.RemoveFollowFromCache(existing)
	}

	return &livefollowdto.UnfollowRes{
		Following:     false,
		FollowerCount: len(livefollowdao.GetFollowersByAnchor(req.AnchorId)),
	}, nil
}

// IsFollowing 查询当前用户是否关注指定主播
func IsFollowing(ctx context.Context, req *livefollowdto.IsFollowingReq) (*livefollowdto.IsFollowingRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if req.AnchorId == 0 {
		return &livefollowdto.IsFollowingRes{Following: false}, nil
	}
	follow := livefollowdao.IsFollowing(userId, req.AnchorId)
	return &livefollowdto.IsFollowingRes{Following: follow}, nil
}

// FollowingList 当前用户关注的主播分页列表
func FollowingList(ctx context.Context, req *livefollowdto.FollowingListReq) (*livefollowdto.FollowingListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizePage(req.Page, req.PageSize)

	all := livefollowdao.GetFollowingsByUser(userId)
	total := len(all)
	start, end := pageRange(total, page, pageSize)
	pageData := all[start:end]

	list := make([]*livefollowdto.AnchorItem, 0, len(pageData))
	for _, f := range pageData {
		item := &livefollowdto.AnchorItem{
			AnchorId:   strconv.FormatUint(f.AnchorId, 10),
			FollowedAt: f.UpdatedAt.Unix(),
			RoomStatus: roomStatus(f.AnchorId),
		}
		if u := userinfodao.GetUserInfoByUserId(f.AnchorId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrl(u.Avatar)
			item.GuildId = strconv.FormatUint(u.GuildId, 10)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}

	return &livefollowdto.FollowingListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// FollowerList 当前用户(作为主播)的粉丝分页列表
func FollowerList(ctx context.Context, req *livefollowdto.FollowerListReq) (*livefollowdto.FollowerListRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizePage(req.Page, req.PageSize)

	all := livefollowdao.GetFollowersByAnchor(anchorId)
	total := len(all)
	start, end := pageRange(total, page, pageSize)
	pageData := all[start:end]

	list := make([]*livefollowdto.FollowerItem, 0, len(pageData))
	for _, f := range pageData {
		item := &livefollowdto.FollowerItem{
			UserId:     strconv.FormatUint(f.UserId, 10),
			FollowedAt: f.UpdatedAt.Unix(),
		}
		if u := userinfodao.GetUserInfoByUserId(f.UserId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrl(u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}

	return &livefollowdto.FollowerListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// roomStatus 返回直播间状态;无直播间返回 2
func roomStatus(anchorId uint64) uint8 {
	if r := liveroomdao.GetRoomById(anchorId); r != nil {
		if r.LiveRecordId > 0 {
			return userstatus.LiveRoomStatusLive
		}
		return userstatus.LiveRoomStatusClosed
	}
	return userstatus.Empty
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

func calcAge(birthday *time.Time) int {
	if birthday == nil || birthday.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	anniversary := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

func pageRange(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}
