package appswitch

import (
	"context"

	"xr-game-server/dto/appswitchdto"
	"xr-game-server/module/livecfg"
)

func GetAppSwitches(_ context.Context, _ *appswitchdto.GetAppSwitchesReq) (*appswitchdto.GetAppSwitchesRes, error) {
	return &appswitchdto.GetAppSwitchesRes{
		PersonalAnchorAlbumEnabled: livecfg.IsPersonalAnchorAlbumEnabled(),
	}, nil
}