package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userinfodto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SubmitReport App端用户举报(请求体原文 syndb 入库)
func SubmitReport(ctx context.Context, _ *userinfodto.AppReportReq) (*userinfodto.AppReportRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	_ = userentity.NewUserReportRecord(userId, readAppSubmissionBody(ctx))
	return &userinfodto.AppReportRes{Success: true}, nil
}
