package entity

import (
	"strings"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFirebasePushCfg db.TbName = "firebase_push_cfgs"
)

// Firebase Push 默认文案/图标(CMS 空值时回填; 推送侧同样回落).
const (
	DefaultFirebasePushOfflineCallTitle = "Incoming video call"
	DefaultFirebasePushOfflineCallBody  = "{nickname} is calling you"
	// DefaultFirebasePushOfflineCallIcon 留空表示走 App 默认图标; CMS 可改成 CDN URL.
	DefaultFirebasePushOfflineCallIcon = ""
)

// FirebasePushCfg Firebase 推送文案等配置(CMS 管理,通常仅一条; FCM 凭证复用 firebase_cfgs 服务账号).
type FirebasePushCfg struct {
	migrate.OneModel
	Enabled           uint8  `gorm:"default:0;comment:是否启用FCM推送(0关1开)" json:"enabled"`
	OfflineCallTitle  string `gorm:"size:128;default:'';comment:1v1主播离线推送标题" json:"offlineCallTitle"`
	OfflineCallBody   string `gorm:"size:512;default:'';comment:1v1主播离线推送文案" json:"offlineCallBody"`
	OfflineCallIcon   string `gorm:"size:512;default:'';comment:1v1主播离线推送图标URL" json:"offlineCallIcon"`
}

func (FirebasePushCfg) TableName() string {
	return string(TbFirebasePushCfg)
}

func (c *FirebasePushCfg) IsEnabled() bool {
	return c != nil && c.ID > 0 && c.Enabled == 1
}

// ApplyFirebasePushDefaults 空字段回填默认值(不改 Enabled).
func (c *FirebasePushCfg) ApplyFirebasePushDefaults() {
	if c == nil {
		return
	}
	if strings.TrimSpace(c.OfflineCallTitle) == "" {
		c.OfflineCallTitle = DefaultFirebasePushOfflineCallTitle
	}
	if strings.TrimSpace(c.OfflineCallBody) == "" {
		c.OfflineCallBody = DefaultFirebasePushOfflineCallBody
	}
	if strings.TrimSpace(c.OfflineCallIcon) == "" {
		c.OfflineCallIcon = DefaultFirebasePushOfflineCallIcon
	}
}

func NewDefaultFirebasePushCfg() *FirebasePushCfg {
	row := &FirebasePushCfg{
		Enabled:          0,
		OfflineCallTitle: DefaultFirebasePushOfflineCallTitle,
		OfflineCallBody:  DefaultFirebasePushOfflineCallBody,
		OfflineCallIcon:  DefaultFirebasePushOfflineCallIcon,
	}
	return row
}

func initFirebasePushCfg() {
	migrate.AutoMigrate(&FirebasePushCfg{})
}
