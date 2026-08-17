package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
	"xr-game-server/module/auth"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/userinfo"
)

const (
	AccountUrl = "/account"
)

type AccountController struct {
}

func initAccountController() {
	httpserver.RegCMS(AccountUrl, &AccountController{})
}

func (a *AccountController) Ban(ctx context.Context, req *accountdto.BanReq) (resp *accountdto.BanRes, e error) {
	return auth.Ban(ctx, req)
}

func (a *AccountController) BanAnchor(ctx context.Context, req *accountdto.BanAnchorReq) (resp *accountdto.BanRes, e error) {
	return liveroom.BanAnchor(ctx, req)
}

func (a *AccountController) UnBanAnchor(ctx context.Context, req *accountdto.UnBanAnchorReq) (bool, error) {
	return liveroom.UnBanAnchor(ctx, req)
}

// SetLiveRoomStatus CMS上架/下架主播直播间
func (a *AccountController) SetLiveRoomStatus(ctx context.Context, req *accountdto.SetLiveRoomStatusReq) (*accountdto.SetLiveRoomStatusRes, error) {
	return liveroom.SetLiveRoomStatus(ctx, req)
}

func (a *AccountController) UnBan(ctx context.Context, req *accountdto.UnBanReq) (bool, error) {
	return auth.UnBan(ctx, req)
}

func (a *AccountController) CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	return auth.CancelUser(ctx, req)
}

func (a *AccountController) UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	return auth.UnCancelUser(ctx, req)
}

func (a *AccountController) QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	return userinfo.QueryUserInfo(ctx, req)
}

func (a *AccountController) GetUserDetail(ctx context.Context, req *accountdto.GetUserDetailReq) (*accountdto.GetUserDetailRes, error) {
	return userinfo.QueryUserDetail(ctx, req)
}

func (a *AccountController) QueryAnchorList(ctx context.Context, req *accountdto.QueryAnchorListReq) (res *httpserver.CMSQueryResp, err error) {
	return liveroom.QueryAnchorList(ctx, req)
}

func (a *AccountController) GetAnchorDetail(ctx context.Context, req *accountdto.GetAnchorDetailReq) (*accountdto.GetAnchorDetailRes, error) {
	return liveroom.QueryAnchorDetail(ctx, req)
}

func (a *AccountController) GetAnchorDailyEffectiveLiveList(ctx context.Context, req *accountdto.GetAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	return liveroom.QueryAnchorDailyEffectiveLiveList(ctx, req)
}

// QueryOffShelfLiveRoomList CMS回收站:查询已下架直播间
func (a *AccountController) QueryOffShelfLiveRoomList(ctx context.Context, req *accountdto.QueryOffShelfLiveRoomListReq) (res *httpserver.CMSQueryResp, err error) {
	return liveroom.QueryOffShelfLiveRoomList(ctx, req)
}

func (a *AccountController) SetAnchor(ctx context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	return liveroom.SetAnchor(ctx, req)
}

func (a *AccountController) SetSeniorAnchor(ctx context.Context, req *accountdto.SetSeniorAnchorReq) (*accountdto.SetSeniorAnchorRes, error) {
	return liveroom.SetSeniorAnchor(ctx, req)
}

func (a *AccountController) BatchSetAnchor(ctx context.Context, req *accountdto.BatchSetAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return liveroom.BatchSetAnchor(ctx, req)
}

func (a *AccountController) BatchSetSeniorAnchor(ctx context.Context, req *accountdto.BatchSetSeniorAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return liveroom.BatchSetSeniorAnchor(ctx, req)
}

func (a *AccountController) SetUserType(ctx context.Context, req *accountdto.SetUserTypeReq) (*accountdto.SetUserTypeRes, error) {
	return userinfo.SetUserType(ctx, req)
}

func (a *AccountController) SetCanRank(ctx context.Context, req *accountdto.SetCanRankReq) (*accountdto.SetCanRankRes, error) {
	return userinfo.SetCanRank(ctx, req)
}

func (a *AccountController) SetRechargeWhitelist(ctx context.Context, req *accountdto.SetRechargeWhitelistReq) (*accountdto.SetRechargeWhitelistRes, error) {
	return userinfo.SetRechargeWhitelist(ctx, req)
}

// ExitGuild CMS主播退出工会(将工会ID置为0)
func (a *AccountController) ExitGuild(ctx context.Context, req *accountdto.ExitGuildReq) (*accountdto.ExitGuildRes, error) {
	return liveroom.ExitGuild(ctx, req)
}
