package liveroom

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/push"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

const (
	serverOnlineNormalUserRefreshInterval = 5 * time.Minute
	serverOnlineNormalUserMax             = 500
	serverOnlineNormalUserLoadBatch       = 500
)

var serverOnlineNormalUserSnapshot = struct {
	sync.RWMutex
	userIds     []uint64
	refreshedAt int64
}{
	userIds: make([]uint64, 0),
}

func initServerOnlineNormalUserList() {
	refreshServerOnlineNormalUserList(gctx.New())
	xrtimer.AddSingleton(gctx.New(), serverOnlineNormalUserRefreshInterval, refreshServerOnlineNormalUserList)
}

// refreshServerOnlineNormalUserList 从当前进程 WebSocket 在线用户中筛选普通用户。
// 用户资料缓存未命中的 ID 按批量查询数据库并回填，最终只保存最多500个用户ID。
func refreshServerOnlineNormalUserList(_ context.Context) {
	onlineUserIds := push.OnlineUserIds()

	normalUserIds := make([]uint64, 0, min(len(onlineUserIds), serverOnlineNormalUserMax))
	for start := 0; start < len(onlineUserIds) && len(normalUserIds) < serverOnlineNormalUserMax; start += serverOnlineNormalUserLoadBatch {
		end := start + serverOnlineNormalUserLoadBatch
		if end > len(onlineUserIds) {
			end = len(onlineUserIds)
		}
		batch := onlineUserIds[start:end]
		userInfoMap := userinfodao.EnsureUserInfosCached(batch)
		for _, userId := range batch {
			user := userInfoMap[userId]
			if user == nil || user.UserType != userentity.UserTypeNormal {
				continue
			}
			normalUserIds = append(normalUserIds, userId)
			if len(normalUserIds) == serverOnlineNormalUserMax {
				break
			}
		}
	}

	serverOnlineNormalUserSnapshot.Lock()
	serverOnlineNormalUserSnapshot.userIds = normalUserIds
	serverOnlineNormalUserSnapshot.refreshedAt = time.Now().UnixMilli()
	serverOnlineNormalUserSnapshot.Unlock()
}

func getServerOnlineNormalUserSnapshot() ([]uint64, int64) {
	serverOnlineNormalUserSnapshot.RLock()
	defer serverOnlineNormalUserSnapshot.RUnlock()
	userIds := append([]uint64(nil), serverOnlineNormalUserSnapshot.userIds...)
	return userIds, serverOnlineNormalUserSnapshot.refreshedAt
}

// GetServerOnlineNormalUserList 返回当前服务器每5分钟生成的在线普通用户快照。
func GetServerOnlineNormalUserList(_ context.Context, _ *liveroomdto.GetServerOnlineNormalUserListReq) (*liveroomdto.GetServerOnlineNormalUserListRes, error) {
	userIds, refreshedAt := getServerOnlineNormalUserSnapshot()
	list := make([]*liveroomdto.ServerOnlineNormalUserItem, 0, len(userIds))
	for _, userId := range userIds {
		user := userinfodao.GetUserInfoFromMemory(userId)
		if user == nil || user.UserType != userentity.UserTypeNormal {
			continue
		}
		list = append(list, &liveroomdto.ServerOnlineNormalUserItem{
			UserId:   strconv.FormatUint(userId, 10),
			Nickname: user.Nickname,
			Avatar:   upload.ResolveAvatarUrlForUser(userId, user.Avatar),
			VipLevel: user.VipLevel,
			Gender:   user.Gender,
			Age:      calcAge(user.Birthday),
			UserType: user.UserType,
		})
	}
	now := time.Now().UnixMilli()
	return &liveroomdto.GetServerOnlineNormalUserListRes{
		Total:       len(list),
		List:        list,
		RefreshedAt: refreshedAt,
		SysTime:     now,
	}, nil
}
