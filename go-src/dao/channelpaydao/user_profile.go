package channelpaydao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
)

var channelPayUserProfileCacheMgr *cache.RowCache[*entity.ChannelPayUserProfile]

func Init() {
	channelPayUserProfileCacheMgr = cache.NewRowCache[*entity.ChannelPayUserProfile]()
}

func loadFromDB(userId uint64) *entity.ChannelPayUserProfile {
	if userId == 0 {
		return nil
	}
	var data *entity.ChannelPayUserProfile
	_ = g.Model(string(entity.TbChannelPayUserProfile)).Unscoped().Where(g.Map{
		string(db.IdName): userId,
	}).Scan(&data)
	if data == nil || data.ID == 0 {
		return nil
	}
	return data
}

// GetExisting 仅读缓存/库；不存在返回 nil，绝不新建（防假 userId 写库）
func GetExisting(userId uint64) *entity.ChannelPayUserProfile {
	if userId == 0 || channelPayUserProfileCacheMgr == nil {
		return nil
	}
	if v, _ := channelPayUserProfileCacheMgr.GetRowCached(gctx.New(), userId); v != nil && v.ID != 0 {
		return v
	}
	row := loadFromDB(userId)
	if row == nil {
		return nil
	}
	channelPayUserProfileCacheMgr.PublishRow(gctx.New(), userId, row)
	return row
}

// Publish 原地修改后刷新缓存
func Publish(data *entity.ChannelPayUserProfile) {
	if data == nil || data.ID == 0 || channelPayUserProfileCacheMgr == nil {
		return
	}
	channelPayUserProfileCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

// UpsertPayerInfo 写入非空 name/email。user 必须在 user_infos 中真实存在，否则不写库。
func UpsertPayerInfo(userId uint64, name, email string) *entity.ChannelPayUserProfile {
	if userId == 0 {
		return nil
	}
	if userinfodao.GetUserInfoFromDB(userId) == nil {
		return nil
	}
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if isPlaceholderPayEmail(email) {
		email = ""
	}
	if name == "" && email == "" {
		return GetExisting(userId)
	}

	row := GetExisting(userId)
	if row == nil {
		row = entity.NewChannelPayUserProfile(userId)
	}
	changed := false
	if name != "" && name != row.Name {
		row.SetName(name)
		changed = true
	}
	if email != "" && email != row.Email {
		row.SetEmail(email)
		changed = true
	}
	if changed {
		Publish(row)
	}
	return row
}

func isPlaceholderPayEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	return strings.HasSuffix(email, "@noreply.local") || strings.HasSuffix(email, "@noreply.sara.live")
}
