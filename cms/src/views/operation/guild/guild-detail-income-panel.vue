<template>
  <el-descriptions :column="1" border>
    <el-descriptions-item :label="t('pages.anchorList.liveIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.giftIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalGiftIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalPaidDanmakuIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomTicketIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalPrivateRoomTicketIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomWatchIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalPrivateRoomWatchIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallTicketIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')"><span class="money-amount">{{ formatWalletBalance(displayData.totalVideoCallBillingIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">{{ formatLiveDurationMinutes(displayData.totalLiveDuration, t) }}</el-descriptions-item>
    <el-descriptions-item v-if="showSettlementShare" :label="t('pages.anchorList.settlementShareAmount')">
      <span class="money-amount">{{ formatWalletBalance(settlementShareAmount) }}</span>
    </el-descriptions-item>
    <el-descriptions-item v-if="updatedAt" :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(updatedAt) }}</el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  data?: LiveRoomIncomeAmounts | null
  settlementShareAmount?: number | null
  showSettlementShare?: boolean
  updatedAt?: string | null
}>()

const {t} = useI18n()

const displayData = computed<Required<LiveRoomIncomeAmounts>>(() => ({
  totalIncome: props.data?.totalIncome ?? 0,
  totalGiftIncome: props.data?.totalGiftIncome ?? 0,
  totalPaidDanmakuIncome: props.data?.totalPaidDanmakuIncome ?? 0,
  totalPrivateRoomTicketIncome: props.data?.totalPrivateRoomTicketIncome ?? 0,
  totalPrivateRoomWatchIncome: props.data?.totalPrivateRoomWatchIncome ?? 0,
  totalVideoCallIncome: props.data?.totalVideoCallIncome ?? 0,
  totalVideoCallTicketIncome: props.data?.totalVideoCallTicketIncome ?? 0,
  totalVideoCallBillingIncome: props.data?.totalVideoCallBillingIncome ?? 0,
  totalLiveDuration: props.data?.totalLiveDuration ?? 0,
}))

const formatDate = (dateString?: string | null) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}
</script>
