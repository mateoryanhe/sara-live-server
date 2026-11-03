package apkpromotionchannel

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dto/apkpromotionchanneldto"
	entity "xr-game-server/entity/cms"
	"xr-game-server/errercode"
	"xr-game-server/module/promotionattribution"
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
	h5AppCode, err := normalizeApkH5AppChannelCode(req.H5AppChannelCode, code, 0)
	if err != nil {
		return nil, err
	}
	if BuildChannelShortURL(code) == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.ApkPromotionChannel{
		ChannelName:      channelName,
		H5AppChannelCode: h5AppCode,
		ChannelCode:      code,
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
	h5AppCode, err := normalizeApkH5AppChannelCode(req.H5AppChannelCode, code, row.ID)
	if err != nil {
		return nil, err
	}
	row.ChannelName = channelName
	row.H5AppChannelCode = h5AppCode
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

func normalizeApkH5AppChannelCode(raw, apkChannelCode string, excludeID uint64) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	h5AppCode, ok := promotionattribution.NormalizeH5AppChannelCode(raw)
	if !ok {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	urlRow := promotionchannelurldao.GetByH5AppChannelCode(h5AppCode)
	if urlRow == nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	if strings.TrimSpace(urlRow.H5AppChannelCode) != h5AppCode {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	// 复合唯一 (h5_app_channel_code, channel_code) 由 DB 约束; 先查同 APK 码是否被其它 h5 占用
	if existing := apkpromotionchanneldao.GetByChannelCode(apkChannelCode); existing != nil && existing.ID != excludeID {
		if strings.TrimSpace(existing.H5AppChannelCode) != h5AppCode {
			return "", errercode.CreateCode(errercode.InvalidParam)
		}
	}
	return h5AppCode, nil
}
