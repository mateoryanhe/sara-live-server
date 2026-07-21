package userinfo

import (
	"context"

	"xr-game-server/dto/userinfodto"
)

// SubmitFeedback App端用户反馈(空实现)
func SubmitFeedback(_ context.Context, _ *userinfodto.AppFeedbackReq) (*userinfodto.AppFeedbackRes, error) {
	return &userinfodto.AppFeedbackRes{Success: true}, nil
}
