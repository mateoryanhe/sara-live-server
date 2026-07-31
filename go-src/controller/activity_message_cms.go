package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/activitymessagedto"
	"xr-game-server/module/message"
)

const ActivityMessageUrl = "/activityMessage"

type ActivityMessageController struct{}

func initActivityMessageController() {
	httpserver.RegCMS(ActivityMessageUrl, &ActivityMessageController{})
}

func (c *ActivityMessageController) ActivityMessageList(ctx context.Context, req *activitymessagedto.ActivityMessageListReq) (*httpserver.CMSQueryResp, error) {
	return message.GetActivityMessageList(ctx, req)
}

func (c *ActivityMessageController) CreateActivityMessage(ctx context.Context, req *activitymessagedto.CreateActivityMessageReq) (*activitymessagedto.CreateActivityMessageRes, error) {
	return message.CreateActivityMessage(ctx, req)
}

func (c *ActivityMessageController) UpdateActivityMessage(ctx context.Context, req *activitymessagedto.UpdateActivityMessageReq) (*activitymessagedto.UpdateActivityMessageRes, error) {
	return message.UpdateActivityMessage(ctx, req)
}

func (c *ActivityMessageController) DeleteActivityMessage(ctx context.Context, req *activitymessagedto.DeleteActivityMessageReq) (*activitymessagedto.DeleteActivityMessageRes, error) {
	return message.DeleteActivityMessage(ctx, req)
}

func (c *ActivityMessageController) PublishActivityMessage(ctx context.Context, req *activitymessagedto.PublishActivityMessageReq) (*activitymessagedto.PublishActivityMessageRes, error) {
	return message.PublishActivityMessage(ctx, req)
}

func (c *ActivityMessageController) UnpublishActivityMessage(ctx context.Context, req *activitymessagedto.UnpublishActivityMessageReq) (*activitymessagedto.UnpublishActivityMessageRes, error) {
	return message.UnpublishActivityMessage(ctx, req)
}
