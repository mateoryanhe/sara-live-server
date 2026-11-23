package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/upload"
)

func handleUploadOneToOneCover(r *ghttp.Request) {
	storedName, err := parseAppImageMultipartFile(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	userId := httpserver.GetAuthId(r.Context())
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil || !info.IsAnchor() {
		upload.DeleteUploadedFile(storedName)
		r.SetError(errercode.CreateCode(errercode.LiveRoomNotAnchor))
		return
	}
	res := liveroom.BuildOneToOneCoverUploadRes(storedName)
	httpserver.SetHandlerResponseData(r, res)
}

// HandleUploadOneToOneCover Sofie 等新前缀复用同一套 1v1 封面上传解析。
func HandleUploadOneToOneCover(r *ghttp.Request) {
	handleUploadOneToOneCover(r)
}