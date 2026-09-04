package recharge

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargeorderdto"
	fiatentity "xr-game-server/entity/fiat"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/fiatcurrency"
	"xr-game-server/module/fxrate"
)

const cmsYhPayTestPackageName = "cms.yhpay.test"

// CreateChannelRechargeOrder App渠道充值建单(yhpay IDR手动入款,无需鉴权,userId由App上报)
func CreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if userinfodao.GetUserInfoByUserId(userId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	packageName := strings.TrimSpace(httpserver.GetPackageNameFromContext(ctx))
	return createChannelRechargeOrder(ctx, userId, packageName, req.CfgId, req.CurrencyCode, false)
}

// CMSCreateChannelRechargeOrder CMS第三方充值测试建单并返回支付URL
func CMSCreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.CMSCreateChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if userinfodao.GetUserInfoByUserId(userId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	packageName := strings.TrimSpace(req.PackageName)
	if packageName == "" {
		packageName = cmsYhPayTestPackageName
	}
	return createChannelRechargeOrder(ctx, userId, packageName, req.CfgId, req.CurrencyCode, true)
}

func createChannelRechargeOrder(ctx context.Context, userId uint64, packageName string, cfgId uint64, currencyCode string, forcePayUrl bool) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	packageName = strings.TrimSpace(packageName)
	if userId == 0 || cfgId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	currencyCode = strings.ToUpper(strings.TrimSpace(currencyCode))
	if currencyCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	fiatCfg := cfgdao.GetFiatCurrencyCfgByCode(currencyCode)
	if fiatCfg == nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyNonExist)
	}
	if fiatCfg.Status != fiatentity.FiatCurrencyStatusEnabled {
		return nil, errercode.CreateCode(errercode.FiatCurrencyDisabled)
	}

	rechargeCfg := GetRechargeCfgFromCacheById(cfgId)
	if rechargeCfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if rechargeCfg.Status != entity.RechargeCfgStatusOnShelf {
		return nil, errercode.CreateCode(errercode.RechargeCfgOffShelf)
	}
	if rechargeCfg.Price <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
	}
	goldAmount := float64(rechargeCfg.Gold)
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}

	if currencyCode != yhPaySupportedCurrency {
		return nil, errercode.CreateCode(errercode.YhPayCurrencyNotSupported)
	}
	if fiatCfg.CurrencyType == fiatentity.FiatCurrencyTypeCrypto {
		return nil, errercode.CreateCode(errercode.YhPayCurrencyNotSupported)
	}

	rate, err := fiatcurrency.ResolveEnabledUsdToQuoteRate(ctx, currencyCode)
	if err != nil {
		return nil, err
	}
	payAmount := fxrate.ConvertUsdToQuote(rechargeCfg.Price, rate)
	if payAmount <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
	}

	order := entity.NewRechargeOrder(userId, rechargeCfg.ID, rechargeCfg.Price, currencyCode, goldAmount, entity.RechargeOrderSourceApp)
	order.SetPayAmount(payAmount)
	order.SetPayChannel(entity.RechargeCfgTypeChannel)
	order.SetPackageName(packageName)
	rechargeorderdao.AddOrderToCache(order)
	ScheduleRechargeOrderTimeout(order.ID, order.CreatedAt)
	event.Pub(gameevent.RechargeOrderCreatedEvent, gameevent.NewRechargeOrderCreatedEventData(order.ID))

	orderIdStr := strconv.FormatUint(order.ID, 10)
	if !forcePayUrl && userinfodao.IsRechargeWhitelist(userId) {
		if _, completeErr := CompleteOrder(order.ID); completeErr == nil {
			if completed := rechargeorderdao.GetById(order.ID); completed != nil {
				order = completed
			}
		}
		return &rechargeorderdto.AppCreateChannelRechargeOrderRes{
			OrderId:   orderIdStr,
			PayUrl:    "",
			Price:     order.Price,
			PayAmount: order.PayAmount,
			Currency:  order.Currency,
			Status:    order.Status,
		}, nil
	}

	if !cfgdao.YhPayEnabled() {
		return nil, errercode.CreateCode(errercode.YhPayNotConfigured)
	}
	cfg, err := getYhPayActiveCfg()
	if err != nil {
		return nil, errercode.CreateCode(errercode.YhPayNotConfigured)
	}

	playerName := strconv.FormatUint(userId, 10)
	if user := userinfodao.GetUserInfoByUserId(userId); user != nil && strings.TrimSpace(user.Nickname) != "" {
		playerName = strings.TrimSpace(user.Nickname)
	}
	playerIP := "127.0.0.1"
	if r := g.RequestFromCtx(ctx); r != nil {
		if ip := strings.TrimSpace(r.GetClientIp()); ip != "" {
			playerIP = ip
		}
	}

	payUrl, err := createYhPayManualDeposit(ctx, cfg, orderIdStr, playerName, playerIP, currencyCode, payAmount)
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay create pay url failed orderId=%s userId=%d currency=%s payAmount=%v err=%v",
			orderIdStr, userId, currencyCode, payAmount, err)
		return nil, errercode.CreateCode(errercode.YhPayCreateFailed)
	}
	xrlog.DetailLog.Infof(ctx, "yhpay create pay url ok orderId=%s userId=%d currency=%s priceUsd=%v payAmount=%v",
		orderIdStr, userId, currencyCode, order.Price, payAmount)

	return &rechargeorderdto.AppCreateChannelRechargeOrderRes{
		OrderId:   orderIdStr,
		PayUrl:    payUrl,
		Price:     order.Price,
		PayAmount: order.PayAmount,
		Currency:  order.Currency,
		Status:    order.Status,
	}, nil
}
