<template>
  <el-empty v-if="!data" :description="t('pages.anchorList.noIncomeData')"/>
  <el-descriptions v-else :column="1" border>
    <el-descriptions-item v-if="settlementSalary != null" :label="t('pages.anchorList.settlementSalary')">
      <span class="money-amount">{{ formatWalletBalance(settlementSalary) }}</span>
    </el-descriptions-item>
    <el-descriptions-item v-if="settlementShareAmount != null" :label="t('pages.anchorList.settlementFlowCommission')">
      <span class="money-amount">{{ formatWalletBalance(settlementShareAmount) }}</span>
    </el-descriptions-item>
    <el-descriptions-item v-if="settlementShareAmountUsd != null" :label="t('pages.anchorList.settlementShareAmountUsd')">
      <span class="money-amount">{{ formatWalletBalance(settlementShareAmountUsd) }}</span>
    </el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.liveIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.giftIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalGiftIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalPaidDanmakuIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallTicketIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalVideoCallBillingIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.shortVideoIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalShortVideoIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.gameIncome')"><span class="money-amount">{{ formatWalletBalance(data.totalGameIncome) }}</span></el-descriptions-item>
    <el-descriptions-item :label="t('pages.anchorList.totalLiveDuration')">{{ formatLiveDurationMinutes(data.totalLiveDuration, t) }}</el-descriptions-item>
    <el-descriptions-item v-if="updatedAt" :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(updatedAt) }}</el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import type {LiveRoomIncomeAmounts} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

defineProps<{
  data?: LiveRoomIncomeAmounts | null
  settlementSalary?: number
  settlementShareAmount?: number
  settlementShareAmountUsd?: number | null
  updatedAt?: string | null
}>()

const {t} = useI18n()
</script>
