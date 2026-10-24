package accountcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetAccountCfgReq struct {
	g.Meta `path:"/getAccountCfg" method:"post" summary:"查询账号配置" tags:"账号配置"`
}

type AccountCfgItem struct {
	ID                         string `json:"id"`
	CancelAccountByCodeEnabled bool   `json:"cancelAccountByCodeEnabled"`
	BlockSimulatorLogin        bool   `json:"blockSimulatorLogin"`
	EnvType                    uint8  `json:"envType"`
	DeviceRegisterRiskEnabled  bool   `json:"deviceRegisterRiskEnabled"`
	DeviceAccountMaxCount      int    `json:"deviceAccountMaxCount"`
	DeviceCancelDailyLimit     int     `json:"deviceCancelDailyLimit"`
	UserWithdrawEnabled        bool    `json:"userWithdrawEnabled"`
	UserWithdrawAndroidVisible bool    `json:"userWithdrawAndroidVisible"`
	UserWithdrawMinUsd         float64 `json:"userWithdrawMinUsd"`
	UserWithdrawDailyMaxCount  int     `json:"userWithdrawDailyMaxCount"`
	CreatedAt                  string  `json:"createdAt"`
	UpdatedAt                  string  `json:"updatedAt"`
}

type GetAccountCfgRes struct {
	Cfg *AccountCfgItem `json:"cfg"`
}

type SaveAccountCfgReq struct {
	g.Meta                     `path:"/saveAccountCfg" method:"post" summary:"保存账号配置" tags:"账号配置"`
	ID                         uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	CancelAccountByCodeEnabled bool   `json:"cancelAccountByCodeEnabled" dc:"注销码销户开关(官网公开接口)"`
	BlockSimulatorLogin        bool   `json:"blockSimulatorLogin" dc:"拦截模拟器登录(默认关闭=不拦截)"`
	EnvType                    uint8  `json:"envType" v:"in:0,1,2#环境类型无效" dc:"环境类型(0正式服,1提审服,2测试服)"`
	DeviceRegisterRiskEnabled  bool   `json:"deviceRegisterRiskEnabled" dc:"设备码注册风控开关"`
	DeviceAccountMaxCount      int    `json:"deviceAccountMaxCount" dc:"同设备最大账号数(含已注销)"`
	DeviceCancelDailyLimit     int     `json:"deviceCancelDailyLimit" dc:"同设备每天最多注销次数"`
	UserWithdrawEnabled        bool    `json:"userWithdrawEnabled" dc:"用户提现开关"`
	UserWithdrawAndroidVisible bool    `json:"userWithdrawAndroidVisible" dc:"安卓端提现入口展示开关(仅客户端展示)"`
	UserWithdrawMinUsd         float64 `json:"userWithdrawMinUsd" dc:"用户最低提现金额(USD)"`
	UserWithdrawDailyMaxCount  int     `json:"userWithdrawDailyMaxCount" dc:"用户每天最多提现次数(0=不限制)"`
}

type SaveAccountCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
