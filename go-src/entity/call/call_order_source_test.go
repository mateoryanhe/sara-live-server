package entity

import "testing"

func TestNormalizeCallOrderSource(t *testing.T) {
	tests := []struct {
		name string
		in   uint8
		want uint8
	}{
		{name: "omitted defaults to live room", in: 0, want: CallOrderSourceLiveRoom},
		{name: "live room", in: CallOrderSourceLiveRoom, want: CallOrderSourceLiveRoom},
		{name: "private message", in: CallOrderSourcePrivateMessage, want: CallOrderSourcePrivateMessage},
		{name: "one-to-one room", in: CallOrderSourceOneToOneRoom, want: CallOrderSourceOneToOneRoom},
		{name: "unknown defaults to live room", in: 99, want: CallOrderSourceLiveRoom},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeCallOrderSource(tt.in); got != tt.want {
				t.Fatalf("NormalizeCallOrderSource(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
