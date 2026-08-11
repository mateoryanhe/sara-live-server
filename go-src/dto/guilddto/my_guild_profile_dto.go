package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetMyGuildProfileReq struct {
	g.Meta `path:"/getMyGuildProfile" method:"post" summary:"获取当前CMS用户关联工会基础信息" tags:"直播工会"`
}

type GetMyGuildProfileRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BankCard    string `json:"bankCard"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updatedAt"`
}

type UpdateMyGuildProfileReq struct {
	g.Meta      `path:"/updateMyGuildProfile" method:"post" summary:"更新当前CMS用户关联工会基础信息" tags:"直播工会"`
	Name        string `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	BankCard    string `json:"bankCard" dc:"银行卡"`
	Contact     string `json:"contact" dc:"联系方式"`
	Description string `json:"description" dc:"工会简介"`
}

type UpdateMyGuildProfileRes struct {
	Success bool `json:"success"`
}
