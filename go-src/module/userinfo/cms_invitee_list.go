package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

// GetCMSInviteeList CMS分页查询某用户邀请的下级用户
func GetCMSInviteeList(ctx context.Context, req *accountdto.CMSInviteeListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := userinfodao.CMSInviteeList(&userinfodao.CMSInviteeListFilter{
		InviterId: req.InviterId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*accountdto.CMSInviteeItem, 0, len(rows))
	for _, row := range rows {
		if item := toCMSInviteeItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toCMSInviteeItem(v *userentity.UserInfo) *accountdto.CMSInviteeItem {
	if v == nil {
		return nil
	}
	createdAt := v.CreatedAt
	return &accountdto.CMSInviteeItem{
		UserId:    v.ID,
		Nickname:  v.Nickname,
		Avatar:    upload.ResolveAvatarUrlForUser(v.ID, v.Avatar),
		CreatedAt: &createdAt,
	}
}
