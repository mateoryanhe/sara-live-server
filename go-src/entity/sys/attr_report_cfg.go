package entity

import (
	"strings"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAttrReportCfg db.TbName = "attr_report_cfgs"
)

// AttrReportCfg 归因数据上报配置(CMS 管理,通常仅一条).
type AttrReportCfg struct {
	migrate.OneModel
	BaseApi   string `gorm:"size:512;default:'';comment:归因上报主机根地址(自动拼/tj/api/accept.php)" json:"baseApi"`
	SecretKey string `gorm:"size:255;default:'';comment:归因上报密钥(api_token)" json:"secretKey"`
}

func (AttrReportCfg) TableName() string {
	return string(TbAttrReportCfg)
}

func (c *AttrReportCfg) IsConfigured() bool {
	return c != nil && c.ID > 0 &&
		strings.TrimSpace(c.BaseApi) != "" &&
		strings.TrimSpace(c.SecretKey) != ""
}

func initAttrReportCfg() {
	migrate.AutoMigrate(&AttrReportCfg{})
}
