package liveroomdao

import (
	"reflect"
	"testing"
	"time"

	"xr-game-server/core/migrate"
	entity "xr-game-server/entity/live"
)

func TestFilterHistoricalTransferredAnchorSettlements(t *testing.T) {
	cutoff := time.Date(2026, time.September, 21, 0, 0, 0, 0, time.Local)
	rows := []*entity.AnchorIncomeSettlementLog{
		{OneModel: migrate.OneModel{ID: 1}, Status: entity.AnchorIncomeSettlementStatusTransferred, CreatedAt: cutoff.Add(-time.Second)},
		{OneModel: migrate.OneModel{ID: 2}, Status: entity.AnchorIncomeSettlementStatusPending, CreatedAt: cutoff.Add(-time.Hour)},
		{OneModel: migrate.OneModel{ID: 3}, Status: entity.AnchorIncomeSettlementStatusTransferred, CreatedAt: cutoff},
		{OneModel: migrate.OneModel{ID: 4}, Status: entity.AnchorIncomeSettlementStatusTransferred, CreatedAt: cutoff.Add(time.Hour)},
	}

	got := filterHistoricalTransferredAnchorSettlements(rows, cutoff)
	gotIDs := make([]uint64, 0, len(got))
	for _, row := range got {
		gotIDs = append(gotIDs, row.ID)
	}
	wantIDs := []uint64{2, 3, 4}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("filtered IDs = %v, want %v", gotIDs, wantIDs)
	}
}
