package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbUserPayoutInfo db.TbName = "user_payout_infos"
)

const UserPayoutAccountTypeBank = "BANK_ACCOUNT"

// UserPayoutInfo 普通用户收款/代付信息(主键 ID=用户 ID)
type UserPayoutInfo struct {
	migrate.OneModel
	CountryCode  string `gorm:"size:8;not null;default:'';comment:国家/地区ISO简码" json:"countryCode"`
	Currency     string `gorm:"size:16;not null;default:'';comment:代付币种" json:"currency"`
	AccountType  string `gorm:"size:32;not null;default:'BANK_ACCOUNT';comment:HaiPay accountType" json:"accountType"`
	PayeeName    string `gorm:"size:128;default:'';comment:收款人姓名" json:"payeeName"`
	Phone        string `gorm:"size:32;default:'';comment:收款人手机" json:"phone"`
	Email        string `gorm:"size:128;default:'';comment:收款人邮箱" json:"email"`
	BankName     string `gorm:"size:128;default:'';comment:支付方式说明" json:"bankName"`
	AccountNo    string `gorm:"size:128;default:'';comment:收款账号" json:"accountNo"`
	BankCode     string `gorm:"size:64;default:'';comment:银行/支付编码" json:"bankCode"`
	IdentifyType string `gorm:"size:64;default:'';comment:PIX/Papara类型或ACH路由" json:"identifyType"`
	Address1     string `gorm:"size:255;default:'';comment:收款地址-街道" json:"address1"`
	Address2     string `gorm:"size:128;default:'';comment:收款地址-城市" json:"address2"`
	Address3     string `gorm:"size:128;default:'';comment:收款地址-省州" json:"address3"`
	PostalCode   string `gorm:"size:32;default:'';comment:收款邮编" json:"postalCode"`
	Remark       string `gorm:"size:255;default:'';comment:备注" json:"remark"`
}

func NewUserPayoutInfo(userId uint64) *UserPayoutInfo {
	now := time.Now()
	return &UserPayoutInfo{
		OneModel: migrate.OneModel{
			ID:        userId,
			CreatedAt: now,
			UpdatedAt: now,
		},
		AccountType: UserPayoutAccountTypeBank,
	}
}

func initUserPayoutInfo() {
	migrate.AutoMigrate(&UserPayoutInfo{})
}
