package country

import "strings"

// HaiPayAppID 是 HaiPay 后台固定创建的业务 ID。
type HaiPayAppID int64

const (
	HaiPayAppIDIDR     HaiPayAppID = 25238
	HaiPayAppIDPHP     HaiPayAppID = 25239
	HaiPayAppIDUSD     HaiPayAppID = 25240
	HaiPayAppIDKRW     HaiPayAppID = 25241
	HaiPayAppIDEUR     HaiPayAppID = 25242
	HaiPayAppIDTWD     HaiPayAppID = 25243
	HaiPayAppIDVND     HaiPayAppID = 25244
	HaiPayAppIDTHB     HaiPayAppID = 25245
	HaiPayAppIDSAR     HaiPayAppID = 25246
	HaiPayAppIDKWD     HaiPayAppID = 25247
	HaiPayAppIDOMR     HaiPayAppID = 25248
	HaiPayAppIDQAR     HaiPayAppID = 25249
	HaiPayAppIDAED     HaiPayAppID = 25250
	HaiPayAppIDBHD     HaiPayAppID = 25251
	HaiPayAppIDUSDT    HaiPayAppID = 25252
	HaiPayAppIDBRL     HaiPayAppID = 25253
	HaiPayAppIDHKD     HaiPayAppID = 25254
	HaiPayAppIDINR     HaiPayAppID = 25255
	HaiPayAppIDMYR     HaiPayAppID = 25256
	HaiPayAppIDTRY     HaiPayAppID = 25257
	HaiPayAppIDPKR     HaiPayAppID = 25258
	HaiPayAppIDSGD     HaiPayAppID = 25259
	HaiPayAppIDRUB     HaiPayAppID = 25260
	HaiPayAppIDEGP     HaiPayAppID = 25261
	HaiPayAppIDJPY     HaiPayAppID = 25262
	HaiPayAppIDGBP     HaiPayAppID = 25263
	HaiPayAppIDPLN     HaiPayAppID = 25264
	HaiPayAppIDCAD     HaiPayAppID = 25265
	HaiPayAppIDJOD     HaiPayAppID = 25266
	HaiPayAppIDIQD     HaiPayAppID = 25267
	HaiPayAppIDBDT     HaiPayAppID = 25268
	HaiPayAppIDUSDC    HaiPayAppID = 25269
	HaiPayAppIDMXN     HaiPayAppID = 25270
	HaiPayAppIDNGN     HaiPayAppID = 25271
	HaiPayAppIDCashier HaiPayAppID = 25272
)

const (
	HaiPayAppIDProdIDR     HaiPayAppID = 7724
	HaiPayAppIDProdUSD     HaiPayAppID = 7810
	HaiPayAppIDProdMYR     HaiPayAppID = 7730
	HaiPayAppIDProdCashier HaiPayAppID = 7545
)

var haiPayAppIDByBusinessCodeTest = map[string]HaiPayAppID{
	"IDR": HaiPayAppIDIDR, "PHP": HaiPayAppIDPHP, "USD": HaiPayAppIDUSD,
	"KRW": HaiPayAppIDKRW, "EUR": HaiPayAppIDEUR, "TWD": HaiPayAppIDTWD,
	"VND": HaiPayAppIDVND, "THB": HaiPayAppIDTHB, "SAR": HaiPayAppIDSAR,
	"KWD": HaiPayAppIDKWD, "OMR": HaiPayAppIDOMR, "QAR": HaiPayAppIDQAR,
	"AED": HaiPayAppIDAED, "BHD": HaiPayAppIDBHD, "USDT": HaiPayAppIDUSDT,
	"BRL": HaiPayAppIDBRL, "HKD": HaiPayAppIDHKD, "INR": HaiPayAppIDINR,
	"MYR": HaiPayAppIDMYR, "TRY": HaiPayAppIDTRY, "PKR": HaiPayAppIDPKR,
	"SGD": HaiPayAppIDSGD, "RUB": HaiPayAppIDRUB, "EGP": HaiPayAppIDEGP,
	"JPY": HaiPayAppIDJPY, "GBP": HaiPayAppIDGBP, "PLN": HaiPayAppIDPLN,
	"CAD": HaiPayAppIDCAD, "JOD": HaiPayAppIDJOD, "IQD": HaiPayAppIDIQD,
	"BDT": HaiPayAppIDBDT, "USDC": HaiPayAppIDUSDC, "MXN": HaiPayAppIDMXN,
	"NGN": HaiPayAppIDNGN, "CASHIER": HaiPayAppIDCashier,
}

var haiPayAppIDByBusinessCodeProd = map[string]HaiPayAppID{
	"IDR": HaiPayAppIDProdIDR, "USD": HaiPayAppIDProdUSD, "MYR": HaiPayAppIDProdMYR,
	"CASHIER": HaiPayAppIDProdCashier,
}

// LookupHaiPayAppID 按币种或 CASHIER 业务编码返回 AppID。
// useProd=false 用测试业务 ID，true 用正式业务 ID（当前仅 IDR/USD/MYR/收银台）。
func LookupHaiPayAppID(businessCode string, useProd bool) (int64, bool) {
	table := haiPayAppIDByBusinessCodeTest
	if useProd {
		table = haiPayAppIDByBusinessCodeProd
	}
	appID, ok := table[strings.ToUpper(strings.TrimSpace(businessCode))]
	return int64(appID), ok
}

// HaiPayCashierAppID 返回全球收银台 AppID。
func HaiPayCashierAppID(useProd bool) int64 {
	appID, ok := LookupHaiPayAppID("CASHIER", useProd)
	if !ok {
		return 0
	}
	return appID
}
