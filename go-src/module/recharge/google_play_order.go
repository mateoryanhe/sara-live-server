package recharge

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"google.golang.org/api/androidpublisher/v3"
	"xr-game-server/constants/currency"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity/recharge"
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
	gmlock.Lock(rechargeThirdOrderLockKey(googleOrderId))
	defer gmlock.Unlock(rechargeThirdOrderLockKey(googleOrderId))

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
	if getActiveGooglePlayCfg() == nil {
		logGooglePlayInfo(ctx, "rtdn skipped: google play cfg missing")
		return nil
	}
	return processGooglePlayPurchase(ctx, packageName, sku, purchaseToken, 0, 0)
}

// ScheduleGooglePlayPurchaseVerify App 上报 purchaseToken 后异步验单并完成充值
func ScheduleGooglePlayPurchaseVerify(ctx context.Context, userId, orderId uint64, packageName, productId, purchaseToken string) {
	xrpool.AddWithRecover(ctx, func(poolCtx context.Context) {
		dbCtx := gctx.New()
		if err := verifyGooglePlayPurchaseByAppReport(dbCtx, userId, orderId, packageName, productId, purchaseToken); err != nil {
			logGooglePlayError(dbCtx, "app async verify failed orderId=%d err=%v", orderId, err)
		}
	})
}

func verifyGooglePlayPurchaseByAppReport(ctx context.Context, userId, orderId uint64, packageName, sku, purchaseToken string) error {
	if !googlePlayEnabled() {
		logGooglePlayInfo(ctx, "app verify skipped: google play disabled orderId=%d", orderId)
		return nil
	}
	if getActiveGooglePlayCfg() == nil {
		logGooglePlayInfo(ctx, "app verify skipped: google play cfg missing orderId=%d", orderId)
		return nil
	}
	order := rechargeorderdao.GetById(orderId)
	if order == nil || order.UserId != userId {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status != entity.RechargeOrderStatusPending {
		return nil
	}
	return processGooglePlayPurchase(ctx, packageName, sku, purchaseToken, orderId, userId)
}

func processGooglePlayPurchase(ctx context.Context, packageName, sku, purchaseToken string, fixedOrderId, fixedUserId uint64) error {
	packageName = strings.TrimSpace(packageName)
	sku = strings.TrimSpace(sku)
	purchaseToken = strings.TrimSpace(purchaseToken)
	if packageName == "" || sku == "" || purchaseToken == "" {
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.InvalidParam)
		}
		logGooglePlayError(ctx, "google play purchase missing params sku=%s", sku)
		return nil
	}

	rechargeCfg := GetGoogleRechargeCfgByProductId(packageName, sku)
	if rechargeCfg == nil {
		logGooglePlayError(ctx, "google play unknown sku=%s pkg=%s", sku, packageName)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.RechargeCfgNonExist)
		}
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
		logGooglePlayInfo(ctx, "google play purchase not completed sku=%s state=%d", sku, state)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.InvalidParam)
		}
		return nil
	}

	googleOrderId := strings.TrimSpace(purchase.OrderId)
	if googleOrderId == "" {
		logGooglePlayError(ctx, "google play missing google orderId sku=%s", sku)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.InvalidParam)
		}
		return nil
	}
	if existing := rechargeorderdao.GetByThirdOrderId(googleOrderId); existing != nil {
		if existing.Status == entity.RechargeOrderStatusCompleted {
			logGooglePlayInfo(ctx, "google play duplicate googleOrderId=%s orderId=%d", googleOrderId, existing.ID)
			return maybeConsumeGooglePurchase(ctx, packageName, sku, purchaseToken, purchase)
		}
	}

	var orderId uint64
	if fixedOrderId > 0 {
		orderId = fixedOrderId
		if obfuscated := strings.TrimSpace(purchase.ObfuscatedExternalAccountId); obfuscated != "" {
			resolved, resolveErr := resolveGooglePlayOrderId(obfuscated)
			if resolveErr != nil || resolved != orderId {
				logGooglePlayError(ctx, "google play obfuscated account mismatch orderId=%d obfuscated=%s", orderId, obfuscated)
				return errercode.CreateCode(errercode.InvalidParam)
			}
		}
	} else {
		orderId, err = resolveGooglePlayOrderId(purchase.ObfuscatedExternalAccountId)
		if err != nil {
			logGooglePlayError(ctx, "google play resolve orderId failed sku=%s googleOrderId=%s err=%v", sku, googleOrderId, err)
			return nil
		}
	}

	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		logGooglePlayError(ctx, "google play order not found orderId=%d sku=%s", orderId, sku)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.RechargeOrderNonExist)
		}
		return nil
	}
	if fixedUserId > 0 && order.UserId != fixedUserId {
		logGooglePlayError(ctx, "google play order user mismatch orderId=%d userId=%d", orderId, fixedUserId)
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if orderPackage := strings.TrimSpace(order.PackageName); orderPackage != "" && orderPackage != packageName {
		logGooglePlayError(ctx, "google play order package mismatch orderId=%d orderPkg=%s pkg=%s", orderId, orderPackage, packageName)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.InvalidParam)
		}
		return nil
	}
	packageName = resolveGooglePlayPackageName(order, packageName)
	if order.Status != entity.RechargeOrderStatusPending {
		if order.Status == entity.RechargeOrderStatusCompleted && order.ThirdOrderId == googleOrderId {
			return maybeConsumeGooglePurchase(ctx, packageName, sku, purchaseToken, purchase)
		}
		logGooglePlayInfo(ctx, "google play order state invalid orderId=%d status=%d", orderId, order.Status)
		return nil
	}
	if order.CfgId > 0 && order.CfgId != rechargeCfg.ID {
		logGooglePlayError(ctx, "google play cfg mismatch orderId=%d orderCfg=%d skuCfg=%d", orderId, order.CfgId, rechargeCfg.ID)
		if fixedOrderId > 0 {
			return errercode.CreateCode(errercode.InvalidParam)
		}
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

func resolveGooglePlayPackageName(order *entity.RechargeOrder, rtdnPackageName string) string {
	rtdnPackageName = strings.TrimSpace(rtdnPackageName)
	if order != nil {
		if orderPackage := strings.TrimSpace(order.PackageName); orderPackage != "" {
			return orderPackage
		}
	}
	return rtdnPackageName
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
