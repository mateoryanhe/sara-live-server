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
	watch.SetUpdatedAt(time.Now())
	//initFreeTime(userId, req.VideoId)
	//clearFreeTime(userId, req.VideoId)
	////检查是否是付费视频
	//if video.AuthorId != userId && video.IsPaid == entity.ShortVideoPaidYes && watch.PaidTime == nil {
	//	//开始检查是否有免费时间
	//	freeTime := watch.GetFreeTime(video.FreeWatchSeconds)
	//	if freeTime == 0 {
	//		//需要购买才可以观看
	//		return nil, errercode.CreateCode(errercode.ShortVideoMustPayToWatch)
	//	} else {
	//		now := time.Now()
	//		watch.SetFreeStartTime(&now)
	//		watch.SetFreeUsing(true)
	//	}
	//}

	// 累计观看次数(不去重,每次开始观看+1)
	if stat := shortvideodao.GetStatByVideoId(req.VideoId); stat != nil {
		stat.AddWatchCount(1)
	}

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

	if video.IsPaid != entity.ShortVideoPaidYes || video.AuthorId == userId {
		shortvideodao.PrependWatchToWatchListCache(userId, watch)
	}

	return &shortvideodto.WatchShortVideoStartRes{}, nil
}

func initFreeTime(userId uint64, videoId uint64) {
	video := shortvideodao.GetShortVideoById(videoId)
	watch := shortvideodao.GetOneShortVideoWatch(userId, videoId)
	if video.IsPaid == entity.ShortVideoPaidYes {
		if watch.FirstWatchTime == nil {
			now := time.Now()
			watch.SetFreeTime(uint64(video.FreeWatchSeconds))
			watch.SetFirstWatchTime(&now)
		}
	}
}

func clearFreeTime(userId uint64, videoId uint64) {
	video := shortvideodao.GetShortVideoById(videoId)
	watch := shortvideodao.GetOneShortVideoWatch(userId, videoId)
	if video.IsPaid == entity.ShortVideoPaidYes {
		if watch.FreeUsing {
			diff := time.Since(*watch.FreeStartTime)
			watch.SetFreeUsing(false)
			watch.SubFreeTime(uint64(diff.Seconds()))
		}
	}
}

// WatchShortVideoEnd App端结束观看短视频
func WatchShortVideoEnd(ctx context.Context, req *shortvideodto.WatchShortVideoEndReq) (*shortvideodto.WatchShortVideoEndRes, error) {
	userId := httpserver.GetAuthId(ctx)
	clearFreeTime(userId, req.VideoId)
	return &shortvideodto.WatchShortVideoEndRes{}, nil
}

// PayShortVideo App端短视频付费观看,按视频付费价格一次性扣除钻石
func PayShortVideo(ctx context.Context, req *shortvideodto.PayShortVideoReq) (*shortvideodto.PayShortVideoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	video := shortvideodao.GetShortVideoById(req.VideoId)
	if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if video.IsPaid != entity.ShortVideoPaidYes {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	price := video.PayDiamond
	if price <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	watch := shortvideodao.GetOneShortVideoWatch(userId, req.VideoId)
	user := userinfodao.GetUserInfoByUserId(userId)
	if watch.PaidTime != nil {
		return &shortvideodto.PayShortVideoRes{
			Deducted: 0,
			Diamond:  user.Diamond,
		}, nil
	}

	diamond, err := wallet.DiamondSub(userId, price, currency.ReasonShortVideoWatch)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	watch.SetPaidTime(&now)

	lockName := fmt.Sprintf("pay_shortvideo_%v", req.VideoId)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)
	if stat := shortvideodao.GetStatByVideoId(req.VideoId); stat != nil {
		stat.AddTotalDiamondIncome(price)
	}

	shortvideodao.PrependWatchToWatchListCache(userId, watch)

	return &shortvideodto.PayShortVideoRes{
		Deducted: price,
		Diamond:  diamond,
	}, nil
}
