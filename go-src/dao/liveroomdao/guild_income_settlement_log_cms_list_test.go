package liveroomdao

import (
	"reflect"
	"testing"
	"time"

	"xr-game-server/core/migrate"
	live "xr-game-server/entity/live"
)

func TestFilterHistoricalTransferredGuildSettlements(t *testing.T) {
	cutoff := time.Date(2026, time.September, 21, 0, 0, 0, 0, time.Local)
	rows := []*live.GuildIncomeSettlementLog{
		{OneModel: migrate.OneModel{ID: 1}, CreatedAt: cutoff.Add(-time.Second), Status: live.GuildIncomeSettlementStatusTransferred},
		{OneModel: migrate.OneModel{ID: 2}, CreatedAt: cutoff.Add(-time.Hour), Status: live.GuildIncomeSettlementStatusPending},
		{OneModel: migrate.OneModel{ID: 3}, CreatedAt: cutoff, Status: live.GuildIncomeSettlementStatusTransferred},
		{OneModel: migrate.OneModel{ID: 4}, CreatedAt: cutoff.Add(time.Hour), Status: live.GuildIncomeSettlementStatusTransferred},
	}

	got := filterHistoricalTransferredGuildSettlements(rows, cutoff)
	gotIds := make([]uint64, 0, len(got))
	for _, row := range got {
		gotIds = append(gotIds, row.ID)
	}
	wantIds := []uint64{2, 3, 4}
	if !reflect.DeepEqual(gotIds, wantIds) {
		t.Fatalf("filtered IDs = %v, want %v", gotIds, wantIds)
	}
}
