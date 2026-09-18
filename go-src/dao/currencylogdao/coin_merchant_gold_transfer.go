package currencylogdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

type CoinMerchantGoldTransferListFilter struct {
	MerchantUserId uint64
	TargetUserId   uint64
	StartTime      int64
	EndTime        int64
	PageIndex      int
	PageSize       int
}

// CoinMerchantGoldTransferListRow 是转账记录与目标 user_info 的联表查询结果。
type CoinMerchantGoldTransferListRow struct {
	ID             uint64    `orm:"id"`
	TargetUserId   uint64    `orm:"target_user_id"`
	TargetNickname string    `orm:"target_nickname"`
	TargetAvatar   string    `orm:"target_avatar"`
	Amount         float64   `orm:"amount"`
	CreatedAt      time.Time `orm:"created_at"`
}

// ListCoinMerchantGoldTransferLogs 直接查询数据库并关联 user_infos，按最新记录倒序分页。
func ListCoinMerchantGoldTransferLogs(filter *CoinMerchantGoldTransferListFilter) (int, []*CoinMerchantGoldTransferListRow, error) {
	rows := make([]*CoinMerchantGoldTransferListRow, 0)
	if filter == nil || filter.MerchantUserId == 0 {
		return 0, rows, nil
	}
	if filter.PageIndex <= 0 {
		filter.PageIndex = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	const transferAlias = "cmgtl"
	const userAlias = "target_user"
	model := g.DB().Model(string(userentity.TbCoinMerchantGoldTransferLog)+" "+transferAlias).Ctx(gctx.New()).
		LeftJoin(string(userentity.TbUserInfo)+" "+userAlias,
			userAlias+".id = "+transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetUserId)).
		Where(transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogMerchantUserId)+" = ?", filter.MerchantUserId)
	if filter.TargetUserId > 0 {
		model = model.Where(transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetUserId)+" = ?", filter.TargetUserId)
	}
	if filter.StartTime > 0 {
		model = model.Where(transferAlias+".created_at >= ?", time.Unix(filter.StartTime, 0))
	}
	if filter.EndTime > 0 {
		model = model.Where(transferAlias+".created_at <= ?", time.Unix(filter.EndTime, 0))
	}

	total, err := model.Clone().Count()
	if err != nil {
		return 0, rows, err
	}
	err = model.Clone().Fields(
		transferAlias+".id",
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetUserId),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogAmount),
		transferAlias+".created_at",
		userAlias+"."+string(userentity.UserInfoNickname)+" AS target_nickname",
		userAlias+"."+string(userentity.UserInfoAvatar)+" AS target_avatar",
	).Order(transferAlias + ".created_at desc, " + transferAlias + ".id desc").
		Limit(filter.PageSize).Offset((filter.PageIndex - 1) * filter.PageSize).
		Scan(&rows)
	return total, rows, err
}
