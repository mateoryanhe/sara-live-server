package accountdto

import "github.com/gogf/gf/v2/frame/g"

type GetAnchorTransferInfoReq struct {
	g.Meta   `path:"/getAnchorTransferInfo" method:"post" summary:"获取平台主播转账信息" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"平台主播用户ID"`
}

type AnchorTransferWalletOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AnchorTransferMethodOption struct {
	AccountType          string   `json:"accountType"`
	BankCode             string   `json:"bankCode"`
	Limit                string   `json:"limit"`
	Description          string   `json:"description"`
	IdentifyTypeRequired bool     `json:"identifyTypeRequired"`
	IdentifyTypeOptions  []string `json:"identifyTypeOptions"`
	CountryRequired      bool     `json:"countryRequired"`
	AddressRequired      bool     `json:"addressRequired"`
}

type AnchorTransferCountryOption struct {
	CountryCode  string                       `json:"countryCode"`
	NameEn       string                       `json:"nameEn"`
	NameZh       string                       `json:"nameZh"`
	Currency     string                       `json:"currency"`
	Region       string                       `json:"region"`
	AccountTypes []string                     `json:"accountTypes"`
	Wallets      []AnchorTransferWalletOption `json:"wallets"`
	Methods      []AnchorTransferMethodOption `json:"methods"`
	Icon         string                       `json:"icon"`
}

type AnchorTransferInfoItem struct {
	AnchorId     string `json:"anchorId"`
	CountryCode  string `json:"countryCode"`
	Currency     string `json:"currency"`
	AccountType  string `json:"accountType"`
	PayeeName    string `json:"payeeName"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	BankName     string `json:"bankName"`
	AccountNo    string `json:"accountNo"`
	BankCode     string `json:"bankCode"`
	IdentifyType string `json:"identifyType"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Address3     string `json:"address3"`
	PostalCode   string `json:"postalCode"`
	Remark       string `json:"remark"`
	UpdatedAt    string `json:"updatedAt"`
}

type GetAnchorTransferInfoRes struct {
	Info      *AnchorTransferInfoItem        `json:"info"`
	Countries []*AnchorTransferCountryOption `json:"countries"`
}

type SaveAnchorTransferInfoReq struct {
	g.Meta       `path:"/saveAnchorTransferInfo" method:"post" summary:"保存平台主播转账信息" tags:"账号"`
	AnchorId     uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"平台主播用户ID"`
	CountryCode  string `json:"countryCode" v:"required|length:2,8#国家不能为空|国家简码无效"`
	Currency     string `json:"currency" v:"max-length:16#币种最长16字符"`
	AccountType  string `json:"accountType" v:"required|max-length:32#accountType不能为空|accountType最长32"`
	PayeeName    string `json:"payeeName" v:"max-length:128#收款人姓名最长128字符"`
	Phone        string `json:"phone" v:"max-length:32#手机号最长32字符"`
	Email        string `json:"email" v:"max-length:128#邮箱最长128字符"`
	BankName     string `json:"bankName" v:"max-length:128#支付方式说明最长128字符"`
	AccountNo    string `json:"accountNo" v:"max-length:128#收款账号最长128字符"`
	BankCode     string `json:"bankCode" v:"required|max-length:64#支付编码不能为空|支付编码最长64字符"`
	IdentifyType string `json:"identifyType" v:"max-length:64#identifyType最长64字符"`
	Address1     string `json:"address1" v:"max-length:255#详细街道最长255字符"`
	Address2     string `json:"address2" v:"max-length:128#城市最长128字符"`
	Address3     string `json:"address3" v:"max-length:128#省州最长128字符"`
	PostalCode   string `json:"postalCode" v:"max-length:32#邮编最长32字符"`
	Remark       string `json:"remark" v:"max-length:255#备注最长255字符"`
}

type SaveAnchorTransferInfoRes struct {
	Success bool `json:"success"`
}
