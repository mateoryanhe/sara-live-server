package recharge

import (
	"context"
	"strings"
	"sync"

	"xr-game-server/errercode"
)

// ChannelPayCreateReq 渠道支付建单请求（业务层 → Provider）
type ChannelPayCreateReq struct {
	OrderID    string
	UserID     uint64
	PlayerName string
	PlayerIP   string
	Email      string
	Region     string // 可选覆盖；空则 Provider 用自身默认(如 HaiPay DefaultRegion)
	Amount     float64
	PriceUsd   float64
	PayChannel uint8
}

// ChannelPayCreateRes Provider 建单结果
type ChannelPayCreateRes struct {
	PayURL       string
	ThirdOrderID string
}

// ChannelPayProvider 可插拔渠道支付实现（HaiPay 等）。Google Play / iOS 不走此接口。
type ChannelPayProvider interface {
	Name() string
	Enabled() bool
	// QuotePay 决定订单落库币种与金额。HaiPay 美金包装固定返回 USD + priceUsd。
	QuotePay(priceUsd float64) (payCurrency string, payAmount float64, err error)
	CreatePay(ctx context.Context, req *ChannelPayCreateReq) (*ChannelPayCreateRes, error)
}

var (
	channelPayMu        sync.RWMutex
	channelPayProviders []ChannelPayProvider
)

// RegisterChannelPayProvider 注册渠道支付实现；后注册的同名实现会覆盖前者。
func RegisterChannelPayProvider(p ChannelPayProvider) {
	if p == nil {
		return
	}
	name := strings.TrimSpace(p.Name())
	channelPayMu.Lock()
	defer channelPayMu.Unlock()
	if name != "" {
		for i, existing := range channelPayProviders {
			if existing != nil && strings.EqualFold(existing.Name(), name) {
				channelPayProviders[i] = p
				return
			}
		}
	}
	channelPayProviders = append(channelPayProviders, p)
}

func resolveActiveChannelPayProvider() (ChannelPayProvider, error) {
	channelPayMu.RLock()
	defer channelPayMu.RUnlock()
	for _, p := range channelPayProviders {
		if p != nil && p.Enabled() {
			return p, nil
		}
	}
	return nil, errercode.CreateCode(errercode.ChannelPayNotConfigured)
}
