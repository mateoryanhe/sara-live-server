package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type UpdateShortVideoReq struct {
	g.Meta           `path:"/updateShortVideo" method:"post" summary:"修改短视频" tags:"短视频"`
	ID               uint64 `json:"id" v:"required#短视频 ID不能为空" dc:"短视频 ID"`
	Title            string `json:"title" v:"required|length:1,64#标题不能为空|标题长度需在1到64之间" dc:"标题"`
	Video            string `json:"video" v:"required|max-length:255#视频不能为空|视频资源名最长255字符" dc:"视频资源名"`
	Cover            string `json:"cover" v:"max-length:255#封面资源名最长255字符" dc:"封面资源名"`
	Sort             int    `json:"sort" dc:"排序值(越大越靠前)"`
	IsPaid           uint8  `json:"isPaid" v:"in:0,1#是否付费取值无效" dc:"是否付费(0免费,1付费)"`
	DiamondPerSecond uint64 `json:"diamondPerSecond" dc:"每秒钻石数(付费时必填)"`
	Description      string `json:"description" v:"max-length:255#描述最长255字符" dc:"描述"`
}

type UpdateShortVideoRes struct {
	Success bool `json:"success"`
}
