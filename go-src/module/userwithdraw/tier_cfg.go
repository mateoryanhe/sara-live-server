package userwithdraw

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/syndb"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/userwithdrawtiercfgdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

const tierCfgSyndbFlushTimeout = 3 * time.Second

func flushTierCfgSyndb() error {
	if !syndb.FlushUntilIdle(tierCfgSyndbFlushTimeout) {
		return errercode.CreateCode(errercode.SysError)
	}
	return nil
}

func normalizeTierUsd(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func normalizeTierGold(v float64) float64 {
	return math.Round(v*100) / 100
}

func ListUserWithdrawTierCfgs(_ context.Context, req *userwithdrawtiercfgdto.UserWithdrawTierCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := cfgdao.GetUserWithdrawTierCfgList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func CreateUserWithdrawTierCfg(_ context.Context, req *userwithdrawtiercfgdto.CreateUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.CreateUserWithdrawTierCfgRes, error) {
	name := strings.TrimSpace(req.Name)
	withdrawUsd := normalizeTierUsd(req.WithdrawUsd)
	minGameConsumeGold := normalizeTierGold(req.MinGameConsumeGold)
	if name == "" || withdrawUsd <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if minGameConsumeGold < 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.GetUserWithdrawTierCfgByName(name) != nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	cfg := &userentity.UserWithdrawTierCfg{
		Name:               name,
		WithdrawUsd:        withdrawUsd,
		MinGameConsumeGold: minGameConsumeGold,
		Status:         userentity.UserWithdrawTierCfgStatusOffShelf,
	}
	if err := cfgdao.CreateUserWithdrawTierCfg(cfg); err != nil {
		return nil, err
	}
	reloadUserWithdrawTierCfgCache()
	return &userwithdrawtiercfgdto.CreateUserWithdrawTierCfgRes{ID: strconv.FormatUint(cfg.ID, 10)}, nil
}

func UpdateUserWithdrawTierCfg(_ context.Context, req *userwithdrawtiercfgdto.UpdateUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.UpdateUserWithdrawTierCfgRes, error) {
	cfg := cfgdao.GetUserWithdrawTierCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	name := strings.TrimSpace(req.Name)
	withdrawUsd := normalizeTierUsd(req.WithdrawUsd)
	minGameConsumeGold := normalizeTierGold(req.MinGameConsumeGold)
	if name == "" || withdrawUsd <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if minGameConsumeGold < 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := cfgdao.GetUserWithdrawTierCfgByName(name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	cfg.SetName(name)
	cfg.SetWithdrawUsd(withdrawUsd)
	cfg.SetMinGameConsumeGold(minGameConsumeGold)
	if err := flushTierCfgSyndb(); err != nil {
		return nil, err
	}
	reloadUserWithdrawTierCfgCache()
	return &userwithdrawtiercfgdto.UpdateUserWithdrawTierCfgRes{Success: true}, nil
}

func DeleteUserWithdrawTierCfg(_ context.Context, req *userwithdrawtiercfgdto.DeleteUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.DeleteUserWithdrawTierCfgRes, error) {
	if cfgdao.GetUserWithdrawTierCfgById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if err := cfgdao.DeleteUserWithdrawTierCfg(req.ID); err != nil {
		return nil, err
	}
	reloadUserWithdrawTierCfgCache()
	return &userwithdrawtiercfgdto.DeleteUserWithdrawTierCfgRes{Success: true}, nil
}

func OnShelfUserWithdrawTierCfg(_ context.Context, req *userwithdrawtiercfgdto.OnShelfUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.OnShelfUserWithdrawTierCfgRes, error) {
	cfg := cfgdao.GetUserWithdrawTierCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != userentity.UserWithdrawTierCfgStatusOnShelf {
		cfg.SetStatus(userentity.UserWithdrawTierCfgStatusOnShelf)
		if err := flushTierCfgSyndb(); err != nil {
			return nil, err
		}
		reloadUserWithdrawTierCfgCache()
	}
	return &userwithdrawtiercfgdto.OnShelfUserWithdrawTierCfgRes{
		Success: true,
		Status:  userentity.UserWithdrawTierCfgStatusOnShelf,
	}, nil
}

func OffShelfUserWithdrawTierCfg(_ context.Context, req *userwithdrawtiercfgdto.OffShelfUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.OffShelfUserWithdrawTierCfgRes, error) {
	cfg := cfgdao.GetUserWithdrawTierCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != userentity.UserWithdrawTierCfgStatusOffShelf {
		cfg.SetStatus(userentity.UserWithdrawTierCfgStatusOffShelf)
		if err := flushTierCfgSyndb(); err != nil {
			return nil, err
		}
		reloadUserWithdrawTierCfgCache()
	}
	return &userwithdrawtiercfgdto.OffShelfUserWithdrawTierCfgRes{
		Success: true,
		Status:  userentity.UserWithdrawTierCfgStatusOffShelf,
	}, nil
}

func AppListWithdrawTiers(_ context.Context, _ *userwithdrawtiercfgdto.AppUserWithdrawTierListReq) (*userwithdrawtiercfgdto.AppUserWithdrawTierListRes, error) {
	list := GetOnShelfUserWithdrawTiers()
	if list == nil {
		list = emptyUserWithdrawTierCfgList
	}
	return &userwithdrawtiercfgdto.AppUserWithdrawTierListRes{List: list}, nil
}
