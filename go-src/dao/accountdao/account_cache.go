package accountdao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/core/phoneutil"
	"xr-game-server/entity/user"
	"xr-game-server/module/usermaxid"
)

// PhoneOpenId 手机号登陆 open_id = 区号_手机号
func PhoneOpenId(areaCode, phone string) string {
	return phoneutil.NormalizeAreaCode(areaCode) + "_" + strings.TrimSpace(phone)
}

// LogicalOpenId 缓存 key 使用的 open_id
func LogicalOpenId(openId string) string {
	return entity.NormalizeOpenId(openId)
}

func accountListCacheKey(channel uint, openId string) string {
	return fmt.Sprintf("%d:%s", channel, LogicalOpenId(openId))
}

func parseAccountListCacheKey(key string) (channel uint, openId string, ok bool) {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 {
		return 0, "", false
	}
	ch, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, "", false
	}
	return uint(ch), parts[1], true
}

func loadAccountListFromDB(openId string, channel uint) []*entity.Account {
	openId = LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return nil
	}
	var list []*entity.Account
	_ = g.Model(string(entity.TbAccount)).Unscoped().Where(
		"channel = ? AND open_id = ?",
		channel, openId,
	).Order("id desc").Scan(&list)
	return list
}

func dbLoaderForOpenId(openId string, channel uint) cache.Loader[[]*entity.Account] {
	return func(ctx context.Context) ([]*entity.Account, error) {
		return loadAccountListFromDB(openId, channel), nil
	}
}

// GetAccountList 按 channel + openId 获取账号列表缓存(含已注销)
func GetAccountList(openId string, channel uint) []*entity.Account {
	openId = LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return nil
	}
	key := accountListCacheKey(channel, openId)
	return accountCacheMgr.MustGetList(gctx.New(), key, dbLoaderForOpenId(openId, channel))
}

func castAccountList(data any) []*entity.Account {
	if data == nil {
		return nil
	}
	list, ok := data.([]*entity.Account)
	if !ok || len(list) == 0 {
		return nil
	}
	return list
}

func findActiveAccount(list []*entity.Account) *entity.Account {
	for _, item := range list {
		if item != nil && !item.Cancel {
			return item
		}
	}
	return nil
}

func findAccountByID(list []*entity.Account, accountId uint64) *entity.Account {
	for _, item := range list {
		if item != nil && item.ID == accountId {
			return item
		}
	}
	return nil
}

// FindActiveAccount 获取未注销账号
func FindActiveAccount(openId string, channel uint) *entity.Account {
	return findActiveAccount(GetAccountList(openId, channel))
}

// FindActivePhoneAccount 获取未注销手机号账号
func FindActivePhoneAccount(areaCode, phone string, channel uint) *entity.Account {
	return FindActiveAccount(PhoneOpenId(areaCode, phone), channel)
}

// GetAccountFromCache 从列表缓存中按 accountId 获取账号(CMS/App 写操作)
func GetAccountFromCache(openId string, channel uint, accountId uint64) *entity.Account {
	list := GetAccountList(openId, channel)
	if acc := findAccountByID(list, accountId); acc != nil {
		return acc
	}
	return nil
}

func appendAccountToListCache(openId string, channel uint, account *entity.Account) {
	openId = LogicalOpenId(openId)
	key := accountListCacheKey(channel, openId)
	list := GetAccountList(openId, channel)
	newList := append(append([]*entity.Account{}, list...), account)
	accountCacheMgr.PublishList(gctx.New(), key, newList)
}

// RegisterAccount 注册新账号并写入列表缓存(ID 由 usermaxid 分配)
func RegisterAccount(openId string, channel uint) *entity.Account {
	return RegisterAccountWithID(usermaxid.NextUserId(), openId, channel)
}

// RegisterAccountWithID 使用指定ID注册新账号(id=0 时回退为 usermaxid 分配)
func RegisterAccountWithID(id uint64, openId string, channel uint) *entity.Account {
	openId = LogicalOpenId(strings.TrimSpace(openId))
	if openId == "" || channel == 0 {
		return nil
	}
	if id == 0 {
		id = usermaxid.NextUserId()
	}
	acc := entity.NewAccountWithID(id, openId, channel)
	appendAccountToListCache(openId, channel, acc)
	return acc
}

// RegisterPhoneAccount 手机号注册(open_id 由区号+手机号组成)
func RegisterPhoneAccount(areaCode, phone string, channel uint) *entity.Account {
	openId := PhoneOpenId(areaCode, phone)
	acc := RegisterAccount(openId, channel)
	if acc == nil {
		return nil
	}
	acc.SetPhoneAreaCode(phoneutil.NormalizeAreaCode(areaCode))
	return acc
}

// RegisterDeviceAccount 设备码注册
func RegisterDeviceAccount(deviceId string, channel uint) *entity.Account {
	return RegisterAccount(deviceId, channel)
}

// GetOrRegisterDeviceAccount 设备码登陆/注册
func GetOrRegisterDeviceAccount(deviceId string, channel uint) (*entity.Account, bool) {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" || channel == 0 {
		return nil, false
	}
	if acc := FindActiveAccount(deviceId, channel); acc != nil {
		return acc, false
	}
	acc := RegisterDeviceAccount(deviceId, channel)
	if acc == nil || acc.ID == 0 {
		return nil, false
	}
	return acc, true
}

// PublishAccountList 账号列表原地修改后刷新缓存.
func PublishAccountList(openId string, channel uint) {
	openId = LogicalOpenId(openId)
	if openId == "" || channel == 0 || accountCacheMgr == nil {
		return
	}
	key := accountListCacheKey(channel, openId)
	list := GetAccountList(openId, channel)
	accountCacheMgr.PublishList(gctx.New(), key, list)
}

// SyncAccountListCache 从数据库刷新列表缓存
func SyncAccountListCache(openId string, channel uint) {
	openId = LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return
	}
	key := accountListCacheKey(channel, openId)
	list := loadAccountListFromDB(openId, channel)
	accountCacheMgr.PublishList(gctx.New(), key, list)
}

// FindAccountInCacheByID 遍历缓存查找账号(App 注销优先走缓存)
func FindAccountInCacheByID(accountId uint64) (*entity.Account, string, uint) {
	if accountId == 0 || accountCacheMgr == nil {
		return nil, "", 0
	}
	ctx := gctx.New()
	keys, err := accountCacheMgr.Keys(ctx)
	if err != nil {
		return nil, "", 0
	}
	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok { continue }
		channel, openId, ok := parseAccountListCacheKey(keyStr)
		if !ok {
			continue
		}
		if acc := findAccountByID(func() []*entity.Account { l, _ := accountCacheMgr.GetListCached(ctx, key); return l }(), accountId); acc != nil {
			return acc, openId, channel
		}
	}
	return nil, "", 0
}

// ResolveAccountForCMS CMS 操作: 按 openId + channel 从列表缓存定位账号
func ResolveAccountForCMS(openId string, channel uint, accountId uint64) *entity.Account {
	return GetAccountFromCache(LogicalOpenId(openId), channel, accountId)
}
