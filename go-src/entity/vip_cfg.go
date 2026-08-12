package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbVipCfg db.TbName = "vip_cfgs"
)

const (
	VipCfgSwitchOff uint8 = 0 // 关(仅App端展示控制)
	VipCfgSwitchOn  uint8 = 1 // 开(仅App端展示控制)
)

// VipCfg VIP等级配置(CMS管理,App端下发见 vipcfgdto.AppVipCfgItem)
type VipCfg struct {
	migrate.OneModel
	Level                 uint32  `gorm:"uniqueIndex;default:0;comment:VIP等级" json:"level"`
	LevelName             string  `gorm:"size:64;default:'';comment:等级名称" json:"levelName"`
	LevelIcon             string  `gorm:"size:255;default:'';comment:VIP等级图标资源文件名" json:"levelIcon"`
	UpgradeRechargeLimit  float64 `gorm:"type:decimal(18,4);default:0;comment:升级充值上限(USD)" json:"upgradeRechargeLimit"`
	WithdrawSwitch        uint8   `gorm:"default:0;comment:提现开关(0关,1开,仅App端使用)" json:"withdrawSwitch"`
	MinWithdrawAmount     float64 `gorm:"type:decimal(18,4);default:0;comment:最低提现金额(USD)" json:"minWithdrawAmount"`
	MaxWithdrawAmount     float64 `gorm:"type:decimal(18,4);default:0;comment:最高提现金额(USD)" json:"maxWithdrawAmount"`
	Fee                   float64 `gorm:"type:decimal(18,4);default:0;comment:手续费" json:"fee"`
	WithdrawIcon          string  `gorm:"size:255;default:'';comment:提现图标资源文件名" json:"withdrawIcon"`
	WithdrawNoticeEn      string  `gorm:"type:text;comment:提现须知(英文)" json:"withdrawNoticeEn"`
	WithdrawNoticeEs      string  `gorm:"type:text;comment:提现须知(西班牙语)" json:"withdrawNoticeEs"`
	WithdrawNoticePt      string  `gorm:"type:text;comment:提现须知(葡萄牙语)" json:"withdrawNoticePt"`
	WithdrawNoticeHi      string  `gorm:"type:text;comment:提现须知(印地语)" json:"withdrawNoticeHi"`
	WithdrawNoticeId      string  `gorm:"type:text;comment:提现须知(印尼语)" json:"withdrawNoticeId"`
	AnimationSwitch       uint8   `gorm:"default:0;comment:进场特效开关(0关,1开,仅App端使用)" json:"animationSwitch"`
	Animation             string  `gorm:"size:255;default:'';comment:进场特效动画资源文件名" json:"animation"`
	AnimationIcon         string  `gorm:"size:255;default:'';comment:进场特效图标资源文件名" json:"animationIcon"`
	AnimationDescEn       string  `gorm:"type:text;comment:进场特效说明(英文)" json:"animationDescEn"`
	AnimationDescEs       string  `gorm:"type:text;comment:进场特效说明(西班牙语)" json:"animationDescEs"`
	AnimationDescPt       string  `gorm:"type:text;comment:进场特效说明(葡萄牙语)" json:"animationDescPt"`
	AnimationDescHi       string  `gorm:"type:text;comment:进场特效说明(印地语)" json:"animationDescHi"`
	AnimationDescId       string  `gorm:"type:text;comment:进场特效说明(印尼语)" json:"animationDescId"`
	CommentEffectSwitch   uint8   `gorm:"default:0;comment:公屏评论特效开关(0关,1开,仅App端使用)" json:"commentEffectSwitch"`
	CommentEffect         string  `gorm:"size:255;default:'';comment:公屏评论特效动画资源文件名" json:"commentEffect"`
	CommentEffectIcon     string  `gorm:"size:255;default:'';comment:公屏评论特效图标资源文件名" json:"commentEffectIcon"`
	CommentEffectDescEn   string  `gorm:"type:text;comment:公屏评论特效说明(英文)" json:"commentEffectDescEn"`
	CommentEffectDescEs   string  `gorm:"type:text;comment:公屏评论特效说明(西班牙语)" json:"commentEffectDescEs"`
	CommentEffectDescPt   string  `gorm:"type:text;comment:公屏评论特效说明(葡萄牙语)" json:"commentEffectDescPt"`
	CommentEffectDescHi   string  `gorm:"type:text;comment:公屏评论特效说明(印地语)" json:"commentEffectDescHi"`
	CommentEffectDescId   string  `gorm:"type:text;comment:公屏评论特效说明(印尼语)" json:"commentEffectDescId"`
	CustomerServiceSwitch uint8   `gorm:"default:0;comment:客服优先开关(0关,1开,仅App端使用)" json:"customerServiceSwitch"`
	CustomerServiceIcon   string  `gorm:"size:255;default:'';comment:客服优先图标资源文件名" json:"customerServiceIcon"`
	CustomerServiceDescEn string  `gorm:"type:text;comment:客服优先说明(英文)" json:"customerServiceDescEn"`
	CustomerServiceDescEs string  `gorm:"type:text;comment:客服优先说明(西班牙语)" json:"customerServiceDescEs"`
	CustomerServiceDescPt string  `gorm:"type:text;comment:客服优先说明(葡萄牙语)" json:"customerServiceDescPt"`
	CustomerServiceDescHi string  `gorm:"type:text;comment:客服优先说明(印地语)" json:"customerServiceDescHi"`
	CustomerServiceDescId string  `gorm:"type:text;comment:客服优先说明(印尼语)" json:"customerServiceDescId"`
}

func initVipCfg() {
	migrate.AutoMigrate(&VipCfg{})
}
