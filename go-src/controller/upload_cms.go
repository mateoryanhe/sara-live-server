package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/uploaddto"
	"xr-game-server/module/upload"
)

const (
	UploadUrl = "/upload"
)

type UploadController struct{}

func initUploadController() {
	httpserver.RegCMSHandler(UploadUrl, "/uploadFile", handleCMSUploadFile)
	httpserver.RegCMS(UploadUrl, &UploadController{})
}

func handleCMSUploadFile(r *ghttp.Request) {
	name, err := upload.UploadCMSFileFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, &uploaddto.UploadCMSFileRes{
		FileName: name,
		FileUrl:  upload.GetUrlByName(name),
	})
}

func (c *UploadController) GetUploadResourceCfg(ctx context.Context, req *uploaddto.GetUploadResourceCfgReq) (*uploaddto.GetUploadResourceCfgRes, error) {
	return upload.GetUploadResourceCfg(ctx, req)
}

func (c *UploadController) SaveUploadResourceCfg(ctx context.Context, req *uploaddto.SaveUploadResourceCfgReq) (*uploaddto.SaveUploadResourceCfgRes, error) {
	return upload.SaveUploadResourceCfg(ctx, req)
}
