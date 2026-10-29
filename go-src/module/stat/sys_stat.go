package stat

import (
	"context"
	"time"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity/stat"
)

// GetCMSSysStat CMS获取系统总数据(默认 5 分钟内存缓存, 过期再组装/查库)
func GetCMSSysStat(_ context.Context, req *statdto.CMSSysStatReq) (*statdto.CMSSysStatRes, error) {
	bypass := req != nil && req.BypassDashboardCache
	res, err := getCachedDashboardSysStat(bypass, buildCMSSysStat)
	if err != nil {
		return nil, err
	}
	overlayCMSSysStatLiveCounts(res)
	return res, nil
}

func buildCMSSysStat() (*statdto.CMSSysStatRes, error) {
	statRow := statdao.ScanSystemTotalStatForCMS()
	now := time.Now()
	today := entity.FormatDailyLoginStatDate(now)
	week := entity.FormatWeeklyLoginStatKey(now)
	todayStat := statdao.ScanDailyLoginStatByDate(today)
	weeklyStat := statdao.ScanWeeklyLoginStatByWeek(week)
	return assembleCMSSysStatRes(statRow, todayStat, weeklyStat, now), nil
}

func overlayCMSSysStatLiveCounts(res *statdto.CMSSysStatRes) {
	if res == nil {
		return
	}
	now := time.Now()
	onlineTotal := uint64(push.OnlineCount())
	liveAnchors := liveroomdao.CountLivingAnchors()
	gameOnline, gameBetUsers := SnapshotGameTrendCountsForMetric(now)
	audienceOnline := uint64(0)
	if onlineTotal > liveAnchors {
		audienceOnline = onlineTotal - liveAnchors
	}
	res.OnlineCount = onlineTotal
	res.GameOnlineCount = gameOnline
	res.GameBetUserCount = gameBetUsers
	res.AudienceOnlineCount = audienceOnline
	res.LiveAnchorCount = liveAnchors
}

func assembleCMSSysStatRes(
	statRow *entity.SystemTotalStat,
	todayStat *entity.DailyLoginStat,
	weeklyStat *entity.WeeklyLoginStat,
	now time.Time,
) *statdto.CMSSysStatRes {
	if statRow == nil {
		statRow = entity.NewSystemTotalStat(entity.SystemTotalStatDefaultID)
	}
	if todayStat == nil {
		todayStat = entity.NewDailyLoginStat(entity.FormatDailyLoginStatDate(now))
	}
	if weeklyStat == nil {
		weeklyStat = entity.NewWeeklyLoginStat(entity.FormatWeeklyLoginStatKey(now))
	}
	res := &statdto.CMSSysStatRes{
		TotalGold:                 statRow.TotalGold,
		TotalGoldConsume:          statRow.TotalGoldConsume,
		TotalDiamondConsume:       statRow.TotalDiamondConsume,
		TotalRecharge:             statRow.TotalRecharge,
		TotalNormalUserRecharge:   statRow.TotalNormalUserRecharge,
		TotalCoinMerchantRecharge: statRow.TotalCoinMerchantRecharge,
		TotalAnchorPayout:         statRow.TotalAnchorPayout,
		WeeklyAnchorPayout:        weeklyStat.AnchorPayoutAmount,
		TotalVirtualRecharge:      statRow.TotalVirtualRecharge,
		TotalWithdraw:             statRow.TotalWithdraw,
		TotalRegisterUser:         statRow.TotalRegisterUser,
		TodayRecharge:             todayStat.RechargeAmount,
		TodayGoldConsume:          todayStat.GoldConsumeAmount,
		TodayDiamondConsume:       todayStat.DiamondConsumeAmount,
		TodayRegisterUser:         todayStat.RegisterCount,
	}
	return res
}
