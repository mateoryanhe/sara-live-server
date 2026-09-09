package dbbackupdto

import "github.com/gogf/gf/v2/frame/g"

type GetDbBackupCfgReq struct {
	g.Meta `path:"/getDbBackupCfg" method:"post" summary:"查询数据库备份配置" tags:"数据库备份"`
}

type DbBackupCfgItem struct {
	ID             string `json:"id"`
	Enabled        bool   `json:"enabled"`
	StoragePrefix  string `json:"storagePrefix"`
	RetainDays     int    `json:"retainDays"`
	LastSuccessAt  string `json:"lastSuccessAt"`
	LastError      string `json:"lastError"`
	LastObjectKey  string `json:"lastObjectKey"`
	LastFileSize   int64  `json:"lastFileSize"`
	LastTargetHint string `json:"lastTargetHint"`
	S3Enabled      bool   `json:"s3Enabled"`
	SourceDatabase string `json:"sourceDatabase"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type GetDbBackupCfgRes struct {
	Cfg *DbBackupCfgItem `json:"cfg"`
}

type SaveDbBackupCfgReq struct {
	g.Meta     `path:"/saveDbBackupCfg" method:"post" summary:"保存数据库备份配置" tags:"数据库备份"`
	ID         uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled    bool   `json:"enabled" dc:"是否启用每日0点备份"`
	RetainDays int    `json:"retainDays" dc:"云端保留天数,默认1"`
}

type SaveDbBackupCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type RunDbBackupNowReq struct {
	g.Meta `path:"/runDbBackupNow" method:"post" summary:"立即执行一次数据库备份" tags:"数据库备份"`
}

type RunDbBackupNowRes struct {
	Success   bool   `json:"success"`
	ObjectKey string `json:"objectKey"`
	FileSize  int64  `json:"fileSize"`
	Message   string `json:"message"`
}

type ListDbBackupFilesReq struct {
	g.Meta `path:"/listDbBackupFiles" method:"post" summary:"列出云端数据库备份文件" tags:"数据库备份"`
}

type DbBackupFileItem struct {
	ObjectKey    string `json:"objectKey"`
	FileName     string `json:"fileName"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
}

type ListDbBackupFilesRes struct {
	List []*DbBackupFileItem `json:"list"`
}

type RestoreDbBackupReq struct {
	g.Meta         `path:"/restoreDbBackup" method:"post" summary:"从云端备份还原到指定数据库" tags:"数据库备份"`
	ObjectKey      string `json:"objectKey" v:"required#请选择备份文件" dc:"备份相对路径"`
	TargetDatabase string `json:"targetDatabase" v:"required#请填写目标数据库名" dc:"目标数据库名"`
}

type RestoreDbBackupRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
