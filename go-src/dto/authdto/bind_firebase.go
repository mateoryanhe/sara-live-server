package authdto

import "github.com/gogf/gf/v2/frame/g"

// BindFirebaseReq 当前登录用户通过 Firebase ID Token 绑定 Firebase 账号。
type BindFirebaseReq struct {
	g.Meta  `path:"/bindFirebase" method:"post" summary:"绑定Firebase账号" tags:"权限"`
	IdToken string `json:"idToken" v:"required|max-length:16384#Firebase ID Token不能为空|Firebase ID Token过长" dc:"Firebase客户端登录后获取的ID Token"`
}

type BindFirebaseRes struct {
	Success bool `json:"success"`
}
