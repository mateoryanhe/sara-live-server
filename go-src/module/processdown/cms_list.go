package processdown

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/processdowndao"
	"xr-game-server/dto/processdowndto"
	sysentity "xr-game-server/entity/sys"
)

func GetCMSList(ctx context.Context, req *processdowndto.CMSListProcessDownLogReq) (*httpserver.CMSQueryResp, error) {
	_ = ctx
	if req == nil {
		req = &processdowndto.CMSListProcessDownLogReq{}
	}
	total, rows := processdowndao.List(&processdowndao.ListFilter{
		Host:      req.Host,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*processdowndto.CMSProcessDownLogItem, 0, len(rows))
	for _, row := range rows {
		if item := toItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toItem(row *sysentity.ProcessDownLog) *processdowndto.CMSProcessDownLogItem {
	if row == nil {
		return nil
	}
	return &processdowndto.CMSProcessDownLogItem{
		ID:           row.ID,
		Host:         row.Host,
		AppName:      row.AppName,
		MissingSince: formatTime(row.MissingSince),
		DetectedAt:   formatTime(row.DetectedAt),
		ColdStartOK:  row.ColdStartOK,
		NewPID:       row.NewPID,
		Detail:       row.Detail,
		CreatedAt:    formatTime(row.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}
