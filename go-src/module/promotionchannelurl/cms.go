package promotionchannelurl

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
	"xr-game-server/module/promotionattribution"
	"xr-game-server/module/upload"
)

func GetList(_ context.Context, req *promotionchannelurldto.PromotionChannelUrlListReq) (*promotionchannelurldto.PromotionChannelUrlListRes, error) {
	total, list := promotionchannelurldao.GetList(req)
	baseHost := H5LiveShareBaseHost()
	for _, item := range list {
		if item == nil {
			continue
		}
		shareURL := BuildShareURL(item.ChannelCode)
		item.ShareURL = shareURL
		item.URL = shareURL
		item.ApkDownloadURL = GetCachedApkURL(item.ChannelCode)
	}
	return &promotionchannelurldto.PromotionChannelUrlListRes{
		Total:         total,
		Data:          list,
		ShareBaseHost: baseHost,
	}, nil
}

func Create(_ context.Context, req *promotionchannelurldto.CreatePromotionChannelUrlReq) (*promotionchannelurldto.CreatePromotionChannelUrlRes, error) {
	channelName := strings.TrimSpace(req.ChannelName)
	if channelName == "" || req.ChannelCode <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if promotionchannelurldao.GetByChannelCode(req.ChannelCode) != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	h5AppCode, err := normalizeUniqueH5AppChannelCode(req.H5AppChannelCode, 0)
	if err != nil {
		return nil, err
	}
	if err := requireURLH5ApkLandingPairIfApkExists(h5AppCode, req.ChannelCode); err != nil {
		return nil, err
	}
	shareURL := BuildShareURL(req.ChannelCode)
	if shareURL == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.PromotionChannelUrl{
		ChannelName:      channelName,
		ChannelCode:      req.ChannelCode,
		H5AppChannelCode: h5AppCode,
		URL:              shareURL,
	}
	if err := promotionchannelurldao.Create(row); err != nil {
		return nil, err
	}
	SetCachedApkURL(row.ChannelCode, "")
	return &promotionchannelurldto.CreatePromotionChannelUrlRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *promotionchannelurldto.UpdatePromotionChannelUrlReq) (*promotionchannelurldto.UpdatePromotionChannelUrlRes, error) {
	row := promotionchannelurldao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	channelName := strings.TrimSpace(req.ChannelName)
	if channelName == "" || req.ChannelCode <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := promotionchannelurldao.GetByChannelCode(req.ChannelCode); existing != nil && existing.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	h5AppCode, err := normalizeUniqueH5AppChannelCode(req.H5AppChannelCode, row.ID)
	if err != nil {
		return nil, err
	}
	if err := validateURLUpdateH5ApkLandingPair(row, h5AppCode, req.ChannelCode); err != nil {
		return nil, err
	}
	shareURL := BuildShareURL(req.ChannelCode)
	if shareURL == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	oldCode := row.ChannelCode
	row.ChannelName = channelName
	row.ChannelCode = req.ChannelCode
	row.H5AppChannelCode = h5AppCode
	row.URL = shareURL
	if oldCode != req.ChannelCode {
		row.ApkFileName = ""
		row.ApkDownloadURL = ""
		RenameCachedChannel(oldCode, req.ChannelCode)
	}
	if err := promotionchannelurldao.Update(row); err != nil {
		return nil, err
	}
	return &promotionchannelurldto.UpdatePromotionChannelUrlRes{Success: true}, nil
}

func Delete(_ context.Context, req *promotionchannelurldto.DeletePromotionChannelUrlReq) (*promotionchannelurldto.DeletePromotionChannelUrlRes, error) {
	row := promotionchannelurldao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	stored := strings.TrimSpace(row.ApkFileName)
	code := row.ChannelCode
	if err := promotionchannelurldao.Delete(req.ID); err != nil {
		return nil, err
	}
	RemoveCachedChannel(code)
	if stored != "" {
		upload.DeletePromotionChannelUrlApkStored(stored)
	}
	return &promotionchannelurldto.DeletePromotionChannelUrlRes{Success: true}, nil
}

func normalizeUniqueH5AppChannelCode(raw string, excludeID uint64) (string, error) {
	code, ok := promotionattribution.NormalizeH5AppChannelCode(raw)
	if !ok {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := promotionchannelurldao.GetByH5AppChannelCode(code); existing != nil && existing.ID != excludeID {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return code, nil
}

// requireURLH5ApkLandingPair 编辑 H5推广渠道: H5-App 渠道码须与落地页推广 (h5_app_channel_code+channel_code) 配对, 且对应该 URL 整数 channel_code.
func requireURLH5ApkLandingPair(h5AppCode string, urlChannelCode int) error {
	if urlChannelCode <= 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	rows := apkpromotionchanneldao.ListByH5AppChannelCode(h5AppCode)
	if len(rows) == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	for _, apk := range rows {
		if apk == nil {
			continue
		}
		ch, ok := promotionattribution.ValidatePair(h5AppCode, apk.ChannelCode)
		if ok && ch == urlChannelCode {
			return nil
		}
	}
	return errercode.CreateCode(errercode.InvalidParam)
}

func requireURLH5ApkLandingPairIfApkExists(h5AppCode string, urlChannelCode int) error {
	if len(apkpromotionchanneldao.ListByH5AppChannelCode(h5AppCode)) == 0 {
		return nil
	}
	return requireURLH5ApkLandingPair(h5AppCode, urlChannelCode)
}

func validateURLUpdateH5ApkLandingPair(row *entity.PromotionChannelUrl, h5AppCode string, urlChannelCode int) error {
	if row == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	oldH5 := strings.TrimSpace(row.H5AppChannelCode)
	h5Changed := oldH5 != h5AppCode
	chChanged := row.ChannelCode != urlChannelCode
	// 仅改渠道名称等: H5-App 码与整数渠道码未变时不做配对校验.
	if !h5Changed && !chChanged {
		return nil
	}
	if h5Changed && oldH5 != "" && len(apkpromotionchanneldao.ListByH5AppChannelCode(oldH5)) > 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return requireURLH5ApkLandingPairIfApkExists(h5AppCode, urlChannelCode)
}
