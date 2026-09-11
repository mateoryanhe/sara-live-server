package activity

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/currency"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

func onInviteRechargeGoldArrived(val any) {
	data, ok := val.(*gameevent.RechargeGoldArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		g.Log().Errorf(gctx.New(), "invite RechargeGoldArrivedEvent payload type error: %T", val)
		return
	}
	tryGrantInviteRechargeReward(data)
}

// tryGrantInviteRechargeReward 被邀请人充值金币到账后,按配置给邀请人返还金币。
// 条件:活动开启、邀请关系存在、被邀请人已完成账号首充、在首充完成后 ValidDays 内。
// 排除:白名单、币商。含人工确认。返还基数为本次最终到账金币 CreditedGold。
func tryGrantInviteRechargeReward(data *gameevent.RechargeGoldArrivedEventData) {
	if data == nil || data.Order == nil {
		return
	}
	order := data.Order
	if order.UserId == 0 {
		return
	}
	if shouldSkipInviteRewardKind(data.Kind, order) {
		return
	}

	creditedGold := data.CreditedGold
	if creditedGold <= 0 {
		creditedGold = order.Gold
	}
	if creditedGold <= 0 {
		return
	}

	snap := getInviteRewardCfgCache()
	if snap == nil || !snap.Enabled || snap.RewardPercent <= 0 || snap.ValidDays <= 0 {
		return
	}

	inviteeExt := userinfodao.GetUserExtByUserId(order.UserId)
	if inviteeExt == nil || inviteeExt.FirstRecharge {
		// 尚未完成账号首充,不触发邀请返还
		return
	}

	inviterId := inviteeExt.InviterId
	if inviterId == 0 {
		info := userinfodao.GetUserInfoByUserId(order.UserId)
		if info != nil {
			inviterId = info.InviterId
		}
		if inviterId > 0 && inviteeExt.InviterId == 0 {
			inviteeExt.SetInviterId(inviterId)
			userinfodao.PublishUserExt(inviteeExt)
		}
	}
	if inviterId == 0 || inviterId == order.UserId {
		return
	}

	anchor := inviteeExt.FirstRechargeAt
	if anchor == nil || anchor.IsZero() {
		info := userinfodao.GetUserInfoByUserId(order.UserId)
		if info == nil || info.CreatedAt.IsZero() {
			return
		}
		t := info.CreatedAt
		anchor = &t
	}
	deadline := anchor.Add(time.Duration(snap.ValidDays) * 24 * time.Hour)
	if time.Now().After(deadline) {
		return
	}

	bonus := creditedGold * snap.RewardPercent / 100
	if bonus <= 0 {
		return
	}
	if _, err := wallet.GoldAdd(inviterId, bonus, currency.ReasonInviteRechargeReward); err != nil {
		g.Log().Errorf(gctx.New(), "invite recharge reward GoldAdd failed inviter=%d invitee=%d order=%d bonus=%v err=%v",
			inviterId, order.UserId, order.ID, bonus, err)
	}
}

func shouldSkipInviteRewardKind(kind gameevent.UsdIncomeKind, order *entity.RechargeOrder) bool {
	switch kind {
	case gameevent.UsdIncomeKindWhitelist, gameevent.UsdIncomeKindCoinMerchant:
		return true
	}
	if order == nil {
		return true
	}
	if order.PayChannel == entity.RechargeCfgTypeCoinMerchant {
		return true
	}
	if userinfodao.IsRechargeWhitelist(order.UserId) {
		return true
	}
	return false
}
