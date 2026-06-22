package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/entity"
)

// AudienceContributionStatRow 观众贡献聚合结果
type AudienceContributionStatRow struct {
	SenderId    uint64  `json:"sender_id"`
	TotalAmount float64 `json:"total_amount"`
}

const contributionRankTopLimit = 500

// SumAudienceContributionByRoom 统计指定直播间在时间范围内、且当前在线观众的贡献(礼物+付费弹幕)
func SumAudienceContributionByRoom(roomId uint64, startTime, endTime time.Time, onlineUserIds []uint64) []*AudienceContributionStatRow {
	list := make([]*AudienceContributionStatRow, 0)
	if roomId == 0 || endTime.Before(startTime) || len(onlineUserIds) == 0 {
		return list
	}
	ctx := gctx.New()
	now := time.Now()
	err := g.DB().Ctx(ctx).Model(string(entity.TbLiveRevenueLog)+" rl").
		InnerJoin(string(entity.TbAccount)+" a", "a.id = rl.sender_id").
		Fields("rl.sender_id, SUM(rl.total_amount) AS total_amount").
		Where("rl.room_id", roomId).
		WhereIn("rl.sender_id", onlineUserIds).
		WhereIn("rl.revenue_type", []uint8{uint8(liverevenueconst.Gift), uint8(liverevenueconst.PaidDanmaku)}).
		WhereGTE("rl.created_at", startTime).
		WhereLTE("rl.created_at", endTime).
		Where("IFNULL(a."+string(entity.AccountCancel)+", 0)", 0).
		Wheref(`(
    IFNULL(a.`+string(entity.AccountBan)+`, 0) = 0
    OR (a.`+string(entity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(entity.AccountBanApplyTime)+` <= ?)
  )`, now).
		Group("rl.sender_id").
		Having("total_amount > 0").
		OrderDesc("total_amount").
		Limit(contributionRankTopLimit).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumAudienceContributionByRoom roomId=%d error: %v", roomId, err)
		return make([]*AudienceContributionStatRow, 0)
	}
	return list
}
