package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbDbBackupCfg db.TbName = "db_backup_cfgs"
)

// DbBackupCfg 数据库定时备份配置(CMS 管理,通常仅一条)
type DbBackupCfg struct {
	migrate.OneModel
	Enabled        bool      `gorm:"default:0;comment:是否启用每日0点备份" json:"enabled"`
	StoragePrefix  string    `gorm:"size:128;default:'';comment:云桶相对目录(首次生成guid/db_backup)" json:"storagePrefix"`
	RetainDays     int       `gorm:"default:1;comment:云端保留天数" json:"retainDays"`
	LastSuccessAt  time.Time `gorm:"comment:最近成功时间" json:"lastSuccessAt"`
	LastError      string    `gorm:"size:512;default:'';comment:最近错误摘要" json:"lastError"`
	LastObjectKey  string    `gorm:"size:255;default:'';comment:最近成功备份相对路径" json:"lastObjectKey"`
	LastFileSize   int64     `gorm:"default:0;comment:最近成功备份字节数" json:"lastFileSize"`
	LastTargetHint string    `gorm:"size:64;default:'';comment:最近备份来源库名" json:"lastTargetHint"`
}

func (DbBackupCfg) TableName() string {
	return string(TbDbBackupCfg)
}

func initDbBackupCfg() {
	migrate.AutoMigrate(&DbBackupCfg{})
}
