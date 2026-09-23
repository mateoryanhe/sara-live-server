package stat

import (
	"context"
	"time"
	"xr-game-server/core/push"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity/stat"
)

// GetCMSSysStat CMS获取系统总数据
func GetCMSSysStat(_ context.Context, _ *statdto.CMSSysStatReq) (*statdto.CMSSysStatRes, error) {
	stat := statdao.GetSysStat()
	if stat == nil {
		stat = entity.NewSystemTotalStat(entity.SystemTotalStatDefaultID)
	}
	now := time.Now()
	todayStat := statdao.GetDailyLoginStatByDate(entity.FormatDailyLoginStatDate(now))
	weeklyStat := statdao.GetWeeklyLoginStatByWeek(entity.FormatWeeklyLoginStatKey(now))
	return &statdto.CMSSysStatRes{
		TotalGold:                 stat.TotalGold,
		TotalGoldConsume:          stat.TotalGoldConsume,
		TotalDiamondConsume:       stat.TotalDiamondConsume,
		TotalRecharge:             stat.TotalRecharge,
		TotalNormalUserRecharge:   stat.TotalNormalUserRecharge,
		TotalCoinMerchantRecharge: stat.TotalCoinMerchantRecharge,
		TotalAnchorPayout:         stat.TotalAnchorPayout,
		WeeklyAnchorPayout:        weeklyStat.AnchorPayoutAmount,
		TotalVirtualRecharge:      stat.TotalVirtualRecharge,
		TotalWithdraw:             stat.TotalWithdraw,
		TotalRegisterUser:         stat.TotalRegisterUser,
		TodayRecharge:             todayStat.RechargeAmount,
		TodayGoldConsume:          todayStat.GoldConsumeAmount,
		TodayDiamondConsume:       todayStat.DiamondConsumeAmount,
		TodayRegisterUser:         todayStat.RegisterCount,
		OnlineCount:               uint64(push.OnlineCount()),
	}, nil
}
