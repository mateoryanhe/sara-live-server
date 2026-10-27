package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAuditCodeCfg db.TbName = "audit_code_cfgs"
)

// AuditCodeCfg App 审核码/主播码(CMS 生成覆盖,通常仅一条)
type AuditCodeCfg struct {
	migrate.OneModel
	Code string `gorm:"size:64;not null;default:'';comment:审核码(guid小写无下划线)" json:"code"`
}

func initAuditCodeCfg() {
	migrate.AutoMigrate(&AuditCodeCfg{})
}
