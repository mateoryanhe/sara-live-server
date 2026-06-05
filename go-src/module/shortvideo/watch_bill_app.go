package shortvideo

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/wallet"
)

const watchBillIntervalSeconds uint32 = 1

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

	if video.IsPaid != entity.ShortVideoPaidYes {
		ensureShortVideoWatch(userId, req.VideoId)
		return &shortvideodto.WatchBillShortVideoRes{
			Deducted:          0,
			Diamond:           user.Diamond,
			BilledSeconds:     0,
			ChargeableSeconds: 0,
			CanContinue:       true,
		}, nil
	}

	watch := ensureShortVideoWatch(userId, req.VideoId)
	billedSeconds := watch.BilledSeconds

	cost := float64(video.DiamondPerSecond)

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
	watch.SetBilledSeconds(newBilledSeconds)
	//shortvideodao.SaveToCache(watch)

	return &shortvideodto.WatchBillShortVideoRes{
		Deducted:          deducted,
		Diamond:           diamond,
		BilledSeconds:     newBilledSeconds,
		ChargeableSeconds: 0,
		CanContinue:       true,
	}, nil
}

// ensureShortVideoWatch 按用户+视频维度获取观看记录,首次观看时创建记录并累加观看人数
func ensureShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	watch := shortvideodao.GetShortVideoWatchByUserVideo(userId, videoId)
	if watch != nil {
		return watch
	}

	lockName := fmt.Sprintf("watch_first_%d_%d", userId, videoId)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	watch = shortvideodao.GetShortVideoWatchByUserVideo(userId, videoId)
	if watch != nil {
		return watch
	}

	watch = entity.NewShortVideoWatch(userId, videoId)

	statLock := fmt.Sprintf("view_shortvideo_%v", videoId)
	gmlock.Lock(statLock)
	defer gmlock.Unlock(statLock)
	stat := shortvideodao.GetStatByVideoId(videoId)
	if stat != nil {
		stat.AddViewCount(1)
	}

	return watch
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
