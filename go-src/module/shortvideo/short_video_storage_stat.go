package shortvideo

import (
	"context"

	"xr-game-server/core/cfg"
	"xr-game-server/core/disk"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
)

// GetShortVideoStorageStat CMS 查询短视频存储与磁盘占用
func GetShortVideoStorageStat(_ context.Context, _ *shortvideodto.ShortVideoStorageStatReq) (*shortvideodto.ShortVideoStorageStatRes, error) {
	root := cfg.GetImageStaticRoot()
	res := &shortvideodto.ShortVideoStorageStatRes{
		TotalCount:   shortvideodao.CountAll(),
		ImageDirPath: root,
	}
	if root == "" {
		return res, nil
	}
	if used, err := disk.DirSize(root); err == nil {
		res.ImageDirUsedBytes = used
	}
	if diskTotal, diskFree, err := disk.GetDiskUsage(root); err == nil && diskTotal > 0 {
		res.DiskTotalBytes = diskTotal
		res.DiskFreeBytes = diskFree
		res.DiskFreeRatio = float64(diskFree) / float64(diskTotal) * 100
	}
	return res, nil
}
