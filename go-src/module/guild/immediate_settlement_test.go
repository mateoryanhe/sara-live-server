package guild

import (
	"reflect"
	"testing"
)

func TestUniqueGuildIds(t *testing.T) {
	got := uniqueGuildIds([]uint64{3, 0, 2, 3, 2, 1})
	want := []uint64{3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueGuildIds() = %v, want %v", got, want)
	}
}

func TestFilterImmediateSettlementGuildIds(t *testing.T) {
	allowed, denied := filterImmediateSettlementGuildIds([]uint64{1, 2, 3}, []uint64{1, 3}, true)
	if !reflect.DeepEqual(allowed, []uint64{1, 3}) || !reflect.DeepEqual(denied, []uint64{2}) {
		t.Fatalf("filter result allowed=%v denied=%v", allowed, denied)
	}
}
