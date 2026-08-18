<template>
  <el-descriptions :column="1" border>
    <el-descriptions-item :label="t('pages.anchorList.liveIncome')">{{ formatNum(displayData.totalIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.giftIncome')">{{ formatNum(displayData.totalGiftIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')">{{ formatNum(displayData.totalPaidDanmakuIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomTicketIncome')">{{ formatNum(displayData.totalPrivateRoomTicketIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomWatchIncome')">{{ formatNum(displayData.totalPrivateRoomWatchIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')">{{ formatNum(displayData.totalVideoCallIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')">{{ formatNum(displayData.totalVideoCallTicketIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')">{{ formatNum(displayData.totalVideoCallBillingIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">{{ formatNum(displayData.totalLiveDuration) }}</el-descriptions-item>
    <el-descriptions-item v-if="showSettlementShare" :label="t('pages.anchorList.settlementShareAmount')">
      {{ formatNum(settlementShareAmount) }}
    </el-descriptions-item>
    <el-descriptions-item v-if="updatedAt" :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(updatedAt) }}</el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

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

const formatNum = (value?: number | null) => formatAmount(value, '0')

const formatDate = (dateString?: string | null) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}
</script>
