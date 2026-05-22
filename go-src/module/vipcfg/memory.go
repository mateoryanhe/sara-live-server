package vipcfg

import (
	"context"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/vipcfgdao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

type vipCfgSnapshot struct {
	byLevel map[uint32]*vipcfgdto.AppVipCfgItem
	list    []*vipcfgdto.AppVipCfgItem
}

var (
	vipCfgCache     atomic.Value // *vipCfgSnapshot
	emptyVipCfgList = make([]*vipcfgdto.AppVipCfgItem, 0)
)

// Init 服务启动时加载VIP配置到内存
func Init() {
	reloadVipCfgMemory()
}

// reloadVipCfgMemory 从DB重新加载并整体替换内存快照
func reloadVipCfgMemory() {
	rows := vipcfgdao.GetAll()
	byLevel := make(map[uint32]*vipcfgdto.AppVipCfgItem, len(rows))
	list := make([]*vipcfgdto.AppVipCfgItem, 0, len(rows))
	for _, row := range rows {
		item := toAppVipCfgItem(row)
		byLevel[row.Level] = item
		list = append(list, item)
	}
	vipCfgCache.Store(&vipCfgSnapshot{
		byLevel: byLevel,
		list:    list,
	})
}

func getVipCfgSnapshot() *vipCfgSnapshot {
	v := vipCfgCache.Load()
	if v == nil {
		return &vipCfgSnapshot{
			byLevel: make(map[uint32]*vipcfgdto.AppVipCfgItem),
			list:    emptyVipCfgList,
		}
	}
	return v.(*vipCfgSnapshot)
}

func toAppVipCfgItem(row *entity.VipCfg) *vipcfgdto.AppVipCfgItem {
	return &vipcfgdto.AppVipCfgItem{
		ID:                   strconv.FormatUint(row.ID, 10),
		Level:                row.Level,
		LevelName:            row.LevelName,
		Status:               row.Status,
		UpgradeRechargeLimit: row.UpgradeRechargeLimit,
		MinWithdrawAmount:    row.MinWithdrawAmount,
		MaxWithdrawAmount:    row.MaxWithdrawAmount,
		Fee:                  row.Fee,
	}
}

// GetVipCfgFromMemoryByLevel 按等级从内存获取VIP配置(供其它模块使用)
func GetVipCfgFromMemoryByLevel(level uint32) *vipcfgdto.AppVipCfgItem {
	return getVipCfgSnapshot().byLevel[level]
}

// GetAllVipCfgFromMemory 获取全部VIP配置(供其它模块使用)
func GetAllVipCfgFromMemory() []*vipcfgdto.AppVipCfgItem {
	return getVipCfgSnapshot().list
}

// GetAppVipCfgByLevel App端按等级查询VIP配置
func GetAppVipCfgByLevel(_ context.Context, req *vipcfgdto.AppVipCfgByLevelReq) (*vipcfgdto.AppVipCfgByLevelRes, error) {
	item := GetVipCfgFromMemoryByLevel(req.Level)
	if item == nil {
		return nil, errercode.CreateCode(errercode.VipCfgNonExist)
	}
	return &vipcfgdto.AppVipCfgByLevelRes{Item: item}, nil
}

// GetAppVipCfgList App端查询全部VIP配置
func GetAppVipCfgList(_ context.Context, _ *vipcfgdto.AppVipCfgListReq) (*vipcfgdto.AppVipCfgListRes, error) {
	return &vipcfgdto.AppVipCfgListRes{List: GetAllVipCfgFromMemory()}, nil
}
