package userwithdrawdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

type UserWithdrawListFilter struct {
	UserId   uint64
	Status   *uint8
	Page     int
	PageSize int
}

func GetUserWithdrawLogById(id uint64) *userentity.UserWithdrawLog {
	if id == 0 {
		return nil
	}
	var row userentity.UserWithdrawLog
	err := g.DB().Model(string(userentity.TbUserWithdrawLog)).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func InsertUserWithdrawLog(row *userentity.UserWithdrawLog) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	_, err := g.DB().Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New()).Data(row).Insert()
	return err
}

func UpdateUserWithdrawLog(row *userentity.UserWithdrawLog) error {
	if row == nil || row.ID == 0 {
		return nil
	}
	_, err := g.DB().Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New()).WherePri(row.ID).Data(row).Update()
	return err
}

func ListUserWithdrawLogs(filter *UserWithdrawListFilter) (int, []*userentity.UserWithdrawLog, error) {
	rows := make([]*userentity.UserWithdrawLog, 0)
	if filter == nil {
		return 0, rows, nil
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	model := g.DB().Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New())
	if filter.UserId > 0 {
		model = model.Where(string(userentity.UserWithdrawLogUserId)+" = ?", filter.UserId)
	}
	if filter.Status != nil {
		model = model.Where(string(userentity.UserWithdrawLogStatus)+" = ?", *filter.Status)
	}
	total, err := model.Count()
	if err != nil {
		return 0, nil, err
	}
	err = model.OrderDesc("id").
		Page(page, pageSize).
		Scan(&rows)
	if err != nil {
		return 0, nil, err
	}
	return total, rows, nil
}
