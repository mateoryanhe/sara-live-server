package liveroomdao

import "testing"

func TestFirstReachDailyLiveThreshold(t *testing.T) {
	needSec := float64(60)
	tests := []struct {
		name      string
		beforeSec float64
		afterSec  float64
		want      bool
	}{
		{name: "first reach", beforeSec: 0, afterSec: 120, want: true},
		{name: "cross exactly", beforeSec: 30, afterSec: 60, want: true},
		{name: "already met before session", beforeSec: 120, afterSec: 240, want: false},
		{name: "still below", beforeSec: 0, afterSec: 59, want: false},
		{name: "late config already above", beforeSec: 5310, afterSec: 5850, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstReachDailyLiveThreshold(tt.beforeSec, tt.afterSec, needSec); got != tt.want {
				t.Fatalf("firstReachDailyLiveThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}
