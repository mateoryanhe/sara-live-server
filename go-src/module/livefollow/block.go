package livefollow

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/livefollowdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// Block 拉黑用户(同一对用户仅一行,通过 Status 切换)
func Block(ctx context.Context, req *livefollowdto.BlockReq) (*livefollowdto.BlockRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.TargetId == 0 || req.TargetId == userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(req.TargetId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	if livefollowdao.IsBlocked(userId, req.TargetId) {
		return &livefollowdto.BlockRes{
			Success: true,
			Blocked: true,
		}, nil
	}

	existing := livefollowdao.GetByUserAnchor(userId, req.TargetId)
	if existing == nil {
		row := entity.NewLiveFollowWithStatus(userId, req.TargetId, entity.LiveFollowStatusBlock)
		livefollowdao.AddFollowToCache(row)
		livefollowdao.PrependBlockedToListCache(row)
	} else {
		wasFollowing := existing.Status == entity.LiveFollowStatusFollow
		existing.SetStatus(entity.LiveFollowStatusBlock)
		livefollowdao.AddFollowToCache(existing)
		livefollowdao.PrependBlockedToListCache(existing)
		if wasFollowing {
			userinfodao.DecFollowCount(userId, req.TargetId)
			livefollowdao.RemoveFollowingFromListCache(userId, req.TargetId)
			livefollowdao.RemoveFollowerFromListCache(req.TargetId, userId)
		}
	}

	return &livefollowdto.BlockRes{
		Success: true,
		Blocked: true,
	}, nil
}

// Unblock 解除拉黑;不存在或已是非拉黑状态都按幂等处理
func Unblock(ctx context.Context, req *livefollowdto.UnblockReq) (*livefollowdto.UnblockRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.TargetId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := livefollowdao.GetByUserAnchor(userId, req.TargetId)
	if existing != nil && existing.Status == entity.LiveFollowStatusBlock {
		existing.SetStatus(entity.LiveFollowStatusUnfollow)
		livefollowdao.AddFollowToCache(existing)
		livefollowdao.RemoveBlockedFromListCache(userId, req.TargetId)
	}

	return &livefollowdto.UnblockRes{
		Success: true,
		Blocked: false,
	}, nil
}

// BlockList 当前用户拉黑列表(分页)
func BlockList(ctx context.Context, req *livefollowdto.BlockListReq) (*livefollowdto.BlockListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	page, pageSize := normalizePage(req.Page, req.PageSize)

	pageData := livefollowdao.GetBlockedListByUser(userId, page, pageSize)
	list := make([]*livefollowdto.BlockListItem, 0, len(pageData))
	for _, row := range pageData {
		item := &livefollowdto.BlockListItem{
			TargetId:  strconv.FormatUint(row.AnchorId, 10),
			BlockedAt: row.UpdatedAt.Unix(),
		}
		if u := userinfodao.GetUserInfoByUserId(row.AnchorId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrlForUser(row.AnchorId, u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}

	return &livefollowdto.BlockListRes{
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}
