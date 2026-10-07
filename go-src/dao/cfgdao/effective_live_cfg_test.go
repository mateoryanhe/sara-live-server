package cfgdao

import (
	"testing"

	"xr-game-server/entity/user"
)

func TestResolveEffectiveLiveMinSessionMinutes(t *testing.T) {
	tests := []struct {
		name string
		cfg  *entity.WalletExchangeCfg
		want int
	}{
		{name: "missing config", cfg: nil, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "zero value", cfg: &entity.WalletExchangeCfg{}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "negative value", cfg: &entity.WalletExchangeCfg{EffectiveLiveMinSessionMinutes: -1}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "too large value", cfg: &entity.WalletExchangeCfg{EffectiveLiveMinSessionMinutes: MaxEffectiveLiveMinSessionMinutes + 1}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "configured value", cfg: &entity.WalletExchangeCfg{EffectiveLiveMinSessionMinutes: 45}, want: 45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveEffectiveLiveMinSessionMinutes(tt.cfg); got != tt.want {
				t.Fatalf("resolveEffectiveLiveMinSessionMinutes() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestResolveEffectiveLiveDailyAccumulatedMinutes(t *testing.T) {
	tests := []struct {
		name string
		cfg  *entity.WalletExchangeCfg
		want int
	}{
		{name: "missing config", cfg: nil, want: DefaultEffectiveLiveDailyAccumulatedMinutes},
		{name: "zero value", cfg: &entity.WalletExchangeCfg{}, want: 0},
		{name: "negative value", cfg: &entity.WalletExchangeCfg{EffectiveLiveDailyAccumulatedMinutes: -1}, want: DefaultEffectiveLiveDailyAccumulatedMinutes},
		{name: "too large value", cfg: &entity.WalletExchangeCfg{EffectiveLiveDailyAccumulatedMinutes: MaxEffectiveLiveDailyAccumulatedMinutes + 1}, want: DefaultEffectiveLiveDailyAccumulatedMinutes},
		{name: "configured value", cfg: &entity.WalletExchangeCfg{EffectiveLiveDailyAccumulatedMinutes: 120}, want: 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveEffectiveLiveDailyAccumulatedMinutes(tt.cfg); got != tt.want {
				t.Fatalf("resolveEffectiveLiveDailyAccumulatedMinutes() = %d, want %d", got, tt.want)
			}
		})
	}
}
