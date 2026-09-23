package liveroom

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	platformAnchorVisibilityTimeLayout = "2006-01-02 15:04:05"
	batchPlatformAnchorVisibilityMax   = 500
)

func ensurePlatformAnchorExists(anchorId uint64) error {
	if anchorId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	room := liveroomdao.GetRoomFromDB(anchorId)
	if room == nil || room.GuildId != 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	user := userinfodao.GetUserInfoFromDB(anchorId)
	if user == nil || (user.UserType != userentity.UserTypeAnchor && user.UserType != userentity.UserTypeSeniorAnchor) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

// GetPlatformAnchorListForVisibility 可见性管理页拉取全部平台主播，不应用当前用户可见性过滤。
func GetPlatformAnchorListForVisibility(_ context.Context, req *accountdto.PlatformAnchorListForVisibilityReq) (*httpserver.CMSQueryResp, error) {
	listReq := &accountdto.QueryAnchorListReq{
		CMSQueryReq:  req.CMSQueryReq,
		Key:          req.Key,
		PlatformOnly: true,
	}
	total, data := queryAnchorListFromMemory(listReq)
	return httpserver.NewCMSQueryResp(total, data), nil
}

// GetPlatformAnchorVisibilityList 查询平台主播已授权的 CMS 用户。
func GetPlatformAnchorVisibilityList(_ context.Context, req *accountdto.PlatformAnchorVisibilityListReq) (*accountdto.PlatformAnchorVisibilityListRes, error) {
	if err := ensurePlatformAnchorExists(req.AnchorId); err != nil {
		return nil, err
	}
	rows := liveroomdao.ListPlatformAnchorVisibilitiesByAnchorId(req.AnchorId)
	profileMap := userinfodao.GetUserProfileMapByUserIds([]uint64{req.AnchorId})
	list := make([]*accountdto.PlatformAnchorVisibilityItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := buildPlatformAnchorVisibilityItem(row.ID, row.AnchorId, row.CmsUserId, row.CreatedAt, profileMap)
		list = append(list, item)
	}
	return &accountdto.PlatformAnchorVisibilityListRes{List: list}, nil
}

// GetPlatformAnchorVisibilityByUserList 查询指定 CMS 用户已授权的平台主播。
func GetPlatformAnchorVisibilityByUserList(_ context.Context, req *accountdto.PlatformAnchorVisibilityByUserListReq) (*accountdto.PlatformAnchorVisibilityByUserListRes, error) {
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := liveroomdao.ListPlatformAnchorVisibilitiesByCmsUserId(req.CmsUserId)
	anchorIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.AnchorId > 0 {
			anchorIds = append(anchorIds, row.AnchorId)
		}
	}
	profileMap := userinfodao.GetUserProfileMapByUserIds(anchorIds)
	list := make([]*accountdto.PlatformAnchorVisibilityItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := buildPlatformAnchorVisibilityItem(row.ID, row.AnchorId, row.CmsUserId, row.CreatedAt, profileMap)
		list = append(list, item)
	}
	return &accountdto.PlatformAnchorVisibilityByUserListRes{List: list}, nil
}

// GrantPlatformAnchorVisibility 授权 CMS 用户查看平台主播。
func GrantPlatformAnchorVisibility(_ context.Context, req *accountdto.GrantPlatformAnchorVisibilityReq) (*accountdto.GrantPlatformAnchorVisibilityRes, error) {
	if err := ensurePlatformAnchorExists(req.AnchorId); err != nil {
		return nil, err
	}
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := liveroomdao.AddPlatformAnchorVisibility(req.AnchorId, req.CmsUserId); err != nil {
		return nil, err
	}
	return &accountdto.GrantPlatformAnchorVisibilityRes{Success: true}, nil
}

