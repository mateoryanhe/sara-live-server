package liveroomdao

import "testing"

func TestCrossedDailyAccumulatedLiveThreshold(t *testing.T) {
	tests := []struct {
		name      string
		beforeSec float64
		afterSec  float64
		needSec   float64
		want      bool
	}{
		{name: "disabled threshold", needSec: 0, beforeSec: 0, afterSec: 3600, want: false},
		{name: "already met", beforeSec: 7200, afterSec: 9000, needSec: 3600, want: false},
		{name: "still below", beforeSec: 1000, afterSec: 2000, needSec: 3600, want: false},
		{name: "cross on this session", beforeSec: 3000, afterSec: 3600, needSec: 3600, want: true},
		{name: "exactly at threshold after", beforeSec: 0, afterSec: 3600, needSec: 3600, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crossedDailyAccumulatedLiveThreshold(tt.beforeSec, tt.afterSec, tt.needSec); got != tt.want {
				t.Fatalf("crossedDailyAccumulatedLiveThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}
