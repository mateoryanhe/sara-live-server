package entity

import (
	"strconv"
	"strings"
)

// DeviceInfo App端上报的设备信息
type DeviceInfo struct {
	DeviceType  string `json:"deviceType" v:"required|in:ios,android,iOS,Android#设备类型不能为空|设备类型无效" dc:"设备类型(ios/android)"`
	DeviceModel string `json:"deviceModel" dc:"设备型号"`
	CpuModel    string `json:"cpuModel" dc:"CPU型号/硬件标识"`
	OsVersion   string `json:"osVersion" v:"required|length:1,64#系统版本号不能为空" dc:"系统版本号"`
	AppVersion  string `json:"appVersion" dc:"App/APK版本号"`
	PackageName string `json:"packageName" dc:"App包名"`
	DeviceId    string `json:"deviceId" dc:"设备唯一标识"`
	// Meta Conversions API 归因(可选,App/Web 从 Meta SDK 或 Cookie 采集)
	MetaFbp string `json:"metaFbp" dc:"Meta _fbp"`
	MetaFbc string `json:"metaFbc" dc:"Meta _fbc"`
	// Ch 推广渠道编码(可选); JSON 无此字段、空字符串或仅空白时不写 user_ext、不计推广日统计
	Ch string `json:"ch,omitempty" dc:"推广渠道编码(非空且为正整数才统计)"`
}

// PromotionChannelCode 解析 deviceInfo.ch; 无字段/空/非正整数时 ok=false,推广统计全部跳过
func (d *DeviceInfo) PromotionChannelCode() (code int, ok bool) {
	if d == nil {
		return 0, false
	}
	ch := strings.TrimSpace(d.Ch)
	if ch == "" {
		return 0, false
	}
	n, err := strconv.Atoi(ch)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
