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

// CoinMerchantGoldTransferListRow 是转账记录与币商、目标用户 user_info 的联表查询结果。
type CoinMerchantGoldTransferListRow struct {
	ID                   uint64    `orm:"id"`
	CoinMerchantUserId   uint64    `orm:"coin_merchant_user_id"`
	CoinMerchantNickname string    `orm:"coin_merchant_nickname"`
	CoinMerchantAvatar   string    `orm:"coin_merchant_avatar"`
	TargetUserId         uint64    `orm:"target_user_id"`
	TargetNickname       string    `orm:"target_nickname"`
	TargetAvatar         string    `orm:"target_avatar"`
	Amount               float64   `orm:"amount"`
	SenderGoldBefore     float64   `orm:"sender_gold_before"`
	SenderGoldAfter      float64   `orm:"sender_gold_after"`
	TargetGoldBefore     float64   `orm:"target_gold_before"`
	TargetGoldAfter      float64   `orm:"target_gold_after"`
	CreatedAt            time.Time `orm:"created_at"`
}

// ListCoinMerchantGoldTransferLogs 直接查询数据库并关联 user_infos，按最新记录倒序分页。
func ListCoinMerchantGoldTransferLogs(filter *CoinMerchantGoldTransferListFilter) (int, []*CoinMerchantGoldTransferListRow, error) {
	rows := make([]*CoinMerchantGoldTransferListRow, 0)
	if filter == nil {
		return 0, rows, nil
	}
	if filter.PageIndex <= 0 {
		filter.PageIndex = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	const transferAlias = "cmgtl"
	const merchantUserAlias = "merchant_user"
	const targetUserAlias = "target_user"
	model := g.DB().Model(string(userentity.TbCoinMerchantGoldTransferLog)+" "+transferAlias).Ctx(gctx.New()).
		LeftJoin(string(userentity.TbUserInfo)+" "+merchantUserAlias,
			merchantUserAlias+".id = "+transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogMerchantUserId)).
		LeftJoin(string(userentity.TbUserInfo)+" "+targetUserAlias,
			targetUserAlias+".id = "+transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetUserId))
	if filter.MerchantUserId > 0 {
		model = model.Where(transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogMerchantUserId)+" = ?", filter.MerchantUserId)
	}
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
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogMerchantUserId),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetUserId),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogAmount),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogSenderGoldBefore),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogSenderGoldAfter),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetGoldBefore),
		transferAlias+"."+string(userentity.CoinMerchantGoldTransferLogTargetGoldAfter),
		transferAlias+".created_at",
		merchantUserAlias+"."+string(userentity.UserInfoNickname)+" AS coin_merchant_nickname",
		merchantUserAlias+"."+string(userentity.UserInfoAvatar)+" AS coin_merchant_avatar",
		targetUserAlias+"."+string(userentity.UserInfoNickname)+" AS target_nickname",
		targetUserAlias+"."+string(userentity.UserInfoAvatar)+" AS target_avatar",
	).Order(transferAlias + ".created_at desc, " + transferAlias + ".id desc").
		Limit(filter.PageSize).Offset((filter.PageIndex - 1) * filter.PageSize).
		Scan(&rows)
	return total, rows, err
}
