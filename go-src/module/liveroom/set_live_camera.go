package liveroom

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
)

// SetLiveCameraEnabled 主播开关摄像头;关闭时不累计有效直播时长(心跳仍保活).
func SetLiveCameraEnabled(ctx context.Context, req *liveroomdto.SetLiveCameraEnabledReq) (*liveroomdto.SetLiveCameraEnabledRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 || req.RoomId == 0 || req.RoomId != userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}
	setLiveCameraEnabled(room.ID, req.Enabled)
	return &liveroomdto.SetLiveCameraEnabledRes{
		Success:       true,
		CameraEnabled: req.Enabled,
	}, nil
}
