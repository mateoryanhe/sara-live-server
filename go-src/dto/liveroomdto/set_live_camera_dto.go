package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// SetLiveCameraEnabledReq 主播开关摄像头(默认开;关时不计有效直播时长).
type SetLiveCameraEnabledReq struct {
	g.Meta  `path:"/setLiveCameraEnabled" method:"post" summary:"开关摄像头" tags:"直播间"`
	RoomId  uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID(须为主播本人)"`
	Enabled bool   `json:"enabled" dc:"true=开摄像头,false=关摄像头"`
}

type SetLiveCameraEnabledRes struct {
	Success       bool `json:"success"`
	CameraEnabled bool `json:"cameraEnabled"`
}
