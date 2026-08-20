package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/module/cmsexport"
)

const CMSExportUrl = "/cmsExport"

type CMSExportController struct{}

func initCMSExportController() {
	httpserver.RegCMS(CMSExportUrl, &CMSExportController{})
}

func (c *CMSExportController) SubmitExportJob(ctx context.Context, req *cmsexportdto.CMSSubmitExportJobReq) (res *cmsexportdto.CMSSubmitExportJobRes, err error) {
	return cmsexport.SubmitExportJob(ctx, req)
}

func (c *CMSExportController) GetExportJob(ctx context.Context, req *cmsexportdto.CMSGetExportJobReq) (res *cmsexportdto.CMSGetExportJobRes, err error) {
	return cmsexport.GetExportJob(ctx, req)
}

func (c *CMSExportController) DeleteExport(ctx context.Context, req *cmsexportdto.CMSDeleteExportReq) (res *cmsexportdto.CMSDeleteExportRes, err error) {
	return cmsexport.DeleteExport(ctx, req)
}
