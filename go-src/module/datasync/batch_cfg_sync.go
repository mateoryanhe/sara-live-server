package datasync

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/entity/live"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/module/banner"
	"xr-game-server/module/gift"
	"xr-game-server/module/recharge"
)

func SyncBanner(_ context.Context, req *datasyncdto.SyncBannerReq) (*datasyncdto.SyncBatchRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}
	rows := cfgdao.GetBannersByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	files, err := buildSyncFiles(collectBannerFileNames(rows))
	if err != nil {
		return nil, err
	}
	payload := &datasyncdto.ReceiveBannerReq{Rows: rows, Files: files}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveBanner", payload, &receiveRes); err != nil {
		return nil, err
	}
	return newSyncBatchRes(&receiveRes, "Banner"), nil
}

func ReceiveBanner(_ context.Context, req *datasyncdto.ReceiveBannerReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}
	fileCount, err := saveSyncFiles(req.Files)
	if err != nil {
		return nil, err
	}
	rowCount, err := saveHomeBanners(req.Rows)
	if err != nil {
		return nil, err
	}
	banner.ReloadBannerMemory()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: rowCount, FileCount: fileCount}, nil
}

func SyncGift(_ context.Context, req *datasyncdto.SyncGiftReq) (*datasyncdto.SyncBatchRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}
	rows := cfgdao.GetGiftsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	files, err := buildSyncFiles(collectGiftFileNames(rows))
	if err != nil {
		return nil, err
	}
	payload := &datasyncdto.ReceiveGiftReq{Rows: rows, Files: files}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveGift", payload, &receiveRes); err != nil {
		return nil, err
	}
	return newSyncBatchRes(&receiveRes, "礼物"), nil
}

// SyncGiftAssets 同步所选礼物的图标/动画文件,并更新目标环境库中的 icon/animation 字段后刷新缓存。
func SyncGiftAssets(_ context.Context, req *datasyncdto.SyncGiftAssetsReq) (*datasyncdto.SyncBatchRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}
	rows := cfgdao.GetGiftsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	files, err := buildSyncFiles(collectGiftFileNames(rows))
	if err != nil {
		return nil, err
	}
	payload := &datasyncdto.ReceiveGiftReq{Rows: rows, Files: files, AssetsOnly: true}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveGift", payload, &receiveRes); err != nil {
		return nil, err
	}
	return &datasyncdto.SyncBatchRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条礼物资源名、%d 个文件", receiveRes.RowCount, receiveRes.FileCount),
	}, nil
}

func ReceiveGift(_ context.Context, req *datasyncdto.ReceiveGiftReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}
	fileCount, err := saveSyncFiles(req.Files)
	if err != nil {
		return nil, err
	}
	var rowCount int
	if req.AssetsOnly {
		rowCount, err = saveLiveGiftAssets(req.Rows)
	} else {
		rowCount, err = saveLiveGifts(req.Rows)
	}
	if err != nil {
		return nil, err
	}
	gift.ReloadGiftCache()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: rowCount, FileCount: fileCount}, nil
}

func SyncRechargeCfg(_ context.Context, req *datasyncdto.SyncRechargeCfgReq) (*datasyncdto.SyncBatchRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}
	rows := cfgdao.GetRechargeCfgsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	files, err := buildSyncFiles(collectRechargeCfgFileNames(rows))
	if err != nil {
		return nil, err
	}
	payload := &datasyncdto.ReceiveRechargeCfgReq{Rows: rows, Files: files}
	var receiveRes datasyncdto.ReceiveBatchRes
	if err := postSyncReceive("/dataSync/receiveRechargeCfg", payload, &receiveRes); err != nil {
		return nil, err
	}
	return newSyncBatchRes(&receiveRes, "充值配置"), nil
}

func ReceiveRechargeCfg(_ context.Context, req *datasyncdto.ReceiveRechargeCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}
	fileCount, err := saveSyncFiles(req.Files)
	if err != nil {
		return nil, err
	}
	rowCount, err := saveRechargeCfgs(req.Rows)
	if err != nil {
		return nil, err
	}
	recharge.ReloadRechargeCfgCache()
	return &datasyncdto.ReceiveBatchRes{Success: true, RowCount: rowCount, FileCount: fileCount}, nil
}

func saveHomeBanners(rows []*entity.HomeBanner) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if _, err := g.DB().Model(string(entity.TbHomeBanner)).Save(row); err != nil {
			return rowCount, fmt.Errorf("save banner id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func saveLiveGifts(rows []*entity.LiveGift) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if _, err := g.DB().Model(string(entity.TbLiveGift)).Save(row); err != nil {
			return rowCount, fmt.Errorf("save gift id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func saveLiveGiftAssets(rows []*entity.LiveGift) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		_, err := g.DB().Model(string(entity.TbLiveGift)).
			WherePri(row.ID).
			Data(g.Map{
				"icon":      row.Icon,
				"animation": row.Animation,
			}).
			Update()
		if err != nil {
			return rowCount, fmt.Errorf("update gift assets id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func saveRechargeCfgs(rows []*rechargeentity.RechargeCfg) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if _, err := g.DB().Model(string(rechargeentity.TbRechargeCfg)).Save(row); err != nil {
			return rowCount, fmt.Errorf("save recharge cfg id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func collectBannerFileNames(rows []*entity.HomeBanner) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		names = appendUniqueFileName(seen, names, row.Image)
	}
	return names
}

func collectGiftFileNames(rows []*entity.LiveGift) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		names = appendUniqueFileName(seen, names, row.Icon)
		names = appendUniqueFileName(seen, names, row.Animation)
	}
	return names
}

func collectRechargeCfgFileNames(rows []*rechargeentity.RechargeCfg) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		names = appendUniqueFileName(seen, names, row.Icon)
	}
	return names
}
