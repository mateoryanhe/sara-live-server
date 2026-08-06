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

// VendorTransferReq 第三方下注转账回调请求.
type VendorTransferReq struct {
	g.Meta                `path:"/transfer" method:"post" summary:"第三方下注转账回调" tags:"游戏第三方回调"`
	OperatorToken         string  `json:"operator_token"`
	OperatorPlayerSession string  `json:"operator_player_session"`
	SecretKey             string  `json:"secret_key"`
	IP                    string  `json:"ip"`
	GameID                string  `json:"game_id"`
	PlayerName            string  `json:"player_name"`
	ParentBetID           string  `json:"parent_bet_id"`
	BetID                 string  `json:"bet_id"`
	BetAmount             float64 `json:"bet_amount"`
	WinAmount             float64 `json:"win_amount"`
	TransferAmount        float64 `json:"transfer_amount"`
	TransactionID         string  `json:"transaction_id"`
	CurrencyCode          string  `json:"currency_code"`
	Platform              string  `json:"platform"`
	WalletType            string  `json:"wallet_type"`
	CreateTime            int64   `json:"create_time"`
	UpdatedTime           int64   `json:"updated_time"`
	IsFeature             string  `json:"is_feature"`
	IsFeatureBuy          string  `json:"is_feature_buy"`
	BetType               int     `json:"bet_type"`
	RealTransferAmount    float64 `json:"real_transfer_amount"`
	InventoryAmount       float64 `json:"inventory_amount"`
	InventoryPoolID       string  `json:"inventory_pool_id"`
}

type VendorTransferData struct {
	Balance      float64 `json:"balance"`
	CurrencyCode string  `json:"currency_code"`
	UpdatedTime  int64   `json:"updated_time"`
}

type VendorTransferError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type VendorTransferRes struct {
	Data  *VendorTransferData  `json:"data"`
	Error *VendorTransferError `json:"error"`
}
