package dbbackup

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/dbbackupdto"
	"xr-game-server/errercode"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/module/upload"
)

func GetCfg(_ context.Context, _ *dbbackupdto.GetDbBackupCfgReq) (*dbbackupdto.GetDbBackupCfgRes, error) {
	return &dbbackupdto.GetDbBackupCfgRes{Cfg: toCfgItem(cfgdao.GetDbBackupCfgCached())}, nil
}

func SaveCfg(_ context.Context, req *dbbackupdto.SaveDbBackupCfgReq) (*dbbackupdto.SaveDbBackupCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := cfgdao.GetDbBackupCfgCached()
	if row == nil || row.ID == 0 {
		row = &sysentity.DbBackupCfg{}
	}
	if req.ID > 0 && row.ID != 0 && req.ID != row.ID {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.Enabled = req.Enabled
	if req.RetainDays > 0 {
		row.RetainDays = req.RetainDays
	}
	ensureStoragePrefix(row)
	if err := persistCfg(row); err != nil {
		return nil, err
	}
	return &dbbackupdto.SaveDbBackupCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func RunNow(_ context.Context, _ *dbbackupdto.RunDbBackupNowReq) (*dbbackupdto.RunDbBackupNowRes, error) {
	row := cfgdao.GetDbBackupCfgCached()
	if row == nil || row.ID == 0 {
		row = &sysentity.DbBackupCfg{RetainDays: 1}
	}
	ensureStoragePrefix(row)
	if err := persistCfg(row); err != nil {
		return nil, err
	}
	objectKey, size, err := runBackupOnce(row, true)
	if err != nil {
		return nil, err
	}
	return &dbbackupdto.RunDbBackupNowRes{
		Success:   true,
		ObjectKey: objectKey,
		FileSize:  size,
		Message:   fmt.Sprintf("备份成功: %s (%d bytes)", objectKey, size),
	}, nil
}

func ListFiles(_ context.Context, _ *dbbackupdto.ListDbBackupFilesReq) (*dbbackupdto.ListDbBackupFilesRes, error) {
	row := cfgdao.GetDbBackupCfgCached()
	if row == nil || strings.TrimSpace(row.StoragePrefix) == "" {
		return &dbbackupdto.ListDbBackupFilesRes{List: []*dbbackupdto.DbBackupFileItem{}}, nil
	}
	if !upload.IsS3Enabled() {
		return nil, fmt.Errorf("云桶未开启,无法列出备份")
	}
	objs, err := upload.ListStoredObjectsByPrefix(row.StoragePrefix)
	if err != nil {
		return nil, err
	}
	list := make([]*dbbackupdto.DbBackupFileItem, 0, len(objs))
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj.StoredName == "" || strings.HasSuffix(obj.StoredName, "/") {
			continue
		}
		name := obj.StoredName
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		list = append(list, &dbbackupdto.DbBackupFileItem{
			ObjectKey:    obj.StoredName,
			FileName:     name,
			Size:         obj.Size,
			LastModified: formatCfgTime(obj.LastModified.UTC()),
		})
	}
	return &dbbackupdto.ListDbBackupFilesRes{List: list}, nil
}

func Restore(_ context.Context, req *dbbackupdto.RestoreDbBackupReq) (*dbbackupdto.RestoreDbBackupRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := cfgdao.GetDbBackupCfgCached()
	if row == nil || strings.TrimSpace(row.StoragePrefix) == "" {
		return nil, fmt.Errorf("备份目录尚未初始化,请先保存配置或执行一次备份")
	}
	if err := restoreFromObject(row, req.ObjectKey, req.TargetDatabase); err != nil {
		return nil, err
	}
	return &dbbackupdto.RestoreDbBackupRes{
		Success: true,
		Message: fmt.Sprintf("已还原到数据库 %s", strings.TrimSpace(req.TargetDatabase)),
	}, nil
}
