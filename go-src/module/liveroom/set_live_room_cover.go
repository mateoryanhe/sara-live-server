package liveroom

import (
	"context"
	"strings"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// SetLiveRoomCover CMS 设置直播间封面(cover 为 uploadFile 返回的文件名,空字符串表示清除)
func SetLiveRoomCover(_ context.Context, req *accountdto.SetLiveRoomCoverReq) (*accountdto.SetLiveRoomCoverRes, error) {
	if req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	room := liveroomdao.ResolveRoom(req.AnchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	cover := strings.TrimSpace(req.Cover)
	room.SetCover(cover)
	liveroomdao.FlushRoomCache(room)
	return &accountdto.SetLiveRoomCoverRes{
		Success: true,
		Cover:   upload.GetUrlByName(room.Cover),
	}, nil
}
