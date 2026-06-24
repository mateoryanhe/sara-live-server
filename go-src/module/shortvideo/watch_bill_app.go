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

const (
	watchBillCallIntervalSeconds uint64 = 1  // App 每秒调用一次
	watchBillTickSeconds         uint64 = 10 // 每累计 10 秒可计费观看时长扣一次费
)

// WatchBillShortVideo App端短视频观看扣费,App 每秒调用一次,每 10 秒扣一次费
func WatchBillShortVideo(ctx context.Context, req *shortvideodto.WatchBillShortVideoReq) (*shortvideodto.WatchBillShortVideoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	video := shortvideodao.GetShortVideoById(req.VideoId)
	if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
		return nil, nil
	}

	watch := shortvideodao.GetOneShortVideoWatch(userId, req.VideoId)

	if video.IsPaid != entity.ShortVideoPaidYes {
		return buildWatchBillShortVideoRes(0, 0, watch, 0, true), nil
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	//只扣费一次
	if isShortVideoBillingCompleted(video.Duration-video.FreeWatchSeconds, watch.BilledSeconds) {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, 0, true), nil
	}
	//判断是否进入扣费环节
	if uint64(video.FreeWatchSeconds) > watch.WatchSeconds {
		watch.AddWatchSeconds(watchBillCallIntervalSeconds)
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, 0, false), nil
	}
	//观看视频叠加
	watch.AddWatchSeconds(watchBillCallIntervalSeconds)

	billSeconds := (watch.WatchSeconds - uint64(video.FreeWatchSeconds)) / watchBillTickSeconds

	if 0 >= billSeconds {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, 0, true), nil
	}

	cost := calcShortVideoWatchCost(video.DiamondPerMinute, billSeconds*watchBillTickSeconds)
	if cost <= 0 {
		return buildWatchBillShortVideoRes(user.Diamond, 0, watch, 0, true), nil
	}

	remaining, err := wallet.DiamondSub(userId, cost, currency.ReasonShortVideoWatch)

	if err != nil {
		return nil, err
	}
	watch.SubWatchSeconds(billSeconds * watchBillTickSeconds)
	if stat := shortvideodao.GetStatByVideoId(req.VideoId); stat != nil {
		stat.AddTotalDiamondIncome(cost)
	}
	watch.AddBilledSeconds(billSeconds * watchBillTickSeconds)

	return buildWatchBillShortVideoRes(remaining, cost, watch, uint32(billSeconds), true), nil
}

// WatchShortVideoStart App端开始观看短视频
func WatchShortVideoStart(ctx context.Context, req *shortvideodto.WatchShortVideoStartReq) (*shortvideodto.WatchShortVideoStartRes, error) {
	video := shortvideodao.GetShortVideoById(req.VideoId)
	if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
		return nil, nil
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
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
		//watch.SetUpdatedAt(time.Now())
	}
	watch.ResetWatchSeconds()
	if video.IsPaid != entity.ShortVideoPaidYes {
		return &shortvideodto.WatchShortVideoStartRes{}, nil
	}
	return &shortvideodto.WatchShortVideoStartRes{}, nil
}

// WatchShortVideoEnd App端结束观看短视频
func WatchShortVideoEnd(_ context.Context, _ *shortvideodto.WatchShortVideoEndReq) (*shortvideodto.WatchShortVideoEndRes, error) {
	return &shortvideodto.WatchShortVideoEndRes{}, nil
}

func isShortVideoBillingCompleted(duration uint32, billedSeconds uint64) bool {
	if duration == 0 {
		return true
	}
	return billedSeconds >= uint64(duration)
}

func calcChargeableWatchSeconds(watchSeconds uint64, freeWatchSeconds uint32) uint64 {
	if freeWatchSeconds == 0 {
		return watchSeconds
	}
	if watchSeconds <= uint64(freeWatchSeconds) {
		return 0
	}
	return watchSeconds - uint64(freeWatchSeconds)
}

// calcPendingShortVideoBillSeconds 计算本次是否达到 10 秒扣费点,返回本次应计费秒数(0 表示暂不扣费)
func calcPendingShortVideoBillSeconds(watchSeconds, billedSeconds uint64, freeWatchSeconds uint32, videoDuration uint32) uint64 {
	chargeableWatchSeconds := calcChargeableWatchSeconds(watchSeconds, freeWatchSeconds)
	if chargeableWatchSeconds <= billedSeconds {
		return 0
	}
	pendingSeconds := chargeableWatchSeconds - billedSeconds
	if pendingSeconds < watchBillTickSeconds {
		return 0
	}

	billSeconds := watchBillTickSeconds
	if videoDuration > 0 {
		remaining := uint64(videoDuration) - billedSeconds
		if remaining == 0 {
			return 0
		}
		if remaining < billSeconds {
			billSeconds = remaining
		}
	}
	return billSeconds
}

func buildWatchBillShortVideoRes(diamond, deducted float64, watch *entity.ShortVideoWatch, chargeableSeconds uint32, canContinue bool) *shortvideodto.WatchBillShortVideoRes {
	return &shortvideodto.WatchBillShortVideoRes{
		Deducted:          deducted,
		Diamond:           diamond,
		BilledSeconds:     uint32(watch.BilledSeconds),
		ChargeableSeconds: chargeableSeconds,
		CanContinue:       canContinue,
	}
}

func calcShortVideoWatchCost(diamondPerMinute float64, billSeconds uint64) float64 {
	if diamondPerMinute <= 0 || billSeconds == 0 {
		return 0
	}
	return diamondPerMinute * float64(billSeconds) / 60
}
