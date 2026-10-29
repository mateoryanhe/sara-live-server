package statdao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/promotionchannelurldto"
)

// ListDailyPromotionChannelStats 按日期范围分页查询推广渠道日统计
func ListDailyPromotionChannelStats(req *promotionchannelurldto.PromotionChannelDailyStatListReq) (int, []*promotionchannelurldto.PromotionChannelDailyStatItem) {
	sql := `select stat_date as statDate, channel_name as channelName, channel_code as channelCode,
register_count as registerCount, order_count as orderCount, pay_user_count as payUserCount,
order_total_amount as orderTotalAmount, recharge_success_amount as rechargeSuccessAmount,
h5_download_apk_count as h5DownloadApkCount
from daily_promotion_channel_stats where stat_date >= ? and stat_date <= ? `
	param := []any{strings.TrimSpace(req.StartDate), strings.TrimSpace(req.EndDate)}
	if req.ChannelCode > 0 {
		sql += ` and channel_code = ? `
		param = append(param, req.ChannelCode)
	}
	sql += ` order by stat_date desc, channel_code asc`
	ctx := gctx.New()
	ret := make([]*promotionchannelurldto.PromotionChannelDailyStatItem, 0)
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
