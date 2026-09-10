package recharge

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/channelpaydao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity/recharge"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/coinmerchant"
)

const cmsChannelPayTestPackageName = "cms.channelpay.test"

// CreateChannelRechargeOrder App渠道充值建单(无需鉴权,userId由App上报)
func CreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := requireExistingAppUser(userId); err != nil {
		return nil, err
	}
	packageName := strings.TrimSpace(httpserver.GetPackageNameFromContext(ctx))
	return createChannelRechargeOrder(ctx, userId, packageName, req.CfgId, req.CurrencyCode, req.PayName, req.PayEmail, false)
}

// CMSCreateChannelRechargeOrder CMS第三方充值测试建单并返回支付URL
func CMSCreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.CMSCreateChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := requireExistingAppUser(userId); err != nil {
		return nil, err
	}
	packageName := strings.TrimSpace(req.PackageName)
	if packageName == "" {
		packageName = cmsChannelPayTestPackageName
	}
	return createChannelRechargeOrder(ctx, userId, packageName, req.CfgId, req.CurrencyCode, req.PayName, req.PayEmail, true)
}

func createChannelRechargeOrder(ctx context.Context, userId uint64, packageName string, cfgId uint64, currencyCode, payName, payEmail string, forcePayUrl bool) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	packageName = strings.TrimSpace(packageName)
	if userId == 0 || cfgId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	region, err := resolveChannelPayRegion(currencyCode)
	if err != nil {
		return nil, err
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

	provider, err := resolveActiveChannelPayProvider()
	if err != nil {
		return nil, err
	}
	payCurrency, payAmount, err := provider.QuotePay(rechargeCfg.Price)
	if err != nil || payAmount <= 0 || strings.TrimSpace(payCurrency) == "" {
		xrlog.DetailLog.Errorf(ctx, "channelPay quote failed provider=%s priceUsd=%v err=%v",
			provider.Name(), rechargeCfg.Price, err)
		return nil, errercode.CreateCode(errercode.ChannelPayCreateFailed)
	}
	payCurrency = strings.ToUpper(strings.TrimSpace(payCurrency))

	order := entity.NewRechargeOrder(userId, rechargeCfg.ID, rechargeCfg.Price, payCurrency, goldAmount, entity.RechargeOrderSourceApp)
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

	return createChannelPayURL(ctx, order, userId, orderIdStr, rechargeCfg.Price, payAmount, entity.RechargeCfgTypeChannel, provider, region, payName, payEmail)
}

// CreateCoinMerchantChannelRechargeOrder 币商App渠道建单(需登录;档位来自币商充值配置缓存)
func CreateCoinMerchantChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateCoinMerchantChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	if user.UserType != userentity.UserTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	packageName := strings.TrimSpace(httpserver.GetPackageNameFromContext(ctx))
	return createCoinMerchantChannelRechargeOrder(ctx, userId, packageName, req.CfgId, req.CurrencyCode, req.PayName, req.PayEmail)
}

func createCoinMerchantChannelRechargeOrder(ctx context.Context, userId uint64, packageName string, cfgId uint64, currencyCode, payName, payEmail string) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	packageName = strings.TrimSpace(packageName)
	if userId == 0 || cfgId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	region, err := resolveChannelPayRegion(currencyCode)
	if err != nil {
		return nil, err
	}

	rechargeCfg := coinmerchant.GetCoinMerchantRechargeCfgFromCacheById(cfgId)
	if rechargeCfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if rechargeCfg.Price <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
	}
	goldAmount := float64(rechargeCfg.Gold)
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}

	provider, err := resolveActiveChannelPayProvider()
	if err != nil {
		return nil, err
	}
	payCurrency, payAmount, err := provider.QuotePay(rechargeCfg.Price)
	if err != nil || payAmount <= 0 || strings.TrimSpace(payCurrency) == "" {
		xrlog.DetailLog.Errorf(ctx, "channelPay coinMerchant quote failed provider=%s priceUsd=%v err=%v",
			provider.Name(), rechargeCfg.Price, err)
		return nil, errercode.CreateCode(errercode.ChannelPayCreateFailed)
	}
	payCurrency = strings.ToUpper(strings.TrimSpace(payCurrency))

	order := entity.NewRechargeOrder(userId, rechargeCfg.ID, rechargeCfg.Price, payCurrency, goldAmount, entity.RechargeOrderSourceApp)
	order.SetPayAmount(payAmount)
	order.SetPayChannel(entity.RechargeCfgTypeCoinMerchant)
	order.SetPackageName(packageName)
	rechargeorderdao.AddOrderToCache(order)
	ScheduleRechargeOrderTimeout(order.ID, order.CreatedAt)
	event.Pub(gameevent.RechargeOrderCreatedEvent, gameevent.NewRechargeOrderCreatedEventData(order.ID))

	orderIdStr := strconv.FormatUint(order.ID, 10)
	if userinfodao.IsRechargeWhitelist(userId) {
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

	return createChannelPayURL(ctx, order, userId, orderIdStr, rechargeCfg.Price, payAmount, entity.RechargeCfgTypeCoinMerchant, provider, region, payName, payEmail)
}

