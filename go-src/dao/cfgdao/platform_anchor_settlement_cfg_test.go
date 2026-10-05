package cfgdao

import (
	"testing"

	"xr-game-server/entity/user"
)

func TestResolvePlatformAnchorMinimumSettlementUsd(t *testing.T) {
	tests := []struct {
		name string
		cfg  *entity.WalletExchangeCfg
		want float64
	}{
		{name: "missing config", cfg: nil, want: DefaultPlatformAnchorMinimumSettlementUsd},
		{name: "zero value", cfg: &entity.WalletExchangeCfg{}, want: DefaultPlatformAnchorMinimumSettlementUsd},
		{name: "negative value", cfg: &entity.WalletExchangeCfg{PlatformAnchorMinimumSettlementUsd: -1}, want: DefaultPlatformAnchorMinimumSettlementUsd},
		{name: "too large value", cfg: &entity.WalletExchangeCfg{PlatformAnchorMinimumSettlementUsd: MaxPlatformAnchorMinimumSettlementUsd + 1}, want: DefaultPlatformAnchorMinimumSettlementUsd},
		{name: "configured value", cfg: &entity.WalletExchangeCfg{PlatformAnchorMinimumSettlementUsd: 12.5}, want: 12.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolvePlatformAnchorMinimumSettlementUsd(tt.cfg); got != tt.want {
				t.Fatalf("resolvePlatformAnchorMinimumSettlementUsd() = %v, want %v", got, tt.want)
			}
		})
	}
}
