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

	watch := shortvideodao.GetOneShortVideoWatch(userId, req.VideoId)

	//记录视频观看人数
	if watch.ViewCounted == entity.ShortVideoWatchViewCountedNo {
		watch.SetViewCounted(entity.ShortVideoWatchViewCountedYes)
		stat := shortvideodao.GetStatByVideoId(watch.VideoId)
		if stat != nil {
			stat.AddViewCount(1)
		}
	}

	if video.IsPaid != entity.ShortVideoPaidYes {
		return &shortvideodto.WatchBillShortVideoRes{
			Deducted:          0,
			Diamond:           user.Diamond,
			BilledSeconds:     0,
			ChargeableSeconds: 0,
			CanContinue:       true,
		}, nil
	}

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

	return &shortvideodto.WatchBillShortVideoRes{
		Deducted:          deducted,
		Diamond:           diamond,
		BilledSeconds:     newBilledSeconds,
		ChargeableSeconds: 0,
		CanContinue:       true,
	}, nil
}

// ensureShortVideoWatch 按用户+视频维度获取观看记录,未计入观看人数时累加统计
func ensureShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	watch := shortvideodao.GetOneShortVideoWatch(userId, videoId)
	if watch != nil {
		tryIncrViewCount(watch)
		return watch
	}

	lockName := fmt.Sprintf("watch_first_%d_%d", userId, videoId)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	tryIncrViewCount(watch)
	return watch
}

func tryIncrViewCount(watch *entity.ShortVideoWatch) {
	if watch == nil || watch.ViewCounted == entity.ShortVideoWatchViewCountedYes {
		return
	}

	lockName := fmt.Sprintf("watch_view_count_%s", watch.ID)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if watch.ViewCounted == entity.ShortVideoWatchViewCountedYes {
		return
	}

	statLock := fmt.Sprintf("view_shortvideo_%v", watch.VideoId)
	gmlock.Lock(statLock)
	stat := shortvideodao.GetStatByVideoId(watch.VideoId)
	if stat != nil {
		stat.AddViewCount(1)
	}
	gmlock.Unlock(statLock)

	watch.SetViewCounted(entity.ShortVideoWatchViewCountedYes)
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
