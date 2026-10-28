package userappsubmission

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userappsubmissiondao"
	"xr-game-server/dto/userappsubmissiondto"
	userentity "xr-game-server/entity/user"
)

func ListFeedbackCMS(ctx context.Context, req *userappsubmissiondto.CMSListUserFeedbackRecordReq) (*httpserver.CMSQueryResp, error) {
	_ = ctx
	if req == nil {
		req = &userappsubmissiondto.CMSListUserFeedbackRecordReq{}
	}
	total, rows := userappsubmissiondao.ListFeedback(&userappsubmissiondao.ListFilter{
		UserId:    req.UserId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*userappsubmissiondto.CMSUserFeedbackRecordItem, 0, len(rows))
	for _, row := range rows {
		if item := feedbackToItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func ListReportCMS(ctx context.Context, req *userappsubmissiondto.CMSListUserReportRecordReq) (*httpserver.CMSQueryResp, error) {
	_ = ctx
	if req == nil {
		req = &userappsubmissiondto.CMSListUserReportRecordReq{}
	}
	total, rows := userappsubmissiondao.ListReport(&userappsubmissiondao.ListFilter{
		UserId:    req.UserId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*userappsubmissiondto.CMSUserReportRecordItem, 0, len(rows))
	for _, row := range rows {
		if item := reportToItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func feedbackToItem(row *userentity.UserFeedbackRecord) *userappsubmissiondto.CMSUserFeedbackRecordItem {
	if row == nil {
		return nil
	}
	return &userappsubmissiondto.CMSUserFeedbackRecordItem{
		ID:        row.ID,
		UserId:    row.UserId,
		Content:   row.Content,
		CreatedAt: formatTime(row.CreatedAt),
	}
}

func reportToItem(row *userentity.UserReportRecord) *userappsubmissiondto.CMSUserReportRecordItem {
	if row == nil {
		return nil
	}
	return &userappsubmissiondto.CMSUserReportRecordItem{
		ID:        row.ID,
		UserId:    row.UserId,
		Content:   row.Content,
		CreatedAt: formatTime(row.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}
