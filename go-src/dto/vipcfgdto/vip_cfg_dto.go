package vipcfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// VipCfgListReq CMS分页查询VIP配置
type VipCfgListReq struct {
	g.Meta `path:"/vipCfgList" method:"post" summary:"获取VIP配置列表" tags:"VIP配置"`
	httpserver.CMSQueryReq
	LevelName            string `json:"levelName" dc:"等级名称(模糊匹配)"`
	WithdrawSwitchFilter int    `json:"withdrawSwitchFilter" dc:"提现开关过滤(0=全部,1=只看关闭,2=只看开启)"`
}

// VipCfgListRes 列表项
type VipCfgListRes struct {
	ID                    string  `json:"id"`
	Level                 uint32  `json:"level"`
	LevelName             string  `json:"levelName"`
	WithdrawSwitch        uint8   `json:"withdrawSwitch" dc:"提现开关(0关,1开,仅App端使用)"`
	AnimationSwitch       uint8   `json:"animationSwitch" dc:"进场特效开关(0关,1开,仅App端使用)"`
	CommentEffectSwitch   uint8   `json:"commentEffectSwitch" dc:"公屏评论特效开关(0关,1开,仅App端使用)"`
	UpgradeRechargeLimit  float64 `json:"upgradeRechargeLimit,string"`
	MinWithdrawAmount     float64 `json:"minWithdrawAmount,string"`
	MaxWithdrawAmount     float64 `json:"maxWithdrawAmount,string"`
	Fee                   float64 `json:"fee,string"`
	Animation             string  `json:"animation" dc:"动画完整URL(列表展示)"`
	AnimationName         string  `json:"animationName" dc:"动画资源文件名(编辑保存用)"`
	AnimationIcon         string  `json:"animationIcon" dc:"进场特效图标完整URL(列表展示)"`
	AnimationIconName     string  `json:"animationIconName" dc:"进场特效图标资源文件名(编辑保存用)"`
	AnimationDescEn       string  `json:"animationDescEn" dc:"进场特效说明(英文)"`
	AnimationDescEs       string  `json:"animationDescEs" dc:"进场特效说明(西班牙语)"`
	AnimationDescPt       string  `json:"animationDescPt" dc:"进场特效说明(葡萄牙语)"`
	AnimationDescHi       string  `json:"animationDescHi" dc:"进场特效说明(印地语)"`
	CommentEffect         string  `json:"commentEffect" dc:"公屏评论特效动画完整URL(列表展示)"`
	CommentEffectName     string  `json:"commentEffectName" dc:"公屏评论特效动画资源文件名(编辑保存用)"`
	CommentEffectIcon     string  `json:"commentEffectIcon" dc:"公屏评论特效图标完整URL(列表展示)"`
	CommentEffectIconName string  `json:"commentEffectIconName" dc:"公屏评论特效图标资源文件名(编辑保存用)"`
	CommentEffectDescEn   string  `json:"commentEffectDescEn" dc:"公屏评论特效说明(英文)"`
	CommentEffectDescEs   string  `json:"commentEffectDescEs" dc:"公屏评论特效说明(西班牙语)"`
	CommentEffectDescPt   string  `json:"commentEffectDescPt" dc:"公屏评论特效说明(葡萄牙语)"`
	CommentEffectDescHi   string  `json:"commentEffectDescHi" dc:"公屏评论特效说明(印地语)"`
	WithdrawIcon          string  `json:"withdrawIcon" dc:"提现图标完整URL(列表展示)"`
	WithdrawIconName      string  `json:"withdrawIconName" dc:"提现图标资源文件名(编辑保存用)"`
	WithdrawNoticeEn      string  `json:"withdrawNoticeEn" dc:"提现须知(英文)"`
	WithdrawNoticeEs      string  `json:"withdrawNoticeEs" dc:"提现须知(西班牙语)"`
	WithdrawNoticePt      string  `json:"withdrawNoticePt" dc:"提现须知(葡萄牙语)"`
	WithdrawNoticeHi      string  `json:"withdrawNoticeHi" dc:"提现须知(印地语)"`
	CreatedAt             string  `json:"createdAt"`
	UpdatedAt             string  `json:"updatedAt"`
}

