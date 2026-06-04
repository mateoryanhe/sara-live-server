package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/uploaddto"
	"xr-game-server/module/upload"
)

const (
	UploadUrl = "/upload"
)

type UploadController struct{}

func initUploadController() {
	httpserver.RegCMS(UploadUrl, &UploadController{})
}

// UploadFile CMS后台上传图片或礼物动画资源,返回文件名
func (c *UploadController) UploadFile(_ context.Context, req *uploaddto.UploadCMSFileReq) (res *uploaddto.UploadCMSFileRes, err error) {
	name, err := upload.UploadCMSFile(req.File)
	if err != nil {
		return nil, err
	}
	return &uploaddto.UploadCMSFileRes{FileName: name}, nil
}

func (c *UploadController) GetUploadResourceCfg(ctx context.Context, req *uploaddto.GetUploadResourceCfgReq) (*uploaddto.GetUploadResourceCfgRes, error) {
	return upload.GetUploadResourceCfg(ctx, req)
}

func (c *UploadController) SaveUploadResourceCfg(ctx context.Context, req *uploaddto.SaveUploadResourceCfgReq) (*uploaddto.SaveUploadResourceCfgRes, error) {
	return upload.SaveUploadResourceCfg(ctx, req)
}
