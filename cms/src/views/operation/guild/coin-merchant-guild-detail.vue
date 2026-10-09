<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.guildList.back') }}</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <el-empty v-if="!guildId" :description="t('pages.guildList.detailNotFound')"/>
        <el-descriptions v-else :column="1" border>
          <el-descriptions-item label="ID">{{ guildId }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildList.guildName')">{{ guildName || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildList.sharePercent')">
            {{ sharePercentText }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.coinMerchantGuildList.giftIncomeCumulative')">
            <span class="money-amount">{{ formatWalletBalance(giftIncomeCumulative) }}</span>
          </el-descriptions-item>
          <el-descriptions-item v-if="incomeUpdatedAt" :label="t('common.updatedAt')">
            {{ formatDate(incomeUpdatedAt) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {GuildDetailIncome} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const incomeData = ref<GuildDetailIncome | null>(null)

const parseQueryValue = (key: string) => {
  const value = route.query[key]
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
}

const guildId = computed(() => parseQueryValue('id'))
const guildName = computed(() => parseQueryValue('name'))
const sharePercentText = computed(() => {
  const raw = parseQueryValue('sharePercent')
  if (raw === '') {
    return '-'
  }
  return raw
})

const pageTitle = computed(() => {
  if (guildName.value) {
    return t('pages.coinMerchantGuildList.detailTitleWithName', {name: guildName.value})
  }
  if (guildId.value) {
    return t('pages.coinMerchantGuildList.detailTitleWithId', {id: guildId.value})
  }
  return t('pages.coinMerchantGuildList.detailTitle')
})

/** 与列表 unsettledTotalIncome 一致：未结算礼物流水 */
const giftIncomeCumulative = computed(() => incomeData.value?.incomeUnsettled?.totalGiftIncome ?? 0)

const incomeUpdatedAt = computed(() => incomeData.value?.incomeUnsettled?.updatedAt ?? '')

const fetchIncome = async () => {
  if (!guildId.value) {
    incomeData.value = null
    return
  }
  loading.value = true
  try {
    incomeData.value = await guildApi.getGuildDetail(guildId.value)
  } catch (error) {
    console.error('load coin merchant guild detail failed:', error)
    incomeData.value = null
    ElMessage.error(t('pages.guildList.detailFetchFailed'))
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'CoinMerchantGuildListManagement'})
}

watch(guildId, (_id, prev) => {
  if (prev !== undefined) {
    fetchIncome()
  }
})

onActivated(() => {
  fetchIncome()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  font-weight: bold;
}
</style>
