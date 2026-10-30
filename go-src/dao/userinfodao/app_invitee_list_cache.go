package userinfodao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	userentity "xr-game-server/entity/user"
)

const (
	// AppInviteeListCachePageSize App 邀请人列表固定分页大小
	AppInviteeListCachePageSize = 20
	appInviteeListCachedPages   = 2
	appInviteeListCacheMaxSize  = AppInviteeListCachePageSize * appInviteeListCachedPages
)

// AppInviteeListRow App 邀请人列表缓存行(被邀请人维度)
type AppInviteeListRow struct {
	UserId            uint64    `orm:"user_id"`
	Nickname          string    `orm:"nickname"`
	Avatar            string    `orm:"avatar"`
	InviterRewardGold float64   `orm:"inviter_reward_gold"`
	CreatedAt         time.Time `orm:"created_at"`
}

var appInviteeListCacheMgr *cache.ListCache[*AppInviteeListRow]

func initAppInviteeListCache() {
	appInviteeListCacheMgr = cache.NewListCache[*AppInviteeListRow]()
}

func appInviteeListCacheKey(inviterId uint64) string {
	return fmt.Sprintf("user_app_invitee_list:%d", inviterId)
}

func getAppInviteeListCache(inviterId uint64) []*AppInviteeListRow {
	if appInviteeListCacheMgr == nil || inviterId == 0 {
		return make([]*AppInviteeListRow, 0)
	}
	list := appInviteeListCacheMgr.MustGetList(gctx.New(), appInviteeListCacheKey(inviterId), func(ctx context.Context) ([]*AppInviteeListRow, error) {
		return loadAppInviteesFromDB(inviterId, 1, appInviteeListCacheMaxSize), nil
	})
	if list == nil {
		return make([]*AppInviteeListRow, 0)
	}
	return list
}

func putAppInviteeListCache(inviterId uint64, list []*AppInviteeListRow) {
	if appInviteeListCacheMgr == nil || inviterId == 0 {
		return
	}
	if list == nil {
		list = make([]*AppInviteeListRow, 0)
	}
	appInviteeListCacheMgr.PublishList(gctx.New(), appInviteeListCacheKey(inviterId), list)
}

func newAppInviteeListRow(info *userentity.UserInfo, rewardGold float64) *AppInviteeListRow {
	if info == nil || info.ID == 0 {
		return nil
	}
	return &AppInviteeListRow{
		UserId:            info.ID,
		Nickname:          info.Nickname,
		Avatar:            info.Avatar,
		InviterRewardGold: rewardGold,
		CreatedAt:         info.CreatedAt,
	}
}

// PrependAppInviteeToListCache 新绑定邀请关系后写入列表缓存头部
func PrependAppInviteeToListCache(inviterId uint64, invitee *userentity.UserInfo, rewardGold float64) {
	if appInviteeListCacheMgr == nil || inviterId == 0 {
		return
	}
	row := newAppInviteeListRow(invitee, rewardGold)
	if row == nil {
		return
	}
	list := getAppInviteeListCache(inviterId)
	newList := make([]*AppInviteeListRow, 0, len(list)+1)
	newList = append(newList, row)
	for _, item := range list {
		if item != nil && item.UserId != row.UserId {
			newList = append(newList, item)
		}
	}
	if len(newList) > appInviteeListCacheMaxSize {
		newList = newList[:appInviteeListCacheMaxSize]
	}
	putAppInviteeListCache(inviterId, newList)
}

// UpdateAppInviteeListRewardGoldInCache 邀请充值返还后更新缓存中的返现金币累计
func UpdateAppInviteeListRewardGoldInCache(inviterId, inviteeId uint64, totalRewardGold float64) {
	if appInviteeListCacheMgr == nil || inviterId == 0 || inviteeId == 0 {
		return
	}
	if _, ok := appInviteeListCacheMgr.GetListCached(gctx.New(), appInviteeListCacheKey(inviterId)); !ok {
		return
	}
	list := getAppInviteeListCache(inviterId)
	updated := false
	for _, row := range list {
		if row != nil && row.UserId == inviteeId {
			row.InviterRewardGold = totalRewardGold
			updated = true
			break
		}
	}
	if updated {
		putAppInviteeListCache(inviterId, list)
	}
}

// GetAppInviteesByInviter 分页获取名下邀请人; pageSize=20 时缓存2页, 前2页走缓存, 第3页起查库
func GetAppInviteesByInviter(inviterId uint64, page, pageSize int) []*AppInviteeListRow {
	list := make([]*AppInviteeListRow, 0)
	if inviterId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = AppInviteeListCachePageSize
	}
	if pageSize != AppInviteeListCachePageSize {
		return loadAppInviteesFromDB(inviterId, page, pageSize)
	}
	if page <= appInviteeListCachedPages {
		cached := getAppInviteeListCache(inviterId)
		start := (page - 1) * pageSize
		if start >= len(cached) {
			return list
		}
		end := start + pageSize
		if end > len(cached) {
			end = len(cached)
		}
		return cached[start:end]
	}
	return loadAppInviteesFromDB(inviterId, page, pageSize)
}

func loadAppInviteesFromDB(inviterId uint64, page, pageSize int) []*AppInviteeListRow {
	list := make([]*AppInviteeListRow, 0)
	if inviterId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = AppInviteeListCachePageSize
	}
	const uiAlias = "ui"
	const ueAlias = "ue"
	_ = g.DB().Model(string(userentity.TbUserInfo)+" "+uiAlias).Ctx(gctx.New()).
		LeftJoin(string(userentity.TbUserExt)+" "+ueAlias, ueAlias+".id = "+uiAlias+".id").
		Fields(
			uiAlias+".id AS user_id",
			uiAlias+".nickname AS nickname",
			uiAlias+".avatar AS avatar",
			uiAlias+".created_at AS created_at",
			"IFNULL("+ueAlias+"."+string(userentity.UserExtInviterInviteRechargeRewardGold)+", 0) AS inviter_reward_gold",
		).
		Where(uiAlias+"."+string(userentity.UserInfoInviterId)+" = ?", inviterId).
		Order(uiAlias + ".id desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}

// NotifyAppInviteeBound 邀请关系首次建立时刷新 App 列表缓存
func NotifyAppInviteeBound(inviterId uint64, invitee *userentity.UserInfo) {
	if inviterId == 0 || invitee == nil || invitee.ID == 0 {
		return
	}
	PrependAppInviteeToListCache(inviterId, invitee, 0)
}
