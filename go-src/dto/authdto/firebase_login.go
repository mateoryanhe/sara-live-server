package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

// FirebaseLoginReq Firebase ID Token 登录(不存在则自动注册).
type FirebaseLoginReq struct {
	g.Meta     `path:"/firebaseLogin" method:"post" summary:"Firebase登录" tags:"权限"`
	IdToken    string             `json:"idToken" v:"required|max-length:16384#Firebase ID Token不能为空|Firebase ID Token过长" dc:"Firebase客户端登录后获取的ID Token"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo" dc:"设备信息(可选)"`
}

type FirebaseLoginRes struct {
	Token     string `json:"token"`
	IsNewUser bool   `json:"isNewUser" dc:"是否首次注册"`
}
