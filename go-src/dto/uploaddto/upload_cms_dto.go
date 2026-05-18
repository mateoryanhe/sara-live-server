package uploaddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadCMSFileReq CMS后台上传文件(图片/礼物动画资源)
type UploadCMSFileReq struct {
	g.Meta `path:"/uploadFile" method:"post" mime:"multipart/form-data" summary:"CMS后台上传文件" tags:"上传管理"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择上传文件" dc:"图片或礼物动画资源"`
}

type UploadCMSFileRes struct {
	FileName string `json:"fileName" dc:"保存后的文件名"`
}
