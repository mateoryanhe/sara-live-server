package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/liveroom"
)

const OneToOneRoomCMSUrl = "/oneToOneRoom"

type OneToOneRoomCMSController struct{}

func initOneToOneRoomCMSController() {
	httpserver.RegCMS(OneToOneRoomCMSUrl, &OneToOneRoomCMSController{})
}

func (c *OneToOneRoomCMSController) List(ctx context.Context, req *liveroomdto.CMSOneToOneRoomListReq) (*httpserver.CMSQueryResp, error) {
	return liveroom.GetCMSOneToOneRoomList(ctx, req)
}

func (c *OneToOneRoomCMSController) Create(ctx context.Context, req *liveroomdto.CMSCreateOneToOneRoomReq) (*liveroomdto.CMSCreateOneToOneRoomRes, error) {
	return liveroom.CMSCreateOneToOneRoom(ctx, req)
}

func (c *OneToOneRoomCMSController) Update(ctx context.Context, req *liveroomdto.CMSUpdateOneToOneRoomReq) (*liveroomdto.CMSUpdateOneToOneRoomRes, error) {
	return liveroom.CMSUpdateOneToOneRoom(ctx, req)
}

func (c *OneToOneRoomCMSController) SetStatus(ctx context.Context, req *liveroomdto.CMSSetOneToOneRoomStatusReq) (*liveroomdto.CMSSetOneToOneRoomStatusRes, error) {
	return liveroom.CMSSetOneToOneRoomStatus(ctx, req)
}
