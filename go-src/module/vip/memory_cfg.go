package vip

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const vipCfgTimeLayout = "2006-01-02 15:04:05"

type vipCfgSnapshot struct {
	byID       map[uint64]*entity.VipCfg
	byLevel    map[uint32]*entity.VipCfg
	allList    []*entity.VipCfg
	appByLevel map[uint32]*vipcfgdto.AppVipCfgItem
	appList    []*vipcfgdto.AppVipCfgItem
}

var (
	vipCfgCache     atomic.Value // *vipCfgSnapshot
	emptyVipCfgList = make([]*vipcfgdto.AppVipCfgItem, 0)
)

// reloadVipCfgMemory 从DB重新加载并整体替换内存快照
func reloadVipCfgMemory() {
	rows := cfgdao.GetAllVipCfg()
	byID := make(map[uint64]*entity.VipCfg, len(rows))
	byLevel := make(map[uint32]*entity.VipCfg, len(rows))
	allList := make([]*entity.VipCfg, 0, len(rows))
	appByLevel := make(map[uint32]*vipcfgdto.AppVipCfgItem, len(rows))
	appList := make([]*vipcfgdto.AppVipCfgItem, 0, len(rows))

	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		byID[row.ID] = row
		byLevel[row.Level] = row
		allList = append(allList, row)
		item := toAppVipCfgItem(row)
		appByLevel[row.Level] = item
		appList = append(appList, item)
	}

	sort.Slice(allList, func(i, j int) bool {
		if allList[i].Level != allList[j].Level {
			return allList[i].Level < allList[j].Level
		}
		return allList[i].CreatedAt.After(allList[j].CreatedAt)
	})

	vipCfgCache.Store(&vipCfgSnapshot{
		byID:       byID,
		byLevel:    byLevel,
		allList:    allList,
		appByLevel: appByLevel,
		appList:    appList,
	})
}

func getVipCfgSnapshot() *vipCfgSnapshot {
	v := vipCfgCache.Load()
	if v == nil {
		return &vipCfgSnapshot{
			byID:       make(map[uint64]*entity.VipCfg),
			byLevel:    make(map[uint32]*entity.VipCfg),
			allList:    make([]*entity.VipCfg, 0),
			appByLevel: make(map[uint32]*vipcfgdto.AppVipCfgItem),
			appList:    emptyVipCfgList,
		}
	}
	return v.(*vipCfgSnapshot)
}

func getVipCfgByIDFromMemory(id uint64) *entity.VipCfg {
	return getVipCfgSnapshot().byID[id]
}

func findVipCfgByLevelFromMemory(level uint32, excludeID uint64) *entity.VipCfg {
	row := getVipCfgSnapshot().byLevel[level]
	if row == nil || row.ID == excludeID {
		return nil
	}
	return row
}

