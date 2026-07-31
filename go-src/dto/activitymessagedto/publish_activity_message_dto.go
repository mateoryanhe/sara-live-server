package activitymessagedto

import "github.com/gogf/gf/v2/frame/g"

type PublishActivityMessageReq struct {
	g.Meta `path:"/publishActivityMessage" method:"post" summary:"发布活动消息" tags:"活动消息"`
	ID     uint64 `json:"id" v:"required#活动消息ID不能为空"`
}

type PublishActivityMessageRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type UnpublishActivityMessageReq struct {
	g.Meta `path:"/unpublishActivityMessage" method:"post" summary:"取消发布活动消息" tags:"活动消息"`
	ID     uint64 `json:"id" v:"required#活动消息ID不能为空"`
}

type UnpublishActivityMessageRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
