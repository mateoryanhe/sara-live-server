package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

// CoinMerchantLoginReq 币商用户名+密码登录(App,不自动注册)
type CoinMerchantLoginReq struct {
	g.Meta     `path:"/coinMerchantLogin" method:"post" summary:"币商登录" tags:"权限"`
	Username   string             `json:"username" v:"required|length:2,32#用户名不能为空|用户名长度2-32" dc:"用户名"`
	Password   string             `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32" dc:"客户端MD5后的密码(与库中哈希比对)"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo" dc:"设备信息"`
}

// CoinMerchantLoginRes 币商登录响应
type CoinMerchantLoginRes struct {
	Token string `json:"token"`
}
