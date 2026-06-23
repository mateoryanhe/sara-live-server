package shortvideo

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/wallet"
)

const watchBillIntervalSeconds = 1

// WatchBillShortVideo App端短视频观看扣费,每次按1秒进度计费
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
		lockName := fmt.Sprintf("watch_shortvideo_%v", req.VideoId)
		gmlock.Lock(lockName)
		defer gmlock.Unlock(lockName)
		stat := shortvideodao.GetStatByVideoId(watch.VideoId)
		if stat != nil {
			stat.AddViewCount(1)
		}
	} else {
		watch.SetUpdatedAt(time.Now())
	}

	watch.AddWatchSeconds(watchBillIntervalSeconds)

	if video.IsPaid != entity.ShortVideoPaidYes {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, true), nil
	}

	if isWithinShortVideoFreeWatch(watch.WatchSeconds) {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, true), nil
	}

	if isShortVideoBillingCompleted(video.Duration, watch.BilledSeconds) {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, true), nil
	}

	cost := video.DiamondPerMinute / watchBillIntervalSeconds

	diamond := user.Diamond
	deducted := float64(0)

	if cost > 0 {
		remaining, err := wallet.DiamondSub(userId, cost, currency.ReasonShortVideoWatch)
		if err != nil {
			return nil, err
		}
		diamond = remaining
		deducted = cost
		if stat := shortvideodao.GetStatByVideoId(req.VideoId); stat != nil {
			stat.AddTotalDiamondIncome(cost)
		}
	}

	watch.AddBilledSeconds(watchBillIntervalSeconds)

	return buildWatchBillShortVideoRes(diamond, deducted, watch, true), nil
}

func isShortVideoBillingCompleted(duration uint32, billedSeconds uint64) bool {
	if duration == 0 {
		return false
	}
	return billedSeconds >= uint64(duration)
}

func isWithinShortVideoFreeWatch(watchSeconds uint64) bool {
	freeWatchSeconds := getShortVideoFreeWatchSeconds()
	if freeWatchSeconds == 0 {
		return false
	}
	return watchSeconds <= uint64(freeWatchSeconds)
}

func buildWatchBillShortVideoRes(diamond, deducted float64, watch *entity.ShortVideoWatch, canContinue bool) *shortvideodto.WatchBillShortVideoRes {
	return &shortvideodto.WatchBillShortVideoRes{
		Deducted:          deducted,
		Diamond:           diamond,
		BilledSeconds:     uint32(watch.BilledSeconds),
		ChargeableSeconds: 0,
		CanContinue:       canContinue,
	}
}
