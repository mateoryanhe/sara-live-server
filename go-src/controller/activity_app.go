package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/activitydto"
	"xr-game-server/module/gameactivity"
)

// ActivityAppUrl App 活动相关接口前缀(对外路径不含游戏字眼).
const ActivityAppUrl = "/activity"

type ActivityAppController struct{}

func initActivityAppController() {
	httpserver.RegAPI(ActivityAppUrl, &ActivityAppController{})
}

// GetRoomActivity 房间活动配置(按钮 icon / 列表 URL / 推荐 itemId)
func (c *ActivityAppController) GetRoomActivity(ctx context.Context, req *activitydto.GetRoomActivityReq) (*activitydto.GetRoomActivityRes, error) {
	return gameactivity.GetRoomActivity(ctx, req)
}
