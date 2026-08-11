package datasync

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/entity"
)

func GetDataSyncCfg(_ context.Context, _ *datasyncdto.GetDataSyncCfgReq) (*datasyncdto.GetDataSyncCfgRes, error) {
	cfg := cfgdao.GetDataSyncCfg()
	if cfg == nil {
		return &datasyncdto.GetDataSyncCfgRes{Cfg: nil}, nil
	}
	return &datasyncdto.GetDataSyncCfgRes{Cfg: toDataSyncCfgItem(cfg)}, nil
}

func SaveDataSyncCfg(_ context.Context, req *datasyncdto.SaveDataSyncCfgReq) (*datasyncdto.SaveDataSyncCfgRes, error) {
	existing := cfgdao.GetDataSyncCfg()
	row := &entity.DataSyncCfg{
		TargetApiBase: strings.TrimSpace(req.TargetApiBase),
		Token:         strings.TrimSpace(req.Token),
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errInvalidParam()
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveDataSyncCfg(row); err != nil {
		return nil, err
	}
	return &datasyncdto.SaveDataSyncCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toDataSyncCfgItem(cfg *entity.DataSyncCfg) *datasyncdto.DataSyncCfgItem {
	if cfg == nil {
		return nil
	}
	return &datasyncdto.DataSyncCfgItem{
		ID:            strconv.FormatUint(cfg.ID, 10),
		TargetApiBase: cfg.TargetApiBase,
		Token:         cfg.Token,
		CreatedAt:     formatDataSyncCfgTime(cfg.CreatedAt),
		UpdatedAt:     formatDataSyncCfgTime(cfg.UpdatedAt),
	}
}

func formatDataSyncCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
