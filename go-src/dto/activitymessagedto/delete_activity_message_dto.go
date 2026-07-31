package activitymessagedto

import "github.com/gogf/gf/v2/frame/g"

type DeleteActivityMessageReq struct {
	g.Meta `path:"/deleteActivityMessage" method:"post" summary:"删除活动消息" tags:"活动消息"`
	ID     uint64 `json:"id" v:"required#活动消息ID不能为空" dc:"活动消息ID"`
}

type DeleteActivityMessageRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
