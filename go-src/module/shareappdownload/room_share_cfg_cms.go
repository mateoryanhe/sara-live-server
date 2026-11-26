package shareappdownload

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetLiveRoomShareCfg(_ context.Context, _ *shareappdownloaddto.GetLiveRoomShareCfgReq) (*shareappdownloaddto.GetLiveRoomShareCfgRes, error) {
	cfg := cfgdao.LoadLiveRoomShareCfg()
	item := toLiveRoomShareCfgItem(cfg)
	return &shareappdownloaddto.GetLiveRoomShareCfgRes{Cfg: item}, nil
}

func SaveLiveRoomShareCfg(_ context.Context, req *shareappdownloaddto.SaveLiveRoomShareCfgReq) (*shareappdownloaddto.SaveLiveRoomShareCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	afLink := strings.TrimSpace(req.AfLink)
	if afLink == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadLiveRoomShareCfg()
	row := &entity.LiveRoomShareCfg{
		ShareTextZhCN: strings.TrimSpace(req.ShareTextZhCN),
		ShareTextZhTW: strings.TrimSpace(req.ShareTextZhTW),
		ShareTextEn:     strings.TrimSpace(req.ShareTextEn),
		ShareTextEs:     strings.TrimSpace(req.ShareTextEs),
		ShareTextPt:     strings.TrimSpace(req.ShareTextPt),
		ShareTextHi:     strings.TrimSpace(req.ShareTextHi),
		ShareTextId:     strings.TrimSpace(req.ShareTextId),
		AfLink:          afLink,
	}
	if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveLiveRoomShareCfg(row); err != nil {
		return nil, err
	}
	reloadLiveRoomShareMemory()
	return &shareappdownloaddto.SaveLiveRoomShareCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func GetLiveRoomShareText(ctx context.Context, req *liveroomdto.GetLiveRoomShareTextReq) (*liveroomdto.GetLiveRoomShareTextRes, error) {
	if req == nil || req.RoomId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	info := userinfodao.FindUserInfoByUserId(userId)
	inviteCode := ""
	if info != nil {
		inviteCode = strings.TrimSpace(info.ShareCode)
	}
	snap := getLiveRoomShareSnapshot()
	shareText := liveRoomShareTextForLang(snap, httpserver.GetLangFromContext(ctx))
	afLink := buildLiveRoomAfLink(snap.AfLink, inviteCode, req.RoomId)
	res := &liveroomdto.GetLiveRoomShareTextRes{
		ShareText: shareText,
		AfLink:    afLink,
	}
	if shareText != "" && afLink != "" {
		res.ShareCode = shareText + "\n" + afLink
	} else if shareText != "" {
		res.ShareCode = shareText
	} else {
		res.ShareCode = afLink
	}
	return res, nil
}

func toLiveRoomShareCfgItem(cfg *entity.LiveRoomShareCfg) *shareappdownloaddto.LiveRoomShareCfgItem {
	bEn, bEs, bPt, bHi, bId := BuiltinLiveRoomShareTexts()
	if cfg == nil {
		return &shareappdownloaddto.LiveRoomShareCfgItem{
			ID:                 "0",
			BuiltinShareTextEn: bEn,
			BuiltinShareTextEs: bEs,
			BuiltinShareTextPt: bPt,
			BuiltinShareTextHi: bHi,
			BuiltinShareTextId: bId,
		}
	}
	return &shareappdownloaddto.LiveRoomShareCfgItem{
		ID:            strconv.FormatUint(cfg.ID, 10),
		ShareTextZhCN: strings.TrimSpace(cfg.ShareTextZhCN),
		ShareTextZhTW: strings.TrimSpace(cfg.ShareTextZhTW),
		ShareTextEn:   strings.TrimSpace(cfg.ShareTextEn),
		ShareTextEs:   strings.TrimSpace(cfg.ShareTextEs),
		ShareTextPt:   strings.TrimSpace(cfg.ShareTextPt),
		ShareTextHi:   strings.TrimSpace(cfg.ShareTextHi),
		ShareTextId:   strings.TrimSpace(cfg.ShareTextId),
		AfLink:        strings.TrimSpace(cfg.AfLink),
		CreatedAt:     formatTime(cfg.CreatedAt),
		UpdatedAt:     formatTime(cfg.UpdatedAt),
		BuiltinShareTextEn: bEn,
		BuiltinShareTextEs: bEs,
		BuiltinShareTextPt: bPt,
		BuiltinShareTextHi: bHi,
		BuiltinShareTextId: bId,
	}
}