package stat

import (
	"context"
	"time"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity"
)

// GetCMSSysStat CMS获取系统总数据
func GetCMSSysStat(_ context.Context, _ *statdto.CMSSysStatReq) (*statdto.CMSSysStatRes, error) {
	stat := statdao.GetSysStat()
	if stat == nil {
		stat = entity.NewSystemTotalStat(entity.SystemTotalStatDefaultID)
	}
	todayStat := dailyloginstatdao.GetByDate(entity.FormatDailyLoginStatDate(time.Now()))
	return &statdto.CMSSysStatRes{
		TotalGold:           stat.TotalGold,
		TotalGoldConsume:    stat.TotalGoldConsume,
		TotalDiamondConsume: stat.TotalDiamondConsume,
		TotalRecharge:       stat.TotalRecharge,
		TotalWithdraw:       stat.TotalWithdraw,
		TotalRegisterUser:   stat.TotalRegisterUser,
		TodayRecharge:       todayStat.RechargeAmount,
		TodayGoldConsume:    todayStat.GoldConsumeAmount,
		TodayDiamondConsume: todayStat.DiamondConsumeAmount,
		TodayRegisterUser:   todayStat.RegisterCount,
	}, nil
}
