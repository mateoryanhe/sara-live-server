package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppUserShortVideoListReq App端查询指定用户短视频列表(仅已上架,按审核时间降序,走内存缓存)
type AppUserShortVideoListReq struct {
	g.Meta   `path:"/appUserShortVideoList" method:"post" summary:"App查询指定用户短视频列表(已上架,按审核时间降序,走内存缓存)" tags:"短视频"`
	UserId   uint64 `json:"userId,string" v:"required#用户ID不能为空" dc:"目标用户ID"`
	Page     int    `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int    `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}
