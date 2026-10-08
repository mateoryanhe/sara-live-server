package liveroom

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrtime"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
)

func initAnchorSettlement() {
	event.Sub(gameevent.WeekEvent, onWeekAnchorSettlement)
}

func onWeekAnchorSettlement(_ any) {
	settlementRunMu.Lock()
	defer settlementRunMu.Unlock()
	settleOnShelfAnchors()
	settleOnShelfGuilds()
}

// settleOnShelfAnchors 周一0点:结算全部上架主播薪资+未结算收益
func settleOnShelfAnchors() {
	tierCfg := loadAnchorWeeklySettlementCfg()
	periodEnd := xrtime.WeekStart(time.Now())
	rooms := liveroomdao.GetAllLiveRoom()
	coinMerchantGuildIds := loadCoinMerchantGuildIdSet()
	ctx := gctx.New()
	liveroomdao.ResetGuildWeeklyAnchorSalary()
	liveroomdao.ResetGuildWeeklyAnchorSettlement()
	g.Log().Infof(ctx, "anchor weekly settlement start, rooms=%d", len(rooms))
	// 先检查整个普通工会的配置；任一主播缺少所需配置，本轮整家工会都不结算。
	for _, room := range rooms {
		if room == nil || room.ID == 0 || room.GuildId == 0 || isCoinMerchantGuildAnchor(room, coinMerchantGuildIds) {
			continue
		}
		if !normalGuildAnchorTieredReady(room, tierCfg, periodEnd) {
			liveroomdao.MarkGuildWeeklyAnchorSettlementBlocked(room.GuildId)
			g.Log().Warningf(ctx, "guild weekly settlement preflight blocked: config missing guildId=%d roomId=%d", room.GuildId, room.ID)
		}
	}
	for _, room := range rooms {
		if room == nil || room.ID == 0 {
			continue
		}
		if room.GuildId == 0 {
			if !settlePlatformAnchorTiered(room, tierCfg, periodEnd) {
				g.Log().Warningf(ctx, "platform anchor weekly settlement retained roomId=%d", room.ID)
			}
			continue
		}
		if isCoinMerchantGuildAnchor(room, coinMerchantGuildIds) {
			settleOneAnchor(room, coinMerchantGuildIds)
			continue
		}
		if liveroomdao.IsGuildWeeklyAnchorSettlementBlocked(room.GuildId) {
			continue
		}
		if !settleNormalGuildAnchorTiered(room, tierCfg, periodEnd) {
			liveroomdao.MarkGuildWeeklyAnchorSettlementBlocked(room.GuildId)
		}
	}
	g.Log().Infof(ctx, "anchor weekly settlement done")
}

func loadCoinMerchantGuildIdSet() map[uint64]struct{} {
	set := make(map[uint64]struct{})
	for _, guild := range guilddao.ListOnShelfGuilds() {
		if guild != nil && guild.GuildType == entity.LiveGuildTypeCoinMerchant {
			set[guild.ID] = struct{}{}
		}
	}
	return set
}

func isCoinMerchantGuildAnchor(room *entity.LiveRoom, coinMerchantGuildIds map[uint64]struct{}) bool {
	if room == nil || room.GuildId == 0 || coinMerchantGuildIds == nil {
		return false
	}
	_, ok := coinMerchantGuildIds[room.GuildId]
	return ok
}

// settleOneAnchor 仅归档币商工会主播的原始流水。
// 币商主播没有底薪，也没有主播个人社交/游戏分佣；实际金额只在工会级结算一次。
func settleOneAnchor(room *entity.LiveRoom, coinMerchantGuildIds map[uint64]struct{}) {
	if !isCoinMerchantGuildAnchor(room, coinMerchantGuildIds) {
		return
	}
	roomId := room.ID
	dailyRows := liveroomdao.ListRecentUnsettledDailyEffectiveLives(roomId)
	unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(roomId)
	if unsettled == nil {
		return
	}
	snap := unsettled.Snapshot()
	if len(dailyRows) == 0 && !snap.HasSettlementFlow() {
		return
	}

	logRow := entity.NewAnchorIncomeSettlementLog(roomId, &snap, unsettled.EffectiveLiveDays, 0, 0, 0, 0)
	if !syndb.FlushUntilIdle(settlementPersistTimeout) || !liveroomdao.VerifyAnchorIncomeSettlementLogPersisted(logRow) {
		g.Log().Errorf(gctx.New(), "coin merchant anchor settlement retained: log persist failed roomId=%d guildId=%d logId=%d", roomId, room.GuildId, logRow.ID)
		return
	}

	unsettled.ConsumeSnapshot(&snap)
	if settled := liveroomdao.GetLiveRoomIncomeSettled(roomId); settled != nil {
		settled.AddAmounts(&snap)
	}
	if len(dailyRows) > 0 {
		liveroomdao.MarkDailyEffectiveLivesSettled(dailyRows)
	}
}

