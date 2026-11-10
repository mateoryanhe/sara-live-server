package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbSimulatorBlacklistLog db.TbName = "simulator_blacklist_logs"

const (
	SimulatorBlacklistLogIP          db.TbCol = "ip"
	SimulatorBlacklistLogDeviceId    db.TbCol = "device_id"
	SimulatorBlacklistLogCpuModel    db.TbCol = "cpu_model"
	SimulatorBlacklistLogDeviceType  db.TbCol = "device_type"
	SimulatorBlacklistLogDeviceModel db.TbCol = "device_model"
	SimulatorBlacklistLogOsVersion   db.TbCol = "os_version"
	SimulatorBlacklistLogAppVersion  db.TbCol = "app_version"
	SimulatorBlacklistLogPackageName db.TbCol = "package_name"
	SimulatorBlacklistLogLoginType   db.TbCol = "login_type"
	SimulatorBlacklistLogEmail       db.TbCol = "email"
	SimulatorBlacklistLogAccountId   db.TbCol = "account_id"
	SimulatorBlacklistLogRegisterAt  db.TbCol = "register_at"
)

// SimulatorBlacklistLog 模拟器登录拦截日志(append-only, syndb 缓冲入库)
type SimulatorBlacklistLog struct {
	migrate.OneModel
	IP          string     `gorm:"size:64;index;default:'';comment:客户端IP" json:"ip"`
	DeviceId    string     `gorm:"size:128;index;default:'';comment:设备码" json:"deviceId"`
	CpuModel    string     `gorm:"size:128;default:'';comment:CPU型号/硬件标识" json:"cpuModel"`
	DeviceType  string     `gorm:"size:32;default:'';comment:设备类型" json:"deviceType"`
	DeviceModel string     `gorm:"size:128;default:'';comment:设备型号" json:"deviceModel"`
	OsVersion   string     `gorm:"size:64;default:'';comment:系统版本" json:"osVersion"`
	AppVersion  string     `gorm:"size:64;default:'';comment:App版本" json:"appVersion"`
	PackageName string     `gorm:"size:128;default:'';comment:包名" json:"packageName"`
	LoginType   string     `gorm:"size:32;index;default:'';comment:登录方式(device/email等)" json:"loginType"`
	Email       string     `gorm:"size:128;index;default:'';comment:邮箱(邮箱登录时)" json:"email"`
	AccountId   uint64     `gorm:"index;default:0;comment:已有账号ID(无则0)" json:"accountId"`
	RegisterAt  *time.Time `gorm:"index;comment:账号注册时间(已有账号时)" json:"registerAt"`
}

func NewSimulatorBlacklistLog(
	ip, deviceId, cpuModel, deviceType, deviceModel, osVersion, appVersion, packageName, loginType, email string,
	accountId uint64,
	registerAt *time.Time,
) *SimulatorBlacklistLog {
	now := time.Now()
	row := &SimulatorBlacklistLog{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetIP(ip)
	row.SetDeviceId(deviceId)
	row.SetCpuModel(cpuModel)
	row.SetDeviceType(deviceType)
	row.SetDeviceModel(deviceModel)
	row.SetOsVersion(osVersion)
	row.SetAppVersion(appVersion)
	row.SetPackageName(packageName)
	row.SetLoginType(loginType)
	row.SetEmail(email)
	row.SetAccountId(accountId)
	row.SetRegisterAt(registerAt)
	return row
}

func (r *SimulatorBlacklistLog) SetIP(v string) {
	r.IP = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogIP, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetDeviceId(v string) {
	r.DeviceId = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetCpuModel(v string) {
	r.CpuModel = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogCpuModel, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetDeviceType(v string) {
	r.DeviceType = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetDeviceModel(v string) {
	r.DeviceModel = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceModel, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetOsVersion(v string) {
	r.OsVersion = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogOsVersion, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetAppVersion(v string) {
	r.AppVersion = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogAppVersion, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetPackageName(v string) {
	r.PackageName = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogPackageName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetLoginType(v string) {
	r.LoginType = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogLoginType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetEmail(v string) {
	r.Email = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogEmail, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetAccountId(v uint64) {
	r.AccountId = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogAccountId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetRegisterAt(v *time.Time) {
	r.RegisterAt = v
	syndb.AddData(TbSimulatorBlacklistLog, SimulatorBlacklistLogRegisterAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbSimulatorBlacklistLog, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *SimulatorBlacklistLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbSimulatorBlacklistLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initSimulatorBlacklistLog() {
	syndb.RegQuick(TbSimulatorBlacklistLog, db.CreatedAtName)
	syndb.RegQuick(TbSimulatorBlacklistLog, db.UpdatedAtName)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogIP)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceId)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogCpuModel)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceType)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogDeviceModel)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogOsVersion)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogAppVersion)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogPackageName)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogLoginType)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogEmail)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogAccountId)
	syndb.RegQuick(TbSimulatorBlacklistLog, SimulatorBlacklistLogRegisterAt)
	migrate.AutoMigrate(&SimulatorBlacklistLog{})
}
