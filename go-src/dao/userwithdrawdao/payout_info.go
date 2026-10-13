package userwithdrawdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	userentity "xr-game-server/entity/user"
)

func GetUserPayoutInfo(userId uint64) *userentity.UserPayoutInfo {
	if userId == 0 {
		return nil
	}
	var row userentity.UserPayoutInfo
	err := g.DB().Model(string(userentity.TbUserPayoutInfo)).WherePri(userId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func SaveUserPayoutInfo(row *userentity.UserPayoutInfo) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	now := time.Now()
	row.UpdatedAt = now
	existing := GetUserPayoutInfo(row.ID)
	if existing == nil {
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		_, err := g.DB().Model(string(userentity.TbUserPayoutInfo)).Data(row).Insert()
		return err
	}
	row.CreatedAt = existing.CreatedAt
	_, err := g.DB().Model(string(userentity.TbUserPayoutInfo)).WherePri(row.ID).Data(row).Update()
	return err
}
