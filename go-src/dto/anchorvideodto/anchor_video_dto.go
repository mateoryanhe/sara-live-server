package anchorvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// VideoListReq App端查询主播视频列表(库拉全量最多1000后再分页)
type VideoListReq struct {
	g.Meta    `path:"/list" method:"post" summary:"查询主播相册视频列表" tags:"主播视频"`
	AnchorId  uint64 `json:"anchorId,string" v:"required#主播ID不能为空" dc:"主播用户ID"`
	PageIndex int    `json:"pageIndex" dc:"页码(从1开始,默认1)"`
	PageSize  int    `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

type VideoItem struct {
	Id          uint64  `json:"id,string" dc:"视频ID"`
	Video       string  `json:"video" dc:"视频完整URL"`
	Cover       string  `json:"cover" dc:"封面(原样返回)"`
	FileType    string  `json:"fileType" dc:"文件类型(原样返回)"`
	Duration    uint32  `json:"duration" dc:"视频秒数"`
	Price       float64 `json:"price" dc:"钻石价格(0=免费)"`
	Description string  `json:"description" dc:"描述(可空)"`
	Paid        bool    `json:"paid" dc:"当前用户是否已付费/可看:价格0或主播本人或已购为true"`
}

type VideoListRes struct {
	List      []*VideoItem `json:"list"`
	Total     int          `json:"total" dc:"有效视频总数(最多按库限制1000统计)"`
	PageIndex int          `json:"pageIndex"`
	PageSize  int          `json:"pageSize"`
}

// UnlockVideoReq App端付费解锁主播相册视频
type UnlockVideoReq struct {
	g.Meta  `path:"/unlock" method:"post" summary:"付费解锁主播相册视频" tags:"主播视频"`
	VideoId uint64 `json:"videoId,string" v:"required#视频ID不能为空" dc:"相册视频ID"`
}

type UnlockVideoRes struct {
	Deducted float64 `json:"deducted" dc:"本次扣除钻石数(已解锁/免费则为0)"`
	Diamond  float64 `json:"diamond" dc:"扣费后钻石余额"`
	Paid     bool    `json:"paid" dc:"是否已可看"`
}

// UploadVideoReq App端上传主播相册视频(仅单个;multipart: file + price + description)
type UploadVideoReq struct {
	g.Meta      `path:"/upload" method:"post" mime:"multipart/form-data" summary:"上传主播相册视频(单个)" tags:"主播视频"`
	File        *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择视频" dc:"视频文件(mp4,仅1个)"`
	Price       float64           `json:"price" v:"required#请选择价格" dc:"钻石价格(0=免费)"`
	Description string            `json:"description" dc:"描述(可空,最多200字)"`
	Cover       string            `json:"cover" dc:"封面(可空,原样存储)"`
	FileType    string            `json:"fileType" dc:"文件类型(可空,原样存储)"`
	Duration    string            `json:"duration" dc:"视频秒数(可空,原样解析为整数,非法则为0)"`
}

// UploadVideoRes 本次上传成功的视频
type UploadVideoRes struct {
	Id          uint64  `json:"id,string" dc:"视频ID"`
	Video       string  `json:"video" dc:"视频完整URL"`
	Cover       string  `json:"cover" dc:"封面(原样返回)"`
	FileType    string  `json:"fileType" dc:"文件类型(原样返回)"`
	Duration    uint32  `json:"duration" dc:"视频秒数"`
	Price       float64 `json:"price" dc:"钻石价格"`
	Description string  `json:"description" dc:"描述(可空)"`
}

// UpdateVideoReq App端修改自己的相册视频价格/描述
type UpdateVideoReq struct {
	g.Meta      `path:"/update" method:"post" summary:"修改主播相册视频价格或描述" tags:"主播视频"`
	VideoId     uint64  `json:"videoId,string" v:"required#视频ID不能为空" dc:"相册视频ID"`
	Price       float64 `json:"price" dc:"钻石价格(0=免费)"`
	Description string  `json:"description" dc:"描述(可空,最多200字)"`
}

type UpdateVideoRes struct {
	Id          uint64  `json:"id,string" dc:"视频ID"`
	Video       string  `json:"video" dc:"视频完整URL"`
	Cover       string  `json:"cover" dc:"封面(原样返回)"`
	FileType    string  `json:"fileType" dc:"文件类型(原样返回)"`
	Duration    uint32  `json:"duration" dc:"视频秒数"`
	Price       float64 `json:"price" dc:"钻石价格"`
	Description string  `json:"description" dc:"描述(可空)"`
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
