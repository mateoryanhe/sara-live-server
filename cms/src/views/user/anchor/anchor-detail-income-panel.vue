<template>
  <el-empty v-if="!data" :description="t('pages.anchorList.noIncomeData')"/>
  <div v-else class="income-panel">
    <section class="detail-section">
      <div class="detail-section-title">{{ t('pages.anchorList.incomeOverviewSection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.totalIncomeConvertedDiamond')">
          <span class="money-amount">{{ formatWalletBalance(data.totalIncome) }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="showAnchorLiveDurationInOverview" :label="t('pages.anchorList.totalLiveDuration')">
          {{ formatLiveDurationMinutes(data.totalLiveDuration, t) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="!forGuild" :label="t('pages.anchorList.gameTotalGoldFlow')">
          <span class="money-amount">{{ formatWalletBalance(data.totalGameIncome) }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </section>

    <section v-if="hasGuildSettlementSummary" class="detail-section">
      <div class="detail-section-title">{{ t('pages.guildList.settlementSummarySection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item v-if="showGuildSettledSalary" :label="t('pages.guildList.settledSalaryTotal')">
          <span class="money-amount">{{ formatWalletBalance(guildSettlementSalary) }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="showGuildSettledShare" :label="t('pages.guildList.settledGuildShareDiamond')">
          <span class="money-amount">{{ formatWalletBalance(guildSettlementShare) }}</span>
        </el-descriptions-item>
        <template v-if="guildSettledView">
          <el-descriptions-item :label="t('pages.guildList.settledGuildReceivableUsdTotal')">
            <span class="money-amount">{{ formatWalletBalance(guildSettlementGuildReceivableUsd) }} USD</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildList.settledAnchorReceivableUsdTotal')">
            <span class="money-amount">{{ formatWalletBalance(guildSettlementAnchorReceivableUsd) }} USD</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildList.settledReceivableUsdTotal')">
            <span class="money-amount">{{ formatWalletBalance(guildSettlementReceivableUsd) }} USD</span>
          </el-descriptions-item>
        </template>
      </el-descriptions>
    </section>

    <section v-if="forGuild" class="detail-section">
      <div class="detail-section-title">{{ t('pages.guildList.gameFlowSection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.gameTotalGoldFlow')">
          <span class="money-amount">{{ formatWalletBalance(data.totalGameIncome) }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </section>

    <section v-if="showUnsettledLiveStats" class="detail-section">
      <div class="detail-section-title">{{ t('pages.anchorList.unsettledLiveStatsSection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">
          {{ formatLiveDurationMinutes(data.totalLiveDuration, t) }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.effectiveLiveDays')">
          {{ effectiveLiveDaysDisplay }}
        </el-descriptions-item>
      </el-descriptions>
    </section>

    <section v-if="hasAnchorSettlementSummary" class="detail-section">
      <div class="detail-section-title">{{ t('pages.anchorList.settlementSummarySection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.settlementSalary')">
          <span class="money-amount">{{ settlementSalary == null ? '-' : formatWalletBalance(settlementSalary) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.settlementFlowCommission')">
          <span class="money-amount">{{ settlementShareAmount == null ? '-' : formatWalletBalance(settlementShareAmount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.settlementShareAmountUsd')">
          <span class="money-amount">{{ settlementShareAmountUsd == null ? '-' : formatWalletBalance(settlementShareAmountUsd) }}</span>
        </el-descriptions-item>
        <el-descriptions-item
            v-if="settlementReceivableUsd != null"
            :label="t('pages.anchorList.settlementReceivableUsd')"
        >
          <span class="money-amount">{{ formatWalletBalance(settlementReceivableUsd) }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-else class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
      </el-descriptions>
    </section>

    <section v-if="showFullSocialBreakdown" class="detail-section">
      <div class="detail-section-title">{{ socialSectionTitle }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.socialTotalDiamondFlow')">
          <span class="money-amount">{{ formatWalletBalance(data.totalSocialIncome) }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.giftIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalGiftIncome) }}</span></el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalPaidDanmakuIncome) }}</span></el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallIncome) }}</span></el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallTicketIncome) }}</span></el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallBillingIncome) }}</span></el-descriptions-item>
        <el-descriptions-item :label="t('pages.anchorList.shortVideoIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalShortVideoIncome) }}</span></el-descriptions-item>
      </el-descriptions>
    </section>

    <section v-else-if="forGuild" class="detail-section">
      <div class="detail-section-title">{{ t('pages.guildList.settledSocialFlowSection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.socialTotalDiamondFlow')">
          <span class="money-amount">{{ formatWalletBalance(data.totalSocialIncome) }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </section>

    <div v-if="updatedAt" class="income-updated-at">
      {{ t('pages.anchorList.roomUpdatedAt') }}：{{ formatDate(updatedAt) }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const props = defineProps<{
  data?: LiveRoomIncomeAmounts | null
  /** 未结算收益 Tab：有效直播天数与直播时长合并展示 */
  effectiveLiveDays?: number
  settlementSalary?: number
  settlementShareAmount?: number
  settlementShareAmountUsd?: number | null
  /** 工会详情已结算/累计 Tab 展示 */
  settlementGuildReceivableUsd?: number | null
  settlementAnchorReceivableUsd?: number | null
  settlementReceivableUsd?: number | null
  /** 工会详情：无有效开播统计，游戏流水独立分区 */
  forGuild?: boolean
  /** 工会已结算/生涯累计：结算汇总与社交明细精简 */
  guildSettledView?: boolean
  updatedAt?: string | null
}>()

const {t} = useI18n()
const detailLabelWidth = 170
const forGuild = computed(() => props.forGuild === true)
const showUnsettledLiveStats = computed(() => !forGuild.value && props.effectiveLiveDays != null)
const showAnchorLiveDurationInOverview = computed(() => !forGuild.value && !showUnsettledLiveStats.value)
const effectiveLiveDaysDisplay = computed(() => {
  const days = Math.max(0, Math.trunc(Number(props.effectiveLiveDays ?? 0)))
  return t('pages.anchorList.effectiveLiveDaysCount', {days})
})
const guildSettledView = computed(() => forGuild.value && props.guildSettledView === true)
const guildSettlementSalary = computed(() => Number(props.settlementSalary ?? 0))
const guildSettlementShare = computed(() => Number(props.settlementShareAmount ?? 0))
const guildSettlementGuildReceivableUsd = computed(() => Number(props.settlementGuildReceivableUsd ?? 0))
const guildSettlementAnchorReceivableUsd = computed(() => Number(props.settlementAnchorReceivableUsd ?? 0))
const guildSettlementReceivableUsd = computed(() => Number(props.settlementReceivableUsd ?? 0))
const showGuildSettledSalary = computed(() => guildSettlementSalary.value !== 0)
const showGuildSettledShare = computed(() => guildSettlementShare.value !== 0)
const hasGuildSettlementSummary = computed(() => guildSettledView.value)
const hasAnchorSettlementSummary = computed(() => !forGuild.value && (
  props.settlementSalary != null
  || props.settlementShareAmount != null
  || props.settlementShareAmountUsd != null
  || props.settlementReceivableUsd != null
))
const showFullSocialBreakdown = computed(() => !guildSettledView.value)
const socialSectionTitle = computed(() => (
  guildSettledView.value ? t('pages.guildList.settledSocialFlowSection') : t('pages.anchorList.socialIncomeSection')
))
</script>

<style scoped>
.income-panel {
  padding-top: 8px;
}

.detail-section + .detail-section {
  margin-top: 20px;
}

.detail-section-title {
  margin-bottom: 10px;
  padding-left: 10px;
  border-left: 3px solid var(--el-color-primary);
  color: var(--el-text-color-primary);
  font-size: 15px;
  font-weight: 600;
  line-height: 20px;
}

.detail-descriptions :deep(.el-descriptions__table) {
  table-layout: fixed;
}

.detail-descriptions :deep(.el-descriptions__content) {
  overflow-wrap: anywhere;
}

.detail-descriptions :deep(.detail-placeholder-cell) {
  background: var(--el-fill-color-blank);
}

.income-updated-at {
  margin-top: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  text-align: right;
}

@media (max-width: 1280px) {
  .detail-descriptions :deep(.el-descriptions__label) {
    width: 140px !important;
  }
}
</style>
