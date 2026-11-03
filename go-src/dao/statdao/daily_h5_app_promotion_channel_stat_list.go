package statdao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/promotionchannelurldto"
)

func ListDailyH5AppPromotionChannelStats(req *promotionchannelurldto.H5AppPromotionDailyStatListReq) (int, []*promotionchannelurldto.H5AppPromotionDailyStatItem) {
	sql := `select s.stat_date as statDate,
coalesce(nullif(u.channel_name, ''), s.channel_name) as channelName,
nullif(trim(u.h5_app_channel_code), '') as h5AppChannelCode,
s.register_count as registerCount,
s.order_count as orderCount, s.pay_user_count as payUserCount,
s.order_total_amount as orderTotalAmount, s.recharge_success_amount as rechargeSuccessAmount
from daily_h5_app_promotion_channel_stats s
inner join promotion_channel_urls u on u.channel_code = s.channel_code
where s.stat_date >= ? and s.stat_date <= ? `
	param := []any{strings.TrimSpace(req.StartDate), strings.TrimSpace(req.EndDate)}
	h5App := strings.TrimSpace(req.H5AppChannelCode)
	if h5App != "" {
		sql += ` and u.h5_app_channel_code = ? `
		param = append(param, h5App)
	} else if req.ChannelCode > 0 {
		sql += ` and s.channel_code = ? `
		param = append(param, req.ChannelCode)
	}
	sql += ` order by s.stat_date desc, u.h5_app_channel_code asc, s.channel_code asc`
	ctx := gctx.New()
	ret := make([]*promotionchannelurldto.H5AppPromotionDailyStatItem, 0)
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
