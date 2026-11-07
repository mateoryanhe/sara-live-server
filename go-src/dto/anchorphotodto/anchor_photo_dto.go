package anchorphotodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// PhotoListReq App端查询主播相册
type PhotoListReq struct {
	g.Meta   `path:"/list" method:"post" summary:"查询主播相册照片列表" tags:"主播相册"`
	AnchorId uint64 `json:"anchorId,string" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

type PhotoItem struct {
	Id    uint64 `json:"id,string" dc:"照片ID"`
	Photo string `json:"photo" dc:"照片完整URL"`
}

type PhotoListRes struct {
	List []*PhotoItem `json:"list"`
}

// UploadPhotoReq App端上传主播相册照片(可一次传多张,form 字段名均为 file)
type UploadPhotoReq struct {
	g.Meta `path:"/upload" method:"post" mime:"multipart/form-data" summary:"上传主播相册照片(可多张)" tags:"主播相册"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择图片" dc:"照片文件,可重复传多个 file 字段"`
}

// UploadPhotoRes 本次成功写入的照片列表
type UploadPhotoRes struct {
	List []*PhotoItem `json:"list" dc:"本次上传成功的照片"`
}

// DeletePhotoReq App端批量删除自己的相册照片
type DeletePhotoReq struct {
	g.Meta   `path:"/delete" method:"post" summary:"批量删除主播相册照片" tags:"主播相册"`
	PhotoIds []uint64 `json:"photoIds,string" v:"required|min-length:1#请至少选择一张照片" dc:"照片ID列表"`
}

type DeletePhotoRes struct {
	Success    bool     `json:"success"`
	DeletedIds []uint64 `json:"deletedIds,string" dc:"实际删除成功的照片ID"`
}
