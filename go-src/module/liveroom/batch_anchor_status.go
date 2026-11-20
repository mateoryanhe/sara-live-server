package liveroom

import (
	"context"
	"strconv"

	"xr-game-server/constants/userstatus"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
)

const batchAnchorStatusMax = 100

// BatchAnchorStatus 按主播ID列表返回开播与 WebSocket 在线状态.
func BatchAnchorStatus(_ context.Context, req *liveroomdto.BatchAnchorStatusReq) (*liveroomdto.BatchAnchorStatusRes, error) {
	if len(req.AnchorIds) > batchAnchorStatusMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := dedupeAnchorIdsInOrder(req.AnchorIds)
	list := make([]*liveroomdto.AnchorStatusItem, 0, len(ids))
	for _, anchorId := range ids {
		list = append(list, buildAnchorStatusItem(anchorId))
	}
	return &liveroomdto.BatchAnchorStatusRes{List: list}, nil
}

func dedupeAnchorIdsInOrder(anchorIds []uint64) []uint64 {
	if len(anchorIds) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(anchorIds))
	out := make([]uint64, 0, len(anchorIds))
	for _, id := range anchorIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func buildAnchorStatusItem(anchorId uint64) *liveroomdto.AnchorStatusItem {
	item := &liveroomdto.AnchorStatusItem{
		AnchorId:     strconv.FormatUint(anchorId, 10),
		Status:       userstatus.LiveRoomStatusClosed,
		OnlineStatus: liveroomdto.OneToOneRoomStatusFilterOffline,
	}
	if push.IsOnline(anchorId) {
		item.OnlineStatus = liveroomdto.OneToOneRoomStatusFilterOnline
	}
	room := liveroomdao.ResolveRoom(anchorId)
	if room != nil && room.LiveRecordId > 0 {
		item.Status = userstatus.LiveRoomStatusLive
	}
	return item
}