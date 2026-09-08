package coinmerchantdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CoinMerchantListReq CMS币商列表
type CoinMerchantListReq struct {
	g.Meta `path:"/coinMerchantList" method:"post" summary:"币商列表" tags:"币商"`
	httpserver.CMSQueryReq
	Key    string `json:"key" dc:"用户ID模糊/用户名精确"`
	Cancel *int   `json:"cancel" dc:"注销状态:1已注销,0正常,不传不过滤"`
}

// CoinMerchantItem 币商列表项
type CoinMerchantItem struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Gold      float64    `json:"gold"`
	Cancel    bool       `json:"cancel"`
	Channel   uint       `json:"channel"`
	CreatedAt *time.Time `json:"createdAt"`
}

// CreateCoinMerchantReq CMS新建币商
type CreateCoinMerchantReq struct {
	g.Meta   `path:"/createCoinMerchant" method:"post" summary:"新建币商" tags:"币商"`
	Username string `json:"username" v:"required|length:2,32#用户名不能为空|用户名长度2-32" dc:"用户名(openId)"`
	Password string `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32" dc:"明文密码(服务端MD5入库)"`
}

// CreateCoinMerchantRes 新建结果
type CreateCoinMerchantRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// ResetCoinMerchantPasswordReq CMS重置币商密码
type ResetCoinMerchantPasswordReq struct {
	g.Meta    `path:"/resetCoinMerchantPassword" method:"post" summary:"重置币商密码" tags:"币商"`
	AccountId uint64 `json:"accountId" v:"required#账号ID不能为空" dc:"账号ID"`
	Password  string `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32" dc:"明文密码(服务端MD5入库)"`
}

// ResetCoinMerchantPasswordRes 重置结果
type ResetCoinMerchantPasswordRes struct {
	Success bool `json:"success"`
}

// CancelCoinMerchantReq CMS注销币商
type CancelCoinMerchantReq struct {
	g.Meta    `path:"/cancelCoinMerchant" method:"post" summary:"注销币商" tags:"币商"`
	AccountId uint64 `json:"accountId" v:"required#账号ID不能为空" dc:"账号ID"`
}

// CancelCoinMerchantRes 注销结果
type CancelCoinMerchantRes struct {
	Success bool `json:"success"`
}
