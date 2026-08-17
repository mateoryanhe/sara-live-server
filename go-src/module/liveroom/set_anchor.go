package liveroom

import (
	"context"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SetAnchor CMS 将用户设为主播(仅允许从未主播变为主播,不可回退)
func SetAnchor(ctx context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	res, err := setUserAsAnchor(req.AccountId, entity.UserTypeAnchor)
	if err != nil {
		return nil, err
	}
	RefreshRoomListCache(ctx)
	return res, nil
}

// SetSeniorAnchor CMS 将用户设为高级主播(仅允许普通用户,不可回退)
func SetSeniorAnchor(ctx context.Context, req *accountdto.SetSeniorAnchorReq) (*accountdto.SetSeniorAnchorRes, error) {
	_, err := setUserAsAnchor(req.AccountId, entity.UserTypeSeniorAnchor)
	if err != nil {
		return nil, err
	}
	RefreshRoomListCache(ctx)
	return &accountdto.SetSeniorAnchorRes{Success: true}, nil
}

// BatchSetAnchor CMS 批量设普通主播
func BatchSetAnchor(ctx context.Context, req *accountdto.BatchSetAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return batchSetAnchor(ctx, req.IDs, entity.UserTypeAnchor)
}

// BatchSetSeniorAnchor CMS 批量设高级主播
func BatchSetSeniorAnchor(ctx context.Context, req *accountdto.BatchSetSeniorAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return batchSetAnchor(ctx, req.IDs, entity.UserTypeSeniorAnchor)
}

func batchSetAnchor(ctx context.Context, ids []uint64, userType uint8) (*accountdto.BatchSetAnchorRes, error) {
	res := &accountdto.BatchSetAnchorRes{}
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if err := setUserAsAnchorBatch(id, userType); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, id)
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount > 0 {
		RefreshRoomListCache(ctx)
	}
	return res, nil
}

func setUserAsAnchor(accountId uint64, userType uint8) (*accountdto.SetAnchorRes, error) {
	if err := setUserAsAnchorBatch(accountId, userType); err != nil {
		return nil, err
	}
	return &accountdto.SetAnchorRes{Success: true}, nil
}

func setUserAsAnchorBatch(accountId uint64, userType uint8) error {
	user := userinfodao.GetUserInfoByUserId(accountId)
	if entity.UserTypeIsAnchor(user.UserType) {
		return nil
	}
	if user.UserType != entity.UserTypeNormal {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	user.SetUserType(userType)
	EnsureAnchorRoom(accountId, 0)
	return nil
}

// SetUserAsAnchorIfNeeded 将普通用户设为主播类型;已是主播则跳过
func SetUserAsAnchorIfNeeded(accountId uint64, userType uint8) error {
	return setUserAsAnchorBatch(accountId, userType)
}

// ExitGuild CMS主播退出工会(将 live_room.guild_id 置为0)
func ExitGuild(ctx context.Context, req *accountdto.ExitGuildReq) (*accountdto.ExitGuildRes, error) {
	user := userinfodao.GetUserInfoByUserId(req.AnchorId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !entity.UserTypeIsAnchor(user.UserType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	room := EnsureAnchorRoom(req.AnchorId, 0)
	room.SetGuildId(0)
	RefreshRoomListCache(ctx)
	return &accountdto.ExitGuildRes{Success: true}, nil
}
