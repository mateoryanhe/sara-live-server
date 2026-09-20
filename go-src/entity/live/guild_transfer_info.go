package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGuildTransferInfo db.TbName = "live_guild_transfer_infos"
)

const (
	// GuildTransferAccountTypeBank 默认银行卡代付(HaiPay accountType)
	GuildTransferAccountTypeBank = "BANK_ACCOUNT"
)

// LiveGuildTransferInfo 工会收款/转账信息(主键ID=工会ID,直写数据库)
// 对齐 HaiPay 代付: 选国家(CountryCode)→推导 Currency; name/phone/email/accountType/bankCode/accountNo
type LiveGuildTransferInfo struct {
	migrate.OneModel
	CountryCode  string `gorm:"size:8;not null;default:'';comment:国家/地区ISO简码(如ID),与country包一致" json:"countryCode"`
	Currency     string `gorm:"size:16;not null;default:'';comment:代付币种(由CountryCode推导,如IDR)" json:"currency"`
	AccountType  string `gorm:"size:32;not null;default:'BANK_ACCOUNT';comment:HaiPay accountType" json:"accountType"`
	PayeeName    string `gorm:"size:128;default:'';comment:收款人姓名(HaiPay name)" json:"payeeName"`
	Phone        string `gorm:"size:32;default:'';comment:收款人手机(HaiPay phone)" json:"phone"`
	Email        string `gorm:"size:128;default:'';comment:收款人邮箱(HaiPay email)" json:"email"`
	BankName     string `gorm:"size:128;default:'';comment:HaiPay支付方式说明(由支付编码自动填充)" json:"bankName"`
	AccountNo    string `gorm:"size:128;default:'';comment:收款账号(HaiPay accountNo)" json:"accountNo"`
	BankCode     string `gorm:"size:64;default:'';comment:银行/支付编码(HaiPay bankCode)" json:"bankCode"`
	IdentifyType string `gorm:"size:64;default:'';comment:PIX/Papara账号类型或ACH银行路由编码(HaiPay identifyType)" json:"identifyType"`
	Address1     string `gorm:"size:255;default:'';comment:收款地址-详细街道(HaiPay address1)" json:"address1"`
	Address2     string `gorm:"size:128;default:'';comment:收款地址-城市(HaiPay address2)" json:"address2"`
	Address3     string `gorm:"size:128;default:'';comment:收款地址-省州(HaiPay address3)" json:"address3"`
	PostalCode   string `gorm:"size:32;default:'';comment:收款邮编(HaiPay postalCode)" json:"postalCode"`
	Remark       string `gorm:"size:255;default:'';comment:备注" json:"remark"`
}

func NewLiveGuildTransferInfo(guildId uint64) *LiveGuildTransferInfo {
	now := time.Now()
	return &LiveGuildTransferInfo{
		OneModel: migrate.OneModel{
			ID:        guildId,
			CreatedAt: now,
			UpdatedAt: now,
		},
		AccountType: GuildTransferAccountTypeBank,
	}
}

func initLiveGuildTransferInfo() {
	migrate.AutoMigrate(&LiveGuildTransferInfo{})
}
