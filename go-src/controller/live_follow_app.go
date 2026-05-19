package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/livefollowdto"
	"xr-game-server/module/livefollow"
)

const (
	LiveFollowAppUrl = "/liveFollow"
)

type LiveFollowAppController struct{}

func initLiveFollowAppController() {
	httpserver.RegAPI(LiveFollowAppUrl, &LiveFollowAppController{})
}

// Follow 关注主播
func (c *LiveFollowAppController) Follow(ctx context.Context, req *livefollowdto.FollowReq) (res *livefollowdto.FollowRes, err error) {
	return livefollow.Follow(ctx, req)
}

// Unfollow 取消关注
func (c *LiveFollowAppController) Unfollow(ctx context.Context, req *livefollowdto.UnfollowReq) (res *livefollowdto.UnfollowRes, err error) {
	return livefollow.Unfollow(ctx, req)
}

// IsFollowing 查询是否关注
func (c *LiveFollowAppController) IsFollowing(ctx context.Context, req *livefollowdto.IsFollowingReq) (res *livefollowdto.IsFollowingRes, err error) {
	return livefollow.IsFollowing(ctx, req)
}

// FollowingList 我关注的主播列表
func (c *LiveFollowAppController) FollowingList(ctx context.Context, req *livefollowdto.FollowingListReq) (res *livefollowdto.FollowingListRes, err error) {
	return livefollow.FollowingList(ctx, req)
}

// FollowerList 我的粉丝列表
func (c *LiveFollowAppController) FollowerList(ctx context.Context, req *livefollowdto.FollowerListReq) (res *livefollowdto.FollowerListRes, err error) {
	return livefollow.FollowerList(ctx, req)
}
