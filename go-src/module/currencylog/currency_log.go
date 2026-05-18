package currencylog

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

// Init 订阅货币变动事件,落流水
func Init() {
	event.Sub(gameevent.CurrencyChangeEvent, onCurrencyChange)
}

func onCurrencyChange(val any) {
	data, ok := val.(*gameevent.CurrencyChangeEventData)
	if !ok || data == nil {
		g.Log().Errorf(gctx.New(), "CurrencyChangeEvent payload type error: %T", val)
		return
	}
	_ = entity.NewCurrencyLog(
		data.UserId,
		data.Type,
		data.Action,
		data.Amount,
		data.Before,
		data.After,
		data.Reason,
	)
	g.Log().Debugf(context.Background(),
		"currency log userId=%d type=%d action=%d amount=%v before=%v after=%v reason=%d(%s)",
		data.UserId, data.Type, data.Action, data.Amount, data.Before, data.After, uint8(data.Reason), data.Reason.String(),
	)
}

// GetByUserId 查询用户的货币流水(分页)
func GetByUserId(ctx context.Context, req *userinfodto.GetCurrencyLogReq) (*userinfodto.GetCurrencyLogRes, error) {
	total, list := currencylogdao.GetByUserId(req.UserId, req.PageIndex, req.PageSize)
	items := make([]*userinfodto.CurrencyLogItem, 0, len(list))
	for _, v := range list {
		items = append(items, &userinfodto.CurrencyLogItem{
			Id:       v.ID,
			UserId:   v.UserId,
			Type:     v.Type,
			Action:   v.Action,
			Amount:   v.Amount,
			Before:   v.Before,
			After:    v.After,
			Reason:   v.Reason,
			CreateAt: v.CreatedAt.Unix(),
		})
	}
	return &userinfodto.GetCurrencyLogRes{
		Total: total,
		List:  items,
	}, nil
}
