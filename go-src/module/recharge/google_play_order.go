package recharge

import (
	"context"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/androidpublisher/v3"
	"xr-game-server/constants/currency"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const (
	googlePurchaseStatePurchased  = int64(0)
	googleOneTimeProductPurchased = 1
)

// CompleteGooglePlayOrder Google Play 验单成功后完成充值订单并发放金币(幂等)
func CompleteGooglePlayOrder(orderId uint64, googleOrderId string) error {
	googleOrderId = strings.TrimSpace(googleOrderId)
	if orderId == 0 || googleOrderId == "" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := rechargeorderdao.GetByThirdOrderId(googleOrderId); existing != nil {
		if existing.ID == orderId && existing.Status == entity.RechargeOrderStatusCompleted {
			return nil
		}
		if existing.ID != orderId {
			return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
		}
	}
	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		if order.ThirdOrderId == "" || order.ThirdOrderId == googleOrderId {
			return nil
		}
		return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	order.SetThirdOrderId(googleOrderId)
	order.SetPayChannel(entity.RechargeCfgTypeGoogle)
	_, err := completeOrder(order, currency.ReasonRecharge)
	return err
}

func handleGoogleOneTimeProductPurchased(ctx context.Context, packageName, sku, purchaseToken string) error {
	if !googlePlayEnabled() {
		logGooglePlayInfo(ctx, "rtdn skipped: google play disabled")
		return nil
	}
	gpCfg := getActiveGooglePlayCfg()
	if gpCfg == nil {
		logGooglePlayInfo(ctx, "rtdn skipped: google play cfg missing")
		return nil
	}
	if packageName != "" && packageName != gpCfg.PackageName {
		logGooglePlayError(ctx, "rtdn package mismatch, got=%s want=%s", packageName, gpCfg.PackageName)
		return nil
	}
	packageName = gpCfg.PackageName

	rechargeCfg := GetGoogleRechargeCfgByProductId(packageName, sku)
	if rechargeCfg == nil {
		logGooglePlayError(ctx, "rtdn unknown sku=%s", sku)
		return nil
	}

	purchase, err := getGoogleProductPurchase(ctx, packageName, sku, purchaseToken)
	if err != nil {
		logGooglePlayError(ctx, "verify purchase failed sku=%s err=%v", sku, err)
		return err
	}
	if purchase == nil || purchase.PurchaseState != googlePurchaseStatePurchased {
		state := int64(-1)
		if purchase != nil {
			state = purchase.PurchaseState
		}
		logGooglePlayInfo(ctx, "rtdn purchase not completed sku=%s state=%d", sku, state)
		return nil
	}

	googleOrderId := strings.TrimSpace(purchase.OrderId)
	if googleOrderId == "" {
		logGooglePlayError(ctx, "rtdn missing google orderId sku=%s", sku)
		return nil
	}
	if existing := rechargeorderdao.GetByThirdOrderId(googleOrderId); existing != nil {
		if existing.Status == entity.RechargeOrderStatusCompleted {
			logGooglePlayInfo(ctx, "rtdn duplicate googleOrderId=%s orderId=%d", googleOrderId, existing.ID)
			return maybeConsumeGooglePurchase(ctx, packageName, sku, purchaseToken, purchase)
		}
	}

	orderId, err := resolveGooglePlayOrderId(purchase.ObfuscatedExternalAccountId)
	if err != nil {
		logGooglePlayError(ctx, "rtdn resolve orderId failed sku=%s googleOrderId=%s err=%v", sku, googleOrderId, err)
		return nil
	}
	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		logGooglePlayError(ctx, "rtdn order not found orderId=%d sku=%s", orderId, sku)
		return nil
	}
	if order.Status != entity.RechargeOrderStatusPending {
		if order.Status == entity.RechargeOrderStatusCompleted && order.ThirdOrderId == googleOrderId {
			return maybeConsumeGooglePurchase(ctx, packageName, sku, purchaseToken, purchase)
		}
		logGooglePlayInfo(ctx, "rtdn order state invalid orderId=%d status=%d", orderId, order.Status)
		return nil
	}
	if order.CfgId > 0 && order.CfgId != rechargeCfg.ID {
		logGooglePlayError(ctx, "rtdn cfg mismatch orderId=%d orderCfg=%d skuCfg=%d", orderId, order.CfgId, rechargeCfg.ID)
		return nil
	}
	if order.CfgId == 0 && rechargeCfg.ID > 0 {
		order.SetCfgId(rechargeCfg.ID)
	}

	if err := CompleteGooglePlayOrder(orderId, googleOrderId); err != nil {
		logGooglePlayError(ctx, "complete order failed orderId=%d googleOrderId=%s err=%v", orderId, googleOrderId, err)
		return err
	}
	logGooglePlayInfo(ctx, "recharge completed orderId=%d googleOrderId=%s sku=%s userId=%d gold=%.0f",
		orderId, googleOrderId, sku, order.UserId, order.Gold)

	if purchase.PurchaseTimeMillis > 0 {
		paidAt := time.UnixMilli(purchase.PurchaseTimeMillis)
		if refreshed := rechargeorderdao.GetById(orderId); refreshed != nil {
			refreshed.SetPaidAt(paidAt)
			rechargeorderdao.FlushOrderCache(refreshed)
		}
	}

	return maybeConsumeGooglePurchase(ctx, packageName, sku, purchaseToken, purchase)
}

func resolveGooglePlayOrderId(obfuscatedAccountId string) (uint64, error) {
	obfuscatedAccountId = strings.TrimSpace(obfuscatedAccountId)
	if obfuscatedAccountId == "" {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	orderId, err := strconv.ParseUint(obfuscatedAccountId, 10, 64)
	if err != nil || orderId == 0 {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return orderId, nil
}

func maybeConsumeGooglePurchase(ctx context.Context, packageName, sku, purchaseToken string, purchase *androidpublisher.ProductPurchase) error {
	if purchase == nil {
		return nil
	}
	if purchase.ConsumptionState == 1 {
		return nil
	}
	if err := consumeGoogleProductPurchase(ctx, packageName, sku, purchaseToken); err != nil {
		logGooglePlayError(ctx, "consume failed sku=%s err=%v", sku, err)
		return err
	}
	logGooglePlayInfo(ctx, "consume success sku=%s", sku)
	return nil
}
