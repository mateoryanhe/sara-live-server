package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/upload"
	"xr-game-server/module/userinfo"
)

// handleUploadAvatar 由控制器逐段解析 multipart，图片流直接交给上传模块。
func handleUploadAvatar(r *ghttp.Request) {
	storedName, err := parseAppImageMultipartFile(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	res, err := userinfo.UpdateAvatarFromStoredFile(r.Context(), storedName)
	if err != nil {
		upload.DeleteUploadedFile(storedName)
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}
