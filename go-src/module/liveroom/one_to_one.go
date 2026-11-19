package liveroom

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/common"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/upload"
)

func IsOneToOneOnShelf(anchorId uint64) bool {
	return liveroomdao.IsOneToOneOnShelf(anchorId)
}

// OneToOneVideoCallBillingPerMinute 1v1 视频通话(source=1、3)每分钟钻石价，统一取自 one_to_one_rooms.billing。
func OneToOneVideoCallBillingPerMinute(anchorId uint64) float64 {
	row := liveroomdao.GetOneToOneRoom(anchorId)
	if row == nil {
		return 0
	}
	return row.Billing
}

// overlayOneToOneListItem 用 1v1 独立资料覆盖列表中的展示字段(计费仍以 one_to_one_rooms 为准)。
func overlayOneToOneListItem(item *liveroomdto.LiveRoomListItem, one *liveentity.OneToOneRoom) {
	if item == nil || one == nil {
		return
	}
	item.Billing = one.Billing
	if one.Title != "" {
		item.Title = one.Title
	}
	if one.Cover != "" {
		item.Cover = upload.GetUrlByName(one.Cover)
	}
	if one.TagId > 0 {
		item.TagId = strconv.FormatUint(one.TagId, 10)
		item.TagName = getRoomTagName(one.TagId)
	} else {
		item.TagId = ""
		item.TagName = ""
	}
}

type oneToOneRoomProfile struct {
	Billing float64
	Title   string
	Cover   string
	TagId   uint64
}

func applyOneToOneRoomProfile(row *liveentity.OneToOneRoom, profile oneToOneRoomProfile, isNew bool) {
	if row == nil {
		return
	}
	row.SetBilling(profile.Billing)
	if isNew || profile.Title != "" {
		row.SetTitle(profile.Title)
	}
	if isNew || profile.Cover != "" {
		row.SetCover(profile.Cover)
	}
	if isNew || profile.TagId > 0 {
		row.SetTagId(profile.TagId)
	}
}

func profileFromCreateReq(req *liveroomdto.CreateOneToOneRoomReq) oneToOneRoomProfile {
	if req == nil {
		return oneToOneRoomProfile{}
	}
	return oneToOneRoomProfile{
		Billing: req.Billing,
		Title:   req.Title,
		Cover:   req.Cover,
		TagId:   req.TagId,
	}
}

func profileFromCMSCreateReq(req *liveroomdto.CMSCreateOneToOneRoomReq) oneToOneRoomProfile {
	if req == nil {
		return oneToOneRoomProfile{}
	}
	return oneToOneRoomProfile{
		Billing: req.Billing,
		Title:   req.Title,
		Cover:   req.Cover,
		TagId:   req.TagId,
	}
}

func profileFromCMSUpdateReq(req *liveroomdto.CMSUpdateOneToOneRoomReq) oneToOneRoomProfile {
	if req == nil {
		return oneToOneRoomProfile{}
	}
	return oneToOneRoomProfile{
		Billing: req.Billing,
		Title:   req.Title,
		Cover:   req.Cover,
		TagId:   req.TagId,
	}
}

// maybeAutoOpenOneToOneFromCreateLiveRoom 完善直播间资料时:通话分钟价>0且接受所有人,且未开通1v1则自动开通并上架;已上架则仅同步计费.
func maybeAutoOpenOneToOneFromCreateLiveRoom(ctx context.Context, anchorId uint64, req *liveroomdto.CreateLiveRoomReq, coverName string, room *liveentity.LiveRoom, cfg *liveentity.LiveRoomCfg) {
	if req == nil || anchorId == 0 || room == nil || cfg == nil {
		return
	}
	if !livecfg.IsOneToOneLiveRoomEnabled() {
		return
	}
	if req.Billing <= common.Zero {
		return
	}
	if cfg.PrivateInviteType != liveentity.LiveRoomPrivateInviteAll {
		return
	}
	profile := oneToOneRoomProfile{
		Billing: req.Billing,
		Title:   room.Title,
		TagId:   cfg.TagId,
	}
	if coverName != "" {
		profile.Cover = coverName
	} else if room.Cover != "" {
		profile.Cover = room.Cover
	}
	existing := liveroomdao.GetOneToOneRoom(anchorId)
	if existing == nil {
		if _, err := enrollOneToOneRoom(anchorId, profile); err != nil {
			g.Log().Warningf(ctx, "auto open one_to_one from createRoom anchorId=%d err=%v", anchorId, err)
		}
		return
	}
	if existing.IsOnShelf() && existing.Billing != profile.Billing {
		existing.SetBilling(profile.Billing)
	}
}

