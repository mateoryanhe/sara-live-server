package vip

import (
	"strconv"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dto/vipdto"
)

// pushVipLevelToApp 推送用户最新VIP等级到 App 端
func pushVipLevelToApp(userId uint64, vipLevel uint32) {
	item := &vipdto.VipLevelPushItem{
		VipLevel: strconv.FormatUint(uint64(vipLevel), 10),
	}
	if cfg := GetVipCfgFromMemoryByLevel(vipLevel); cfg != nil {
		item.LevelName = cfg.LevelName
	}
	push.Data(userId, cmd.VipLevelPush, item)
}
