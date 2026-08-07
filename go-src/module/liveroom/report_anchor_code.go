package liveroom

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const appNormalAnchorCode = "971200"

// ReportAnchorCode App 上报主播码,校验通过后将当前用户设为普通主播.
func ReportAnchorCode(ctx context.Context, req *liveroomdto.ReportAnchorCodeReq) (*liveroomdto.ReportAnchorCodeRes, error) {
	if req == nil || strings.TrimSpace(req.AnchorCode) != appNormalAnchorCode {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if _, err := setUserAsAnchor(userId, entity.UserTypeAnchor); err != nil {
		return nil, err
	}
	return &liveroomdto.ReportAnchorCodeRes{Success: true}, nil
}