// CreateVipCfgReq 创建VIP配置
type CreateVipCfgReq struct {
	g.Meta               `path:"/createVipCfg" method:"post" summary:"创建VIP配置" tags:"VIP配置"`
	Level                uint32  `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
	LevelName            string  `json:"levelName" v:"required|length:1,64#等级名称不能为空|等级名称长度需在1到64之间" dc:"等级名称"`
	WithdrawSwitch       uint8   `json:"withdrawSwitch" v:"in:0,1#提现开关无效" dc:"提现开关(0关,1开,仅App端使用)"`
	AnimationSwitch      uint8   `json:"animationSwitch" v:"in:0,1#进场特效开关无效" dc:"进场特效开关(0关,1开,仅App端使用)"`
	CommentEffectSwitch  uint8   `json:"commentEffectSwitch" v:"in:0,1#公屏评论特效开关无效" dc:"公屏评论特效开关(0关,1开,仅App端使用)"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit" dc:"升级充值上限(USD,保留4位小数)"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount" dc:"最低提现金额(USD,保留4位小数)"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount" dc:"最高提现金额(USD,保留4位小数)"`
	Fee                  float64 `json:"fee" dc:"手续费(保留4位小数)"`
	Animation            string  `json:"animation" dc:"进场特效动画资源文件名(mp4)"`
	AnimationIcon        string  `json:"animationIcon" dc:"进场特效图标资源文件名"`
	AnimationDescEn      string  `json:"animationDescEn" dc:"进场特效说明(英文)"`
	AnimationDescEs      string  `json:"animationDescEs" dc:"进场特效说明(西班牙语)"`
	AnimationDescPt      string  `json:"animationDescPt" dc:"进场特效说明(葡萄牙语)"`
	AnimationDescHi      string  `json:"animationDescHi" dc:"进场特效说明(印地语)"`
	CommentEffect        string  `json:"commentEffect" dc:"公屏评论特效动画资源文件名(mp4)"`
	CommentEffectIcon    string  `json:"commentEffectIcon" dc:"公屏评论特效图标资源文件名"`
	CommentEffectDescEn  string  `json:"commentEffectDescEn" dc:"公屏评论特效说明(英文)"`
	CommentEffectDescEs  string  `json:"commentEffectDescEs" dc:"公屏评论特效说明(西班牙语)"`
	CommentEffectDescPt  string  `json:"commentEffectDescPt" dc:"公屏评论特效说明(葡萄牙语)"`
	CommentEffectDescHi  string  `json:"commentEffectDescHi" dc:"公屏评论特效说明(印地语)"`
	WithdrawIcon         string  `json:"withdrawIcon" dc:"提现图标资源文件名"`
	WithdrawNoticeEn     string  `json:"withdrawNoticeEn" dc:"提现须知(英文)"`
	WithdrawNoticeEs     string  `json:"withdrawNoticeEs" dc:"提现须知(西班牙语)"`
	WithdrawNoticePt     string  `json:"withdrawNoticePt" dc:"提现须知(葡萄牙语)"`
	WithdrawNoticeHi     string  `json:"withdrawNoticeHi" dc:"提现须知(印地语)"`
}

type CreateVipCfgRes struct {
	ID string `json:"id"`
}

// UpdateVipCfgReq 更新VIP配置
type UpdateVipCfgReq struct {
	g.Meta               `path:"/updateVipCfg" method:"post" summary:"修改VIP配置" tags:"VIP配置"`
	ID                   uint64  `json:"id" v:"required#ID不能为空" dc:"配置ID"`
	Level                uint32  `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
	LevelName            string  `json:"levelName" v:"required|length:1,64#等级名称不能为空|等级名称长度需在1到64之间" dc:"等级名称"`
	WithdrawSwitch       uint8   `json:"withdrawSwitch" v:"in:0,1#提现开关无效" dc:"提现开关(0关,1开,仅App端使用)"`
	AnimationSwitch      uint8   `json:"animationSwitch" v:"in:0,1#进场特效开关无效" dc:"进场特效开关(0关,1开,仅App端使用)"`
	CommentEffectSwitch  uint8   `json:"commentEffectSwitch" v:"in:0,1#公屏评论特效开关无效" dc:"公屏评论特效开关(0关,1开,仅App端使用)"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit" dc:"升级充值上限(USD,保留4位小数)"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount" dc:"最低提现金额(USD,保留4位小数)"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount" dc:"最高提现金额(USD,保留4位小数)"`
	Fee                  float64 `json:"fee" dc:"手续费(保留4位小数)"`
	Animation            string  `json:"animation" dc:"进场特效动画资源文件名(mp4)"`
	AnimationIcon        string  `json:"animationIcon" dc:"进场特效图标资源文件名"`
	AnimationDescEn      string  `json:"animationDescEn" dc:"进场特效说明(英文)"`
	AnimationDescEs      string  `json:"animationDescEs" dc:"进场特效说明(西班牙语)"`
	AnimationDescPt      string  `json:"animationDescPt" dc:"进场特效说明(葡萄牙语)"`
	AnimationDescHi      string  `json:"animationDescHi" dc:"进场特效说明(印地语)"`
	CommentEffect        string  `json:"commentEffect" dc:"公屏评论特效动画资源文件名(mp4)"`
	CommentEffectIcon    string  `json:"commentEffectIcon" dc:"公屏评论特效图标资源文件名"`
	CommentEffectDescEn  string  `json:"commentEffectDescEn" dc:"公屏评论特效说明(英文)"`
	CommentEffectDescEs  string  `json:"commentEffectDescEs" dc:"公屏评论特效说明(西班牙语)"`
	CommentEffectDescPt  string  `json:"commentEffectDescPt" dc:"公屏评论特效说明(葡萄牙语)"`
	CommentEffectDescHi  string  `json:"commentEffectDescHi" dc:"公屏评论特效说明(印地语)"`
	WithdrawIcon         string  `json:"withdrawIcon" dc:"提现图标资源文件名"`
	WithdrawNoticeEn     string  `json:"withdrawNoticeEn" dc:"提现须知(英文)"`
	WithdrawNoticeEs     string  `json:"withdrawNoticeEs" dc:"提现须知(西班牙语)"`
	WithdrawNoticePt     string  `json:"withdrawNoticePt" dc:"提现须知(葡萄牙语)"`
	WithdrawNoticeHi     string  `json:"withdrawNoticeHi" dc:"提现须知(印地语)"`
}

