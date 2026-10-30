package userinfo

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	inviteeListDefaultPageSize = 20
	inviteeListMaxPageSize     = 100
)

// InviteeList 当前用户名下邀请人列表(分页)
func InviteeList(ctx context.Context, req *userinfodto.InviteeListReq) (*userinfodto.InviteeListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	page, pageSize := normalizeInviteeListPage(req.Page, req.PageSize)
	pageData := userinfodao.GetAppInviteesByInviter(userId, page, pageSize)
	list := make([]*userinfodto.InviteeListItem, 0, len(pageData))
	for _, row := range pageData {
		if row == nil {
			continue
		}
		resolved := appInviteeListRowWithMemoryProfile(row)
		list = append(list, &userinfodto.InviteeListItem{
			UserId:            strconv.FormatUint(resolved.UserId, 10),
			Nickname:          resolved.Nickname,
			Avatar:            upload.ResolveAvatarUrlForUser(resolved.UserId, resolved.Avatar),
			InviterRewardGold: resolved.InviterRewardGold,
		})
	}
	return &userinfodto.InviteeListRes{
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func normalizeInviteeListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = inviteeListDefaultPageSize
	}
	if pageSize > inviteeListMaxPageSize {
		pageSize = inviteeListMaxPageSize
	}
	return page, pageSize
}

func appInviteeListRowWithMemoryProfile(row *userinfodao.AppInviteeListRow) *userinfodao.AppInviteeListRow {
	if row == nil {
		return row
	}
	resolved := *row
	if user := userinfodao.GetUserInfoFromMemory(row.UserId); user != nil {
		resolved.Nickname = user.Nickname
		resolved.Avatar = user.Avatar
	}
	return &resolved
}
