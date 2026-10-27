package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/auditcodedto"
	"xr-game-server/module/auditcode"
)

const AuditCodeCMSUrl = "/auditCode"

type AuditCodeCMSController struct{}

func initAuditCodeCMSController() {
	httpserver.RegCMS(AuditCodeCMSUrl, &AuditCodeCMSController{})
}

func (c *AuditCodeCMSController) GetAuditCodeCfg(ctx context.Context, req *auditcodedto.GetAuditCodeCfgReq) (*auditcodedto.GetAuditCodeCfgRes, error) {
	return auditcode.GetAuditCodeCfg(ctx, req)
}

func (c *AuditCodeCMSController) GenerateAuditCodeCfg(ctx context.Context, req *auditcodedto.GenerateAuditCodeCfgReq) (*auditcodedto.GenerateAuditCodeCfgRes, error) {
	return auditcode.GenerateAuditCodeCfg(ctx, req)
}
