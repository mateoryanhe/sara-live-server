package message

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/activitymessagedto"
	"xr-game-server/entity/message"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func GetActivityMessageList(_ context.Context, req *activitymessagedto.ActivityMessageListReq) (*httpserver.CMSQueryResp, error) {
	total, list := messagedao.GetActivityMessageList(req)
	fillActivityMessageListAssets(list)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func fillActivityMessageListAssets(list []*activitymessagedto.ActivityMessageListRes) {
	for _, row := range list {
		if row == nil {
			continue
		}
		row.IconEnName = row.IconEn
		row.IconEsName = row.IconEs
		row.IconPtName = row.IconPt
		row.IconHiName = row.IconHi
		row.IconIdName = row.IconId
		row.IconEn = upload.GetUrlByName(row.IconEnName)
		row.IconEs = upload.GetUrlByName(row.IconEsName)
		row.IconPt = upload.GetUrlByName(row.IconPtName)
		row.IconHi = upload.GetUrlByName(row.IconHiName)
		row.IconId = upload.GetUrlByName(row.IconIdName)

		row.BgEnName = row.BgEn
		row.BgEsName = row.BgEs
		row.BgPtName = row.BgPt
		row.BgHiName = row.BgHi
		row.BgIdName = row.BgId
		row.BgEn = upload.GetUrlByName(row.BgEnName)
		row.BgEs = upload.GetUrlByName(row.BgEsName)
		row.BgPt = upload.GetUrlByName(row.BgPtName)
		row.BgHi = upload.GetUrlByName(row.BgHiName)
		row.BgId = upload.GetUrlByName(row.BgIdName)
	}
}

func CreateActivityMessage(_ context.Context, req *activitymessagedto.CreateActivityMessageReq) (*activitymessagedto.CreateActivityMessageRes, error) {
	row := &entity.ActivityMessage{
		IconEn:    req.IconEn,
		IconEs:    req.IconEs,
		IconPt:    req.IconPt,
		IconHi:    req.IconHi,
		IconId:    req.IconId,
		BgEn:      req.BgEn,
		BgEs:      req.BgEs,
		BgPt:      req.BgPt,
		BgHi:      req.BgHi,
		BgId:      req.BgId,
		TitleEn:   req.TitleEn,
		TitleEs:   req.TitleEs,
		TitlePt:   req.TitlePt,
		TitleHi:   req.TitleHi,
		TitleId:   req.TitleId,
		ContentEn: req.ContentEn,
		ContentEs: req.ContentEs,
		ContentPt: req.ContentPt,
		ContentHi: req.ContentHi,
		ContentId: req.ContentId,
		UrlEn:     req.UrlEn,
		UrlEs:     req.UrlEs,
		UrlPt:     req.UrlPt,
		UrlHi:     req.UrlHi,
		UrlId:     req.UrlId,
		Status:    entity.ActivityMessageStatusUnpublished,
	}
	if err := messagedao.CreateActivityMessage(row); err != nil {
		return nil, err
	}
	return &activitymessagedto.CreateActivityMessageRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateActivityMessage(_ context.Context, req *activitymessagedto.UpdateActivityMessageReq) (*activitymessagedto.UpdateActivityMessageRes, error) {
	row := messagedao.GetActivityMessageById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ActivityMessageNonExist)
	}
	row.IconEn = req.IconEn
	row.IconEs = req.IconEs
	row.IconPt = req.IconPt
	row.IconHi = req.IconHi
	row.IconId = req.IconId
	row.BgEn = req.BgEn
	row.BgEs = req.BgEs
	row.BgPt = req.BgPt
	row.BgHi = req.BgHi
	row.BgId = req.BgId
	row.TitleEn = req.TitleEn
	row.TitleEs = req.TitleEs
	row.TitlePt = req.TitlePt
	row.TitleHi = req.TitleHi
	row.TitleId = req.TitleId
	row.ContentEn = req.ContentEn
	row.ContentEs = req.ContentEs
	row.ContentPt = req.ContentPt
	row.ContentHi = req.ContentHi
	row.ContentId = req.ContentId
	row.UrlEn = req.UrlEn
	row.UrlEs = req.UrlEs
	row.UrlPt = req.UrlPt
	row.UrlHi = req.UrlHi
	row.UrlId = req.UrlId
	if err := messagedao.UpdateActivityMessage(row); err != nil {
		return nil, err
	}
	return &activitymessagedto.UpdateActivityMessageRes{Success: true}, nil
}

func DeleteActivityMessage(_ context.Context, req *activitymessagedto.DeleteActivityMessageReq) (*activitymessagedto.DeleteActivityMessageRes, error) {
	if row := messagedao.GetActivityMessageById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.ActivityMessageNonExist)
	}
	if err := messagedao.DeleteActivityMessage(req.ID); err != nil {
		return nil, err
	}
	removeUnpublishedActivityMessageFromCachedUsers(req.ID)
	if err := messagedao.RemoveUserActivityMessageByActivityMessageId(req.ID); err != nil {
		return nil, err
	}
	return &activitymessagedto.DeleteActivityMessageRes{Success: true}, nil
}

func PublishActivityMessage(_ context.Context, req *activitymessagedto.PublishActivityMessageReq) (*activitymessagedto.PublishActivityMessageRes, error) {
	row := messagedao.GetActivityMessageById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ActivityMessageNonExist)
	}
	if row.Status != entity.ActivityMessageStatusPublished {
		now := time.Now()
		row.Status = entity.ActivityMessageStatusPublished
		row.PublishedAt = &now
		if err := messagedao.UpdateActivityMessage(row); err != nil {
			return nil, err
		}
		prependPublishedActivityMessageToCachedUsers(row)
		push.BroadcastNonData(cmd.ActivityMessagePush)
	}
	return &activitymessagedto.PublishActivityMessageRes{Success: true, Status: entity.ActivityMessageStatusPublished}, nil
}

func UnpublishActivityMessage(_ context.Context, req *activitymessagedto.UnpublishActivityMessageReq) (*activitymessagedto.UnpublishActivityMessageRes, error) {
	row := messagedao.GetActivityMessageById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ActivityMessageNonExist)
	}
	if row.Status != entity.ActivityMessageStatusUnpublished {
		row.Status = entity.ActivityMessageStatusUnpublished
		row.PublishedAt = nil
		if err := messagedao.UpdateActivityMessage(row); err != nil {
			return nil, err
		}
		removeUnpublishedActivityMessageFromCachedUsers(req.ID)
		if err := messagedao.RemoveUserActivityMessageByActivityMessageId(req.ID); err != nil {
			return nil, err
		}
	}
	return &activitymessagedto.UnpublishActivityMessageRes{Success: true, Status: entity.ActivityMessageStatusUnpublished}, nil
}

// GetPublishedActivityMessages 获取已发布活动消息(读 gcache 全量后过滤)
func GetPublishedActivityMessages() []*entity.ActivityMessage {
	rows := messagedao.GetAllActivityMessagesCached()
	if len(rows) == 0 {
		return nil
	}
	ret := make([]*entity.ActivityMessage, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.Status != entity.ActivityMessageStatusPublished {
			continue
		}
		ret = append(ret, row)
	}
	return ret
}
