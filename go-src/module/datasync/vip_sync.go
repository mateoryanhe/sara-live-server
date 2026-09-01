package datasync

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/entity/live"
	"xr-game-server/module/upload"
	"xr-game-server/module/vip"
)

// SyncVipCfg 从当前环境读取指定 VIP 配置及资源文件,推送到目标环境
func SyncVipCfg(_ context.Context, req *datasyncdto.SyncVipCfgReq) (*datasyncdto.SyncVipCfgRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}

	rows := cfgdao.GetVipCfgsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	fileNames := collectVipCfgFileNames(rows)
	files := make([]*datasyncdto.SyncFileItem, 0, len(fileNames))
	for _, name := range fileNames {
		content, err := upload.ReadUploadedFileBytes(name)
		if err != nil {
			return nil, fmt.Errorf("read resource file %s: %w", name, err)
		}
		files = append(files, &datasyncdto.SyncFileItem{
			Name: name,
			Data: base64.StdEncoding.EncodeToString(content),
		})
	}

	payload := &datasyncdto.ReceiveVipCfgReq{
		Rows:  rows,
		Files: files,
	}
	var receiveRes datasyncdto.ReceiveVipCfgRes
	if err := postSyncReceive("/dataSync/receiveVipCfg", payload, &receiveRes); err != nil {
		return nil, err
	}

	return &datasyncdto.SyncVipCfgRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条配置、%d 个资源文件", receiveRes.RowCount, receiveRes.FileCount),
	}, nil
}

// SyncVipCfgAssets 同步所选 VIP 配置的图标/动画等资源文件,并更新目标环境库中资源名字段后刷新缓存。
func SyncVipCfgAssets(_ context.Context, req *datasyncdto.SyncVipCfgAssetsReq) (*datasyncdto.SyncVipCfgRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errInvalidParam()
	}

	rows := cfgdao.GetVipCfgsByIDs(req.IDs)
	if len(rows) == 0 {
		return nil, errInvalidParam()
	}
	fileNames := collectVipCfgFileNames(rows)
	files := make([]*datasyncdto.SyncFileItem, 0, len(fileNames))
	for _, name := range fileNames {
		content, err := upload.ReadUploadedFileBytes(name)
		if err != nil {
			return nil, fmt.Errorf("read resource file %s: %w", name, err)
		}
		files = append(files, &datasyncdto.SyncFileItem{
			Name: name,
			Data: base64.StdEncoding.EncodeToString(content),
		})
	}

	payload := &datasyncdto.ReceiveVipCfgReq{Rows: rows, Files: files, AssetsOnly: true}
	var receiveRes datasyncdto.ReceiveVipCfgRes
	if err := postSyncReceive("/dataSync/receiveVipCfg", payload, &receiveRes); err != nil {
		return nil, err
	}

	return &datasyncdto.SyncVipCfgRes{
		Success:   receiveRes.Success,
		RowCount:  receiveRes.RowCount,
		FileCount: receiveRes.FileCount,
		Message:   fmt.Sprintf("已同步 %d 条VIP资源名、%d 个文件", receiveRes.RowCount, receiveRes.FileCount),
	}, nil
}

// ReceiveVipCfg 接收 VIP 配置同步,按主键 upsert 并写入资源文件
func ReceiveVipCfg(_ context.Context, req *datasyncdto.ReceiveVipCfgReq) (*datasyncdto.ReceiveVipCfgRes, error) {
	if req == nil {
		return nil, errInvalidParam()
	}

	fileCount := 0
	for _, item := range req.Files {
		if item == nil || strings.TrimSpace(item.Name) == "" || item.Data == "" {
			continue
		}
		content, err := base64.StdEncoding.DecodeString(item.Data)
		if err != nil {
			return nil, fmt.Errorf("decode file %s: %w", item.Name, err)
		}
		if err := upload.SaveUploadedFileBytes(item.Name, content); err != nil {
			return nil, fmt.Errorf("save file %s: %w", item.Name, err)
		}
		fileCount++
	}

	var rowCount int
	var err error
	if req.AssetsOnly {
		rowCount, err = saveVipCfgAssets(req.Rows)
	} else {
		rowCount, err = saveVipCfgs(req.Rows)
	}
	if err != nil {
		return nil, err
	}

	vip.ReloadVipCfgMemory()

	return &datasyncdto.ReceiveVipCfgRes{
		Success:   true,
		RowCount:  rowCount,
		FileCount: fileCount,
	}, nil
}

func saveVipCfgs(rows []*entity.VipCfg) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if err := cfgdao.CreateVipCfg(row); err != nil {
			return rowCount, fmt.Errorf("save vip cfg id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func saveVipCfgAssets(rows []*entity.VipCfg) (int, error) {
	rowCount := 0
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		_, err := g.DB().Model(string(entity.TbVipCfg)).
			WherePri(row.ID).
			Data(g.Map{
				"level_icon":              row.LevelIcon,
				"animation":               row.Animation,
				"animation_icon":          row.AnimationIcon,
				"comment_effect":          row.CommentEffect,
				"comment_effect_icon":     row.CommentEffectIcon,
				"customer_service_icon":   row.CustomerServiceIcon,
			}).
			Update()
		if err != nil {
			return rowCount, fmt.Errorf("update vip cfg assets id=%d: %w", row.ID, err)
		}
		rowCount++
	}
	return rowCount, nil
}

func collectVipCfgFileNames(rows []*entity.VipCfg) []string {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		for _, name := range []string{
			row.LevelIcon,
			row.Animation,
			row.AnimationIcon,
			row.CommentEffect,
			row.CommentEffectIcon,
			row.CustomerServiceIcon,
		} {
			name = strings.TrimSpace(name)
			if name == "" || strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
				continue
			}
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
	}
	return names
}
