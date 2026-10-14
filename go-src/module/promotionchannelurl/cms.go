package promotionchannelurl

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
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
	shareURL := BuildShareURL(req.ChannelCode)
	if shareURL == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.PromotionChannelUrl{
		ChannelName: channelName,
		ChannelCode: req.ChannelCode,
		URL:         shareURL,
	}
	if err := promotionchannelurldao.Create(row); err != nil {
		return nil, err
	}
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
	shareURL := BuildShareURL(req.ChannelCode)
	if shareURL == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.ChannelName = channelName
	row.ChannelCode = req.ChannelCode
	row.URL = shareURL
	if err := promotionchannelurldao.Update(row); err != nil {
		return nil, err
	}
	return &promotionchannelurldto.UpdatePromotionChannelUrlRes{Success: true}, nil
}

func Delete(_ context.Context, req *promotionchannelurldto.DeletePromotionChannelUrlReq) (*promotionchannelurldto.DeletePromotionChannelUrlRes, error) {
	if promotionchannelurldao.GetById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := promotionchannelurldao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &promotionchannelurldto.DeletePromotionChannelUrlRes{Success: true}, nil
}
