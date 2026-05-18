package liveroom

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/globalcfg"
)

// CreateRoom 创建直播间
// 业务规则:
//  1. 调用者必须已通过工会审批(UserInfo.GuildId != 0),即已成为主播
//  2. 同一主播只能拥有一个直播间(再次调用直接返回已有信息)
func CreateRoom(ctx context.Context, req *liveroomdto.CreateLiveRoomReq) (res *liveroomdto.CreateLiveRoomRes, err error) {
	anchorId := httpserver.GetAuthId(ctx)

	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil || user.GuildId == 0 {
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
		Status:  room.Status,
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
	if room.Status != entity.LiveRoomStatusLive {
		room.SetStatus(entity.LiveRoomStatusLive)
		//liveroomdao.FlushRoomCache(room)
	}
	return &liveroomdto.StartLiveRes{Status: room.Status}, nil
}

// StopLive 下播
func StopLive(ctx context.Context, _ *liveroomdto.StopLiveReq) (*liveroomdto.StopLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.Status != entity.LiveRoomStatusClosed {
		room.SetStatus(entity.LiveRoomStatusClosed)
		//liveroomdao.FlushRoomCache(room)
	}
	return &liveroomdto.StopLiveRes{Status: room.Status}, nil
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

// GetRoom 查询直播间(公开接口,任意登录用户可调用)
func GetRoom(_ context.Context, req *liveroomdto.GetLiveRoomReq) (*liveroomdto.GetLiveRoomRes, error) {
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	return &liveroomdto.GetLiveRoomRes{
		RoomId:   strconv.FormatUint(room.ID, 10),
		GuildId:  strconv.FormatUint(room.GuildId, 10),
		Title:    room.Title,
		Cover:    globalcfg.BuildResourceUrl(room.Cover),
		Notice:   room.Notice,
		Status:   room.Status,
		CreateAt: room.CreatedAt.Unix(),
	}, nil
}
