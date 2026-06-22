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

var (
	commonOnlineMap = gmap.NewKVMap[uint64, []uint64](true)
	vipOnlineMap    = gmap.NewKVMap[uint64, []uint64](true)
)

type onlineSortKey struct {
	userId      uint64
	totalReward float64
	vipLevel    uint32
}

func buildOnlineSortKeys(roomId uint64, userIds []uint64) []onlineSortKey {
	keys := make([]onlineSortKey, 0, len(userIds))
	for _, userId := range userIds {
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
	return keys
}

func sortOnlineUserIds(keys []onlineSortKey) []uint64 {
	if len(keys) == 0 {
		return nil
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

	sorted := make([]uint64, 0, len(keys))
	for _, key := range rewarded {
		sorted = append(sorted, key.userId)
	}
	for _, key := range noReward {
		sorted = append(sorted, key.userId)
	}
	return sorted
}

func flushCommonOnlineMap(roomId uint64) {
	all := getOnline(roomId)
	commonOnlineMap.Set(roomId, sortOnlineUserIds(buildOnlineSortKeys(roomId, all)))
}

func flushVipOnlineMap(roomId uint64) {
	all := getOnline(roomId)
	keys := buildOnlineSortKeys(roomId, all)
	vipKeys := make([]onlineSortKey, 0, len(keys))
	for _, key := range keys {
		if key.vipLevel > 0 {
			vipKeys = append(vipKeys, key)
		}
	}
	vipOnlineMap.Set(roomId, sortOnlineUserIds(vipKeys))
}

func flushOnlineLists(roomId uint64) {
	flushCommonOnlineMap(roomId)
	flushVipOnlineMap(roomId)
}

func clearOnlineLists(roomId uint64) {
	commonOnlineMap.Remove(roomId)
	vipOnlineMap.Remove(roomId)
}

func normalizeOnlineListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func buildOnlineUserItems(roomId uint64, userIds []uint64) []*liveroomdto.OnlineUserItem {
	list := make([]*liveroomdto.OnlineUserItem, 0, len(userIds))
	for _, userId := range userIds {
		onlineId := entity.BuildLiveRoomOnlineId(userId, roomId)
		onlineData := liveroomdao.GetOnlineById(onlineId, userId, roomId)
		joinTime := ""
		var joinedUnix int64
		if onlineData != nil && onlineData.JoinTime != nil {
			joinTime = onlineData.JoinTime.Format("2006-01-02 15:04:05")
			joinedUnix = onlineData.JoinTime.UnixMilli()
		}
		item := &liveroomdto.OnlineUserItem{
			UserId:     strconv.FormatUint(userId, 10),
			JoinedAt:   joinTime,
			JoinedUnix: joinedUnix,
		}
		if onlineData != nil {
			item.Muted = onlineData.Muted
			item.TotalReward = onlineData.TotalReward
		}
		if u := userinfodao.GetUserInfoByUserId(userId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrl(u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}
	return list
}

func queryOnlineUserListPage(roomId uint64, page, pageSize int, cached *gmap.KVMap[uint64, []uint64]) *liveroomdto.GetOnlineUserListRes {
	page, pageSize = normalizeOnlineListPage(page, pageSize)

	all := cached.Get(roomId)
	if all == nil {
		all = make([]uint64, 0)
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

	return &liveroomdto.GetOnlineUserListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     buildOnlineUserItems(roomId, all[start:end]),
		SysTime:  time.Now().UnixMilli(),
	}
}

// GetOnlineUserList 分页查询直播间在线玩家(附用户基础信息)
func GetOnlineUserList(_ context.Context, req *liveroomdto.GetOnlineUserListReq) (*liveroomdto.GetOnlineUserListRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	return queryOnlineUserListPage(req.RoomId, req.Page, req.PageSize, commonOnlineMap), nil
}
