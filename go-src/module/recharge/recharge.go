package recharge

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
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

// toCMSItem 将 entity 转换为 CMS 列表 DTO(含用户昵称)
func toCMSItem(o *entity.RechargeOrder) *rechargeorderdto.CMSRechargeOrderListItem {
	item := &rechargeorderdto.CMSRechargeOrderListItem{
		ID:           strconv.FormatUint(o.ID, 10),
		UserId:       strconv.FormatUint(o.UserId, 10),
		CfgId:        strconv.FormatUint(o.CfgId, 10),
		Price:        o.Price,
		Currency:     o.Currency,
		Gold:         o.Gold,
		Status:       o.Status,
		Source:       o.Source,
		PayChannel:   o.PayChannel,
		ThirdOrderId: o.ThirdOrderId,
		Remark:       o.Remark,
		OperatorId:   strconv.FormatUint(o.OperatorId, 10),
		CreatedAt:    o.CreatedAt.Unix(),
	}
	if o.OperatorId == 0 {
		item.OperatorId = ""
	}
	if o.CfgId == 0 {
		item.CfgId = ""
	}
	if !o.PaidAt.IsZero() {
		item.PaidAt = o.PaidAt.Unix()
	}
	if u := userinfodao.GetUserInfoByUserId(o.UserId); u != nil {
		item.Nickname = u.Nickname
	}
	return item
}

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// ===== App =====

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
		UserId:    parseUint64Filter(req.UserId),
		OrderId:   parseUint64Filter(req.OrderId),
		StatusVal: resolveStatusFilter(req.StatusFilter),
		Source:    req.Source,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*rechargeorderdto.CMSRechargeOrderListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, toCMSItem(r))
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}
