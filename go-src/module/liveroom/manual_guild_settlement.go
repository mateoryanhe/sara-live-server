package liveroom

import (
	"context"
	"errors"
	"sync"
	"time"

	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrtime"
	"xr-game-server/dao/anchorsalarycfgdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
)

var (
	settlementRunMu      sync.Mutex
	ErrSettlementRunning = errors.New("income settlement is already running")
)

// BatchGuildSettlementResult 批量立即结算的业务结果。
type BatchGuildSettlementResult struct {
	SettledCount int
	NoDataCount  int
	FailGuildIds []uint64
}

type immediateGuildSettlementOutcome uint8

const (
	immediateGuildSettlementFailed immediateGuildSettlementOutcome = iota
	immediateGuildSettlementNoData
	immediateGuildSettlementCreated
)

// BatchSettleGuildsImmediately 立即结算所选上架工会。
// 同一时刻只允许一套手工或周结算流程运行，防止重复消费未结算快照。
func BatchSettleGuildsImmediately(ctx context.Context, guildIds []uint64) (*BatchGuildSettlementResult, error) {
	if !settlementRunMu.TryLock() {
		return nil, ErrSettlementRunning
	}
	defer settlementRunMu.Unlock()

	result := &BatchGuildSettlementResult{FailGuildIds: make([]uint64, 0)}
	tierCfg := loadAnchorWeeklySettlementCfg()
	legacySalaryCfgs := anchorsalarycfgdao.ListAllOrderBySalaryDesc()
	periodEnd := xrtime.NextWeekStart(time.Now())
	xrlog.DetailLog.Infof(ctx, "batch immediate guild settlement start guildIds=%v periodEnd=%s", guildIds, periodEnd.Format(time.RFC3339))

	for _, guildId := range guildIds {
		guild := guilddao.GetGuildById(guildId)
		if guild == nil {
			result.FailGuildIds = append(result.FailGuildIds, guildId)
			continue
		}
		switch settleOneGuildImmediately(guild, tierCfg, legacySalaryCfgs, periodEnd) {
		case immediateGuildSettlementCreated:
			result.SettledCount++
		case immediateGuildSettlementNoData:
			result.NoDataCount++
		default:
			result.FailGuildIds = append(result.FailGuildIds, guildId)
		}
	}

	xrlog.DetailLog.Infof(ctx, "batch immediate guild settlement done settled=%d noData=%d failed=%d failGuildIds=%v",
		result.SettledCount, result.NoDataCount, len(result.FailGuildIds), result.FailGuildIds)
	return result, nil
}

func settleOneGuildImmediately(
	guild *entity.LiveGuild,
	tierCfg *anchorWeeklySettlementCfg,
	legacySalaryCfgs []*entity.AnchorSalaryCfg,
	periodEnd time.Time,
) immediateGuildSettlementOutcome {
	if guild == nil || guild.ID == 0 {
		return immediateGuildSettlementFailed
	}

	rooms := liveroomdao.ListRoomsByGuild(guild.ID)
	liveroomdao.PrepareGuildWeeklyAnchorSettlement(guild.ID)
	liveroomdao.ResetGuildWeeklyAnchorLegacySettlement(guild.ID)

	if guild.GuildType == entity.LiveGuildTypeCoinMerchant {
		coinMerchantGuildIds := map[uint64]struct{}{guild.ID: {}}
		for _, room := range rooms {
			if room == nil || room.ID == 0 {
				continue
			}
			settleOneAnchor(room, legacySalaryCfgs, coinMerchantGuildIds)
		}
	} else {
		// 与周结算一致：先检查整家工会，任一主播缺配置时不消费任何流水。
		for _, room := range rooms {
			if room == nil || room.ID == 0 {
				continue
			}
			if !normalGuildAnchorTieredReady(room, tierCfg, periodEnd) {
				liveroomdao.MarkGuildWeeklyAnchorSettlementBlocked(guild.ID)
				return immediateGuildSettlementFailed
			}
		}
		for _, room := range rooms {
			if room == nil || room.ID == 0 {
				continue
			}
			if !settleNormalGuildAnchorTiered(room, tierCfg, periodEnd) {
				liveroomdao.MarkGuildWeeklyAnchorSettlementBlocked(guild.ID)
				return immediateGuildSettlementFailed
			}
		}
	}

	created := settleOneGuild(guild)
	if liveroomdao.IsGuildWeeklyAnchorSettlementBlocked(guild.ID) {
		return immediateGuildSettlementFailed
	}
	if created {
		return immediateGuildSettlementCreated
	}
	return immediateGuildSettlementNoData
}
