package xrtime

import (
	"testing"
	"time"
)

func TestNextWeekStart(t *testing.T) {
	location := time.Local
	tests := []struct {
		name string
		at   time.Time
		want time.Time
	}{
		{
			name: "monday starts next monday",
			at:   time.Date(2026, time.September, 21, 0, 0, 0, 0, location),
			want: time.Date(2026, time.September, 28, 0, 0, 0, 0, location),
		},
		{
			name: "sunday expires the following day",
			at:   time.Date(2026, time.September, 27, 23, 59, 59, 0, location),
			want: time.Date(2026, time.September, 28, 0, 0, 0, 0, location),
		},
		{
			name: "crosses year boundary",
			at:   time.Date(2026, time.December, 31, 12, 0, 0, 0, location),
			want: time.Date(2027, time.January, 4, 0, 0, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextWeekStart(tt.at); !got.Equal(tt.want) {
				t.Fatalf("NextWeekStart(%v) = %v, want %v", tt.at, got, tt.want)
			}
		})
	}
}

func TestWeekStart(t *testing.T) {
	location := time.Local
	at := time.Date(2026, time.September, 23, 16, 30, 45, 0, location)
	want := time.Date(2026, time.September, 21, 0, 0, 0, 0, location)
	if got := WeekStart(at); !got.Equal(want) {
		t.Fatalf("WeekStart(%v) = %v, want %v", at, got, want)
	}
}