type UpdateVipCfgRes struct {
	Success bool `json:"success"`
}

// DeleteVipCfgReq 删除VIP配置
type DeleteVipCfgReq struct {
	g.Meta `path:"/deleteVipCfg" method:"post" summary:"删除VIP配置" tags:"VIP配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"配置ID"`
}

// DeleteVipCfgRes 删除VIP配置响应
type DeleteVipCfgRes struct {
	Success bool `json:"success"`
}

// ===== App =====

const (
	AppVipPrivilegeTypeWithdraw      uint8 = 1 // 提现
	AppVipPrivilegeTypeEntryEffect   uint8 = 2 // 进场特效
	AppVipPrivilegeTypeCommentEffect uint8 = 3 // 公屏评论特效
)

// AppVipPrivilegeItem App端VIP特权项(按开关动态组装)
type AppVipPrivilegeItem struct {
	PrivilegeType     uint8   `json:"privilegeType" dc:"特权类型(1=提现,2=进场特效,3=公屏评论特效)"`
	Icon              string  `json:"icon,omitempty" dc:"图标完整URL"`
	Animation         string  `json:"animation,omitempty" dc:"特效动画完整URL"`
	DescEn            string  `json:"descEn,omitempty" dc:"特效说明(英文)"`
	DescEs            string  `json:"descEs,omitempty" dc:"特效说明(西班牙语)"`
	DescPt            string  `json:"descPt,omitempty" dc:"特效说明(葡萄牙语)"`
	DescHi            string  `json:"descHi,omitempty" dc:"特效说明(印地语)"`
	MinWithdrawAmount float64 `json:"minWithdrawAmount,string,omitempty" dc:"最低提现金额(USD,保留4位小数)"`
	MaxWithdrawAmount float64 `json:"maxWithdrawAmount,string,omitempty" dc:"最高提现金额(USD,保留4位小数)"`
	Fee               float64 `json:"fee,string,omitempty" dc:"提现手续费(保留4位小数)"`
	NoticeEn          string  `json:"noticeEn,omitempty" dc:"提现须知(英文)"`
	NoticeEs          string  `json:"noticeEs,omitempty" dc:"提现须知(西班牙语)"`
	NoticePt          string  `json:"noticePt,omitempty" dc:"提现须知(葡萄牙语)"`
	NoticeHi          string  `json:"noticeHi,omitempty" dc:"提现须知(印地语)"`
}

// AppVipCfgItem App端VIP等级权益配置
type AppVipCfgItem struct {
	Level                uint32                 `json:"level" dc:"VIP等级"`
	LevelName            string                 `json:"levelName" dc:"等级名称"`
	UpgradeRechargeLimit float64                `json:"upgradeRechargeLimit,string" dc:"升级所需累计充值上限(USD,保留4位小数)"`
	PrivilegeList        []*AppVipPrivilegeItem `json:"privilegeList" dc:"已开启的特权列表(按提现/进场/公屏评论顺序)"`
}

// AppVipCfgByLevelReq App端按等级查询VIP配置
type AppVipCfgByLevelReq struct {
	g.Meta `path:"/getVipCfgByLevel" method:"post" summary:"按等级查询VIP配置" tags:"VIP配置"`
	Level  uint32 `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
}

// AppVipCfgByLevelRes App端按等级查询VIP配置响应
type AppVipCfgByLevelRes struct {
	Item *AppVipCfgItem `json:"item" dc:"VIP等级权益配置,不存在时为null"`
}

// AppVipCfgListReq App端查询全部VIP配置
type AppVipCfgListReq struct {
	g.Meta `path:"/vipCfgListForApp" method:"post" summary:"查询全部VIP配置" tags:"VIP配置"`
}

// AppVipCfgListRes App端查询全部VIP配置响应
type AppVipCfgListRes struct {
	List []*AppVipCfgItem `json:"list" dc:"VIP等级权益配置列表(按等级升序)"`
}
