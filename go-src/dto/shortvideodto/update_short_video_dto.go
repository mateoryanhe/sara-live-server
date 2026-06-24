package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type UpdateShortVideoReq struct {
	g.Meta           `path:"/updateShortVideo" method:"post" summary:"修改短视频" tags:"短视频"`
	ID               uint64  `json:"id" v:"required#短视频 ID不能为空" dc:"短视频 ID"`
	Title            string  `json:"title" v:"required|length:1,64#标题不能为空|标题长度需在1到64之间" dc:"标题"`
	Cover            string  `json:"cover" v:"max-length:255#封面资源名最长255字符" dc:"封面资源名"`
	Sort             int     `json:"sort" dc:"排序值(越大越靠前)"`
	IsPaid           uint8   `json:"isPaid" v:"in:0,1#是否付费取值无效" dc:"是否付费(0免费,1付费)"`
	DiamondPerMinute float64 `json:"diamondPerMinute" dc:"每分钟钻石数(付费时必填)"`
	CategoryId       int     `json:"categoryId" dc:"视频分类ID"`
	Source           uint8   `json:"source" v:"required|in:1,2,3#视频来源不能为空|视频来源取值无效" dc:"视频来源(1原创,2转发,3AI生成)"`
	FreeWatchSeconds uint32  `json:"freeWatchSeconds" v:"min:0#免费观看时长不能小于0" dc:"免费观看时长(秒)"`
}

type UpdateShortVideoRes struct {
	Success bool `json:"success"`
}
