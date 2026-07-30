package livefollow

import (
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/livefollowdto"
)

func pushFollowCountToUser(userId uint64) {
	if userId == 0 {
		return
	}
	push.Data(userId, cmd.LiveFollowCountPush, &livefollowdto.FollowCountPushItem{
		FollowCount:   userinfodao.GetFollowCount(userId),
		FollowerCount: userinfodao.GetFollowerCount(userId),
	})
}

func pushFollowCountChange(userId, anchorId uint64) {
	pushFollowCountToUser(userId)
	if anchorId != userId {
		pushFollowCountToUser(anchorId)
	}
}
