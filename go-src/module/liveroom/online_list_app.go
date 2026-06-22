package liveroom

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

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

	all := getOnline(req.RoomId)
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
		onlineId := entity.BuildLiveRoomOnlineId(o, req.RoomId)
		onlineData := liveroomdao.GetOnlineById(onlineId, o, req.RoomId)
		joinTime := ""
		if onlineData.JoinTime != nil {
			joinTime = onlineData.JoinTime.Format("2006-01-02 15:04:05")
		}
		item := &liveroomdto.OnlineUserItem{
			UserId:     strconv.FormatUint(o, 10),
			JoinedAt:   joinTime,
			JoinedUnix: onlineData.JoinTime.UnixMilli(),
			Muted:      onlineData.Muted,
		}
		if u := userinfodao.GetUserInfoByUserId(o); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrl(u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}

	return &liveroomdto.GetOnlineUserListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
		SysTime:  time.Now().UnixMilli(),
	}, nil
}
