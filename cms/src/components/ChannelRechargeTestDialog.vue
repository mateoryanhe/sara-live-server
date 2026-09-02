<template>
  <el-dialog
      :model-value="modelValue"
      :title="t('pages.rechargeOrderList.channelRechargeTest')"
      width="760px"
      @closed="resetChannelTest"
      @update:model-value="emit('update:modelValue', $event)"
  >
    <el-form ref="channelTestFormRef" :model="channelTestForm" :rules="channelTestRules" label-width="90px">
      <el-form-item :label="t('pages.rechargeOrderList.playerId')" prop="userId">
        <el-input
            v-model="channelTestForm.userId"
            clearable
            :disabled="!!lockedUserId"
            :placeholder="t('pages.rechargeOrderList.enterPlayerId')"
        />
      </el-form-item>
      <el-form-item :label="t('pages.rechargeOrderList.currencyCode')" prop="currencyCode">
        <el-select
            v-model="channelTestForm.currencyCode"
            clearable
            filterable
            :loading="channelTestCurrencyLoading"
            :placeholder="t('pages.rechargeOrderList.selectCurrencyPlaceholder')"
            style="width: 100%"
            @change="loadChannelTestCfgList"
        >
          <el-option
              v-for="item in channelTestCurrencyList"
              :key="item.currencyCode"
              :label="currencyOptionLabel(item)"
              :value="item.currencyCode"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <el-table
        v-loading="channelTestCfgLoading || channelTestCreating"
        :data="channelTestCfgList"
        highlight-current-row
        max-height="420"
        style="width: 100%"
        @row-click="handleChannelTestCfgPick"
    >
      <el-table-column :label="t('pages.rechargeOrderList.cfgId')" prop="id" width="100"/>
      <el-table-column :label="t('pages.rechargeOrderList.rechargeCfgName')" min-width="120" prop="name"/>
      <el-table-column :label="t('pages.rechargeOrderList.priceUsd')" width="100">
        <template #default="{ row }">{{ formatAmount(row.price) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.priceLocal')" width="140">
        <template #default="{ row }">
          {{ formatAmount(row.displayPrice) }}
          {{ channelTestForm.currencyCode || '' }}
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.gold')" width="120">
        <template #default="{ row }">{{ formatAmount(cfgGold(row)) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click.stop="handleChannelTestCfgPick(row)">
            {{ t('pages.rechargeOrderList.openPayUrl') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {fiatCurrencyApi, rechargeCfgApi, rechargeOrderApi} from '@/api'
import type {FiatCurrency, RechargeCfg} from '@/types/api.ts'
import {formatAmount} from '@/utils/number-format'

const props = defineProps<{
  modelValue: boolean
  userId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  created: []
}>()

const {t} = useI18n()
const channelTestCurrencyLoading = ref(false)
const channelTestCfgLoading = ref(false)
const channelTestCreating = ref(false)
const channelTestCurrencyList = ref<FiatCurrency[]>([])
const channelTestCfgList = ref<Array<RechargeCfg & {displayPrice: number}>>([])
const channelTestFormRef = ref<FormInstance>()
const lockedUserId = computed(() => String(props.userId || '').trim())
const channelTestForm = reactive({
  userId: '',
  currencyCode: '',
  currencyType: 1,
  rate: 1,
})

const channelTestRules = computed<FormRules>(() => ({
  userId: [{required: true, message: t('pages.rechargeOrderList.playerIdRequired'), trigger: 'blur'}],
  currencyCode: [{required: true, message: t('pages.rechargeOrderList.currencyRequired'), trigger: 'change'}],
}))

const cfgGold = (row: RechargeCfg) => {
  const anyRow = row as RechargeCfg & {gold?: number; diamond?: number}
  return anyRow.gold ?? anyRow.diamond ?? 0
}

const currencyOptionLabel = (item: FiatCurrency) => {
  const typeLabel = item.currencyType === 2
      ? t('pages.rechargeOrderList.currencyTypeCrypto')
      : t('pages.rechargeOrderList.currencyTypeFiat')
  return `${item.currencyCode} · ${item.name || '-'} (${typeLabel})`
}

const loadCurrencies = async () => {
  channelTestCurrencyLoading.value = true
  try {
    const response = await fiatCurrencyApi.getList({
      pageIndex: 1,
      pageSize: 200,
      statusFilter: 2,
    })
    channelTestCurrencyList.value = (response.data || []).slice().sort((a, b) => (b.sort || 0) - (a.sort || 0))
  } catch (error) {
    console.error('load fiat currency failed:', error)
    ElMessage.error(t('pages.rechargeOrderList.loadCurrencyFailed'))
    channelTestCurrencyList.value = []
  } finally {
    channelTestCurrencyLoading.value = false
  }
}

const resetChannelTest = () => {
  channelTestForm.userId = lockedUserId.value
  channelTestForm.currencyCode = ''
  channelTestForm.currencyType = 1
  channelTestForm.rate = 1
  channelTestCfgList.value = []
  channelTestFormRef.value?.clearValidate()
}

const loadChannelTestCfgList = async () => {
  const code = channelTestForm.currencyCode
  if (!code) {
    channelTestCfgList.value = []
    return
  }
  const currency = channelTestCurrencyList.value.find((item) => item.currencyCode === code)
  if (!currency) {
    channelTestCfgList.value = []
    return
  }
  channelTestForm.currencyType = currency.currencyType
  channelTestForm.rate = 1
  channelTestCfgLoading.value = true
  try {
    if (currency.currencyType !== 2) {
      const rateRes = await fiatCurrencyApi.getExchangeRate(currency.currencyCode)
      channelTestForm.rate = rateRes?.rate || 0
      if (!channelTestForm.rate) {
        ElMessage.error(t('pages.rechargeOrderList.loadRateFailed'))
        channelTestCfgList.value = []
        return
      }
    }
    const response = await rechargeCfgApi.getRechargeCfgList({
      pageIndex: 1,
      pageSize: 200,
      statusFilter: 2,
    })
    const list = (response.data || []).slice().sort((a, b) => (a.price || 0) - (b.price || 0))
    channelTestCfgList.value = list.map((item) => ({
      ...item,
      displayPrice: currency.currencyType === 2 ? item.price : item.price * channelTestForm.rate,
    }))
  } catch (error) {
    console.error('load channel test cfg failed:', error)
    ElMessage.error(t('pages.rechargeOrderList.loadRechargeCfgFailed'))
    channelTestCfgList.value = []
  } finally {
    channelTestCfgLoading.value = false
  }
}

const handleChannelTestCfgPick = async (row: RechargeCfg & {displayPrice?: number}) => {
  if (!row?.id || channelTestCreating.value) return
  if (!channelTestFormRef.value) return
  try {
    await channelTestFormRef.value.validate()
  } catch {
    return
  }
  channelTestCreating.value = true
  try {
    const res = await rechargeOrderApi.createChannelRechargeOrderTest({
      userId: channelTestForm.userId.trim(),
      cfgId: Number(row.id),
      currencyCode: channelTestForm.currencyCode,
    })
    if (!res?.payUrl) {
      ElMessage.error(t('pages.rechargeOrderList.openPayUrlFailed'))
      return
    }
    ElMessage.success(t('pages.rechargeOrderList.channelTestCreated', {
      orderId: res.orderId,
      price: formatAmount(res.price),
      currency: res.currency,
    }))
    window.open(res.payUrl, '_blank')
    emit('update:modelValue', false)
    emit('created')
  } catch (error) {
    console.error('channel recharge test failed:', error)
  } finally {
    channelTestCreating.value = false
  }
}

watch(
    () => props.modelValue,
    async (visible) => {
      if (!visible) return
      channelTestForm.userId = lockedUserId.value
      channelTestForm.currencyCode = ''
      channelTestForm.currencyType = 1
      channelTestForm.rate = 1
      channelTestCfgList.value = []
      await loadCurrencies()
    },
)
</script>
