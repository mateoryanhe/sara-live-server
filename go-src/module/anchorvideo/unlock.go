package anchorvideo

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/common"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorvideodao"
	"xr-game-server/dao/anchorvideopurchasedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorvideodto"
	entity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/wallet"
)

// UnlockVideo App 付费解锁主播相册视频;免费/已付费/主播本人不扣费;成功后刷新 userId+anchorId 购买缓存.
func UnlockVideo(ctx context.Context, req *anchorvideodto.UnlockVideoReq) (*anchorvideodto.UnlockVideoRes, error) {
	if req == nil || req.VideoId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	video := anchorvideodao.GetById(req.VideoId)
	if video == nil || video.ID == 0 || video.IsDeleted {
		return nil, errercode.CreateCode(errercode.AnchorVideoNonExist)
	}

	// 免费或主播本人
	if video.Price <= common.Zero || video.AnchorId == userId {
		return &anchorvideodto.UnlockVideoRes{
			Deducted: 0,
			Diamond:  user.Diamond,
			Paid:     true,
		}, nil
	}

	lockName := fmt.Sprintf("anchor_video_unlock_%d_%d", userId, req.VideoId)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if anchorvideopurchasedao.HasPurchased(userId, video.AnchorId, req.VideoId) {
		user = userinfodao.GetUserInfoByUserId(userId)
		diamond := float64(0)
		if user != nil {
			diamond = user.Diamond
		}
		return &anchorvideodto.UnlockVideoRes{
			Deducted: 0,
			Diamond:  diamond,
			Paid:     true,
		}, nil
	}

	price := video.Price
	diamond, err := wallet.DiamondSubWithGoldExchange(userId, price, currency.ReasonAnchorVideoUnlock)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	row := entity.NewAnchorVideoPurchase(userId, video.ID, video.AnchorId, price, now)
	anchorvideopurchasedao.SavePurchase(row)
	RecordAnchorVideoPayIncome(video.AnchorId, price)

	return &anchorvideodto.UnlockVideoRes{
		Deducted: price,
		Diamond:  diamond,
		Paid:     true,
	}, nil
}
