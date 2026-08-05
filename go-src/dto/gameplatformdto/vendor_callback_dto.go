package gameplatformdto

import "github.com/gogf/gf/v2/frame/g"

// VendorVerifyReq 第三方身份验证回调请求.
type VendorVerifyReq struct {
	g.Meta        `path:"/verify" method:"post" summary:"第三方身份验证回调" tags:"游戏第三方回调"`
	OperatorToken string `json:"operator_token"`
	Ops           string `json:"ops"`
	Sign          string `json:"sign"`
	Timestamp     int64  `json:"timestamp"`
}

type VendorVerifyData struct {
	PlayerName string  `json:"player_name"`
	Currency   string  `json:"currency"`
	Balance    float64 `json:"balance"`
}

type VendorVerifyRes struct {
	Code    int               `json:"code"`
	Message string            `json:"message,omitempty"`
	Data    *VendorVerifyData `json:"data,omitempty"`
}

// VendorBalanceReq 第三方获取余额回调请求.
type VendorBalanceReq struct {
	g.Meta        `path:"/balance" method:"post" summary:"第三方获取余额回调" tags:"游戏第三方回调"`
	OperatorToken string `json:"operator_token"`
	PlayerName    string `json:"player_name"`
	Sign          string `json:"sign"`
	Timestamp     int64  `json:"timestamp"`
}

type VendorBalanceData struct {
	Balance      float64 `json:"balance"`
	CurrencyCode string  `json:"currency_code,omitempty"`
}

type VendorBalanceRes struct {
	Code    int                `json:"code"`
	Message string             `json:"message,omitempty"`
	Data    *VendorBalanceData `json:"data,omitempty"`
}
