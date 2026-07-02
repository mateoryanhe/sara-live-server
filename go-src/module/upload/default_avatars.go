package upload

import "math/rand"

// 默认头像列表(文件位于 upload/images,由 pub-tool/avatars 上传)
var defaultAvatarNames = []string{
	"demo_avatar_1.jpg",
	"demo_avatar_2.jpg",
	"demo_avatar_3.jpg",
	"demo_avatar_4.jpg",
	"demo_avatar_5.jpg",
	"demo_avatar_6.jpg",
	"demo_avatar_7.jpg",
	"demo_avatar_8.jpg",
	"demo_avatar_9.jpg",
	"demo_avatar_10.jpg",
}

// PickDefaultAvatarName 按用户 ID 从默认头像列表中取一个(同一用户固定,未上传头像时使用)
func PickDefaultAvatarName(userId uint64) string {
	n := len(defaultAvatarNames)
	if n == 0 {
		return ""
	}
	if userId == 0 {
		return PickRandomDefaultAvatarName()
	}
	return defaultAvatarNames[userId%uint64(n)]
}

// PickRandomDefaultAvatarName 从默认头像列表中随机取一个
func PickRandomDefaultAvatarName() string {
	n := len(defaultAvatarNames)
	if n == 0 {
		return ""
	}
	return defaultAvatarNames[rand.Intn(n)]
}

// defaultAvatarUrl 默认头像完整 URL
func defaultAvatarUrl(name string) string {
	if name == "" {
		return ""
	}
	return GetUrlByName(name)
}
