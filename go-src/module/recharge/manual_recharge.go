package recharge

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargecfgdao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// ManualRecharge 后台手动给玩家充值到账(创建已完成订单 + 立即发放金币)
// 支持两种方式:
//  1. CfgId>0: 从配置取 Price/Gold/Currency(配置必须存在,允许下架状态以便补单)
//  2. CfgId=0: 使用入参 Price 与 Gold(均必须>0)
func ManualRecharge(ctx context.Context, req *rechargeorderdto.CMSManualRechargeReq) (*rechargeorderdto.CMSManualRechargeRes, error) {
	operatorId := httpserver.GetAuthId(ctx)

	var (
		price      uint64
		cur        string
		goldAmount float64
	)
	if req.CfgId > 0 {
		cfg := rechargecfgdao.GetById(req.CfgId)
		if cfg == nil {
			return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
		}
		price = cfg.Price
		cur = cfg.Currency
		goldAmount = float64(cfg.Gold + cfg.ExtraGold)
	} else {
		price = req.Price
		cur = req.Currency
		goldAmount = req.Gold
	}
	if cur == "" {
		cur = defaultCurrency
	}
	if price == 0 {
		return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
	}
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}

	order := entity.NewRechargeOrder(req.UserId, req.CfgId, price, cur, goldAmount, entity.RechargeOrderSourceManual)
	order.SetPayChannel("manual")
	order.SetOperatorId(operatorId)
	if req.Remark != "" {
		order.SetRemark(req.Remark)
	}

	after, err := completeOrder(order, currency.ReasonGmAdjust)
	if err != nil {
		// 发放失败,标记订单为已取消,避免遗留"待支付"脏数据
		order.SetStatus(entity.RechargeOrderStatusCancelled)
		order.SetUpdatedAt(time.Now())
		return nil, err
	}

	return &rechargeorderdto.CMSManualRechargeRes{
		OrderId: strconv.FormatUint(order.ID, 10),
		Gold:    goldAmount,
		After:   after,
		Success: true,
	}, nil
}
