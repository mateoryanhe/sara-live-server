package cfgdao

import (
	"testing"

	"xr-game-server/entity/live"
)

func TestResolveEffectiveLiveMinSessionMinutes(t *testing.T) {
	tests := []struct {
		name string
		cfg  *entity.EffectiveLiveCfg
		want int
	}{
		{name: "missing config", cfg: nil, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "zero value", cfg: &entity.EffectiveLiveCfg{}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "negative value", cfg: &entity.EffectiveLiveCfg{MinSessionMinutes: -1}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "too large value", cfg: &entity.EffectiveLiveCfg{MinSessionMinutes: MaxEffectiveLiveMinSessionMinutes + 1}, want: DefaultEffectiveLiveMinSessionMinutes},
		{name: "configured value", cfg: &entity.EffectiveLiveCfg{MinSessionMinutes: 45}, want: 45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveEffectiveLiveMinSessionMinutes(tt.cfg); got != tt.want {
				t.Fatalf("resolveEffectiveLiveMinSessionMinutes() = %d, want %d", got, tt.want)
			}
		})
	}
}
