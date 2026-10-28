package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userappsubmissiondto"
	"xr-game-server/module/userappsubmission"
)

const UserFeedbackRecordCMSUrl = "/userFeedbackRecord"
const UserReportRecordCMSUrl = "/userReportRecord"

type UserFeedbackRecordCMSController struct{}

type UserReportRecordCMSController struct{}

func initUserAppSubmissionCMSController() {
	httpserver.RegCMS(UserFeedbackRecordCMSUrl, &UserFeedbackRecordCMSController{})
	httpserver.RegCMS(UserReportRecordCMSUrl, &UserReportRecordCMSController{})
}

func (c *UserFeedbackRecordCMSController) ListUserFeedbackRecord(ctx context.Context, req *userappsubmissiondto.CMSListUserFeedbackRecordReq) (*httpserver.CMSQueryResp, error) {
	return userappsubmission.ListFeedbackCMS(ctx, req)
}

func (c *UserReportRecordCMSController) ListUserReportRecord(ctx context.Context, req *userappsubmissiondto.CMSListUserReportRecordReq) (*httpserver.CMSQueryResp, error) {
	return userappsubmission.ListReportCMS(ctx, req)
}
