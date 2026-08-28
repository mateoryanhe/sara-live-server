<template>
  <el-empty v-if="!data" :description="t('pages.anchorList.noIncomeData')"/>
  <el-descriptions v-else :column="1" border>
    <el-descriptions-item :label="t('pages.anchorList.liveIncome')">{{ formatWalletBalance(data.totalIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.giftIncome')">{{ formatAmount(data.totalGiftIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')">{{ formatAmount(data.totalPaidDanmakuIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomTicketIncome')">{{ formatAmount(data.totalPrivateRoomTicketIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.privateRoomWatchIncome')">{{ formatAmount(data.totalPrivateRoomWatchIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')">{{ formatAmount(data.totalVideoCallIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')">{{ formatAmount(data.totalVideoCallTicketIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')">{{ formatAmount(data.totalVideoCallBillingIncome) }}</el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">{{ formatLiveDurationMinutes(data.totalLiveDuration, t) }}</el-descriptions-item>
    <el-descriptions-item v-if="settlementSalary != null" :label="t('pages.anchorList.settlementSalary')">
      {{ formatAmount(settlementSalary) }}
    </el-descriptions-item>
    <el-descriptions-item v-if="settlementShareAmount != null" :label="t('pages.anchorList.settlementShareAmount')">
      {{ formatAmount(settlementShareAmount) }}
    </el-descriptions-item>
    <el-descriptions-item v-if="updatedAt" :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(updatedAt) }}</el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

defineProps<{
  data?: LiveRoomIncomeAmounts | null
  settlementSalary?: number
  settlementShareAmount?: number
  updatedAt?: string | null
}>()

const {t} = useI18n()

const formatDate = (dateString?: string | null) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}
</script>
