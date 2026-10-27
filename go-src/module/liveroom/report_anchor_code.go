package liveroom

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/auditcode"
)

// ReportAnchorCode App 上报主播码/审核码,校验通过后将当前用户设为普通主播.
// 审核码每次从库读取(无缓存),由 CMS「审核码生成」覆盖写入.
func ReportAnchorCode(ctx context.Context, req *liveroomdto.ReportAnchorCodeReq) (*liveroomdto.ReportAnchorCodeRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	got := strings.TrimSpace(req.AnchorCode)
	want := auditcode.CurrentAuditCode()
	if got == "" || got != want {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if _, err := setUserAsAnchor(userId, entity.UserTypeAnchor); err != nil {
		return nil, err
	}
	return &liveroomdto.ReportAnchorCodeRes{Success: true}, nil
}
