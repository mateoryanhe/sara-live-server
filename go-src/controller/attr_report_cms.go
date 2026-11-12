package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/attrreportdto"
	"xr-game-server/module/attrreport"
)

const AttrReportCMSUrl = "/attrReportCfg"

type AttrReportCMSController struct{}

func initAttrReportCMSController() {
	httpserver.RegCMS(AttrReportCMSUrl, &AttrReportCMSController{})
}

func (c *AttrReportCMSController) GetAttrReportCfg(ctx context.Context, req *attrreportdto.GetAttrReportCfgReq) (*attrreportdto.GetAttrReportCfgRes, error) {
	return attrreport.GetAttrReportCfg(ctx, req)
}

func (c *AttrReportCMSController) SaveAttrReportCfg(ctx context.Context, req *attrreportdto.SaveAttrReportCfgReq) (*attrreportdto.SaveAttrReportCfgRes, error) {
	return attrreport.SaveAttrReportCfg(ctx, req)
}
