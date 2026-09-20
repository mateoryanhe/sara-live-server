package liveroomdao

import (
	"testing"

	"xr-game-server/core/migrate"
	"xr-game-server/entity/live"
)

func TestMergeGuildIncomeSettlementLogsFromCache(t *testing.T) {
	previous := guildIncomeSettlementLogCacheMgr
	initGuildIncomeSettlementLogDao()
	t.Cleanup(func() {
		guildIncomeSettlementLogCacheMgr = previous
	})

	databaseRow := &entity.GuildIncomeSettlementLog{
		OneModel: migrate.OneModel{ID: 25238},
		Status:   entity.GuildIncomeSettlementStatusPending,
	}
	cachedRow := &entity.GuildIncomeSettlementLog{
		OneModel:        migrate.OneModel{ID: databaseRow.ID},
		Status:          entity.GuildIncomeSettlementStatusApproved,
		TransferOrderId: "25238",
	}
	PublishGuildIncomeSettlementLog(cachedRow)

	rows := mergeGuildIncomeSettlementLogsFromCache([]*entity.GuildIncomeSettlementLog{databaseRow})
	if len(rows) != 1 {
		t.Fatalf("unexpected row count: %d", len(rows))
	}
	if rows[0] != cachedRow {
		t.Fatal("list row did not prefer the cached settlement")
	}
	if rows[0].Status != entity.GuildIncomeSettlementStatusApproved {
		t.Fatalf("unexpected cached status: %d", rows[0].Status)
	}
}

func TestMergeGuildIncomeSettlementLogsKeepsDatabaseRowOnCacheMiss(t *testing.T) {
	previous := guildIncomeSettlementLogCacheMgr
	initGuildIncomeSettlementLogDao()
	t.Cleanup(func() {
		guildIncomeSettlementLogCacheMgr = previous
	})

	databaseRow := &entity.GuildIncomeSettlementLog{OneModel: migrate.OneModel{ID: 25239}}
	rows := mergeGuildIncomeSettlementLogsFromCache([]*entity.GuildIncomeSettlementLog{databaseRow})
	if len(rows) != 1 || rows[0] != databaseRow {
		t.Fatal("cache miss must preserve the database row")
	}
}
