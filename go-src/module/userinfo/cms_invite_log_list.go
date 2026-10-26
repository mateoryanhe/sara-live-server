package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

// GetCMSInviteLogList CMS分页查询全局邀请日志
func GetCMSInviteLogList(ctx context.Context, req *accountdto.CMSInviteLogListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := userinfodao.CMSInviteLogList(&userinfodao.CMSInviteLogListFilter{
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	inviterIds := collectInviteLogInviterIds(rows)
	inviterMap := userinfodao.GetUserProfileMapByUserIds(inviterIds)
	list := make([]*accountdto.CMSInviteLogItem, 0, len(rows))
	for _, row := range rows {
		if item := toCMSInviteLogItem(row, inviterMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func collectInviteLogInviterIds(rows []*userentity.UserInfo) []uint64 {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.InviterId == 0 {
			continue
		}
		if _, ok := seen[row.InviterId]; ok {
			continue
		}
		seen[row.InviterId] = struct{}{}
		ids = append(ids, row.InviterId)
	}
	return ids
}

func toCMSInviteLogItem(v *userentity.UserInfo, inviterMap map[uint64]*userentity.UserInfo) *accountdto.CMSInviteLogItem {
	if v == nil {
		return nil
	}
	createdAt := v.CreatedAt
	item := &accountdto.CMSInviteLogItem{
		UserId:    v.ID,
		Nickname:  v.Nickname,
		Avatar:    upload.ResolveAvatarUrlForUser(v.ID, v.Avatar),
		InviterId: v.InviterId,
		CreatedAt: &createdAt,
	}
	if inviterMap != nil {
		if inv := inviterMap[v.InviterId]; inv != nil {
			item.InviterNickname = inv.Nickname
			item.InviterAvatar = upload.ResolveAvatarUrlForUser(inv.ID, inv.Avatar)
		}
	}
	return item
}
