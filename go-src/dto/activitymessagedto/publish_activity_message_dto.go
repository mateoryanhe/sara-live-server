package activitymessagedto

import "github.com/gogf/gf/v2/frame/g"

type PublishActivityMessageReq struct {
	g.Meta `path:"/publishActivityMessage" method:"post" summary:"发布活动消息" tags:"活动消息"`
	ID     uint64 `json:"id" v:"required#活动消息ID不能为空" dc:"活动消息ID"`
}

type PublishActivityMessageRes struct {
	Success bool  `json:"success" dc:"是否成功"`
	Status  uint8 `json:"status" dc:"发布状态(0未发布,1已发布)"`
}

type UnpublishActivityMessageReq struct {
	g.Meta `path:"/unpublishActivityMessage" method:"post" summary:"取消发布活动消息" tags:"活动消息"`
	ID     uint64 `json:"id" v:"required#活动消息ID不能为空" dc:"活动消息ID"`
}

type UnpublishActivityMessageRes struct {
	Success bool  `json:"success" dc:"是否成功"`
	Status  uint8 `json:"status" dc:"发布状态(0未发布,1已发布)"`
}
