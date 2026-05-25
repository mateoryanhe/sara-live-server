package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadShortVideoCMSReq CMS后台上传短视频文件
type UploadShortVideoCMSReq struct {
	g.Meta `path:"/uploadShortVideo" method:"post" mime:"multipart/form-data" summary:"CMS上传短视频文件" tags:"短视频"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择短视频文件" dc:"短视频文件"`
}

type UploadShortVideoCMSRes struct {
	FileName string `json:"fileName" dc:"视频资源文件名"`
	Url      string `json:"url" dc:"视频完整URL"`
}
