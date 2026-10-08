<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorPayoutList') }} - {{ t('pages.guildTransferList.viewDetail') }}</span>
          <el-button @click="backToList">{{ t('pages.guildTransferList.backToList') }}</el-button>
        </div>
      </template>

      <template v-if="item">
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="t('pages.guildTransferList.basicInfo')" name="basic">
            <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
              <el-descriptions-item :label="t('pages.guildTransferList.settlementId')">{{ item.id }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.status')">
                <el-tag :type="statusTagType(item.status)">{{ statusLabel(item.status) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="settlementText('settlementRuleType')">
                <el-tag :type="Number(item.settlementRuleType) === 1 ? 'success' : 'info'">
                  {{ settlementRuleLabel }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorIncomeSettlementLogList.roomId')">
                {{ item.roomId || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorIncomeSettlementLogList.roomNickname')">
                <el-button
                    v-if="item.roomId && item.roomNickname"
                    link
                    type="primary"
                    @click="openAnchorDetail"
                >
                  {{ item.roomNickname }}
                </el-button>
                <span v-else>{{ item.roomNickname || '-' }}</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.avatar')">
                <el-image
                    v-if="item.roomAvatar"
                    :preview-src-list="[item.roomAvatar]"
                    :src="item.roomAvatar"
                    fit="cover"
                    hide-on-click-modal
                    preview-teleported
                    class="anchor-avatar"
                />
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.settlementReceivableUsd')">
                <strong>{{ amount(item.settlementReceivableUsd) }} USD</strong>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.createdAt')">{{ formatDate(item.createdAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.updatedAt')">{{ formatDate(item.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildTransferList.incomeSnapshot')" lazy name="income">
            <section class="detail-section">
            <div class="detail-section-title">{{ settlementText('snapshotSummary') }}（钻石）</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('totalIncome')">{{ amount(item.totalIncome) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('liveStatsSection') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('totalLiveDuration')">
                  {{ formatLiveDurationMinutes(item.totalLiveDuration, t) }}
                </el-descriptions-item>
                <el-descriptions-item :label="settlementText('effectiveLiveDays')">
                  {{ formatEffectiveLiveDays(item.effectiveLiveDays) }}
                </el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
            <div class="detail-section-title">{{ settlementText('socialIncomeSection') }}（钻石）</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('totalSocialIncome')">{{ amount(item.totalSocialIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalGiftIncome')">{{ amount(item.totalGiftIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalPaidDanmakuIncome')">{{ amount(item.totalPaidDanmakuIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalVideoCallIncome')">{{ amount(item.totalVideoCallIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalVideoCallTicketIncome')">{{ amount(item.totalVideoCallTicketIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalVideoCallBillingIncome')">{{ amount(item.totalVideoCallBillingIncome) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('totalShortVideoIncome')">{{ amount(item.totalShortVideoIncome) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
            <div class="detail-section-title">{{ settlementText('gameIncomeSection') }}（金币）</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('totalGameIncome')">{{ amount(item.totalGameIncome) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildTransferList.settlementCalculation')" lazy name="calculation">
            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('calculationSummary') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="guildSettlementText('totalSettlementDiamond')">{{ amount(item.totalSettlementDiamond) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('settlementShareAmountUsd')">{{ amount(item.settlementReceivableUsd) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('salarySection') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('settlementSalary')">{{ amount(item.settlementSalary) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('hasSalary')">
                  <el-tag :type="item.hasSalary ? 'success' : 'info'">
                    {{ item.hasSalary ? t('common.yes') : t('common.no') }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('socialShareSection') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('anchorSocialShareAmount')">{{ amount(item.anchorSocialShareAmount) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('anchorSocialSharePercent')">{{ percent(item.anchorSocialSharePercent) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('gameShareSection') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="settlementText('anchorGameShareAmountGold')">{{ amount(item.anchorGameShareAmountGold) }}</el-descriptions-item>
                <el-descriptions-item :label="settlementText('anchorGameSharePercent')">{{ percent(item.anchorGameSharePercent) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>

            <section class="detail-section">
              <div class="detail-section-title">{{ settlementText('conversionSection') }}</div>
              <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                <el-descriptions-item :label="guildSettlementText('goldToDiamondRate')">{{ item.goldToDiamondRate || 0 }}</el-descriptions-item>
                <el-descriptions-item :label="guildSettlementText('usdToGoldRate')">{{ item.usdToGoldRate || 0 }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                <el-descriptions-item :label="guildSettlementText('gameShareAmountDiamond')">{{ amount(item.gameShareAmountDiamond) }}</el-descriptions-item>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
              </el-descriptions>
            </section>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildTransferList.payoutInfo')" lazy name="payout">
            <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
              <el-descriptions-item :label="t('pages.guildTransferList.transferAt')">{{ formatDate(item.transferAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferOrderId')">{{ item.transferOrderId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferPlatformNo')">{{ item.transferPlatformNo || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferCurrency')">
                {{ item.transferCurrency || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferLocalAmount')">
                {{ item.transferLocalAmount ? amount(item.transferLocalAmount) : '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferPayeeName')">{{ item.transferPayeeName || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferBankName')">{{ item.transferBankName || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferAccountNo')">{{ item.transferAccountNo || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferBankCode')">{{ item.transferBankCode || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildTransferList.transferFailMsg')" :span="3">
                {{ item.transferFailMsg || '-' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>
        </el-tabs>
      </template>

      <el-empty v-else-if="!loading" :description="t('pages.guildTransferList.noData')"/>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const loading = ref(false)
const activeTab = ref('basic')
const item = ref<AnchorIncomeSettlementLogItem | null>(null)
const detailLabelWidth = 180

const settlementText = (key: string) => t(`pages.anchorIncomeSettlementLogList.${key}`)
const guildSettlementText = (key: string) => t(`pages.guildIncomeSettlementLogList.${key}`)
const amount = (value: number | null | undefined) => formatWalletBalance(Number(value || 0))
const percent = (value: number | null | undefined) => `${Number(value || 0)}%`
const formatEffectiveLiveDays = (value: number | null | undefined) => {
  const days = Math.max(0, Math.trunc(Number(value ?? 0)))
  return t('pages.anchorIncomeSettlementLogList.effectiveLiveDaysCount', {days})
}
const settlementRuleLabel = computed(() => Number(item.value?.settlementRuleType) === 1
    ? settlementText('settlementRuleTiered')
    : settlementText('settlementRuleLegacy'))

const statusLabel = (status: number | undefined) => {
  const value = Number(status)
  if (value === 1) return t('pages.guildTransferList.statusApproved')
  if (value === 2) return t('pages.guildTransferList.statusTransferred')
  if (value === 3) return t('pages.guildTransferList.statusTransferring')
  return t('pages.guildTransferList.statusPending')
}

const statusTagType = (status: number | undefined) => {
  const value = Number(status)
  if (value === 1) return 'success'
  if (value === 2) return 'info'
  if (value === 3) return ''
  return 'warning'
}

const fetchDetail = async () => {
  const settlementId = String(route.params.id || '').trim()
  if (!settlementId) return
  loading.value = true
  try {
    const response = await anchorIncomeSettlementLogApi.getList({
      settlementId,
      directPayout: true,
      includeTransferInfo: true,
      pageIndex: 1,
      pageSize: 1,
    })
    item.value = response.data?.[0] || null
  } catch (error) {
    console.error('Failed to load platform anchor settlement detail:', error)
    ElMessage.error(t('pages.anchorIncomeSettlementLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const backToList = () => router.push({name: 'PlatformAnchorPayoutList'})

const openAnchorDetail = () => {
  if (!item.value?.roomId) return
  router.push({
    name: 'AnchorDetail',
    query: {id: String(item.value.roomId)},
  })
}

onMounted(fetchDetail)
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.anchor-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
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

.detail-descriptions :deep(.detail-placeholder-cell) {
  background: var(--el-fill-color-blank);
}
</style>
