package userwithdrawtiercfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type UserWithdrawTierCfgListReq struct {
	g.Meta `path:"/userWithdrawTierCfgList" method:"post" summary:"用户提现档位列表" tags:"用户提现档位"`
	httpserver.CMSQueryReq
	Name         string `json:"name" dc:"名称模糊"`
	StatusFilter int    `json:"statusFilter" dc:"0全部,1只看下架,2只看上架"`
}

type UserWithdrawTierCfgListRes struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	WithdrawUsd    float64 `json:"withdrawUsd"`
	MinGameConsumeGold float64 `json:"minGameConsumeGold"`
	Status         uint8   `json:"status"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

type CreateUserWithdrawTierCfgReq struct {
	g.Meta         `path:"/createUserWithdrawTierCfg" method:"post" summary:"新建用户提现档位" tags:"用户提现档位"`
	Name           string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64"`
	WithdrawUsd    float64 `json:"withdrawUsd" v:"required|min:0.01#提现USD不能为空|提现USD需大于0"`
	MinGameConsumeGold float64 `json:"minGameConsumeGold" dc:"游戏消费累计金币门槛(仅配置)"`
}

type CreateUserWithdrawTierCfgRes struct {
	ID string `json:"id"`
}

type UpdateUserWithdrawTierCfgReq struct {
	g.Meta         `path:"/updateUserWithdrawTierCfg" method:"post" summary:"修改用户提现档位" tags:"用户提现档位"`
	ID             uint64  `json:"id" v:"required#ID不能为空"`
	Name           string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64"`
	WithdrawUsd    float64 `json:"withdrawUsd" v:"required|min:0.01#提现USD不能为空|提现USD需大于0"`
	MinGameConsumeGold float64 `json:"minGameConsumeGold" dc:"游戏消费累计金币门槛(仅配置)"`
}

type UpdateUserWithdrawTierCfgRes struct {
	Success bool `json:"success"`
}

type DeleteUserWithdrawTierCfgReq struct {
	g.Meta `path:"/deleteUserWithdrawTierCfg" method:"post" summary:"删除用户提现档位" tags:"用户提现档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type DeleteUserWithdrawTierCfgRes struct {
	Success bool `json:"success"`
}

type OnShelfUserWithdrawTierCfgReq struct {
	g.Meta `path:"/onShelfUserWithdrawTierCfg" method:"post" summary:"上架用户提现档位" tags:"用户提现档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OnShelfUserWithdrawTierCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfUserWithdrawTierCfgReq struct {
	g.Meta `path:"/offShelfUserWithdrawTierCfg" method:"post" summary:"下架用户提现档位" tags:"用户提现档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OffShelfUserWithdrawTierCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type AppUserWithdrawTierListReq struct {
	g.Meta `path:"/listWithdrawTiers" method:"post" summary:"App查询提现档位(已上架)" tags:"用户提现"`
}

type AppUserWithdrawTierItem struct {
	ID             uint64  `json:"id,string"`
	Name           string  `json:"name"`
	WithdrawUsd    float64 `json:"withdrawUsd"`
	MinGameConsumeGold float64 `json:"minGameConsumeGold"`
}

type AppUserWithdrawTierListRes struct {
	List []*AppUserWithdrawTierItem `json:"list"`
}
