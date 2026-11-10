package cmsexport

import (
	"context"
	"encoding/json"
	"strconv"

	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/module/userinfo"
)

func exportUserListCSV(ctx context.Context, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	var req cmsexportdto.CMSExportUserListPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	return streamCSVExport(ctx, req.Headers, defaultExportPageSize, func(pageIndex, pageSize int) (int, [][]string) {
		query := &accountdto.QueryUserInfoReq{
			Key:               req.Key,
			StartTime:         req.StartTime,
			EndTime:           req.EndTime,
			RechargeWhitelist: req.RechargeWhitelist,
			IsAnchor:          req.IsAnchor,
		}
		query.PageIndex = pageIndex
		query.PageSize = pageSize
		resp, err := userinfo.QueryUserInfo(ctx, query)
		if err != nil || resp == nil {
			return 0, nil
		}
		list, _ := resp.Data.([]*accountdto.UserInfoDto)
		csvRows := make([][]string, 0, len(list))
		for _, row := range list {
			csvRows = append(csvRows, userInfoToCSVRow(row, &req))
		}
		return resp.Total, csvRows
	}, onProgress)
}

func userInfoToCSVRow(row *accountdto.UserInfoDto, req *cmsexportdto.CMSExportUserListPayload) []string {
	if row == nil {
		return nil
	}
	inviter := ""
	if row.InviterId > 0 {
		inviter = formatCSVUint(row.InviterId)
	}
	yesNo := req.YesNoLabels
	return []string{
		formatCSVUint(row.ID),
		inviter,
		formatCSVTimePtr(row.CreatedAt),
		formatCSVTimePtr(row.LastLoginTime),
		row.Nickname,
		formatCSVFloat(row.Gold),
		formatCSVFloat(row.Diamond),
		strconv.FormatUint(uint64(row.VipLevel), 10),
		lookupExportLabel(req.UserTypeLabels, row.UserType),
		lookupBoolLabel(yesNo, row.IsAnchor),
		dashIfEmpty(row.IP),
		dashIfEmpty(row.LoginCountry),
		dashIfEmpty(row.RegisterIp),
		dashIfEmpty(row.RegisterCountry),
		dashIfEmpty(row.DeviceType),
		dashIfEmpty(row.PackageName),
		dashIfEmpty(row.AppVersion),
		lookupBoolLabel(req.BanLabels, row.Ban),
		lookupBoolLabel(req.CancelLabels, row.Cancel),
		lookupBoolLabel(yesNo, row.RechargeWhitelist),
		formatCSVTimePtr(row.BanApplyTime),
		row.OpenId,
		strconv.FormatUint(uint64(row.Channel), 10),
		dashIfEmpty(row.Remark),
	}
}

func lookupBoolLabel(labels map[string]string, v bool) string {
	key := "0"
	if v {
		key = "1"
	}
	if labels != nil {
		if text, ok := labels[key]; ok && text != "" {
			return text
		}
	}
	return key
}

func dashIfEmpty(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
