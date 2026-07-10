package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/calldto"
	"xr-game-server/module/call"
)

const CallAppUrl = "/call"

type CallAppController struct{}

func initCallAppController() {
	httpserver.RegAPI(CallAppUrl, &CallAppController{})
}

// LiveRoomCall 直播间通话呼叫
func (c *CallAppController) LiveRoomCall(ctx context.Context, req *calldto.LiveRoomCallReq) (*calldto.LiveRoomCallRes, error) {
	return call.LiveRoomCall(ctx, req)
}

// AnchorRejectCall 拒接通话
func (c *CallAppController) AnchorRejectCall(ctx context.Context, req *calldto.AnchorRejectCallReq) (*calldto.AnchorRejectCallRes, error) {
	return call.AnchorRejectCall(ctx, req)
}

// AcceptCall 同意接听通话
func (c *CallAppController) AcceptCall(ctx context.Context, req *calldto.AcceptCallReq) (*calldto.AcceptCallRes, error) {
	return call.AcceptCall(ctx, req)
}

// ConfirmCall 通话应答确认
func (c *CallAppController) ConfirmCall(ctx context.Context, req *calldto.ConfirmCallReq) (*calldto.ConfirmCallRes, error) {
	return call.ConfirmCall(ctx, req)
}

// GetCallConfirmStatus 查询对方应答确认状态
func (c *CallAppController) GetCallConfirmStatus(ctx context.Context, req *calldto.GetCallConfirmStatusReq) (*calldto.GetCallConfirmStatusRes, error) {
	return call.GetCallConfirmStatus(ctx, req)
}

// CallTimeout 呼叫超时
func (c *CallAppController) CallTimeout(ctx context.Context, req *calldto.CallTimeoutReq) (*calldto.CallTimeoutRes, error) {
	return call.CallTimeout(ctx, req)
}

// CallHeart 通话心跳
func (c *CallAppController) CallHeart(ctx context.Context, req *calldto.CallHeartReq) (*calldto.CallHeartRes, error) {
	return call.CallHeart(ctx, req)
}

// EndCall 结束通话
func (c *CallAppController) EndCall(ctx context.Context, req *calldto.EndCallReq) (*calldto.EndCallRes, error) {
	return call.EndCall(ctx, req)
}
