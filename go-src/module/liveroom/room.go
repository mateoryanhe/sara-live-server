package liveroom

import (
	"context"
	"strconv"
	"xr-game-server/constants/roomstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// CreateRoom 创建直播间
// 业务规则:
//  1. 调用者必须已是主播(UserInfo.IsAnchor == true)
//  2. 同一主播只能拥有一个直播间(再次调用直接返回已有信息)
func CreateRoom(ctx context.Context, req *liveroomdto.CreateLiveRoomReq) (res *liveroomdto.CreateLiveRoomRes, err error) {
	anchorId := httpserver.GetAuthId(ctx)

	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil || !user.IsAnchor {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}

	// 同一主播仅允许一个直播间(roomId == anchorId)
	if existing := liveroomdao.GetRoomById(anchorId); existing != nil {
		return nil, errercode.CreateCode(errercode.LiveRoomExist)
	}

	// 通过 syndb 异步入库,不直接 INSERT;LiveRoom.ID 复用主播用户ID
	room := entity.NewLiveRoom(
		anchorId,
		user.GuildId,
		req.Title,
		req.Cover,
		req.Notice,
	)
	liveroomdao.AddRoomToCache(room)

	return &liveroomdto.CreateLiveRoomRes{
		RoomId:  strconv.FormatUint(room.ID, 10),
		GuildId: strconv.FormatUint(room.GuildId, 10),
	}, nil
}

// loadOwnRoom 获取调用者(主播)自己的直播间;不存在则返回 LiveRoomNotExist
func loadOwnRoom(ctx context.Context) (*entity.LiveRoom, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	return room, nil
}

// StartLive 开播
func StartLive(ctx context.Context, _ *liveroomdto.StartLiveReq) (*liveroomdto.StartLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.LiveRecordId > 0 {
		//重置
	}
	liveRecordId := snowflake.GetId()
	liveroomdao.GetLiveRecordById(liveRecordId)
	room.SetLiveRecordId(liveRecordId)
	return &liveroomdto.StartLiveRes{}, nil
}

// StopLive 下播
func StopLive(ctx context.Context, _ *liveroomdto.StopLiveReq) (*liveroomdto.StopLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}

	room.SetLiveRecordId(0)

	return &liveroomdto.StopLiveRes{}, nil
}

// UpdateCover 修改封面
func UpdateCover(ctx context.Context, req *liveroomdto.UpdateCoverReq) (*liveroomdto.UpdateCoverRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.Cover != req.Cover {
		room.SetCover(req.Cover)
	}
	return &liveroomdto.UpdateCoverRes{Success: true}, nil
}

// UpdateNotice 修改公告
func UpdateNotice(ctx context.Context, req *liveroomdto.UpdateNoticeReq) (*liveroomdto.UpdateNoticeRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.Notice != req.Notice {
		room.SetNotice(req.Notice)
		//liveroomdao.FlushRoomCache(room)
	}
	return &liveroomdto.UpdateNoticeRes{Success: true}, nil
}

// JoinRoom 玩家加入直播间,记录状态置为 Online
func JoinRoom(ctx context.Context, req *liveroomdto.JoinRoomReq) (*liveroomdto.JoinRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)

	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(userId, req.RoomId)
	existing := liveroomdao.GetOnlineById(onlineId)
	if existing == nil {
		// 历史无记录,首次加入
		o := entity.NewLiveRoomOnline(userId, req.RoomId)
		liveroomdao.AddOnlineToRoomCache(o)
	} else if existing.Status != entity.LiveRoomOnlineStatusOnline {
		// 历史有记录(此前离开过),状态切回 Online
		existing.SetStatus(entity.LiveRoomOnlineStatusOnline)
		liveroomdao.AddOnlineToRoomCache(existing)
	}

	return &liveroomdto.JoinRoomRes{
		OnlineId:    onlineId,
		OnlineCount: len(liveroomdao.GetOnlinesByRoom(req.RoomId)),
	}, nil
}

// LeaveRoom 玩家离开直播间,记录状态置为 Offline(保留行)
func LeaveRoom(ctx context.Context, req *liveroomdto.LeaveRoomReq) (*liveroomdto.LeaveRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)
	onlineId := entity.BuildLiveRoomOnlineId(userId, req.RoomId)

	if existing := liveroomdao.GetOnlineById(onlineId); existing != nil &&
		existing.Status != entity.LiveRoomOnlineStatusOffline {
		existing.SetStatus(entity.LiveRoomOnlineStatusOffline)
		liveroomdao.RemoveOnlineFromRoomCache(existing)
	}

	return &liveroomdto.LeaveRoomRes{
		OnlineCount: len(liveroomdao.GetOnlinesByRoom(req.RoomId)),
	}, nil
}

// GetOnlineUserList 分页查询直播间在线玩家(附用户基础信息)
func GetOnlineUserList(_ context.Context, req *liveroomdto.GetOnlineUserListReq) (*liveroomdto.GetOnlineUserListRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	all := liveroomdao.GetOnlinesByRoom(req.RoomId)
	total := len(all)

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	pageData := all[start:end]

	list := make([]*liveroomdto.OnlineUserItem, 0, len(pageData))
	for _, o := range pageData {
		item := &liveroomdto.OnlineUserItem{
			UserId:   strconv.FormatUint(o.UserId, 10),
			JoinedAt: o.UpdatedAt.Unix(),
		}
		if u := userinfodao.GetUserInfoByUserId(o.UserId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.GetUrlByName(u.Avatar)
		}
		list = append(list, item)
	}

	return &liveroomdto.GetOnlineUserListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// GetRoom 查询直播间(公开接口,任意登录用户可调用)
func GetRoom(_ context.Context, req *liveroomdto.GetLiveRoomReq) (*liveroomdto.GetLiveRoomRes, error) {
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	status := roomstatus.LiveRoomStatusClosed
	if room.LiveRecordId > 0 {
		status = roomstatus.LiveRoomStatusLive
	}

	return &liveroomdto.GetLiveRoomRes{
		RoomId:   strconv.FormatUint(room.ID, 10),
		GuildId:  strconv.FormatUint(room.GuildId, 10),
		Title:    room.Title,
		Cover:    upload.GetUrlByName(room.Cover),
		Notice:   room.Notice,
		Status:   status,
		CreateAt: room.CreatedAt.Unix(),
	}, nil
}
