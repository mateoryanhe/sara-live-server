<template>
  <el-empty v-if="!data && !showSettlementShare" :description="t('pages.anchorList.noIncomeData')"/>
  <el-descriptions v-else :column="1" border>
    <template v-if="showSettlementShare">
      <el-descriptions-item :label="t('pages.anchorList.settlementSalary')">
        <span class="money-amount">{{ formatWalletBalance(settlementSalaryValue) }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="t('pages.anchorList.settlementFlowCommission')">
        <span class="money-amount">{{ formatWalletBalance(settlementShareAmountValue) }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="t('pages.anchorList.settlementShareAmountUsd')">
        <span class="money-amount">{{ formatWalletBalance(settlementShareAmountUsdValue) }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="t('pages.anchorList.settlementReceivableUsd')">
        <span class="money-amount">{{ formatWalletBalance(settlementReceivableUsdValue) }}</span>
      </el-descriptions-item>
    </template>
    <el-descriptions-item :label="t('pages.anchorList.liveIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.giftIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalGiftIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalPaidDanmakuIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallTicketIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallBillingIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.shortVideoIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalShortVideoIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.gameIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalGameIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">{{ formatLiveDurationMinutes(displayData.totalLiveDuration, t) }}</el-descriptions-item>
    <el-descriptions-item v-if="updatedAt" :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(updatedAt) }}</el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

type GuildIncomeDetail = LiveRoomIncomeAmounts & {
  settlementSalary?: number
  settlementShareAmount?: number
  settlementShareAmountUsd?: number
  settlementReceivableUsd?: number
}

const props = defineProps<{
  data?: GuildIncomeDetail | null
  settlementShareAmount?: number | null
  settlementSalary?: number | null
  settlementShareAmountUsd?: number | null
  settlementReceivableUsd?: number | null
  showSettlementShare?: boolean
  updatedAt?: string | null
}>()

const {t} = useI18n()

const displayData = computed<Required<LiveRoomIncomeAmounts>>(() => ({
  totalIncome: props.data?.totalIncome ?? 0,
  totalGiftIncome: props.data?.totalGiftIncome ?? 0,
  totalPaidDanmakuIncome: props.data?.totalPaidDanmakuIncome ?? 0,
  totalVideoCallIncome: props.data?.totalVideoCallIncome ?? 0,
  totalVideoCallTicketIncome: props.data?.totalVideoCallTicketIncome ?? 0,
  totalVideoCallBillingIncome: props.data?.totalVideoCallBillingIncome ?? 0,
  totalShortVideoIncome: props.data?.totalShortVideoIncome ?? 0,
  totalGameIncome: props.data?.totalGameIncome ?? 0,
  totalLiveDuration: props.data?.totalLiveDuration ?? 0,
}))

const settlementSalaryValue = computed(() => props.settlementSalary ?? props.data?.settlementSalary ?? 0)
const settlementShareAmountValue = computed(() => props.settlementShareAmount ?? props.data?.settlementShareAmount ?? 0)
const settlementShareAmountUsdValue = computed(() => props.settlementShareAmountUsd ?? props.data?.settlementShareAmountUsd ?? 0)
const settlementReceivableUsdValue = computed(() => props.settlementReceivableUsd ?? props.data?.settlementReceivableUsd ?? 0)
</script>
