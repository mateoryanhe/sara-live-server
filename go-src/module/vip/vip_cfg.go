package vip

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *vipcfgdto.VipCfgListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := queryVipCfgListFromMemory(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *vipcfgdto.CreateVipCfgReq) (*vipcfgdto.CreateVipCfgRes, error) {

	if err := validateWithdrawRange(req.MinWithdrawAmount, req.MaxWithdrawAmount); err != nil {
		return nil, err
	}
	if findVipCfgByLevelFromMemory(req.Level, 0) != nil {
		return nil, errercode.CreateCode(errercode.VipCfgExist)
	}
	row := &entity.VipCfg{
		Level:                req.Level,
		LevelName:            req.LevelName,
		WithdrawSwitch:       req.WithdrawSwitch,
		AnimationSwitch:      req.AnimationSwitch,
		CommentEffectSwitch:  req.CommentEffectSwitch,
		UpgradeRechargeLimit: req.UpgradeRechargeLimit,
		MinWithdrawAmount:    req.MinWithdrawAmount,
		MaxWithdrawAmount:    req.MaxWithdrawAmount,
		Fee:                  req.Fee,
		Animation:            req.Animation,
		AnimationIcon:        req.AnimationIcon,
		AnimationDescEn:      req.AnimationDescEn,
		AnimationDescEs:      req.AnimationDescEs,
		AnimationDescPt:      req.AnimationDescPt,
		AnimationDescHi:      req.AnimationDescHi,
		CommentEffect:        req.CommentEffect,
		CommentEffectIcon:    req.CommentEffectIcon,
		CommentEffectDescEn:  req.CommentEffectDescEn,
		CommentEffectDescEs:  req.CommentEffectDescEs,
		CommentEffectDescPt:  req.CommentEffectDescPt,
		CommentEffectDescHi:  req.CommentEffectDescHi,
		WithdrawIcon:         req.WithdrawIcon,
		WithdrawNoticeEn:     req.WithdrawNoticeEn,
		WithdrawNoticeEs:     req.WithdrawNoticeEs,
		WithdrawNoticePt:     req.WithdrawNoticePt,
		WithdrawNoticeHi:     req.WithdrawNoticeHi,
	}
	if err := cfgdao.CreateVipCfg(row); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.CreateVipCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *vipcfgdto.UpdateVipCfgReq) (*vipcfgdto.UpdateVipCfgRes, error) {

	if err := validateWithdrawRange(req.MinWithdrawAmount, req.MaxWithdrawAmount); err != nil {
		return nil, err
	}
	row := getVipCfgByIDFromMemory(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.VipCfgNonExist)
	}
	if findVipCfgByLevelFromMemory(req.Level, req.ID) != nil {
		return nil, errercode.CreateCode(errercode.VipCfgExist)
	}
	updated := *row
	updated.Level = req.Level
	updated.LevelName = req.LevelName
	updated.WithdrawSwitch = req.WithdrawSwitch
	updated.AnimationSwitch = req.AnimationSwitch
	updated.CommentEffectSwitch = req.CommentEffectSwitch
	updated.UpgradeRechargeLimit = req.UpgradeRechargeLimit
	updated.MinWithdrawAmount = req.MinWithdrawAmount
	updated.MaxWithdrawAmount = req.MaxWithdrawAmount
	updated.Fee = req.Fee
	updated.Animation = req.Animation
	updated.AnimationIcon = req.AnimationIcon
	updated.AnimationDescEn = req.AnimationDescEn
	updated.AnimationDescEs = req.AnimationDescEs
	updated.AnimationDescPt = req.AnimationDescPt
	updated.AnimationDescHi = req.AnimationDescHi
	updated.CommentEffect = req.CommentEffect
	updated.CommentEffectIcon = req.CommentEffectIcon
	updated.CommentEffectDescEn = req.CommentEffectDescEn
	updated.CommentEffectDescEs = req.CommentEffectDescEs
	updated.CommentEffectDescPt = req.CommentEffectDescPt
	updated.CommentEffectDescHi = req.CommentEffectDescHi
	updated.WithdrawIcon = req.WithdrawIcon
	updated.WithdrawNoticeEn = req.WithdrawNoticeEn
	updated.WithdrawNoticeEs = req.WithdrawNoticeEs
	updated.WithdrawNoticePt = req.WithdrawNoticePt
	updated.WithdrawNoticeHi = req.WithdrawNoticeHi
	if err := cfgdao.UpdateVipCfg(&updated); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.UpdateVipCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *vipcfgdto.DeleteVipCfgReq) (*vipcfgdto.DeleteVipCfgRes, error) {
	if getVipCfgByIDFromMemory(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.VipCfgNonExist)
	}
	if err := cfgdao.DeleteVipCfg(req.ID); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.DeleteVipCfgRes{Success: true}, nil
}

func validateWithdrawRange(minAmount, maxAmount float64) error {
	if maxAmount > 0 && minAmount > maxAmount {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}
