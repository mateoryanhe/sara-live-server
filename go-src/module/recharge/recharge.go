package recharge

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargecfgdao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/gold"
)

// defaultCurrency 默认结算货币
const defaultCurrency = "USD"

// toItem 将 entity 转换为对外 DTO
func toItem(o *entity.RechargeOrder) *rechargeorderdto.RechargeOrderItem {
	item := &rechargeorderdto.RechargeOrderItem{
		ID:           strconv.FormatUint(o.ID, 10),
		UserId:       o.UserId,
		CfgId:        o.CfgId,
		Price:        o.Price,
		Currency:     o.Currency,
		Gold:         o.Gold,
		Status:       o.Status,
		Source:       o.Source,
		PayChannel:   o.PayChannel,
		ThirdOrderId: o.ThirdOrderId,
		Remark:       o.Remark,
		OperatorId:   o.OperatorId,
		CreatedAt:    o.CreatedAt.Unix(),
	}
	if !o.PaidAt.IsZero() {
		item.PaidAt = o.PaidAt.Unix()
	}
	return item
}

// completeOrder 内部统一的"订单完成 → 发放金币"逻辑
// 幂等:已经是已完成状态的订单不会重复发放
// 返回(发放后玩家金币余额, 错误)
func completeOrder(o *entity.RechargeOrder, reason currency.Reason) (float64, error) {
	if o.Status == entity.RechargeOrderStatusCompleted {
		return 0, errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if o.Gold <= 0 {
		return 0, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}
	after, err := gold.Add(o.UserId, o.Gold, reason)
	if err != nil {
		return 0, err
	}
	o.SetStatus(entity.RechargeOrderStatusCompleted)
	o.SetPaidAt(time.Now())
	o.SetUpdatedAt(time.Now())
	return after, nil
}

// CompleteOrder 对外:支付回调成功时调用此函数,完成订单并发放金币
// (本次需求不开发回调路由,该函数保留以便后续接入第三方支付回调)
func CompleteOrder(orderId uint64) (float64, error) {
	o := rechargeorderdao.GetById(orderId)
	if o == nil {
		return 0, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	return completeOrder(o, currency.ReasonRecharge)
}

// ===== App =====

// CreateOrder App 端创建充值订单(按 cfgId 从配置取价格与金币)
// 订单初始状态=待支付;不立即发金币,等支付回调或后台手动完成
func CreateOrder(ctx context.Context, req *rechargeorderdto.AppCreateRechargeOrderReq) (*rechargeorderdto.AppCreateRechargeOrderRes, error) {
	userId := httpserver.GetAuthId(ctx)
	cfg := rechargecfgdao.GetById(req.CfgId)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.RechargeCfgStatusOnShelf {
		return nil, errercode.CreateCode(errercode.RechargeCfgOffShelf)
	}
	if cfg.Price == 0 {
		return nil, errercode.CreateCode(errercode.RechargeAmountInvalid)
	}
	goldAmount := float64(cfg.Diamond + cfg.ExtraDiamond)
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}
	cur := cfg.Currency
	if cur == "" {
		cur = defaultCurrency
	}

	order := entity.NewRechargeOrder(userId, cfg.ID, cfg.Price, cur, goldAmount, entity.RechargeOrderSourceApp)
	if req.PayChannel != "" {
		order.SetPayChannel(req.PayChannel)
	}

	return &rechargeorderdto.AppCreateRechargeOrderRes{
		OrderId:  strconv.FormatUint(order.ID, 10),
		Price:    order.Price,
		Currency: order.Currency,
		Status:   order.Status,
	}, nil
}

// resolveStatusFilter 将外部 StatusFilter (0=全部, 1=待支付, 2=已完成, 3=已取消)
// 转换为 DAO 使用的 statusVal (<0=不过滤,>=0=实际状态枚举值)
func resolveStatusFilter(f int) int {
	switch f {
	case 1:
		return int(entity.RechargeOrderStatusPending)
	case 2:
		return int(entity.RechargeOrderStatusCompleted)
	case 3:
		return int(entity.RechargeOrderStatusCancelled)
	default:
		return -1
	}
}

// GetMyOrderList App 端查询本人充值订单分页列表
// statusFilter: 0=全部(默认), 1=待支付, 2=已完成, 3=已取消
func GetMyOrderList(ctx context.Context, req *rechargeorderdto.AppMyRechargeOrderListReq) (*rechargeorderdto.AppMyRechargeOrderListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	total, rows := rechargeorderdao.ListByUserId(userId, resolveStatusFilter(req.StatusFilter), req.PageIndex, req.PageSize)
	list := make([]*rechargeorderdto.RechargeOrderItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, toItem(r))
	}
	return &rechargeorderdto.AppMyRechargeOrderListRes{Total: total, List: list}, nil
}

// ===== CMS =====

// GetCMSList 后台分页查询充值订单
// statusFilter: 0=全部, 1=待支付, 2=已完成, 3=已取消;source: 0=全部, 1=App, 2=后台手动
func GetCMSList(_ context.Context, req *rechargeorderdto.CMSRechargeOrderListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := rechargeorderdao.CMSList(&rechargeorderdao.CMSListFilter{
		UserId:    req.UserId,
		OrderId:   req.OrderId,
		StatusVal: resolveStatusFilter(req.StatusFilter),
		Source:    req.Source,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*rechargeorderdto.RechargeOrderItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, toItem(r))
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

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
		goldAmount = float64(cfg.Diamond + cfg.ExtraDiamond)
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
