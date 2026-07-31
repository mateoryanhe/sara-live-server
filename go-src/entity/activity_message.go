package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbActivityMessage db.TbName = "activity_messages"
)

const (
	ActivityMessageStatusUnpublished uint8 = 0 // 未发布
	ActivityMessageStatusPublished   uint8 = 1 // 已发布
)

// ActivityMessage 活动消息(CMS 管理)
type ActivityMessage struct {
	migrate.OneModel
	IconEn      string     `gorm:"size:255;default:'';comment:图标资源名(英文)" json:"iconEn"`
	IconEs      string     `gorm:"size:255;default:'';comment:图标资源名(西班牙语)" json:"iconEs"`
	IconPt      string     `gorm:"size:255;default:'';comment:图标资源名(葡萄牙语)" json:"iconPt"`
	IconHi      string     `gorm:"size:255;default:'';comment:图标资源名(印地语)" json:"iconHi"`
	BgEn        string     `gorm:"size:255;default:'';comment:背景图资源名(英文)" json:"bgEn"`
	BgEs        string     `gorm:"size:255;default:'';comment:背景图资源名(西班牙语)" json:"bgEs"`
	BgPt        string     `gorm:"size:255;default:'';comment:背景图资源名(葡萄牙语)" json:"bgPt"`
	BgHi        string     `gorm:"size:255;default:'';comment:背景图资源名(印地语)" json:"bgHi"`
	TitleEn     string     `gorm:"size:128;default:'';comment:标题(英文)" json:"titleEn"`
	TitleEs     string     `gorm:"size:128;default:'';comment:标题(西班牙语)" json:"titleEs"`
	TitlePt     string     `gorm:"size:128;default:'';comment:标题(葡萄牙语)" json:"titlePt"`
	TitleHi     string     `gorm:"size:128;default:'';comment:标题(印地语)" json:"titleHi"`
	ContentEn   string     `gorm:"type:text;comment:内容(英文)" json:"contentEn"`
	ContentEs   string     `gorm:"type:text;comment:内容(西班牙语)" json:"contentEs"`
	ContentPt   string     `gorm:"type:text;comment:内容(葡萄牙语)" json:"contentPt"`
	ContentHi   string     `gorm:"type:text;comment:内容(印地语)" json:"contentHi"`
	UrlEn       string     `gorm:"size:512;default:'';comment:跳转链接(英文)" json:"urlEn"`
	UrlEs       string     `gorm:"size:512;default:'';comment:跳转链接(西班牙语)" json:"urlEs"`
	UrlPt       string     `gorm:"size:512;default:'';comment:跳转链接(葡萄牙语)" json:"urlPt"`
	UrlHi       string     `gorm:"size:512;default:'';comment:跳转链接(印地语)" json:"urlHi"`
	Status      uint8      `gorm:"default:0;comment:发布状态(0-未发布,1-已发布)" json:"status"`
	PublishedAt *time.Time `gorm:"comment:发布时间" json:"publishedAt"`
}

func initActivityMessage() {
	migrate.AutoMigrate(&ActivityMessage{})
}
