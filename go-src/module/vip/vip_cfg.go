package vip

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	coremath "xr-game-server/core/math"
	"xr-game-server/dao/vipcfgdao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *vipcfgdto.VipCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := vipcfgdao.GetList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func Create(_ context.Context, req *vipcfgdto.CreateVipCfgReq) (*vipcfgdto.CreateVipCfgRes, error) {
	normalizeVipCfgAmounts(
		&req.UpgradeRechargeLimit,
		&req.MinWithdrawAmount,
		&req.MaxWithdrawAmount,
		&req.Fee,
	)
	if err := validateWithdrawRange(req.MinWithdrawAmount, req.MaxWithdrawAmount); err != nil {
		return nil, err
	}
	if existing := vipcfgdao.GetByLevel(req.Level); existing != nil {
		return nil, errercode.CreateCode(errercode.VipCfgExist)
	}
	row := &entity.VipCfg{
		Level:                req.Level,
		LevelName:            req.LevelName,
		Status:               req.Status,
		UpgradeRechargeLimit: req.UpgradeRechargeLimit,
		MinWithdrawAmount:    req.MinWithdrawAmount,
		MaxWithdrawAmount:    req.MaxWithdrawAmount,
		Fee:                  req.Fee,
	}
	if err := vipcfgdao.Create(row); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.CreateVipCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *vipcfgdto.UpdateVipCfgReq) (*vipcfgdto.UpdateVipCfgRes, error) {
	normalizeVipCfgAmounts(
		&req.UpgradeRechargeLimit,
		&req.MinWithdrawAmount,
		&req.MaxWithdrawAmount,
		&req.Fee,
	)
	if err := validateWithdrawRange(req.MinWithdrawAmount, req.MaxWithdrawAmount); err != nil {
		return nil, err
	}
	row := vipcfgdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.VipCfgNonExist)
	}
	if existing := vipcfgdao.GetByLevel(req.Level); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.VipCfgExist)
	}
	row.Level = req.Level
	row.LevelName = req.LevelName
	row.Status = req.Status
	row.UpgradeRechargeLimit = req.UpgradeRechargeLimit
	row.MinWithdrawAmount = req.MinWithdrawAmount
	row.MaxWithdrawAmount = req.MaxWithdrawAmount
	row.Fee = req.Fee
	if err := vipcfgdao.Update(row); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.UpdateVipCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *vipcfgdto.DeleteVipCfgReq) (*vipcfgdto.DeleteVipCfgRes, error) {
	if row := vipcfgdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.VipCfgNonExist)
	}
	if err := vipcfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.DeleteVipCfgRes{Success: true}, nil
}

func normalizeVipCfgAmounts(values ...*float64) {
	for _, v := range values {
		if v == nil {
			continue
		}
		*v = coremath.RoundFloat64(*v)
	}
}

func validateWithdrawRange(minAmount, maxAmount float64) error {
	if maxAmount > 0 && minAmount > maxAmount {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}
