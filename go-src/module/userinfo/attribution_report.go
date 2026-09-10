package userinfo

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userinfodto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// ReportAttribution App端上报归因信息(syndb 缓冲入库)
func ReportAttribution(ctx context.Context, req *userinfodto.ReportAttributionReq) (*userinfodto.ReportAttributionRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	enabled := false
	if req.AttributionEnabled != nil {
		enabled = *req.AttributionEnabled
	}
	_ = userentity.NewUserAttributionReport(
		userId,
		httpserver.GetPackageNameFromContext(ctx),
		enabled,
		strings.TrimSpace(req.AttributionProvider),
		strings.TrimSpace(req.AppsFlyerDevKey),
		strings.TrimSpace(req.AppsFlyerAppId),
	)
	return &userinfodto.ReportAttributionRes{Success: true}, nil
}