func listVipCfgFromMemory(levelName string, withdrawSwitchFilter int) []*entity.VipCfg {
	keyword := strings.ToLower(strings.TrimSpace(levelName))
	rows := getVipCfgSnapshot().allList
	filtered := make([]*entity.VipCfg, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if keyword != "" && !strings.Contains(strings.ToLower(row.LevelName), keyword) {
			continue
		}
		switch withdrawSwitchFilter {
		case 1:
			if row.WithdrawSwitch != entity.VipCfgSwitchOff {
				continue
			}
		case 2:
			if row.WithdrawSwitch != entity.VipCfgSwitchOn {
				continue
			}
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func paginateVipCfgList(rows []*entity.VipCfg, pageIndex, pageSize int) ([]*entity.VipCfg, int) {
	total := len(rows)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return []*entity.VipCfg{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return rows[start:end], total
}

func formatVipCfgTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(vipCfgTimeLayout)
}

func toVipCfgListRes(row *entity.VipCfg) *vipcfgdto.VipCfgListRes {
	if row == nil {
		return nil
	}
	return &vipcfgdto.VipCfgListRes{
		ID:                      strconv.FormatUint(row.ID, 10),
		Level:                   row.Level,
		LevelName:               row.LevelName,
		LevelIconName:           row.LevelIcon,
		LevelIcon:               upload.GetUrlByName(row.LevelIcon),
		WithdrawSwitch:          row.WithdrawSwitch,
		AnimationSwitch:         row.AnimationSwitch,
		CommentEffectSwitch:     row.CommentEffectSwitch,
		CustomerServiceSwitch:   row.CustomerServiceSwitch,
		UpgradeRechargeLimit:    row.UpgradeRechargeLimit,
		MinWithdrawAmount:       row.MinWithdrawAmount,
		MaxWithdrawAmount:       row.MaxWithdrawAmount,
		Fee:                     row.Fee,
		AnimationName:           row.Animation,
		Animation:               upload.GetUrlByName(row.Animation),
		AnimationIconName:       row.AnimationIcon,
		AnimationIcon:           upload.GetUrlByName(row.AnimationIcon),
		AnimationDescEn:         row.AnimationDescEn,
		AnimationDescEs:         row.AnimationDescEs,
		AnimationDescPt:         row.AnimationDescPt,
		AnimationDescHi:         row.AnimationDescHi,
		CommentEffectName:       row.CommentEffect,
		CommentEffect:           upload.GetUrlByName(row.CommentEffect),
		CommentEffectIconName:   row.CommentEffectIcon,
		CommentEffectIcon:       upload.GetUrlByName(row.CommentEffectIcon),
		CommentEffectDescEn:     row.CommentEffectDescEn,
		CommentEffectDescEs:     row.CommentEffectDescEs,
		CommentEffectDescPt:     row.CommentEffectDescPt,
		CommentEffectDescHi:     row.CommentEffectDescHi,
		WithdrawIconName:        row.WithdrawIcon,
		WithdrawIcon:            upload.GetUrlByName(row.WithdrawIcon),
		WithdrawNoticeEn:        row.WithdrawNoticeEn,
		WithdrawNoticeEs:        row.WithdrawNoticeEs,
		WithdrawNoticePt:        row.WithdrawNoticePt,
		WithdrawNoticeHi:        row.WithdrawNoticeHi,
		CustomerServiceIconName: row.CustomerServiceIcon,
		CustomerServiceIcon:     upload.GetUrlByName(row.CustomerServiceIcon),
		CustomerServiceDescEn:   row.CustomerServiceDescEn,
		CustomerServiceDescEs:   row.CustomerServiceDescEs,
		CustomerServiceDescPt:   row.CustomerServiceDescPt,
		CustomerServiceDescHi:   row.CustomerServiceDescHi,
		CreatedAt:               formatVipCfgTime(row.CreatedAt),
		UpdatedAt:               formatVipCfgTime(row.UpdatedAt),
	}
}

func queryVipCfgListFromMemory(req *vipcfgdto.VipCfgListReq) (int, []*vipcfgdto.VipCfgListRes) {
	if req == nil {
		return 0, []*vipcfgdto.VipCfgListRes{}
	}
	rows, total := paginateVipCfgList(
		listVipCfgFromMemory(req.LevelName, req.WithdrawSwitchFilter),
		req.PageIndex,
		req.PageSize,
	)
	list := make([]*vipcfgdto.VipCfgListRes, 0, len(rows))
	for _, row := range rows {
		list = append(list, toVipCfgListRes(row))
	}
	return total, list
}

func toAppVipCfgItem(row *entity.VipCfg) *vipcfgdto.AppVipCfgItem {
	if row == nil {
		return nil
	}
	return &vipcfgdto.AppVipCfgItem{
		Level:                row.Level,
		LevelName:            row.LevelName,
		LevelIcon:            upload.GetUrlByName(row.LevelIcon),
		UpgradeRechargeLimit: row.UpgradeRechargeLimit,
		PrivilegeList:        buildAppVipPrivilegeList(row),
	}
}

func buildAppVipPrivilegeList(row *entity.VipCfg) []*vipcfgdto.AppVipPrivilegeItem {
	if row == nil {
		return []*vipcfgdto.AppVipPrivilegeItem{}
	}
	list := make([]*vipcfgdto.AppVipPrivilegeItem, 0, 4)
	if row.WithdrawSwitch == entity.VipCfgSwitchOn {
		list = append(list, &vipcfgdto.AppVipPrivilegeItem{
			PrivilegeType:     vipcfgdto.AppVipPrivilegeTypeWithdraw,
			Icon:              upload.GetUrlByName(row.WithdrawIcon),
			Desc:              appVipPrivilegeDesc(row.WithdrawNoticeEn, row.WithdrawNoticeEs, row.WithdrawNoticePt, row.WithdrawNoticeHi),
			MinWithdrawAmount: row.MinWithdrawAmount,
			MaxWithdrawAmount: row.MaxWithdrawAmount,
			Fee:               row.Fee,
		})
	}
	if row.AnimationSwitch == entity.VipCfgSwitchOn {
		list = append(list, &vipcfgdto.AppVipPrivilegeItem{
			PrivilegeType: vipcfgdto.AppVipPrivilegeTypeEntryEffect,
			Icon:          upload.GetUrlByName(row.AnimationIcon),
			Desc:          appVipPrivilegeDesc(row.AnimationDescEn, row.AnimationDescEs, row.AnimationDescPt, row.AnimationDescHi),
			Animation:     upload.GetUrlByName(row.Animation),
		})
	}
	if row.CommentEffectSwitch == entity.VipCfgSwitchOn {
		list = append(list, &vipcfgdto.AppVipPrivilegeItem{
			PrivilegeType: vipcfgdto.AppVipPrivilegeTypeCommentEffect,
			Icon:          upload.GetUrlByName(row.CommentEffectIcon),
			Desc:          appVipPrivilegeDesc(row.CommentEffectDescEn, row.CommentEffectDescEs, row.CommentEffectDescPt, row.CommentEffectDescHi),
			Animation:     upload.GetUrlByName(row.CommentEffect),
		})
	}
	if row.CustomerServiceSwitch == entity.VipCfgSwitchOn {
		list = append(list, &vipcfgdto.AppVipPrivilegeItem{
			PrivilegeType: vipcfgdto.AppVipPrivilegeTypeCustomerService,
			Icon:          upload.GetUrlByName(row.CustomerServiceIcon),
			Desc:          appVipPrivilegeDesc(row.CustomerServiceDescEn, row.CustomerServiceDescEs, row.CustomerServiceDescPt, row.CustomerServiceDescHi),
		})
	}
	return list
}

func appVipPrivilegeDesc(en, es, pt, hi string) vipcfgdto.AppVipPrivilegeI18nText {
	return vipcfgdto.AppVipPrivilegeI18nText{
		En: en,
		Es: es,
		Pt: pt,
		Hi: hi,
	}
}

// GetVipCfgFromMemoryByLevel 按等级从内存获取VIP配置(供其它模块使用)
func GetVipCfgFromMemoryByLevel(level uint32) *vipcfgdto.AppVipCfgItem {
	return getVipCfgSnapshot().appByLevel[level]
}

// GetAllVipCfgFromMemory 获取全部VIP配置(供其它模块使用)
func GetAllVipCfgFromMemory() []*vipcfgdto.AppVipCfgItem {
	return getVipCfgSnapshot().appList
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
