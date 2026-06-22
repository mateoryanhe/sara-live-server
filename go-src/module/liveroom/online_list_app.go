package liveroom

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

var commonOnlineMap = gmap.NewKVMap[uint64, []uint64](true)

func flushCommonOnlineMap(roomId uint64) {
	all := getOnline(roomId)
	if len(all) == 0 {
		commonOnlineMap.Set(roomId, all)
		return
	}

	type onlineSortKey struct {
		userId      uint64
		totalReward float64
		vipLevel    uint32
	}

	keys := make([]onlineSortKey, 0, len(all))
	for _, userId := range all {
		onlineId := entity.BuildLiveRoomOnlineId(userId, roomId)
		online := liveroomdao.GetOnlineById(onlineId, userId, roomId)
		key := onlineSortKey{userId: userId}
		if online != nil {
			key.totalReward = online.TotalReward
		}
		if u := userinfodao.GetUserInfoByUserId(userId); u != nil {
			key.vipLevel = u.VipLevel
		}
		keys = append(keys, key)
	}

	rewarded := make([]onlineSortKey, 0, len(keys))
	noReward := make([]onlineSortKey, 0, len(keys))
	for _, key := range keys {
		if key.totalReward > 0 {
			rewarded = append(rewarded, key)
		} else {
			noReward = append(noReward, key)
		}
	}

	sort.Slice(rewarded, func(i, j int) bool {
		if rewarded[i].totalReward != rewarded[j].totalReward {
			return rewarded[i].totalReward > rewarded[j].totalReward
		}
		return rewarded[i].userId < rewarded[j].userId
	})
	sort.Slice(noReward, func(i, j int) bool {
		if noReward[i].vipLevel != noReward[j].vipLevel {
			return noReward[i].vipLevel > noReward[j].vipLevel
		}
		return noReward[i].userId < noReward[j].userId
	})

	sorted := make([]uint64, 0, len(all))
	for _, key := range rewarded {
		sorted = append(sorted, key.userId)
	}
	for _, key := range noReward {
		sorted = append(sorted, key.userId)
	}
	commonOnlineMap.Set(roomId, sorted)
}

func clearCommonOnlineMap(roomId uint64) {
	commonOnlineMap.Remove(roomId)
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

	all := commonOnlineMap.Get(req.RoomId)
	if all == nil {
		flushCommonOnlineMap(req.RoomId)
		all = commonOnlineMap.Get(req.RoomId)
	}
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
