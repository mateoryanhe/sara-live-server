package cfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/userwithdrawtiercfgdto"
	userentity "xr-game-server/entity/user"
)

func GetUserWithdrawTierCfgById(id uint64) *userentity.UserWithdrawTierCfg {
	var cfg userentity.UserWithdrawTierCfg
	err := g.DB().Model(string(userentity.TbUserWithdrawTierCfg)).Where("id = ?", id).Scan(&cfg)
	if err != nil || cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func GetUserWithdrawTierCfgByName(name string) *userentity.UserWithdrawTierCfg {
	var cfg userentity.UserWithdrawTierCfg
	err := g.DB().Model(string(userentity.TbUserWithdrawTierCfg)).Where("name = ?", name).Scan(&cfg)
	if err != nil || cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func CreateUserWithdrawTierCfg(cfg *userentity.UserWithdrawTierCfg) error {
	_, err := g.DB().Model(string(userentity.TbUserWithdrawTierCfg)).Save(cfg)
	return err
}

func DeleteUserWithdrawTierCfg(id uint64) error {
	_, err := g.DB().Model(string(userentity.TbUserWithdrawTierCfg)).WherePri(id).Delete()
	return err
}

func GetOnShelfUserWithdrawTierCfgs() []*userentity.UserWithdrawTierCfg {
	ret := make([]*userentity.UserWithdrawTierCfg, 0)
	_ = g.DB().Model(string(userentity.TbUserWithdrawTierCfg)).
		Where("status = ?", userentity.UserWithdrawTierCfgStatusOnShelf).
		Order("withdraw_usd asc, id asc").
		Scan(&ret)
	return ret
}

func GetUserWithdrawTierCfgList(req *userwithdrawtiercfgdto.UserWithdrawTierCfgListReq) (int, []*userwithdrawtiercfgdto.UserWithdrawTierCfgListRes) {
	sql := `select id, name, withdraw_usd, min_game_consume_gold, status, created_at, updated_at
		from user_withdraw_tier_cfgs where 1=1`
	param := make([]any, 0)
	if req.Name != "" {
		sql += ` and name like ?`
		param = append(param, "%"+req.Name+"%")
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = 0`
	case 2:
		sql += ` and status = 1`
	}
	sql += ` order by withdraw_usd asc, id desc`
	ctx := gctx.New()
	total, _ := g.DB().GetCount(ctx, str.GetCountSQL(sql), param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())

	type row struct {
		ID             uint64  `json:"id"`
		Name           string  `json:"name"`
		WithdrawUsd    float64 `json:"withdraw_usd"`
		MinGameConsumeGold float64 `json:"min_game_consume_gold"`
		Status         uint8   `json:"status"`
		CreatedAt      string  `json:"created_at"`
		UpdatedAt      string  `json:"updated_at"`
	}
	rows := make([]*row, 0)
	_ = g.DB().GetScan(ctx, &rows, sql, param)
	list := make([]*userwithdrawtiercfgdto.UserWithdrawTierCfgListRes, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		list = append(list, &userwithdrawtiercfgdto.UserWithdrawTierCfgListRes{
			ID:             strconv.FormatUint(r.ID, 10),
			Name:           r.Name,
			WithdrawUsd:    r.WithdrawUsd,
			MinGameConsumeGold: r.MinGameConsumeGold,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}
	return total, list
}
