package country

import "strings"

// HaiPayCollectionPaymentMethod 描述官方代收方式总表中的可用支付方式。
type HaiPayCollectionPaymentMethod struct {
	CurrencyCode string
	PayType      string
	InBankCode   string
	MinAmount    string
	MaxAmount    string
	Description  string
	Available    bool
}

// haiPayGlobalCollectionPaymentMethods 是官网“全球”分组，适用于支持 USD 的国家/地区。
var haiPayGlobalCollectionPaymentMethods = []HaiPayCollectionPaymentMethod{
	{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "CREDIT_CARD", MinAmount: "0.99", MaxAmount: "1000", Description: "VISA / MasterCard / JCB", Available: true},
	{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "GOOGLE_PAY", MinAmount: "0.99", MaxAmount: "1000", Description: "Google Pay", Available: true},
	{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "0.99", MaxAmount: "1000", Description: "Apple Pay", Available: true},
}

// haiPayCollectionPaymentMethodsByCountry 只来自 HaiPay“代收支付方式列表”总表。
// 目录只保留状态为“可用”的条目；维护项不会进入 CMS 可选列表。
var haiPayCollectionPaymentMethodsByCountry = map[string][]HaiPayCollectionPaymentMethod{
	"AT": {
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "EPS_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "奥地利 EPS", Available: true},
		{CurrencyCode: "EUR", PayType: "BANK_TRANSFER", InBankCode: "EPS_EUR", MinAmount: "0.99", MaxAmount: "1000", Description: "奥地利 EPS", Available: true},
	},
	"BE": {
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "BANCONTACT_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "比利时 Bancontact", Available: true},
		{CurrencyCode: "EUR", PayType: "BANK_TRANSFER", InBankCode: "BANCONTACT_EUR", MinAmount: "0.99", MaxAmount: "1000", Description: "比利时 Bancontact", Available: true},
	},
	"IT": {
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "BANCOMATPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "意大利 Bancomat Pay", Available: true},
		{CurrencyCode: "EUR", PayType: "EWALLET", InBankCode: "BANCOMATPAY_EUR", MinAmount: "0.99", MaxAmount: "1000", Description: "意大利 Bancomat Pay", Available: true},
	},
	"NL": {
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "IDEAL_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "荷兰 IDEAL", Available: true},
		{CurrencyCode: "EUR", PayType: "BANK_TRANSFER", InBankCode: "IDEAL_EUR", MinAmount: "0.99", MaxAmount: "1000", Description: "荷兰 IDEAL", Available: true},
	},
	"PL": {
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "BLIK_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "波兰 BLIK", Available: true},
	},
	"BR": {
		{CurrencyCode: "BRL", PayType: "QR", InBankCode: "PIX", MinAmount: "5", MaxAmount: "20000", Description: "PIX二维码", Available: true},
	},
	"HK": {
		{CurrencyCode: "HKD", PayType: "EWALLET", InBankCode: "WXPAY_SCANCODE", MinAmount: "1", MaxAmount: "1000", Description: "香港微信扫码支付", Available: true},
		{CurrencyCode: "HKD", PayType: "EWALLET", InBankCode: "ALIPAY_HKD", MinAmount: "1", MaxAmount: "5000", Description: "AlipayHK", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "HK_WXPAY_SCANCODE_USD", MinAmount: "0.99", MaxAmount: "120", Description: "香港微信扫码支付", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "ALIPAY_USD", MinAmount: "0.99", MaxAmount: "499.99", Description: "AlipayHK", Available: true},
	},
	"SG": {
		{CurrencyCode: "SGD", PayType: "EWALLET", InBankCode: "PAYNOW_SGD", MinAmount: "1", MaxAmount: "3000", Description: "PayNow扫码支付", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "PAYNOW_USD", MinAmount: "1", MaxAmount: "3000", Description: "PayNow扫码支付", Available: true},
	},
	"TW": {
		{CurrencyCode: "TWD", PayType: "QR", InBankCode: "JKO_PAY_QR", MinAmount: "30", MaxAmount: "20000", Description: "JkoPay QR", Available: true},
		{CurrencyCode: "TWD", PayType: "EWALLET", InBankCode: "JKO_PAY_TWD", MinAmount: "30", MaxAmount: "20000", Description: "JkoPay", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "ATM", MinAmount: "1000", MaxAmount: "20000", Description: "ATM收银台", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_IBON", MinAmount: "100", MaxAmount: "20000", Description: "7-11超商", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_HILIFEET", MinAmount: "100", MaxAmount: "20000", Description: "萊爾富Life-ET", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_OKGO", MinAmount: "100", MaxAmount: "20000", Description: "OK超商OK-go", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_FAMI", MinAmount: "100", MaxAmount: "20000", Description: "全家FAMIPORT", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_MBC_HILIFEET", MinAmount: "100", MaxAmount: "20000", Description: "超商行動掃碼(萊爾富Life-ET)", Available: true},
		{CurrencyCode: "TWD", PayType: "BANK_TRANSFER", InBankCode: "STORE_MBC_IBON", MinAmount: "100", MaxAmount: "20000", Description: "7-11超商行動掃碼", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "JKO_PAY", MinAmount: "0.99", MaxAmount: "499.99", Description: "JkoPay", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "LINE_PAY_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "LinePay", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "ATM_USD", MinAmount: "32.99", MaxAmount: "499.99", Description: "ATM收银台", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_IBON_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "7-11超商", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_HILIFEET_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "萊爾富Life-ET", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_OKGO_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "OK超商OK-go", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_FAMI_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "全家FAMIPORT", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_MBC_HILIFEET_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "超商行動掃碼(萊爾富Life-ET)", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "STORE_MBC_IBON_USD", MinAmount: "2.99", MaxAmount: "499.99", Description: "超商行動掃碼(7-11)", Available: true},
	},
	"JP": {
		{CurrencyCode: "JPY", PayType: "EWALLET", InBankCode: "AU_PAY", MinAmount: "100", MaxAmount: "150000", Description: "AuPay", Available: true},
		{CurrencyCode: "JPY", PayType: "EWALLET", InBankCode: "MERPAY", MinAmount: "100", MaxAmount: "150000", Description: "MerPay", Available: true},
		{CurrencyCode: "JPY", PayType: "BANK_TRANSFER", InBankCode: "GMO_AOZORA", MinAmount: "100", MaxAmount: "150000", Description: "GMO Aozora VA", Available: true},
		{CurrencyCode: "JPY", PayType: "BANK_TRANSFER", InBankCode: "KONBINI_FAMILYMART", MinAmount: "300", MaxAmount: "150000", Description: "全家便利店", Available: true},
		{CurrencyCode: "JPY", PayType: "BANK_TRANSFER", InBankCode: "KONBINI_LAWSON", MinAmount: "300", MaxAmount: "150000", Description: "罗森便利店", Available: true},
		{CurrencyCode: "JPY", PayType: "BANK_TRANSFER", InBankCode: "KONBINI_SEVEN_ELEVEN", MinAmount: "300", MaxAmount: "150000", Description: "7-11便利店", Available: true},
		{CurrencyCode: "JPY", PayType: "BANK_TRANSFER", InBankCode: "KONBINI_MINISTOP", MinAmount: "300", MaxAmount: "150000", Description: "迷你岛便利店", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "JP_AU_PAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "AuPay", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "JP_MERPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "MerPay", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "JP_GMO_AOZORA_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "GMO Aozora VA", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "JP_KONBINI_FAMILYMART_USD", MinAmount: "1.99", MaxAmount: "1000", Description: "全家便利店", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "JP_KONBINI_LAWSON_USD", MinAmount: "1.99", MaxAmount: "1000", Description: "罗森便利店", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "JP_KONBINI_SEVEN_ELEVEN_USD", MinAmount: "1.99", MaxAmount: "1000", Description: "7-11便利店", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "JP_KONBINI_MINISTOP_USD", MinAmount: "1.99", MaxAmount: "1000", Description: "迷你岛便利店", Available: true},
	},
	"KR": {
		{CurrencyCode: "KRW", PayType: "BANK_TRANSFER", InBankCode: "KAKAOPAY_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "Kakao Pay钱包", Available: true},
		{CurrencyCode: "KRW", PayType: "BANK_TRANSFER", InBankCode: "NAVERPAY_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "Naver Pay钱包", Available: true},
		{CurrencyCode: "KRW", PayType: "BANK_TRANSFER", InBankCode: "TOSS_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "Toss Pay钱包", Available: true},
		{CurrencyCode: "KRW", PayType: "EWALLET", InBankCode: "SAMSUNGPAY_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "Samsung Pay钱包", Available: true},
		{CurrencyCode: "KRW", PayType: "EWALLET", InBankCode: "PAYCO_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "PayCo钱包", Available: true},
		{CurrencyCode: "KRW", PayType: "BANK_TRANSFER", InBankCode: "BANK_KRW", MinAmount: "1000", MaxAmount: "1500000", Description: "银行账号授权支付(支持韩国所有银行)", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "KAKAOPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "Kakao Pay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "NAVERPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "Naver Pay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "TOSS_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "Toss Pay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "SAMSUNGPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "Samsung Pay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "PAYCO_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "PayCo钱包", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "KR_BANK_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "银行账号授权支付(支持韩国所有银行)", Available: true},
	},
	"PH": {
		{CurrencyCode: "PHP", PayType: "QR", InBankCode: "PH_QRPH_DYNAMIC", MinAmount: "100", MaxAmount: "50000", Description: "QRPH动态码", Available: true},
		{CurrencyCode: "PHP", PayType: "PAYMENT_GATEWAY", InBankCode: "GCASH_QR", MinAmount: "100", MaxAmount: "50000", Description: "Gcash钱包(唤醒)", Available: true},
		{CurrencyCode: "PHP", PayType: "PAYMENT_GATEWAY", InBankCode: "GCASH_URL", MinAmount: "100", MaxAmount: "50000", Description: "Gcash钱包(直连)", Available: true},
		{CurrencyCode: "PHP", PayType: "PAYMENT_GATEWAY", InBankCode: "GRPY_URL", MinAmount: "100", MaxAmount: "50000", Description: "GrabPay钱包", Available: true},
		{CurrencyCode: "PHP", PayType: "PAYMENT_GATEWAY", InBankCode: "PAYMAYA_URL", MinAmount: "100", MaxAmount: "50000", Description: "PayMaya钱包", Available: true},
	},
	"TH": {
		{CurrencyCode: "THB", PayType: "QR", InBankCode: "QR", MinAmount: "1", MaxAmount: "8000", Description: "PromptPay QR二维码", Available: true},
		{CurrencyCode: "THB", PayType: "PAYMENT_GATEWAY", InBankCode: "TM", MinAmount: "1", MaxAmount: "8000", Description: "TrueMoney钱包", Available: true},
		{CurrencyCode: "USD", PayType: "QR", InBankCode: "TH_QR_USD", MinAmount: "0.99", MaxAmount: "220", Description: "PromptPay QR二维码", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "TH_TM_USD", MinAmount: "0.99", MaxAmount: "220", Description: "TrueMoney钱包", Available: true},
	},
	"VN": {
		{CurrencyCode: "VND", PayType: "QR", InBankCode: "QR", MinAmount: "10000", MaxAmount: "20000000", Description: "QR二维码", Available: true},
		{CurrencyCode: "VND", PayType: "BANK_TRANSFER", InBankCode: "BANK", MinAmount: "10000", MaxAmount: "20000000", Description: "银行转账", Available: true},
		{CurrencyCode: "VND", PayType: "EWALLET", InBankCode: "MOMO_VND", MinAmount: "100000", MaxAmount: "10000000", Description: "MOMO钱包", Available: true},
		{CurrencyCode: "USD", PayType: "QR", InBankCode: "VN_QR_USD", MinAmount: "0.99", MaxAmount: "2000", Description: "QR二维码", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "VN_BANK_USD", MinAmount: "0.99", MaxAmount: "2000", Description: "银行转账", Available: true},
	},
	"ID": {
		{CurrencyCode: "IDR", PayType: "QR", InBankCode: "dynamic", MinAmount: "10000", MaxAmount: "10000000", Description: "QRIS动态码", Available: true},
		{CurrencyCode: "IDR", PayType: "EWALLET", InBankCode: "DANA", MinAmount: "10000", MaxAmount: "5000000", Description: "Dana钱包", Available: true},
		{CurrencyCode: "IDR", PayType: "EWALLET", InBankCode: "SHOPEEPAY", MinAmount: "10000", MaxAmount: "20000000", Description: "ShopeePay钱包", Available: true},
		{CurrencyCode: "IDR", PayType: "EWALLET", InBankCode: "LINKAJA", MinAmount: "10000", MaxAmount: "20000000", Description: "LinkAja钱包", Available: true},
		{CurrencyCode: "IDR", PayType: "CASHIER", InBankCode: "OVO", MinAmount: "10000", MaxAmount: "20000000", Description: "OVO钱包", Available: true},
		{CurrencyCode: "USD", PayType: "QR", InBankCode: "ID_QRIS_USD", MinAmount: "0.99", MaxAmount: "500", Description: "QRIS动态码", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "ID_DANA_USD", MinAmount: "0.99", MaxAmount: "250", Description: "Dana钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "ID_SHOPEEPAY_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "ShopeePay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "ID_LINKAJA_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "LinkAja钱包", Available: true},
		{CurrencyCode: "USD", PayType: "CASHIER", InBankCode: "ID_OVO_USD", MinAmount: "0.99", MaxAmount: "1000", Description: "OVO钱包", Available: true},
	},
	"IN": {
		{CurrencyCode: "INR", PayType: "PAYMENT_GATEWAY", InBankCode: "UPI", MinAmount: "100", MaxAmount: "25000", Description: "UPI", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "IN_UPI_USD", MinAmount: "0.99", MaxAmount: "200", Description: "UPI", Available: true},
	},
	"MY": {
		{CurrencyCode: "MYR", PayType: "EWALLET", InBankCode: "TNG", MinAmount: "1", MaxAmount: "3000", Description: "Touch N Go", Available: true},
		{CurrencyCode: "MYR", PayType: "EWALLET", InBankCode: "GRAB", MinAmount: "1", MaxAmount: "3000", Description: "GrabPay", Available: true},
		{CurrencyCode: "MYR", PayType: "EWALLET", InBankCode: "BOOST", MinAmount: "1", MaxAmount: "3000", Description: "Boost", Available: true},
		{CurrencyCode: "MYR", PayType: "EWALLET", InBankCode: "SHOPEE", MinAmount: "1", MaxAmount: "3000", Description: "ShopeePay", Available: true},
		{CurrencyCode: "MYR", PayType: "EWALLET", InBankCode: "MCash", MinAmount: "1", MaxAmount: "3000", Description: "MCash", Available: true},
		{CurrencyCode: "MYR", PayType: "BANK_TRANSFER", InBankCode: "FPX_MYR", MinAmount: "1", MaxAmount: "3000", Description: "FPX银行转账", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "TNG_USD", MinAmount: "0.99", MaxAmount: "499.99", Description: "Touch N Go", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "BOOST_USD", MinAmount: "0.99", MaxAmount: "499.99", Description: "Boost", Available: true},
		{CurrencyCode: "USD", PayType: "BANK_TRANSFER", InBankCode: "FPX_USD", MinAmount: "0.99", MaxAmount: "499.99", Description: "FPX银行转账", Available: true},
	},
	"PK": {
		{CurrencyCode: "PKR", PayType: "PAYMENT_GATEWAY", InBankCode: "CASHIER", MinAmount: "10", MaxAmount: "500000", Description: "收银台模式", Available: true},
		{CurrencyCode: "PKR", PayType: "EWALLET", InBankCode: "JAZZCASH", MinAmount: "10", MaxAmount: "500000", Description: "JazzCash钱包", Available: true},
		{CurrencyCode: "PKR", PayType: "EWALLET", InBankCode: "EASYPAISA", MinAmount: "10", MaxAmount: "500000", Description: "EasyPaisa钱包", Available: true},
		{CurrencyCode: "PKR", PayType: "QR", InBankCode: "QR", MinAmount: "100", MaxAmount: "200000", Description: "QR", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "PK_JAZZCASH_USD", MinAmount: "0.5", MaxAmount: "1500", Description: "JazzCash钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "PK_EASYPAISA_USD", MinAmount: "0.5", MaxAmount: "1500", Description: "EasyPaisa钱包", Available: true},
		{CurrencyCode: "USD", PayType: "QR", InBankCode: "PK_QR_USD", MinAmount: "0.5", MaxAmount: "500", Description: "QR", Available: true},
	},
	"BD": {
		{CurrencyCode: "BDT", PayType: "EWALLET", InBankCode: "BDT_BKASH", MinAmount: "100", MaxAmount: "50000", Description: "bkash", Available: true},
		{CurrencyCode: "BDT", PayType: "EWALLET", InBankCode: "BDT_NAGAD", MinAmount: "100", MaxAmount: "50000", Description: "nagad", Available: true},
	},
	"EG": {
		{CurrencyCode: "EGP", PayType: "EWALLET", InBankCode: "REFERENCE_CODE", MinAmount: "10", MaxAmount: "30000", Description: "Fawry钱包", Available: true},
		{CurrencyCode: "EGP", PayType: "EWALLET", InBankCode: "VODAFONE", MinAmount: "5", MaxAmount: "30000", Description: "Vodafone钱包", Available: true},
		{CurrencyCode: "EGP", PayType: "EWALLET", InBankCode: "ORANGE", MinAmount: "5", MaxAmount: "30000", Description: "Orange Cash钱包", Available: true},
		{CurrencyCode: "EGP", PayType: "EWALLET", InBankCode: "ETISALAT", MinAmount: "5", MaxAmount: "30000", Description: "Etisalat Cash钱包", Available: true},
		{CurrencyCode: "EGP", PayType: "EWALLET", InBankCode: "WEPAY", MinAmount: "5", MaxAmount: "30000", Description: "WE Pay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "EG_REFERENCE_CODE_USD", MinAmount: "0.2", MaxAmount: "500", Description: "Fawry钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "EG_VODAFONE_USD", MinAmount: "0.1", MaxAmount: "500", Description: "Vodafone钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "EG_ORANGE_USD", MinAmount: "0.1", MaxAmount: "500", Description: "Orange Cash钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "EG_ETISALAT_USD", MinAmount: "0.1", MaxAmount: "500", Description: "Etisalat Cash钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "EG_WEPAY_USD", MinAmount: "0.1", MaxAmount: "500", Description: "WE Pay钱包", Available: true},
	},
	"SA": {
		{CurrencyCode: "SAR", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "2", MaxAmount: "10000", Description: "支持 VISA, MasterCard, Mada", Available: true},
		{CurrencyCode: "SAR", PayType: "PAYMENT_GATEWAY", InBankCode: "STCPAY", MinAmount: "2", MaxAmount: "10000", Description: "STCPay钱包", Available: true},
		{CurrencyCode: "SAR", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "2", MaxAmount: "10000", Description: "Apple Pay", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "SA_CARD_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "支持 VISA,MasterCard,Mada", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "SA_STCPAY_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "STCPay钱包", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "SA_APPLEPAY_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "Apple Pay", Available: true},
	},
	"AE": {
		{CurrencyCode: "AED", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "2", MaxAmount: "10000", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "AED", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "2", MaxAmount: "10000", Description: "Apple Pay", Available: true},
		{CurrencyCode: "AED", PayType: "VA", InBankCode: "BANK_TRANSFER", MinAmount: "10", MaxAmount: "100000000", Description: "VA BANK TRANSFER", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "AE_CARD_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "AE_APPLEPAY_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "Apple Pay", Available: true},
	},
	"KW": {
		{CurrencyCode: "KWD", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "0.5", MaxAmount: "800", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "KWD", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "0.5", MaxAmount: "800", Description: "Apple Pay", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "4.5", MaxAmount: "2500", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "4.5", MaxAmount: "2500", Description: "Apple Pay", Available: true},
	},
	"QA": {
		{CurrencyCode: "QAR", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "2", MaxAmount: "10000", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "QAR", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "2", MaxAmount: "10000", Description: "Apple Pay", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "QA_CARD_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "QA_APPLEPAY_USD", MinAmount: "0.5", MaxAmount: "2500", Description: "Apple Pay", Available: true},
	},
	"OM": {
		{CurrencyCode: "OMR", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "1.5", MaxAmount: "1000", Description: "支持VISA, MasterCard", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "OM_CARD_USD", MinAmount: "4", MaxAmount: "2500", Description: "支持VISA, MasterCard", Available: true},
	},
	"BH": {
		{CurrencyCode: "BHD", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "1.5", MaxAmount: "1000", Description: "支持 VISA, MasterCard, BENEFIT", Available: true},
		{CurrencyCode: "BHD", PayType: "EWALLET", InBankCode: "APPLE_PAY", MinAmount: "1.5", MaxAmount: "1000", Description: "Apple Pay", Available: true},
		{CurrencyCode: "USD", PayType: "PAYMENT_GATEWAY", InBankCode: "BH_CARD_USD", MinAmount: "4", MaxAmount: "2500", Description: "支持 VISA, MasterCard, BENEFIT", Available: true},
		{CurrencyCode: "USD", PayType: "EWALLET", InBankCode: "BH_APPLE_PAY_USD", MinAmount: "4", MaxAmount: "2500", Description: "Apple Pay", Available: true},
	},
	"JO": {
		{CurrencyCode: "JOD", PayType: "PAYMENT_GATEWAY", InBankCode: "CARD", MinAmount: "0.5", MaxAmount: "1999", Description: "支持 VISA, MasterCard", Available: true},
		{CurrencyCode: "JOD", PayType: "PAYMENT_GATEWAY", InBankCode: "CLIQ", MinAmount: "0.5", MaxAmount: "1999", Description: "CliQ", Available: true},
	},
	"IQ": {
		{CurrencyCode: "IQD", PayType: "EWALLET", InBankCode: "ZAIN", MinAmount: "250", MaxAmount: "10000000", Description: "", Available: true},
	},
}

// ListHaiPayCollectionPaymentMethods 返回指定地区支付方式目录的副本。
// 官网“全球”USD 方式会追加到所有支持 USD 的国家/地区。
func ListHaiPayCollectionPaymentMethods(countryCode string) []HaiPayCollectionPaymentMethod {
	countryCode = normalizeCode(countryCode)
	methods := haiPayCollectionPaymentMethodsByCountry[countryCode]
	out := append([]HaiPayCollectionPaymentMethod(nil), methods...)
	currencies := HaiPayCollectionCurrencies(countryCode)
	if len(currencies) == 0 {
		currencies = HaiPayGlobalCashierCurrencies(countryCode)
	}
	for _, currency := range currencies {
		if currency == "USD" {
			out = append(out, haiPayGlobalCollectionPaymentMethods...)
			break
		}
	}
	return out
}

// HasHaiPayCollectionPaymentMethodCatalog 表示该地区是否存在官网可用代收方式。
func HasHaiPayCollectionPaymentMethodCatalog(countryCode string) bool {
	return len(ListHaiPayCollectionPaymentMethods(countryCode)) > 0
}

// ResolveHaiPayCollectionPaymentMethodCode 校验币种、支付类型、支付编码组合，并返回文档中的规范编码。
// 支付类型和编码可同时留空；只选择支付类型时，验证该类型在当前币种下至少有一种可用方式。
func ResolveHaiPayCollectionPaymentMethodCode(countryCode, currencyCode, payType, inBankCode string) (string, bool) {
	currencyCode = strings.ToUpper(strings.TrimSpace(currencyCode))
	payType = strings.ToUpper(strings.TrimSpace(payType))
	inBankCode = strings.TrimSpace(inBankCode)
	if payType == "" {
		return "", inBankCode == ""
	}
	for _, method := range ListHaiPayCollectionPaymentMethods(countryCode) {
		if !method.Available || method.CurrencyCode != currencyCode || method.PayType != payType {
			continue
		}
		if inBankCode == "" {
			return "", true
		}
		if strings.EqualFold(method.InBankCode, inBankCode) {
			return method.InBankCode, true
		}
	}
	return "", false
}

// IsHaiPayCollectionPaymentMethod 判断币种、支付类型和支付编码是否为可用组合。
func IsHaiPayCollectionPaymentMethod(countryCode, currencyCode, payType, inBankCode string) bool {
	_, ok := ResolveHaiPayCollectionPaymentMethodCode(countryCode, currencyCode, payType, inBankCode)
	return ok
}
