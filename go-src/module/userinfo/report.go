package userinfo

import (
	"context"

	"xr-game-server/dto/userinfodto"
)

// SubmitReport App端用户举报(空实现)
func SubmitReport(_ context.Context, _ *userinfodto.AppReportReq) (*userinfodto.AppReportRes, error) {
	return &userinfodto.AppReportRes{Success: true}, nil
}
