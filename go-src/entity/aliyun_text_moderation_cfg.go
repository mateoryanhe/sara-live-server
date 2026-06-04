package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAliyunTextModerationCfg db.TbName = "aliyun_text_moderation_cfgs"
)

// AliyunTextModerationCfg 阿里云文本审核(敏感词)配置，CMS 管理，通常仅一条
type AliyunTextModerationCfg struct {
	migrate.OneModel
	Enabled         bool   `gorm:"default:0;comment:是否开启文本敏感词过滤" json:"enabled"`
	AccessKeyId     string `gorm:"size:128;default:'';comment:AccessKeyId" json:"accessKeyId"`
	AccessKeySecret string `gorm:"size:256;default:'';comment:AccessKeySecret" json:"accessKeySecret"`
	RegionId        string `gorm:"size:32;default:'cn-shanghai';comment:地域" json:"regionId"`
	Endpoint        string `gorm:"size:128;default:'green-cip.cn-shanghai.aliyuncs.com';comment:接入点" json:"endpoint"`
	ChatService     string `gorm:"size:64;default:'chat_detection';comment:公聊/私信场景 Service" json:"chatService"`
	NicknameService string `gorm:"size:64;default:'nickname_detection';comment:昵称场景 Service" json:"nicknameService"`
	CommentService  string `gorm:"size:64;default:'comment_detection';comment:评论/公告场景 Service" json:"commentService"`
}

func initAliyunTextModerationCfg() {
	migrate.AutoMigrate(&AliyunTextModerationCfg{})
}
