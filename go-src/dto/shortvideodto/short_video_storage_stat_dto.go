package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type ShortVideoStorageStatReq struct {
	g.Meta `path:"/shortVideoStorageStat" method:"post" summary:"短视频存储统计" tags:"短视频"`
}

type ShortVideoStorageStatRes struct {
	TotalCount        int     `json:"totalCount" dc:"短视频总记录数"`
	ImageDirPath      string  `json:"imageDirPath" dc:"图片/短视频静态目录物理路径"`
	ImageDirUsedBytes uint64  `json:"imageDirUsedBytes" dc:"静态目录占用字节数"`
	DiskTotalBytes    uint64  `json:"diskTotalBytes" dc:"磁盘总容量(字节)"`
	DiskFreeBytes     uint64  `json:"diskFreeBytes" dc:"磁盘空闲容量(字节)"`
	DiskFreeRatio     float64 `json:"diskFreeRatio" dc:"磁盘空闲比例(0-100)"`
}
