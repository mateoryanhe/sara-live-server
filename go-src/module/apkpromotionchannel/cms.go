package apkpromotionchannel

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dto/apkpromotionchanneldto"
	entity "xr-game-server/entity/cms"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func GetList(_ context.Context, req *apkpromotionchanneldto.ApkPromotionChannelListReq) (*apkpromotionchanneldto.ApkPromotionChannelListRes, error) {
	total, list := apkpromotionchanneldao.GetList(req)
	baseHost := OfficialSiteShareBaseHost()
	for _, item := range list {
		if item == nil {
			continue
		}
		item.ShareURL = BuildChannelShortURL(item.ChannelCode)
		item.ApkDownloadURL = GetCachedApkURL(item.ChannelCode)
	}
	return &apkpromotionchanneldto.ApkPromotionChannelListRes{
		Total:         total,
		Data:          list,
		ShareBaseHost: baseHost,
	}, nil
}

func Create(_ context.Context, req *apkpromotionchanneldto.CreateApkPromotionChannelReq) (*apkpromotionchanneldto.CreateApkPromotionChannelRes, error) {
	channelName := strings.TrimSpace(req.ChannelName)
	code, ok := normalizeChannelCode(req.ChannelCode)
	if channelName == "" || !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if apkpromotionchanneldao.GetByChannelCode(code) != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if BuildChannelShortURL(code) == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.ApkPromotionChannel{
		ChannelName: channelName,
		ChannelCode: code,
	}
	if err := apkpromotionchanneldao.Create(row); err != nil {
		return nil, err
	}
	SetCachedApkURL(code, "")
	return &apkpromotionchanneldto.CreateApkPromotionChannelRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *apkpromotionchanneldto.UpdateApkPromotionChannelReq) (*apkpromotionchanneldto.UpdateApkPromotionChannelRes, error) {
	row := apkpromotionchanneldao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	channelName := strings.TrimSpace(req.ChannelName)
	code, ok := normalizeChannelCode(req.ChannelCode)
	if channelName == "" || !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := apkpromotionchanneldao.GetByChannelCode(code); existing != nil && existing.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.ChannelName = channelName
	oldCode := row.ChannelCode
	row.ChannelCode = code
	if oldCode != code {
		row.ApkFileName = ""
		row.ApkDownloadURL = ""
		RenameCachedChannel(oldCode, code)
	}
	if err := apkpromotionchanneldao.Update(row); err != nil {
		return nil, err
	}
	return &apkpromotionchanneldto.UpdateApkPromotionChannelRes{Success: true}, nil
}

func Delete(_ context.Context, req *apkpromotionchanneldto.DeleteApkPromotionChannelReq) (*apkpromotionchanneldto.DeleteApkPromotionChannelRes, error) {
	row := apkpromotionchanneldao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	stored := strings.TrimSpace(row.ApkFileName)
	if err := apkpromotionchanneldao.Delete(req.ID); err != nil {
		return nil, err
	}
	RemoveCachedChannel(row.ChannelCode)
	if stored != "" {
		upload.DeleteApkPromotionChannelStored(stored)
	}
	return &apkpromotionchanneldto.DeleteApkPromotionChannelRes{Success: true}, nil
}
