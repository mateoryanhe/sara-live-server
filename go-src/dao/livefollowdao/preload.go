package livefollowdao

// PreloadLiveFollowToCache 批量预热最近活跃用户的关注/粉丝/拉黑列表及单条记录缓存
func PreloadLiveFollowToCache(userIds []uint64) {
	if len(userIds) == 0 {
		return
	}
	for _, userId := range userIds {
		if userId == 0 {
			continue
		}
		preloadUserLiveFollowCaches(userId)
	}
}

func preloadUserLiveFollowCaches(userId uint64) {
	following := loadFollowingsFromDB(userId, 1, followingListCacheMaxSize)
	putFollowingListCache(userId, following)
	for _, row := range following {
		AddFollowToCache(row)
	}

	blocked := loadBlockedListFromDB(userId, 1, followingListCacheMaxSize)
	putBlockListCache(userId, blocked)
	for _, row := range blocked {
		AddFollowToCache(row)
	}

	followers := loadFollowersFromDB(userId, 1, followingListCacheMaxSize)
	putFollowerListCache(userId, followers)
	for _, row := range followers {
		AddFollowToCache(row)
	}
}
