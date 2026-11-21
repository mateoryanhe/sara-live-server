package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAudienceOnlineAnchorGreetCfg db.TbName = "audience_online_anchor_greet_cfgs"
)

const (
	// DefaultAudienceOnlineAnchorGreetAnchorCount 未配置时默认由几名在线主播打招呼
	DefaultAudienceOnlineAnchorGreetAnchorCount uint8 = 5
	// AudienceOnlineAnchorGreetAnchorCountMax CMS 可设上限
	AudienceOnlineAnchorGreetAnchorCountMax uint8 = 20

	DefaultAudienceOnlineAnchorGreetIntervalMinSec uint32 = 30
	DefaultAudienceOnlineAnchorGreetIntervalMaxSec uint32 = 30
	// AudienceOnlineAnchorGreetIntervalSecMax 排队模式下两条招呼之间的最大间隔(秒)
	AudienceOnlineAnchorGreetIntervalSecMax uint32 = 3600
)

// AudienceOnlineAnchorGreetCfg 观众登录成功时由在线主播发送招呼私信(CMS 单条配置)
type AudienceOnlineAnchorGreetCfg struct {
	migrate.OneModel
	Enabled     bool   `gorm:"default:0;comment:是否开启" json:"enabled"`
	AnchorCount uint8  `gorm:"default:5;comment:每次最多几名在线主播打招呼(按1v1通话次数选取)" json:"anchorCount"`
	// IntervalMinSec/IntervalMaxSec 排队模式下相邻两条私信间隔秒数(含端点随机)
	IntervalMinSec uint32 `gorm:"default:30;comment:主播排队招呼间隔下限秒" json:"intervalMinSec"`
	IntervalMaxSec uint32 `gorm:"default:30;comment:主播排队招呼间隔上限秒" json:"intervalMaxSec"`
	Content     string `gorm:"size:1024;default:'';comment:招呼正文(历史单语,回落英)" json:"content"`
	ContentZhCN string `gorm:"size:1024;default:'';comment:招呼正文(简体中文)" json:"contentZhCN"`
	ContentZhTW string `gorm:"size:1024;default:'';comment:招呼正文(繁体中文)" json:"contentZhTW"`
	ContentEn   string `gorm:"size:1024;default:'';comment:招呼正文(英文,兼容单条)" json:"contentEn"`
	ContentEs   string `gorm:"size:1024;default:'';comment:招呼正文(西班牙语,兼容)" json:"contentEs"`
	ContentPt   string `gorm:"size:1024;default:'';comment:招呼正文(葡萄牙语,兼容)" json:"contentPt"`
	ContentHi   string `gorm:"size:1024;default:'';comment:招呼正文(印地语,兼容)" json:"contentHi"`
	ContentId   string `gorm:"size:1024;default:'';comment:招呼正文(印尼语,兼容)" json:"contentId"`
	ContentsEn  string `gorm:"type:text;default:'';comment:英文招呼列表JSON" json:"contentsEn"`
	ContentsEs  string `gorm:"type:text;default:'';comment:西班牙语招呼列表JSON" json:"contentsEs"`
	ContentsPt  string `gorm:"type:text;default:'';comment:葡萄牙语招呼列表JSON" json:"contentsPt"`
	ContentsHi  string `gorm:"type:text;default:'';comment:印地语招呼列表JSON" json:"contentsHi"`
	ContentsId  string `gorm:"type:text;default:'';comment:印尼语招呼列表JSON" json:"contentsId"`
}

// EffectiveAnchorCount 每次最多打招呼的主播人数(含默认与上限).
func (c *AudienceOnlineAnchorGreetCfg) EffectiveAnchorCount() uint8 {
	if c == nil || c.AnchorCount == 0 {
		return DefaultAudienceOnlineAnchorGreetAnchorCount
	}
	if c.AnchorCount > AudienceOnlineAnchorGreetAnchorCountMax {
		return AudienceOnlineAnchorGreetAnchorCountMax
	}
	return c.AnchorCount
}

// EffectiveGreetIntervalSec 排队招呼随机间隔上下限(秒); max 不小于 min.
func (c *AudienceOnlineAnchorGreetCfg) EffectiveGreetIntervalSec() (minSec, maxSec uint32) {
	minSec = DefaultAudienceOnlineAnchorGreetIntervalMinSec
	maxSec = DefaultAudienceOnlineAnchorGreetIntervalMaxSec
	if c == nil {
		return minSec, maxSec
	}
	if c.IntervalMinSec > 0 || c.IntervalMaxSec > 0 {
		minSec = c.IntervalMinSec
		maxSec = c.IntervalMaxSec
		if minSec == 0 && maxSec > 0 {
			minSec = maxSec
		}
		if maxSec == 0 && minSec > 0 {
			maxSec = minSec
		}
	}
	if minSec > AudienceOnlineAnchorGreetIntervalSecMax {
		minSec = AudienceOnlineAnchorGreetIntervalSecMax
	}
	if maxSec > AudienceOnlineAnchorGreetIntervalSecMax {
		maxSec = AudienceOnlineAnchorGreetIntervalSecMax
	}
	if maxSec < minSec {
		maxSec = minSec
	}
	return minSec, maxSec
}

func initAudienceOnlineAnchorGreetCfg() {
	migrate.AutoMigrate(&AudienceOnlineAnchorGreetCfg{})
}