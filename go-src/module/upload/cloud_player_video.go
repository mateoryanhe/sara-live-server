package upload

import "strings"

// ResolveCloudPlayerVideoUrl 云播放器MP4文件名转可访问URL;已是完整URL则原样返回
func ResolveCloudPlayerVideoUrl(video string) string {
	if video == "" {
		return ""
	}
	if strings.HasPrefix(video, "http://") || strings.HasPrefix(video, "https://") {
		return video
	}
	return GetUrlByName(video)
}
