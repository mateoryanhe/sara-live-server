package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbHostTlsMapping db.TbName = "host_tls_mappings"
)

// HostTlsMapping 纯 Host→证书映射(无静态目录),供 API 多域名 HTTPS/SNI.
type HostTlsMapping struct {
	migrate.OneModel
	Domain   string `gorm:"size:253;uniqueIndex;default:'';comment:Host,不含协议与端口" json:"domain"`
	CertFile string `gorm:"size:1024;default:'';comment:HTTPS证书文件路径" json:"certFile"`
	KeyFile  string `gorm:"size:1024;default:'';comment:HTTPS私钥文件路径" json:"keyFile"`
}

func initHostTlsMapping() {
	migrate.AutoMigrate(&HostTlsMapping{})
}
