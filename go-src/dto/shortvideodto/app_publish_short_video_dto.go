package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AppPublishShortVideoReq struct {
	g.Meta           `path:"/appPublishShortVideo" method:"post" mime:"multipart/form-data" summary:"App上传并发布短视频" tags:"短视频"`
	File             *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择短视频文件" dc:"短视频文件"`
	Cover            *ghttp.UploadFile `json:"cover" type:"file" dc:"封面图片(可选)"`
	Title            string            `json:"title" v:"required|length:1,64#标题不能为空|标题长度需在1到64之间" dc:"标题"`
	IsPaid           uint8             `json:"isPaid" v:"in:0,1#是否付费取值无效" dc:"是否付费(0免费,1付费)"`
	DiamondPerMinute float64           `json:"diamondPerMinute" dc:"每分钟钻石数(付费时必填)"`
	CategoryId       int               `json:"categoryId" dc:"视频分类ID"`
	Source           uint8             `json:"source" v:"required|in:1,2,3#视频来源不能为空|视频来源取值无效" dc:"视频来源(1原创,2转发,3AI生成)"`
}

type AppPublishShortVideoRes struct {
	ID    string `json:"id" dc:"短视频 ID"`
	Video string `json:"video" dc:"视频完整URL"`
	Cover string `json:"cover" dc:"封面完整URL"`
}
