package entity

// DeviceInfo App端上报的设备信息
type DeviceInfo struct {
	DeviceType  string `json:"deviceType" v:"required|in:ios,android,iOS,Android#设备类型不能为空|设备类型无效" dc:"设备类型(ios/android)"`
	DeviceModel string `json:"deviceModel" dc:"设备型号"`
	OsVersion   string `json:"osVersion" v:"required|length:1,64#系统版本号不能为空" dc:"系统版本号"`
	AppVersion  string `json:"appVersion" dc:"App/APK版本号"`
	PackageName string `json:"packageName" dc:"App包名"`
	DeviceId    string `json:"deviceId" dc:"设备唯一标识(可选)"`
}
