package simulatorblacklistlog

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/simulatorblacklistlogdao"
	"xr-game-server/dto/simulatorblacklistlogdto"
	userentity "xr-game-server/entity/user"
)

func GetCMSList(ctx context.Context, req *simulatorblacklistlogdto.CMSListSimulatorBlacklistLogReq) (*httpserver.CMSQueryResp, error) {
	_ = ctx
	if req == nil {
		req = &simulatorblacklistlogdto.CMSListSimulatorBlacklistLogReq{}
	}
	total, rows := simulatorblacklistlogdao.List(&simulatorblacklistlogdao.ListFilter{
		Key:       req.Key,
		LoginType: req.LoginType,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*simulatorblacklistlogdto.CMSSimulatorBlacklistLogItem, 0, len(rows))
	for _, row := range rows {
		if item := toItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toItem(row *userentity.SimulatorBlacklistLog) *simulatorblacklistlogdto.CMSSimulatorBlacklistLogItem {
	if row == nil {
		return nil
	}
	accountId := ""
	if row.AccountId > 0 {
		accountId = strconv.FormatUint(row.AccountId, 10)
	}
	return &simulatorblacklistlogdto.CMSSimulatorBlacklistLogItem{
		ID:          row.ID,
		IP:          row.IP,
		DeviceId:    row.DeviceId,
		CpuModel:    row.CpuModel,
		DeviceType:  row.DeviceType,
		DeviceModel: row.DeviceModel,
		OsVersion:   row.OsVersion,
		AppVersion:  row.AppVersion,
		PackageName: row.PackageName,
		LoginType:   row.LoginType,
		Email:       row.Email,
		AccountId:   accountId,
		RegisterAt:  formatTimePtr(row.RegisterAt),
		CreatedAt:   formatTime(row.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