func enrollOneToOneRoom(anchorId uint64, profile oneToOneRoomProfile) (*liveentity.OneToOneRoom, error) {
	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil || !user.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	if existing := liveroomdao.GetOneToOneRoom(anchorId); existing != nil {
		return nil, errercode.CreateCode(errercode.OneToOneRoomExist)
	}
	EnsureAnchorRoom(anchorId, liveroomdao.GetAnchorGuildId(anchorId))
	row := liveroomdao.CreateOneToOneRoom(anchorId, profile.Billing)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	applyOneToOneRoomProfile(row, profile, true)
	return row, nil
}

func CreateOneToOneRoom(ctx context.Context, req *liveroomdto.CreateOneToOneRoomReq) (*liveroomdto.CreateOneToOneRoomRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	profile := profileFromCreateReq(req)
	if existing := liveroomdao.GetOneToOneRoom(anchorId); existing != nil {
		applyOneToOneRoomProfile(existing, profile, false)
		return &liveroomdto.CreateOneToOneRoomRes{
			RoomId: strconv.FormatUint(existing.ID, 10),
		}, nil
	}
	row, err := enrollOneToOneRoom(anchorId, profile)
	if err != nil {
		return nil, err
	}
	return &liveroomdto.CreateOneToOneRoomRes{
		RoomId: strconv.FormatUint(row.ID, 10),
	}, nil
}

// GetOneToOneRoom App 按主播ID查询 one_to_one_rooms 配置。
func GetOneToOneRoom(_ context.Context, req *liveroomdto.GetOneToOneRoomReq) (*liveroomdto.GetOneToOneRoomRes, error) {
	if req == nil || req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := liveroomdao.GetOneToOneRoom(req.AnchorId)
	if row == nil {
		return nil, errercode.CreateCode(errercode.OneToOneRoomNonExist)
	}
	onlineStatus := uint8(liveroomdto.OneToOneRoomStatusFilterOffline)
	if push.IsOnline(req.AnchorId) {
		onlineStatus = liveroomdto.OneToOneRoomStatusFilterOnline
	}
	res := &liveroomdto.GetOneToOneRoomRes{
		RoomId:       strconv.FormatUint(row.ID, 10),
		Status:       row.Status,
		Billing:      row.Billing,
		Title:        row.Title,
		OnlineStatus: onlineStatus,
	}
	if row.Cover != "" {
		res.Cover = upload.GetUrlByName(row.Cover)
	}
	if row.TagId > 0 {
		res.TagId = strconv.FormatUint(row.TagId, 10)
		res.TagName = getRoomTagName(row.TagId)
	}
	return res, nil
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
	row, err := enrollOneToOneRoom(req.UserId, profileFromCMSCreateReq(req))
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
	applyOneToOneRoomProfile(row, profileFromCMSUpdateReq(req), false)
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
		UserId:   strconv.FormatUint(row.ID, 10),
		Status:   row.Status,
		Billing:  row.Billing,
		Title:   row.Title,
	}
	if row.Cover != "" {
		item.Cover = upload.GetUrlByName(row.Cover)
	}
	if row.TagId > 0 {
		item.TagId = strconv.FormatUint(row.TagId, 10)
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
		item.LiveRoomStatus = room.Status
	}
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
	if strings.Contains(strings.ToLower(item.Title), key) {
		return true
	}
	return strings.Contains(strings.ToLower(item.Phone), key)
}