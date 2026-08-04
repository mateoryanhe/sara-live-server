package shortvideo

import (
	"xr-game-server/core/cfg"
	"xr-game-server/core/disk"
	"xr-game-server/errercode"
)

const shortVideoMinDiskFreeRatioPercent = 20

// ensureShortVideoUploadDiskSpace 上传前检查静态资源目录所在磁盘空闲比是否大于 20%
func ensureShortVideoUploadDiskSpace() error {
	root := cfg.GetImageStaticRoot()
	if root == "" {
		return nil
	}
	total, free, err := disk.GetDiskUsage(root)
	if err != nil || total == 0 {
		return nil
	}
	ratio := float64(free) / float64(total) * 100
	if ratio <= shortVideoMinDiskFreeRatioPercent {
		return errercode.CreateCode(errercode.ShortVideoDiskSpaceInsufficient)
	}
	return nil
}
