package statdao

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/apkpromotionchanneldto"
)

func ListDailyApkPromotionChannelStats(req *apkpromotionchanneldto.ApkPromotionChannelDailyStatListReq) (int, []*apkpromotionchanneldto.ApkPromotionChannelDailyStatItem) {
	sql := `select s.stat_date as statDate,
coalesce(nullif(c.channel_name, ''), s.channel_name) as channelName,
s.channel_code as channelCode,
s.register_count as registerCount, s.download_click_count as downloadClickCount,
s.order_count as orderCount, s.pay_user_count as payUserCount,
s.order_total_amount as orderTotalAmount, s.recharge_success_amount as rechargeSuccessAmount
from daily_apk_promotion_channel_stats s
inner join apk_promotion_channels c on c.channel_code = s.channel_code
where s.stat_date >= ? and s.stat_date <= ? `
	param := []any{strings.TrimSpace(req.StartDate), strings.TrimSpace(req.EndDate)}
	channelCode := strings.TrimSpace(req.ChannelCode)
	if channelCode != "" {
		sql += ` and s.channel_code = ? `
		param = append(param, channelCode)
	}
	sql += ` order by s.stat_date desc, s.channel_code asc`
	ctx := gctx.New()
	ret := make([]*apkpromotionchanneldto.ApkPromotionChannelDailyStatItem, 0)
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
