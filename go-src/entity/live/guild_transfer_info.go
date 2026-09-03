package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGuildTransferInfo db.TbName = "live_guild_transfer_infos"
)

// LiveGuildTransferInfo 工会收款/转账信息(主键ID=工会ID,直写数据库)
type LiveGuildTransferInfo struct {
	migrate.OneModel
	Currency  string `gorm:"size:16;not null;default:'';comment:收款币种(如IDR)" json:"currency"`
	PayeeName string `gorm:"size:128;default:'';comment:收款人姓名" json:"payeeName"`
	BankName  string `gorm:"size:128;default:'';comment:银行名称" json:"bankName"`
	AccountNo string `gorm:"size:128;default:'';comment:收款账号" json:"accountNo"`
	BankCode  string `gorm:"size:64;default:'';comment:银行代码" json:"bankCode"`
	Remark    string `gorm:"size:255;default:'';comment:备注" json:"remark"`
}

func NewLiveGuildTransferInfo(guildId uint64) *LiveGuildTransferInfo {
	now := time.Now()
	return &LiveGuildTransferInfo{
		OneModel: migrate.OneModel{
			ID:        guildId,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func initLiveGuildTransferInfo() {
	migrate.AutoMigrate(&LiveGuildTransferInfo{})
}
