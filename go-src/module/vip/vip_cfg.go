package vip

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity/live"
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
	if findVipCfgByLevelFromMemory(req.Level, 0) != nil {
		return nil, errercode.CreateCode(errercode.VipCfgExist)
	}
	row := &entity.VipCfg{
		Level:                  req.Level,
		LevelName:              req.LevelName,
		LevelIcon:              req.LevelIcon,
		AnimationSwitch:        req.AnimationSwitch,
		CommentEffectSwitch:    req.CommentEffectSwitch,
		CustomerServiceSwitch:  req.CustomerServiceSwitch,
		UpgradeRechargeLimit:   req.UpgradeRechargeLimit,
		Animation:              req.Animation,
		AnimationIcon:          req.AnimationIcon,
		AnimationTitleEn:       req.AnimationTitleEn,
		AnimationTitleEs:       req.AnimationTitleEs,
		AnimationTitlePt:       req.AnimationTitlePt,
		AnimationTitleHi:       req.AnimationTitleHi,
		AnimationTitleId:       req.AnimationTitleId,
		AnimationDescEn:        req.AnimationDescEn,
		AnimationDescEs:        req.AnimationDescEs,
		AnimationDescPt:        req.AnimationDescPt,
		AnimationDescHi:        req.AnimationDescHi,
		AnimationDescId:        req.AnimationDescId,
		CommentEffect:          req.CommentEffect,
		CommentEffectIcon:      req.CommentEffectIcon,
		CommentEffectTitleEn:   req.CommentEffectTitleEn,
		CommentEffectTitleEs:   req.CommentEffectTitleEs,
		CommentEffectTitlePt:   req.CommentEffectTitlePt,
		CommentEffectTitleHi:   req.CommentEffectTitleHi,
		CommentEffectTitleId:   req.CommentEffectTitleId,
		CommentEffectDescEn:    req.CommentEffectDescEn,
		CommentEffectDescEs:    req.CommentEffectDescEs,
		CommentEffectDescPt:    req.CommentEffectDescPt,
		CommentEffectDescHi:    req.CommentEffectDescHi,
		CommentEffectDescId:    req.CommentEffectDescId,
		CustomerServiceIcon:    req.CustomerServiceIcon,
		CustomerServiceTitleEn: req.CustomerServiceTitleEn,
		CustomerServiceTitleEs: req.CustomerServiceTitleEs,
		CustomerServiceTitlePt: req.CustomerServiceTitlePt,
		CustomerServiceTitleHi: req.CustomerServiceTitleHi,
		CustomerServiceTitleId: req.CustomerServiceTitleId,
		CustomerServiceDescEn:  req.CustomerServiceDescEn,
		CustomerServiceDescEs:  req.CustomerServiceDescEs,
		CustomerServiceDescPt:  req.CustomerServiceDescPt,
		CustomerServiceDescHi:  req.CustomerServiceDescHi,
		CustomerServiceDescId:  req.CustomerServiceDescId,
	}
	if err := cfgdao.CreateVipCfg(row); err != nil {
		return nil, err
	}
	reloadVipCfgMemory()
	return &vipcfgdto.CreateVipCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *vipcfgdto.UpdateVipCfgReq) (*vipcfgdto.UpdateVipCfgRes, error) {
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
	updated.LevelIcon = req.LevelIcon
	updated.AnimationSwitch = req.AnimationSwitch
	updated.CommentEffectSwitch = req.CommentEffectSwitch
	updated.CustomerServiceSwitch = req.CustomerServiceSwitch
	updated.UpgradeRechargeLimit = req.UpgradeRechargeLimit
	updated.Animation = req.Animation
	updated.AnimationIcon = req.AnimationIcon
	updated.AnimationTitleEn = req.AnimationTitleEn
	updated.AnimationTitleEs = req.AnimationTitleEs
	updated.AnimationTitlePt = req.AnimationTitlePt
	updated.AnimationTitleHi = req.AnimationTitleHi
	updated.AnimationTitleId = req.AnimationTitleId
	updated.AnimationDescEn = req.AnimationDescEn
	updated.AnimationDescEs = req.AnimationDescEs
	updated.AnimationDescPt = req.AnimationDescPt
	updated.AnimationDescHi = req.AnimationDescHi
	updated.AnimationDescId = req.AnimationDescId
	updated.CommentEffect = req.CommentEffect
	updated.CommentEffectIcon = req.CommentEffectIcon
	updated.CommentEffectTitleEn = req.CommentEffectTitleEn
	updated.CommentEffectTitleEs = req.CommentEffectTitleEs
	updated.CommentEffectTitlePt = req.CommentEffectTitlePt
	updated.CommentEffectTitleHi = req.CommentEffectTitleHi
	updated.CommentEffectTitleId = req.CommentEffectTitleId
	updated.CommentEffectDescEn = req.CommentEffectDescEn
	updated.CommentEffectDescEs = req.CommentEffectDescEs
	updated.CommentEffectDescPt = req.CommentEffectDescPt
	updated.CommentEffectDescHi = req.CommentEffectDescHi
	updated.CommentEffectDescId = req.CommentEffectDescId
	updated.CustomerServiceIcon = req.CustomerServiceIcon
	updated.CustomerServiceTitleEn = req.CustomerServiceTitleEn
	updated.CustomerServiceTitleEs = req.CustomerServiceTitleEs
	updated.CustomerServiceTitlePt = req.CustomerServiceTitlePt
	updated.CustomerServiceTitleHi = req.CustomerServiceTitleHi
	updated.CustomerServiceTitleId = req.CustomerServiceTitleId
	updated.CustomerServiceDescEn = req.CustomerServiceDescEn
	updated.CustomerServiceDescEs = req.CustomerServiceDescEs
	updated.CustomerServiceDescPt = req.CustomerServiceDescPt
	updated.CustomerServiceDescHi = req.CustomerServiceDescHi
	updated.CustomerServiceDescId = req.CustomerServiceDescId
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
