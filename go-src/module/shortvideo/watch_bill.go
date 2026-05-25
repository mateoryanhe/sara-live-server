package shortvideo

import (
	"context"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/shortvideowatchdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/wallet"
)

const watchBillIntervalSeconds uint32 = 3

// WatchBillShortVideo App端短视频观看扣费,每次按5秒进度计费
func WatchBillShortVideo(ctx context.Context, req *shortvideodto.WatchBillShortVideoReq) (*shortvideodto.WatchBillShortVideoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	video := shortvideodao.GetShortVideoById(req.VideoId)
	if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}

	user := userinfodao.GetUserInfoByUserId(userId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	if video.IsPaid != entity.ShortVideoPaidYes || video.DiamondPerSecond == 0 {
		return &shortvideodto.WatchBillShortVideoRes{
			Deducted:          0,
			Diamond:           user.Diamond,
			BilledSeconds:     0,
			ChargeableSeconds: 0,
			CanContinue:       true,
		}, nil
	}

	watch := shortvideowatchdao.GetByUserVideo(userId, req.VideoId)
	billedSeconds := uint32(0)
	if watch != nil {
		billedSeconds = watch.BilledSeconds
	}

	freeWatchSeconds := getAppShortVideoCfgCache().FreeWatchSeconds
	chargeableSeconds := calcChargeableSeconds(billedSeconds, watchBillIntervalSeconds, freeWatchSeconds)
	cost := float64(chargeableSeconds) * float64(video.DiamondPerSecond)

	diamond := user.Diamond
	deducted := float64(0)

	if cost > 0 {
		remaining, err := wallet.DiamondSub(userId, cost, currency.ReasonShortVideoWatch)
		if err != nil {
			return nil, err
		}
		diamond = remaining
		deducted = cost
	}

	newBilledSeconds := billedSeconds + watchBillIntervalSeconds
	if watch == nil {
		watch = entity.NewShortVideoWatch(userId, req.VideoId)
	}
	watch.SetBilledSeconds(newBilledSeconds)
	shortvideowatchdao.SaveToCache(watch)

	return &shortvideodto.WatchBillShortVideoRes{
		Deducted:          deducted,
		Diamond:           diamond,
		BilledSeconds:     newBilledSeconds,
		ChargeableSeconds: chargeableSeconds,
		CanContinue:       true,
	}, nil
}

func calcChargeableSeconds(billedSeconds, interval, freeWatchSeconds uint32) uint32 {
	end := billedSeconds + interval
	if end <= freeWatchSeconds {
		return 0
	}
	if billedSeconds >= freeWatchSeconds {
		return interval
	}
	return end - freeWatchSeconds
}
