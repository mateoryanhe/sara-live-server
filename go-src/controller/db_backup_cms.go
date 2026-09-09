package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/dbbackupdto"
	"xr-game-server/module/dbbackup"
)

const DbBackupCMSUrl = "/dbBackup"

type DbBackupCMSController struct{}

func initDbBackupCMSController() {
	httpserver.RegCMS(DbBackupCMSUrl, &DbBackupCMSController{})
}

func (c *DbBackupCMSController) GetDbBackupCfg(ctx context.Context, req *dbbackupdto.GetDbBackupCfgReq) (*dbbackupdto.GetDbBackupCfgRes, error) {
	return dbbackup.GetCfg(ctx, req)
}

func (c *DbBackupCMSController) SaveDbBackupCfg(ctx context.Context, req *dbbackupdto.SaveDbBackupCfgReq) (*dbbackupdto.SaveDbBackupCfgRes, error) {
	return dbbackup.SaveCfg(ctx, req)
}

func (c *DbBackupCMSController) RunDbBackupNow(ctx context.Context, req *dbbackupdto.RunDbBackupNowReq) (*dbbackupdto.RunDbBackupNowRes, error) {
	return dbbackup.RunNow(ctx, req)
}

func (c *DbBackupCMSController) ListDbBackupFiles(ctx context.Context, req *dbbackupdto.ListDbBackupFilesReq) (*dbbackupdto.ListDbBackupFilesRes, error) {
	return dbbackup.ListFiles(ctx, req)
}

func (c *DbBackupCMSController) RestoreDbBackup(ctx context.Context, req *dbbackupdto.RestoreDbBackupReq) (*dbbackupdto.RestoreDbBackupRes, error) {
	return dbbackup.Restore(ctx, req)
}
