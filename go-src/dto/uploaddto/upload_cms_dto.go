package uploaddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UploadCMSFileReq CMS后台上传文件;路由由 RegCMSHandler 注册
type UploadCMSFileReq struct {
	g.Meta `path:"/uploadFile" method:"post" mime:"multipart/form-data" summary:"CMS后台上传文件" tags:"上传管理"`
}

type UploadCMSFileRes struct {
	FileName string `json:"fileName" dc:"保存后的文件名"`
	FileUrl  string `json:"fileUrl" dc:"文件访问URL"`
}
