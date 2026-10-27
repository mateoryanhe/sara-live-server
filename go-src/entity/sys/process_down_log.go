package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbProcessDownLog db.TbName = "process_down_logs"

// ProcessDownLog 探活 gw 判定直播 Go 进程挂掉后的记录(由 live-hot-restart-gw 写入).
type ProcessDownLog struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Host         string     `gorm:"size:128;index;comment:主机名" json:"host"`
	AppName      string     `gorm:"size:64;comment:进程名" json:"appName"`
	MissingSince time.Time  `gorm:"index;comment:连续缺失起点" json:"missingSince"`
	DetectedAt   time.Time  `gorm:"index;comment:达到阈值判定挂掉时间" json:"detectedAt"`
	ColdStartOK  bool       `gorm:"comment:冷启是否成功" json:"coldStartOK"`
	NewPID       int        `gorm:"comment:冷启后PID,失败为0" json:"newPid"`
	Detail       string     `gorm:"size:1024;comment:结果说明" json:"detail"`
	CreatedAt    time.Time  `gorm:"index;comment:入库时间" json:"createdAt"`
}

func initProcessDownLog() {
	migrate.AutoMigrate(&ProcessDownLog{})
}
