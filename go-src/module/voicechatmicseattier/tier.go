package voicechatmicseattier

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/syndb"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/voicechatmicseattierdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

const tierSyndbFlushTimeout = 3 * time.Second

func flushTierSyndb() error {
	if !syndb.FlushUntilIdle(tierSyndbFlushTimeout) {
		return errercode.CreateCode(errercode.SysError)
	}
	return nil
}

func validateSeatCount(seatCount int) bool {
	return seatCount >= liveentity.MinVoiceChatMicSeatCount && seatCount <= liveentity.MaxVoiceChatMicSeatCount
}

func ListVoiceChatMicSeatTiers(_ context.Context, req *voicechatmicseattierdto.VoiceChatMicSeatTierListReq) (*httpserver.CMSQueryResp, error) {
	total, list := cfgdao.GetVoiceChatMicSeatTierList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func CreateVoiceChatMicSeatTier(_ context.Context, req *voicechatmicseattierdto.CreateVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.CreateVoiceChatMicSeatTierRes, error) {
	if !validateSeatCount(req.SeatCount) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.GetVoiceChatMicSeatTierBySeatCount(req.SeatCount) != nil {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierSeatCountExist)
	}
	row := &liveentity.VoiceChatMicSeatTier{
		SeatCount: req.SeatCount,
		Sort:      req.Sort,
		Status:    liveentity.VoiceChatMicSeatTierStatusOffShelf,
	}
	if err := cfgdao.CreateVoiceChatMicSeatTier(row); err != nil {
		return nil, err
	}
	reloadTierCache()
	return &voicechatmicseattierdto.CreateVoiceChatMicSeatTierRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateVoiceChatMicSeatTier(_ context.Context, req *voicechatmicseattierdto.UpdateVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.UpdateVoiceChatMicSeatTierRes, error) {
	row := cfgdao.GetVoiceChatMicSeatTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierNonExist)
	}
	if !validateSeatCount(req.SeatCount) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := cfgdao.GetVoiceChatMicSeatTierBySeatCount(req.SeatCount); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierSeatCountExist)
	}
	row.SetSeatCount(req.SeatCount)
	row.SetSort(req.Sort)
	if err := flushTierSyndb(); err != nil {
		return nil, err
	}
	reloadTierCache()
	return &voicechatmicseattierdto.UpdateVoiceChatMicSeatTierRes{Success: true}, nil
}

func DeleteVoiceChatMicSeatTier(_ context.Context, req *voicechatmicseattierdto.DeleteVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.DeleteVoiceChatMicSeatTierRes, error) {
	if cfgdao.GetVoiceChatMicSeatTierById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierNonExist)
	}
	if err := cfgdao.DeleteVoiceChatMicSeatTier(req.ID); err != nil {
		return nil, err
	}
	reloadTierCache()
	return &voicechatmicseattierdto.DeleteVoiceChatMicSeatTierRes{Success: true}, nil
}

func OnShelfVoiceChatMicSeatTier(_ context.Context, req *voicechatmicseattierdto.OnShelfVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.OnShelfVoiceChatMicSeatTierRes, error) {
	row := cfgdao.GetVoiceChatMicSeatTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierNonExist)
	}
	if row.Status != liveentity.VoiceChatMicSeatTierStatusOnShelf {
		row.SetStatus(liveentity.VoiceChatMicSeatTierStatusOnShelf)
		if err := flushTierSyndb(); err != nil {
			return nil, err
		}
		reloadTierCache()
	}
	return &voicechatmicseattierdto.OnShelfVoiceChatMicSeatTierRes{
		Success: true,
		Status:  liveentity.VoiceChatMicSeatTierStatusOnShelf,
	}, nil
}

func OffShelfVoiceChatMicSeatTier(_ context.Context, req *voicechatmicseattierdto.OffShelfVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.OffShelfVoiceChatMicSeatTierRes, error) {
	row := cfgdao.GetVoiceChatMicSeatTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatTierNonExist)
	}
	if row.Status != liveentity.VoiceChatMicSeatTierStatusOffShelf {
		row.SetStatus(liveentity.VoiceChatMicSeatTierStatusOffShelf)
		if err := flushTierSyndb(); err != nil {
			return nil, err
		}
		reloadTierCache()
	}
	return &voicechatmicseattierdto.OffShelfVoiceChatMicSeatTierRes{
		Success: true,
		Status:  liveentity.VoiceChatMicSeatTierStatusOffShelf,
	}, nil
}

func AppTierList(_ context.Context, _ *voicechatmicseattierdto.AppTierListReq) (*voicechatmicseattierdto.AppTierListRes, error) {
	list := GetOnShelfTiers()
	if list == nil {
		list = emptyTierList
	}
	return &voicechatmicseattierdto.AppTierListRes{List: list}, nil
}
