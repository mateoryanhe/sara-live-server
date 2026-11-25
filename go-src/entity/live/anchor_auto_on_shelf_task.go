package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorAutoOnShelfTask db.TbName = "anchor_auto_on_shelf_tasks"
)

const (
	AnchorAutoOnShelfTaskStatusPending   uint8 = 0 // 待自动上架
	AnchorAutoOnShelfTaskStatusDone      uint8 = 1 // 已完成
	AnchorAutoOnShelfTaskStatusCancelled uint8 = 2 // 已取消(被新任务覆盖或人工上架)
)

const (
	AnchorAutoOnShelfTaskAnchorId      db.TbCol = "anchor_id"
	AnchorAutoOnShelfTaskOffShelfMinutes db.TbCol = "off_shelf_minutes"
	AnchorAutoOnShelfTaskOffShelfAt    db.TbCol = "off_shelf_at"
	AnchorAutoOnShelfTaskRestoreAt     db.TbCol = "restore_at"
	AnchorAutoOnShelfTaskStatus        db.TbCol = "status"
	AnchorAutoOnShelfTaskOperatorId    db.TbCol = "operator_id"
	AnchorAutoOnShelfTaskCompletedAt   db.TbCol = "completed_at"
)

// AnchorAutoOnShelfTask CMS 下架主播后定时自动恢复上架
type AnchorAutoOnShelfTask struct {
	migrate.OneModel
	AnchorId      uint64     `gorm:"index:idx_anchor_auto_on_shelf_anchor_status,priority:1;default:0;comment:主播用户ID" json:"anchorId"`
	OffShelfMinutes uint32   `gorm:"default:0;comment:计划下架时长(分钟)" json:"offShelfMinutes"`
	OffShelfAt    time.Time  `gorm:"comment:下架时间" json:"offShelfAt"`
	RestoreAt     time.Time  `gorm:"index;comment:计划自动上架时间" json:"restoreAt"`
	Status        uint8      `gorm:"index:idx_anchor_auto_on_shelf_anchor_status,priority:2;default:0;comment:状态(0待执行,1已完成,2已取消)" json:"status"`
	OperatorId    uint64     `gorm:"default:0;comment:CMS操作人ID" json:"operatorId"`
	CompletedAt   *time.Time `gorm:"comment:完成或取消时间" json:"completedAt,omitempty"`
}

func initAnchorAutoOnShelfTask() {
	migrate.AutoMigrate(&AnchorAutoOnShelfTask{})
}