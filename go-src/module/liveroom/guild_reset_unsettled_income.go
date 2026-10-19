package liveroom

import (
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
)

// BatchGuildResetUnsettledResult 批量清零工会未结算收益结果
type BatchGuildResetUnsettledResult struct {
	ResetCount   int
	NoDataCount  int
	FailGuildIds []uint64
}

type guildResetUnsettledOutcome uint8

const (
	guildResetUnsettledFailed guildResetUnsettledOutcome = iota
	guildResetUnsettledNoData
	guildResetUnsettledReset
)

// BatchResetGuildsUnsettledIncome 批量清零工会及名下主播未结算收益(不写任何业务日志)
func BatchResetGuildsUnsettledIncome(guildIds []uint64) *BatchGuildResetUnsettledResult {
	result := &BatchGuildResetUnsettledResult{FailGuildIds: make([]uint64, 0)}
	for _, guildId := range guildIds {
		switch resetGuildUnsettledIncome(guildId) {
		case guildResetUnsettledReset:
			result.ResetCount++
		case guildResetUnsettledNoData:
			result.NoDataCount++
		default:
			result.FailGuildIds = append(result.FailGuildIds, guildId)
		}
	}
	return result
}

func resetGuildUnsettledIncome(guildId uint64) guildResetUnsettledOutcome {
	if guildId == 0 || guilddao.GetGuildById(guildId) == nil {
		return guildResetUnsettledFailed
	}
	clearedAny := false
	allZero := true

	switch liveroomdao.ClearGuildIncomeUnsettledForCMS(guildId) {
	case liveroomdao.ClearUnsettledIncomeCleared:
		clearedAny = true
		allZero = false
	case liveroomdao.ClearUnsettledIncomeAlreadyZero:
	default:
		return guildResetUnsettledFailed
	}

	for _, room := range liveroomdao.ListRoomsByGuildFromDB(guildId) {
		if room == nil || room.ID == 0 {
			continue
		}
		switch liveroomdao.ClearLiveRoomIncomeUnsettledForCMS(room.ID) {
		case liveroomdao.ClearUnsettledIncomeCleared:
			clearedAny = true
			allZero = false
		case liveroomdao.ClearUnsettledIncomeAlreadyZero:
		default:
			return guildResetUnsettledFailed
		}
	}

	if clearedAny {
		return guildResetUnsettledReset
	}
	if allZero {
		return guildResetUnsettledNoData
	}
	return guildResetUnsettledFailed
}
