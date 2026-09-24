package liveroom

import (
	"reflect"
	"testing"
)

func TestUniquePlatformAnchorIds(t *testing.T) {
	got := uniquePlatformAnchorIds([]uint64{3, 0, 2, 3, 2, 1})
	want := []uint64{3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("uniquePlatformAnchorIds() = %v, want %v", got, want)
	}
}

func TestFilterImmediateSettlementPlatformAnchorIds(t *testing.T) {
	allowed, denied := filterImmediateSettlementPlatformAnchorIds([]uint64{1, 2, 3}, []uint64{1, 3}, true)
	if !reflect.DeepEqual(allowed, []uint64{1, 3}) || !reflect.DeepEqual(denied, []uint64{2}) {
		t.Fatalf("filter result allowed=%v denied=%v", allowed, denied)
	}
}
