package recharge

import (
	"context"
	"strconv"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// CreateOrder App 端创建充值订单(按 cfgId 从配置取价格与金币)
// 订单初始状态=待支付;不立即发金币,等支付回调或后台手动完成
func CreateOrder(ctx context.Context, req *rechargeorderdto.AppCreateRechargeOrderReq) (*rechargeorderdto.AppCreateRechargeOrderRes, error) {
	userId := httpserver.GetAuthId(ctx)
	var order *entity.RechargeOrder
	//充值金额来自配置
	if req.CfgId > 0 {
		cfg := GetRechargeCfgFromCacheById(req.CfgId)
		if cfg == nil {
			return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
		}
		if cfg.Status != entity.RechargeCfgStatusOnShelf {
			return nil, errercode.CreateCode(errercode.RechargeCfgOffShelf)
		}
		if cfg.Price == 0 {
			return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
		}
		goldAmount := float64(cfg.Gold + cfg.ExtraGold)
		if goldAmount <= 0 {
			return nil, errercode.CreateCode(errercode.RechargeGoldInvalid)
		}
		cur := cfg.Currency
		if cur == "" {
			cur = defaultCurrency
		}
		order = entity.NewRechargeOrder(userId, cfg.ID, cfg.Price, cur, goldAmount, entity.RechargeOrderSourceApp)
		order.SetPayChannel(cfg.CfgType)
	} else {
		//自定义金额
		order = entity.NewRechargeOrder(userId, 0, req.Amount, defaultCurrency, 0, entity.RechargeOrderSourceApp)
		order.SetPayChannel(req.PayChannel)
	}

	rechargeorderdao.AddOrderToCache(order)

	event.Pub(gameevent.RechargeOrderCreatedEvent, gameevent.NewRechargeOrderCreatedEventData(order.ID))

	return &rechargeorderdto.AppCreateRechargeOrderRes{
		OrderId:  strconv.FormatUint(order.ID, 10),
		Price:    order.Price,
		Currency: order.Currency,
		Status:   order.Status,
	}, nil
}
