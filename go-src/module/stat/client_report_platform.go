package stat

import (
	"xr-game-server/dao/userinfodao"
	"xr-game-server/module/promotionattribution"
)

// resolveUserClientReportPlatform 仪表盘客户端 App/H5 分桶.
// 仅按注册时请求头 X-Platform 落库的 register_client_platform:
// web → h5；android/ios → app；空/未知不计入.
func resolveUserClientReportPlatform(userId uint64) string {
	if userId == 0 {
		return ""
	}
	return promotionattribution.MapClientXPlatformToRegisterPlatform(userinfodao.GetRegisterClientPlatform(userId))
}
