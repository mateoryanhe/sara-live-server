package guild

import (
	"context"
	"time"

	"xr-game-server/core/xrtime"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// BatchImportSalaryAnchors CMS批量标记有底薪主播，不区分工会。
func BatchImportSalaryAnchors(_ context.Context, req *guilddto.BatchImportSalaryAnchorsReq) (*guilddto.BatchImportSalaryAnchorsRes, error) {
	if req == nil || len(req.IDs) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &guilddto.BatchImportSalaryAnchorsRes{
		FailIds: make([]uint64, 0),
	}
	now := time.Now()
	salaryEffectiveStartTime := xrtime.WeekStart(now)
	salaryEffectiveEndTime := xrtime.NextWeekStart(now)
	seen := make(map[uint64]struct{}, len(req.IDs))
	for _, userID := range req.IDs {
		if userID == 0 {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}

		user := userinfodao.GetUserInfoByUserId(userID)
		room := liveroomdao.ResolveRoom(userID)
		if user == nil || !userentity.UserTypeIsAnchor(user.UserType) || room == nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, userID)
			continue
		}
		room.SetHasSalary(true)
		room.SetSalaryEffectiveStartTime(&salaryEffectiveStartTime)
		room.SetSalaryEffectiveEndTime(&salaryEffectiveEndTime)
		liveroomdao.FlushRoomCache(room)
		res.SuccessCount++
	}
	return res, nil
}
