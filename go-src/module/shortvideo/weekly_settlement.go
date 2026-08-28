package shortvideo

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/shortvideo"
	"xr-game-server/gameevent"
	"xr-game-server/module/liverevenuesharecfg"
	"xr-game-server/module/wallet"
)

func initWeeklySettlement() {
	event.Sub(gameevent.WeekEvent, onWeekShortVideoAuthorSettlement)
}

func onWeekShortVideoAuthorSettlement(_ any) {
	settleNonAnchorAuthorShortVideoIncome()
}

// settleNonAnchorAuthorShortVideoIncome 周一0点:结算非主播作者短视频未结算收益到钻石钱包
func settleNonAnchorAuthorShortVideoIncome() {
	userIds := userinfodao.CollectUserExtWithShortVideoUnsettledIncomeUserIds()
	ctx := gctx.New()
	anchorSharePercent := liverevenuesharecfg.ResolveAnchorSharePercent()
	g.Log().Infof(ctx, "short video author weekly settlement start, users=%d, anchorSharePercent=%.2f", len(userIds), anchorSharePercent)

	var settledCount int
	for _, userId := range userIds {
		if settleOneNonAnchorAuthorShortVideoIncome(userId, anchorSharePercent) {
			settledCount++
		}
	}
	g.Log().Infof(ctx, "short video author weekly settlement done, settled=%d", settledCount)
}

func settleOneNonAnchorAuthorShortVideoIncome(userId uint64, anchorSharePercent float64) bool {
	if userId == 0 {
		return false
	}
	author := userinfodao.GetUserInfoByUserId(userId)
	if author == nil {
		return false
	}
	if author.IsAnchor() {
		g.Log().Warningf(gctx.New(), "skip short video author settlement for anchor userId=%d", userId)
		return false
	}

	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil {
		return false
	}
	unsettledIncome := ext.ClearShortVideoUnsettledIncome()
	if unsettledIncome <= 0 {
		return false
	}
	userinfodao.PublishUserExt(ext)

	settlementDiamond := liverevenuesharecfg.CalcSettlementShareAmount(0, unsettledIncome)
	if settlementDiamond > 0 {
		if _, err := wallet.DiamondAdd(userId, settlementDiamond, currency.ReasonShortVideoAuthorSettlement); err != nil {
			g.Log().Errorf(gctx.New(), "short video author settlement diamond add failed userId=%d income=%.4f diamond=%.4f err=%v",
				userId, unsettledIncome, settlementDiamond, err)
			ext.AddShortVideoUnsettledIncome(unsettledIncome)
			userinfodao.PublishUserExt(ext)
			return false
		}
	}

	_ = entity.NewShortVideoAuthorSettlementLog(userId, unsettledIncome, settlementDiamond, anchorSharePercent)
	return true
}
