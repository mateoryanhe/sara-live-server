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
      <el-table-column :label="t('pages.rechargeOrderList.priceUsd')" width="120">
        <template #default="{ row }">{{ formatNumberDisplay(row.price, '-', RECHARGE_PRICE_DECIMALS) }}</template>
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
import {rechargeCfgApi, rechargeOrderApi} from '@/api'
import type {RechargeCfg} from '@/types/api.ts'
import {formatAmount, formatNumberDisplay} from '@/utils/number-format'

/** 与充值档位页一致：测试档位可能是 0.0001 这类小金额 */
const RECHARGE_PRICE_DECIMALS = 4

const props = defineProps<{
  modelValue: boolean
  userId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  created: []
}>()

const {t} = useI18n()
const channelTestCfgLoading = ref(false)
const channelTestCreating = ref(false)
const channelTestCfgList = ref<RechargeCfg[]>([])
const channelTestFormRef = ref<FormInstance>()
const lockedUserId = computed(() => String(props.userId || '').trim())
const channelTestForm = reactive({
  userId: '',
})

const channelTestRules = computed<FormRules>(() => ({
  userId: [{required: true, message: t('pages.rechargeOrderList.playerIdRequired'), trigger: 'blur'}],
}))

const cfgGold = (row: RechargeCfg) => {
  const anyRow = row as RechargeCfg & {gold?: number; diamond?: number}
  return anyRow.gold ?? anyRow.diamond ?? 0
}

const resetChannelTest = () => {
  channelTestForm.userId = lockedUserId.value
  channelTestCfgList.value = []
  channelTestFormRef.value?.clearValidate()
}

const loadChannelTestCfgList = async () => {
  channelTestCfgLoading.value = true
  try {
    const response = await rechargeCfgApi.getRechargeCfgList({
      pageIndex: 1,
      pageSize: 200,
      statusFilter: 2,
    })
    channelTestCfgList.value = (response.data || []).slice().sort((a, b) => (a.price || 0) - (b.price || 0))
  } catch (error) {
    console.error('load channel test cfg failed:', error)
    ElMessage.error(t('pages.rechargeOrderList.loadRechargeCfgFailed'))
    channelTestCfgList.value = []
  } finally {
    channelTestCfgLoading.value = false
  }
}

const handleChannelTestCfgPick = async (row: RechargeCfg) => {
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
    })
    if (!res?.payUrl) {
      ElMessage.error(t('pages.rechargeOrderList.openPayUrlFailed'))
      return
    }
    ElMessage.success(t('pages.rechargeOrderList.channelTestCreated', {
      orderId: res.orderId,
      price: formatAmount(res.payAmount ?? res.price),
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
      channelTestCfgList.value = []
      await loadChannelTestCfgList()
    },
)
</script>