func createChannelPayURL(ctx context.Context, order *entity.RechargeOrder, userId uint64, orderIdStr string, priceUsd, payAmount float64, payChannel uint8, provider ChannelPayProvider, region, payName, payEmail string) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	if provider == nil {
		var err error
		provider, err = resolveActiveChannelPayProvider()
		if err != nil {
			return nil, err
		}
	}

	playerName, email := resolveChannelPayPayer(userId, payName, payEmail)
	playerName = normalizePayDisplayName(playerName)
	playerIP := "127.0.0.1"
	if r := g.RequestFromCtx(ctx); r != nil {
		if ip := strings.TrimSpace(r.GetClientIp()); ip != "" {
			playerIP = ip
		}
	}

	payRes, err := provider.CreatePay(ctx, &ChannelPayCreateReq{
		OrderID:    orderIdStr,
		UserID:     userId,
		PlayerName: playerName,
		PlayerIP:   playerIP,
		Email:      email,
		Region:     region,
		Amount:     payAmount,
		PriceUsd:   priceUsd,
		PayChannel: payChannel,
	})
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "channelPay create failed provider=%s orderId=%s userId=%d priceUsd=%v err=%v",
			provider.Name(), orderIdStr, userId, priceUsd, err)
		return nil, errercode.CreateCode(errercode.ChannelPayCreateFailed)
	}
	if payRes == nil || strings.TrimSpace(payRes.PayURL) == "" {
		xrlog.DetailLog.Errorf(ctx, "channelPay empty payUrl provider=%s orderId=%s", provider.Name(), orderIdStr)
		return nil, errercode.CreateCode(errercode.ChannelPayCreateFailed)
	}
	if third := strings.TrimSpace(payRes.ThirdOrderID); third != "" {
		order.SetThirdOrderId(third)
	}
	if isLikelyRealPayName(playerName) || isLikelyRealPayEmail(email) {
		if requireExistingAppUser(userId) == nil {
			channelpaydao.UpsertPayerInfo(userId, playerName, email)
		}
	}

	xrlog.DetailLog.Infof(ctx, "channelPay create ok provider=%s orderId=%s userId=%d settle=%s payAmount=%v priceUsd=%v",
		provider.Name(), orderIdStr, userId, order.Currency, order.PayAmount, priceUsd)

	return &rechargeorderdto.AppCreateChannelRechargeOrderRes{
		OrderId:   orderIdStr,
		PayUrl:    strings.TrimSpace(payRes.PayURL),
		Price:     order.Price,
		PayAmount: order.PayAmount,
		Currency:  order.Currency,
		Status:    order.Status,
	}, nil
}

// resolveChannelPayRegion 将 App 上报的 currencyCode 映射为 HaiPay region；空则留给 Provider 用 CMS 默认。
func resolveChannelPayRegion(currencyCode string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(currencyCode))
	if code == "" {
		return "", nil
	}
	region, err := haiPayResolveRegion(code)
	if err != nil {
		return "", errercode.CreateCode(errercode.ChannelPayCurrencyNotSupported)
	}
	return region, nil
}
