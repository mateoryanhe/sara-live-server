package attrreport

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/gameevent"
)

func initEvents() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
	event.Sub(gameevent.RechargeGoldArrivedEvent, onRechargeGoldArrivedEvent)
}

// onRegisterEvent 有 registerAttr: 入库 + accept; 无上报: 只写负缓存(不上库),后续 Find 不打库.
func onRegisterEvent(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	userId := val.UserId
	if !val.HasRegisterUrlAttr {
		userinfodao.MarkRegisterUrlAttrAbsent(userId)
		return
	}
	typ := val.RegisterUrlAttrType
	ch := val.RegisterUrlAttrCh
	tjUserId := val.RegisterUrlAttrTjUserId
	registeredAt := val.RegisteredAt
	if registeredAt.IsZero() {
		registeredAt = time.Now()
	}
	clientIP := val.Attribution.ClientIPAddress
	clientUA := val.Attribution.ClientUserAgent
	eventSourceURL := val.Attribution.EventSourceURL
	deviceModel := val.DeviceModel
	osVersion := val.OsVersion

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		userinfodao.SaveRegisterUrlAttr(userId, typ, ch, tjUserId)

		ipv4, ipv6 := SplitClientIPs(clientIP)
		PostAccept(ctx, &AcceptPayload{
			Type:            typ,
			Ch:              ch,
			TjUserId:        tjUserId,
			UserId:          strconv.FormatUint(userId, 10),
			ClientIPAddress: clientIP,
			ClientIPv4:      ipv4,
			ClientIPv6:      ipv6,
			EventID:         fmt.Sprintf("%d-%d", userId, registeredAt.Unix()),
			EventName:       EventNameCompleteRegistration,
			EventTime:       registeredAt.Unix(), // Unix 秒 = UTC/GMT
			EventSourceURL:  eventSourceURL,
			Em:              ResolveHashedEmail(userId),
			Ph:              ResolveHashedPhone(userId),
			ActionSource:    DefaultActionSource,
			ClientUserAgent: clientUA,
			PhoneModel:      deviceModel,
			SystemVersion:   osVersion,
			Country:         ResolveCountryForUser(userId, clientIP),
		})
	})
}

// onRechargeGoldArrivedEvent 普通用户账号首充到账后上报 Purchase;
// 必须库中已有有效归因记录(Find 不建空行、空缓存会清掉).
func onRechargeGoldArrivedEvent(data any) {
	val, ok := data.(*gameevent.RechargeGoldArrivedEventData)
	if !ok || val == nil || val.Order == nil || val.Order.UserId == 0 {
		return
	}
	if !val.IsAccountFirst {
		return
	}
	// 排除币商 / 白名单虚拟美金
	if val.Kind == gameevent.UsdIncomeKindCoinMerchant || val.Kind == gameevent.UsdIncomeKindWhitelist {
		return
	}
	order := val.Order
	userId := order.UserId
	orderId := order.ID
	paidAt := order.PaidAt
	if paidAt.IsZero() {
		paidAt = time.Now()
	}
	// FB Purchase: currency+value 有则带上; Price 为配置美金金额, Currency 空且有金额时默认 USD
	currency := strings.TrimSpace(order.Currency)
	value := order.Price
	if currency == "" && value > 0 {
		currency = "USD"
	}

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		attr := userinfodao.FindUserRegisterUrlAttrByUserId(userId)
		if attr == nil || !attr.HasValues() {
			// 无归因记录不上报;勿 New/MustGet 造空行
			return
		}
		deviceModel, osVersion := "", ""
		// 直查库,勿 MustGet 造空登录设备行
		if dev := userlogindevicedao.GetUserLoginDeviceFromDB(userId); dev != nil {
			deviceModel = strings.TrimSpace(dev.DeviceModel)
			osVersion = strings.TrimSpace(dev.OsVersion)
		}
		PostAccept(ctx, &AcceptPayload{
			Type:            attr.Type,
			Ch:              attr.Ch,
			TjUserId:        attr.TjUserId,
			UserId:          strconv.FormatUint(userId, 10),
			EventID:         fmt.Sprintf("%d-%d", orderId, paidAt.Unix()),
			EventName:       EventNamePurchase,
			EventTime:       paidAt.Unix(), // Unix 秒 = UTC/GMT
			Em:              ResolveHashedEmail(userId),
			Ph:              ResolveHashedPhone(userId),
			ActionSource:    DefaultActionSource,
			PhoneModel:      deviceModel,
			SystemVersion:   osVersion,
			Country:         ResolveCountryForUser(userId, ""),
			Currency:        currency,
			Value:           value,
		})
	})
}
