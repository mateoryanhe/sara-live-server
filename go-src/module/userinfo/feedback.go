package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userinfodto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SubmitFeedback App端用户反馈(请求体原文 syndb 入库)
func SubmitFeedback(ctx context.Context, _ *userinfodto.AppFeedbackReq) (*userinfodto.AppFeedbackRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	_ = userentity.NewUserFeedbackRecord(userId, readAppSubmissionBody(ctx))
	return &userinfodto.AppFeedbackRes{Success: true}, nil
}
