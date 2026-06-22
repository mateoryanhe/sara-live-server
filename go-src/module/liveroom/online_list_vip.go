package liveroom

import (
	"context"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
)

// GetVipOnlineUserList 分页查询直播间VIP在线玩家(VIP等级>0)
func GetVipOnlineUserList(_ context.Context, req *liveroomdto.GetVipOnlineUserListReq) (*liveroomdto.GetVipOnlineUserListRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	return queryOnlineUserListPage(req.RoomId, req.Page, req.PageSize, vipOnlineMap), nil
}
