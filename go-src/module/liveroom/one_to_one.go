package liveroom

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func IsOneToOneOnShelf(anchorId uint64) bool {
	return liveroomdao.IsOneToOneOnShelf(anchorId)
}

func enrollOneToOneRoom(anchorId uint64, billing float64) (*liveentity.OneToOneRoom, error) {
	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil || !user.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	if existing := liveroomdao.GetOneToOneRoom(anchorId); existing != nil {
		return nil, errercode.CreateCode(errercode.OneToOneRoomExist)
	}
	EnsureAnchorRoom(anchorId, liveroomdao.GetAnchorGuildId(anchorId))
	row := liveroomdao.CreateOneToOneRoom(anchorId, billing)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return row, nil
}

func CreateOneToOneRoom(ctx context.Context, req *liveroomdto.CreateOneToOneRoomReq) (*liveroomdto.CreateOneToOneRoomRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	billing := float64(0)
	if req != nil {
		billing = req.Billing
	}
	row, err := enrollOneToOneRoom(anchorId, billing)
	if err != nil {
		return nil, err
	}
	return &liveroomdto.CreateOneToOneRoomRes{
		RoomId: strconv.FormatUint(row.ID, 10),
	}, nil
}

func GetCMSOneToOneRoomList(_ context.Context, req *liveroomdto.CMSOneToOneRoomListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return httpserver.NewCMSQueryResp(0, []*liveroomdto.CMSOneToOneRoomItem{}), nil
	}
	key := strings.ToLower(strings.TrimSpace(req.Key))
	items := make([]*liveroomdto.CMSOneToOneRoomItem, 0)
	for _, row := range liveroomdao.GetAllOneToOneRooms() {
		if row == nil {
			continue
		}
		if req.Status != nil && row.Status != *req.Status {
			continue
		}
		item := buildCMSOneToOneRoomItem(row)
		if key != "" && !matchCMSOneToOneRoomItem(item, key) {
			continue
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		ai, aj := items[i], items[j]
		if ai == nil || ai.UpdatedAt == nil {
			return false
		}
		if aj == nil || aj.UpdatedAt == nil {
			return true
		}
		return ai.UpdatedAt.Unix() > aj.UpdatedAt.Unix()
	})
	total := len(items)
	start := req.PageOffset()
	if start > total {
		start = total
	}
	end := start + req.PageSize
	if req.PageSize <= 0 {
		end = total
	}
	if end > total {
		end = total
	}
	return httpserver.NewCMSQueryResp(total, items[start:end]), nil
}

func CMSCreateOneToOneRoom(_ context.Context, req *liveroomdto.CMSCreateOneToOneRoomReq) (*liveroomdto.CMSCreateOneToOneRoomRes, error) {
	if req == nil || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row, err := enrollOneToOneRoom(req.UserId, req.Billing)
	if err != nil {
		return nil, err
	}
	return &liveroomdto.CMSCreateOneToOneRoomRes{
		UserId: strconv.FormatUint(row.ID, 10),
	}, nil
}

func CMSUpdateOneToOneRoom(_ context.Context, req *liveroomdto.CMSUpdateOneToOneRoomReq) (*liveroomdto.CMSUpdateOneToOneRoomRes, error) {
	if req == nil || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := liveroomdao.GetOneToOneRoom(req.UserId)
	if row == nil {
		return nil, errercode.CreateCode(errercode.OneToOneRoomNonExist)
	}
	row.SetBilling(req.Billing)
	return &liveroomdto.CMSUpdateOneToOneRoomRes{}, nil
}

func CMSSetOneToOneRoomStatus(_ context.Context, req *liveroomdto.CMSSetOneToOneRoomStatusReq) (*liveroomdto.CMSSetOneToOneRoomStatusRes, error) {
	if req == nil || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := liveroomdao.GetOneToOneRoom(req.UserId)
	if row == nil {
		return nil, errercode.CreateCode(errercode.OneToOneRoomNonExist)
	}
	if row.Status != req.Status {
		row.SetStatus(req.Status)
	}
	return &liveroomdto.CMSSetOneToOneRoomStatusRes{}, nil
}

func buildCMSOneToOneRoomItem(row *liveentity.OneToOneRoom) *liveroomdto.CMSOneToOneRoomItem {
	item := &liveroomdto.CMSOneToOneRoomItem{
		UserId: strconv.FormatUint(row.ID, 10),
		Status: row.Status,
	}
	updated := row.UpdatedAt
	if !updated.IsZero() {
		item.UpdatedAt = &updated
	}
	user := userinfodao.GetUserInfoByUserId(row.ID)
	if user != nil {
		item.Nickname = user.Nickname
		item.Phone = user.Phone
		item.Avatar = upload.ResolveAvatarUrlForUser(row.ID, user.Avatar)
	}
	room := liveroomdao.ResolveRoom(row.ID)
	if room != nil {
		item.GuildId = strconv.FormatUint(room.GuildId, 10)
		item.RoomTitle = room.Title
		item.LiveRoomStatus = room.Status
	}
	item.Billing = row.Billing
	return item
}

func matchCMSOneToOneRoomItem(item *liveroomdto.CMSOneToOneRoomItem, key string) bool {
	if item == nil {
		return false
	}
	if strings.Contains(strings.ToLower(item.UserId), key) {
		return true
	}
	if strings.Contains(strings.ToLower(item.Nickname), key) {
		return true
	}
	return strings.Contains(strings.ToLower(item.Phone), key)
}