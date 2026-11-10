package anchorvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// VideoListReq App端查询主播视频列表
type VideoListReq struct {
	g.Meta   `path:"/list" method:"post" summary:"查询主播相册视频列表" tags:"主播视频"`
	AnchorId uint64 `json:"anchorId,string" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

type VideoItem struct {
	Id    uint64 `json:"id,string" dc:"视频ID"`
	Video string `json:"video" dc:"视频完整URL"`
}

type VideoListRes struct {
	List []*VideoItem `json:"list"`
}

// UploadVideoReq App端上传主播相册视频(可一次传多个,form 字段名均为 file;仅 mp4)
type UploadVideoReq struct {
	g.Meta `path:"/upload" method:"post" mime:"multipart/form-data" summary:"上传主播相册视频(可多个)" tags:"主播视频"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择视频" dc:"视频文件(mp4),可重复传多个 file 字段"`
}

// UploadVideoRes 本次成功写入的视频列表
type UploadVideoRes struct {
	List []*VideoItem `json:"list" dc:"本次上传成功的视频"`
}

// DeleteVideoReq App端批量删除自己的相册视频
type DeleteVideoReq struct {
	g.Meta   `path:"/delete" method:"post" summary:"批量删除主播相册视频" tags:"主播视频"`
	VideoIds []uint64 `json:"videoIds,string" v:"required|min-length:1#请至少选择一个视频" dc:"视频ID列表"`
}

type DeleteVideoRes struct {
	Success    bool     `json:"success"`
	DeletedIds []uint64 `json:"deletedIds,string" dc:"实际删除成功的视频ID"`
}
