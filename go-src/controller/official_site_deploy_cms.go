package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/officialsitedeploydto"
	"xr-game-server/module/officialsitedeploy"
)

const (
	OfficialSiteDeployCMSUrl         = "/officialSiteDeploy"
	ThirdPayOfficialSiteDeployCMSUrl = "/thirdPayOfficialSiteDeploy"
)

type OfficialSiteDeployCMSController struct{}
type ThirdPayOfficialSiteDeployCMSController struct{}

func initOfficialSiteDeployCMSController() {
	httpserver.RegCMSHandler(OfficialSiteDeployCMSUrl, "/deployZip", handleOfficialSiteDeployZip)
	httpserver.RegCMSHandler(OfficialSiteDeployCMSUrl, "/uploadFileChunk", handleOfficialSiteUploadFileChunk)
	httpserver.RegCMS(OfficialSiteDeployCMSUrl, &OfficialSiteDeployCMSController{})
	httpserver.RegCMSHandler(ThirdPayOfficialSiteDeployCMSUrl, "/deployZip", handleThirdPayOfficialSiteDeployZip)
	httpserver.RegCMSHandler(ThirdPayOfficialSiteDeployCMSUrl, "/uploadFileChunk", handleThirdPayOfficialSiteUploadFileChunk)
	httpserver.RegCMS(ThirdPayOfficialSiteDeployCMSUrl, &ThirdPayOfficialSiteDeployCMSController{})
}

func handleOfficialSiteUploadFileChunk(r *ghttp.Request) {
	res, err := officialsitedeploy.UploadOfficialSiteFileChunkFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func handleOfficialSiteDeployZip(r *ghttp.Request) {
	res, err := officialsitedeploy.DeployOfficialSiteZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func handleThirdPayOfficialSiteDeployZip(r *ghttp.Request) {
	res, err := officialsitedeploy.DeployThirdPayOfficialSiteZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func handleThirdPayOfficialSiteUploadFileChunk(r *ghttp.Request) {
	res, err := officialsitedeploy.UploadThirdPayOfficialSiteFileChunkFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *OfficialSiteDeployCMSController) GetOfficialSiteDeployInfo(ctx context.Context, req *officialsitedeploydto.GetOfficialSiteDeployInfoReq) (*officialsitedeploydto.GetOfficialSiteDeployInfoRes, error) {
	return officialsitedeploy.GetOfficialSiteDeployInfo(ctx, req)
}

func (c *OfficialSiteDeployCMSController) SaveOfficialSiteDeployCfg(ctx context.Context, req *officialsitedeploydto.SaveOfficialSiteDeployCfgReq) (*officialsitedeploydto.SaveOfficialSiteDeployCfgRes, error) {
	return officialsitedeploy.SaveOfficialSiteDeployCfg(ctx, req)
}

func (c *OfficialSiteDeployCMSController) InitOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.InitOfficialSiteFileUploadReq) (*officialsitedeploydto.InitOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.InitOfficialSiteFileUpload(ctx, req)
}

func (c *OfficialSiteDeployCMSController) CompleteOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.CompleteOfficialSiteFileUploadReq) (*officialsitedeploydto.CompleteOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.CompleteOfficialSiteFileUpload(ctx, req)
}

func (c *OfficialSiteDeployCMSController) AbortOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.AbortOfficialSiteFileUploadReq) (*officialsitedeploydto.AbortOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.AbortOfficialSiteFileUpload(ctx, req)
}

func (c *ThirdPayOfficialSiteDeployCMSController) GetThirdPayOfficialSiteDeployInfo(ctx context.Context, req *officialsitedeploydto.GetThirdPayOfficialSiteDeployInfoReq) (*officialsitedeploydto.GetThirdPayOfficialSiteDeployInfoRes, error) {
	return officialsitedeploy.GetThirdPayOfficialSiteDeployInfo(ctx, req)
}

func (c *ThirdPayOfficialSiteDeployCMSController) SaveThirdPayOfficialSiteDeployCfg(ctx context.Context, req *officialsitedeploydto.SaveThirdPayOfficialSiteDeployCfgReq) (*officialsitedeploydto.SaveThirdPayOfficialSiteDeployCfgRes, error) {
	return officialsitedeploy.SaveThirdPayOfficialSiteDeployCfg(ctx, req)
}

func (c *ThirdPayOfficialSiteDeployCMSController) InitThirdPayOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.InitThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.InitThirdPayOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.InitThirdPayOfficialSiteFileUpload(ctx, req)
}

func (c *ThirdPayOfficialSiteDeployCMSController) CompleteThirdPayOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.CompleteThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.CompleteThirdPayOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.CompleteThirdPayOfficialSiteFileUpload(ctx, req)
}

func (c *ThirdPayOfficialSiteDeployCMSController) AbortThirdPayOfficialSiteFileUpload(ctx context.Context, req *officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadReq) (*officialsitedeploydto.AbortThirdPayOfficialSiteFileUploadRes, error) {
	return officialsitedeploy.AbortThirdPayOfficialSiteFileUpload(ctx, req)
}