// BatchGrantPlatformAnchorVisibility 批量授权 CMS 用户查看平台主播。
func BatchGrantPlatformAnchorVisibility(_ context.Context, req *accountdto.BatchGrantPlatformAnchorVisibilityReq) (*accountdto.BatchGrantPlatformAnchorVisibilityRes, error) {
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil || len(req.AnchorIds) == 0 || len(req.AnchorIds) > batchPlatformAnchorVisibilityMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	granted := 0
	seen := make(map[uint64]struct{}, len(req.AnchorIds))
	for _, anchorId := range req.AnchorIds {
		if anchorId == 0 {
			continue
		}
		if _, ok := seen[anchorId]; ok {
			continue
		}
		seen[anchorId] = struct{}{}
		if ensurePlatformAnchorExists(anchorId) != nil || liveroomdao.HasPlatformAnchorVisibility(anchorId, req.CmsUserId) {
			continue
		}
		if err := liveroomdao.AddPlatformAnchorVisibility(anchorId, req.CmsUserId); err != nil {
			return nil, err
		}
		granted++
	}
	return &accountdto.BatchGrantPlatformAnchorVisibilityRes{Success: true, GrantedCount: granted}, nil
}

// RevokePlatformAnchorVisibility 撤销 CMS 用户的平台主播可见性。
func RevokePlatformAnchorVisibility(_ context.Context, req *accountdto.RevokePlatformAnchorVisibilityReq) (*accountdto.RevokePlatformAnchorVisibilityRes, error) {
	if err := ensurePlatformAnchorExists(req.AnchorId); err != nil {
		return nil, err
	}
	if err := liveroomdao.RemovePlatformAnchorVisibility(req.AnchorId, req.CmsUserId); err != nil {
		return nil, err
	}
	return &accountdto.RevokePlatformAnchorVisibilityRes{Success: true}, nil
}

// BatchRevokePlatformAnchorVisibility 批量撤销 CMS 用户的平台主播可见性。
func BatchRevokePlatformAnchorVisibility(_ context.Context, req *accountdto.BatchRevokePlatformAnchorVisibilityReq) (*accountdto.BatchRevokePlatformAnchorVisibilityRes, error) {
	if len(req.AnchorIds) == 0 || len(req.AnchorIds) > batchPlatformAnchorVisibilityMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	revoked := 0
	seen := make(map[uint64]struct{}, len(req.AnchorIds))
	for _, anchorId := range req.AnchorIds {
		if anchorId == 0 {
			continue
		}
		if _, ok := seen[anchorId]; ok {
			continue
		}
		seen[anchorId] = struct{}{}
		if !liveroomdao.HasPlatformAnchorVisibility(anchorId, req.CmsUserId) {
			continue
		}
		if err := liveroomdao.RemovePlatformAnchorVisibility(anchorId, req.CmsUserId); err != nil {
			return nil, err
		}
		revoked++
	}
	return &accountdto.BatchRevokePlatformAnchorVisibilityRes{Success: true, RevokedCount: revoked}, nil
}

func buildPlatformAnchorVisibilityItem(id, anchorId, cmsUserId uint64, createdAt time.Time, profileMap map[uint64]*userentity.UserInfo) *accountdto.PlatformAnchorVisibilityItem {
	item := &accountdto.PlatformAnchorVisibilityItem{
		Id:        strconv.FormatUint(id, 10),
		AnchorId:  strconv.FormatUint(anchorId, 10),
		CmsUserId: strconv.FormatUint(cmsUserId, 10),
		CreatedAt: formatPlatformAnchorVisibilityTime(createdAt),
	}
	if user := profileMap[anchorId]; user != nil {
		item.AnchorName = user.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrlForUser(anchorId, user.Avatar)
	}
	if cmsUser := cmsuserdao.GetCMSUserById(cmsUserId); cmsUser != nil {
		item.CmsUserName = cmsUser.Name
	}
	return item
}

func formatPlatformAnchorVisibilityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(platformAnchorVisibilityTimeLayout)
}
