<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('pages.guildTransferList.detailTitle') }}</span>
          <el-button @click="backToList">{{ t('pages.guildTransferList.backToList') }}</el-button>
        </div>
      </template>

      <template v-if="item">
        <section-title :title="t('pages.guildTransferList.basicInfo')"/>
        <el-descriptions :column="3" border>
          <el-descriptions-item :label="t('pages.guildTransferList.settlementId')">{{ item.id }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.guildId')">
            <el-button v-if="can('viewGuildDetail')" link type="primary" @click="openGuildDetail">
              {{ item.guildId }}
            </el-button>
            <span v-else>{{ item.guildId || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.guildName')">{{ item.guildName || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.settlementReceivableUsd')">
            <strong>{{ amount(item.settlementReceivableUsd) }} USD</strong>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.status')">
            <el-tag :type="statusTagType(item.status)">{{ statusLabel(item.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('common.createdAt')">{{ formatDate(item.createdAt) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.updatedAt')">{{ formatDate(item.updatedAt) }}</el-descriptions-item>
        </el-descriptions>

        <section-title :title="t('pages.guildTransferList.incomeSnapshot')"/>
        <el-descriptions :column="3" border>
          <el-descriptions-item :label="settlementText('totalIncome')">{{ amount(item.totalIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalSocialIncome')">{{ amount(item.totalSocialIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalGiftIncome')">{{ amount(item.totalGiftIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalPaidDanmakuIncome')">{{ amount(item.totalPaidDanmakuIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalVideoCallIncome')">{{ amount(item.totalVideoCallIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalVideoCallTicketIncome')">{{ amount(item.totalVideoCallTicketIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalVideoCallBillingIncome')">{{ amount(item.totalVideoCallBillingIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalShortVideoIncome')">{{ amount(item.totalShortVideoIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalGameIncome')">{{ amount(item.totalGameIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalLiveDuration')">{{ amount(item.totalLiveDuration) }}</el-descriptions-item>
        </el-descriptions>

        <section-title :title="t('pages.guildTransferList.settlementCalculation')"/>
        <el-descriptions :column="3" border>
          <el-descriptions-item :label="settlementText('settlementRuleType')">{{ settlementRuleLabel }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('settlementSalary')">{{ amount(item.settlementSalary) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('settlementShareAmount')">{{ amount(item.settlementShareAmount) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.settlementShareAmountUsd')">{{ amount(item.settlementShareAmountUsd) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('guildSharePercent')">{{ amount(item.guildSharePercent) }}%</el-descriptions-item>
          <el-descriptions-item :label="settlementText('anchorSocialShareAmount')">{{ amount(item.anchorSocialShareAmount) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('guildSocialShareAmount')">{{ amount(item.guildSocialShareAmount) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('anchorGameShareAmountGold')">{{ amount(item.anchorGameShareAmountGold) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('guildGameShareAmountGold')">{{ amount(item.guildGameShareAmountGold) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('goldToDiamondRate')">{{ item.goldToDiamondRate || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('usdToGoldRate')">{{ item.usdToGoldRate || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('gameShareAmountDiamond')">{{ amount(item.gameShareAmountDiamond) }}</el-descriptions-item>
          <el-descriptions-item :label="settlementText('totalSettlementDiamond')">{{ amount(item.totalSettlementDiamond) }}</el-descriptions-item>
        </el-descriptions>

        <section-title :title="t('pages.guildTransferList.payoutInfo')"/>
        <el-descriptions :column="3" border>
          <el-descriptions-item :label="t('pages.guildTransferList.transferAt')">{{ formatDate(item.transferAt) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.transferOrderId')">{{ item.transferOrderId || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.transferPlatformNo')">{{ item.transferPlatformNo || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildTransferList.transferCurrency')">
            <el-button v-if="item.transferCurrency && can('transferInfo')" link type="primary" @click="openTransferInfo">
              {{ item.transferCurrency }}
            </el-button>
            <span v-else>{{ item.transferCurrency || '-' }}</span>
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
      </template>

      <el-empty v-else-if="!loading" :description="t('pages.guildTransferList.noData')"/>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, defineComponent, h, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const SectionTitle = defineComponent({
  props: {title: {type: String, required: true}},
  setup: props => () => h('div', {class: 'section-title'}, props.title),
})

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const isCoinMerchant = Number(route.query.guildType) === 1
const pagePermission = isCoinMerchant
    ? 'CoinMerchantGuildTransferManagement'
    : 'GuildTransferManagement'
const {can} = usePagePermission(pagePermission)
const loading = ref(false)
const item = ref<GuildIncomeSettlementLogItem | null>(null)

const settlementText = (key: string) => t(`pages.guildIncomeSettlementLogList.${key}`)
const amount = (value: number | null | undefined) => formatWalletBalance(Number(value || 0))
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
  const id = String(route.params.id || '').trim()
  if (!id) return
  loading.value = true
  try {
    const response = await guildIncomeSettlementLogApi.getDetail({id})
    item.value = response.item || null
  } catch (error) {
    console.error('Failed to load guild settlement detail:', error)
    ElMessage.error(t('pages.guildTransferList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const backToList = () => router.push({
  name: isCoinMerchant ? 'CoinMerchantGuildTransferManagement' : 'GuildTransferManagement',
})

const openGuildDetail = () => {
  if (!item.value?.guildId) return
  router.push({
    name: 'GuildDetail',
    query: {id: item.value.guildId, name: item.value.guildName || ''},
  })
}

const openTransferInfo = () => {
  if (!item.value?.guildId) return
  router.push({
    name: 'GuildTransferInfoEdit',
    params: {guildId: item.value.guildId},
    query: {
      guildName: item.value.guildName || '',
      from: isCoinMerchant ? 'coinMerchantGuildTransfer' : 'guildTransfer',
    },
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

.section-title {
  margin: 22px 0 12px;
  padding-left: 10px;
  border-left: 3px solid var(--el-color-primary);
  font-size: 16px;
  font-weight: 600;
}

.section-title:first-child {
  margin-top: 0;
}
</style>
