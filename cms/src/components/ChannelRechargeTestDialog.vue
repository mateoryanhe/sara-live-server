<template>
  <el-dialog
      :model-value="modelValue"
      :title="dialogTitle"
      width="820px"
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

    <!-- Step 1: 选择 HaiPay region（App 侧 currencyCode=国家简码） -->
    <template v-if="step === 'region'">
      <el-table
          v-loading="regionLoading"
          :data="regionList"
          highlight-current-row
          max-height="420"
          style="width: 100%"
          @row-click="handleRegionPick"
      >
        <el-table-column :label="t('pages.rechargeOrderList.currencyCode')" prop="code" width="90"/>
        <el-table-column :label="t('pages.rechargeOrderList.fiatCurrencyCode')" prop="currency" width="180"/>
        <el-table-column :label="t('pages.rechargeOrderList.currencyName')" min-width="220">
          <template #default="{ row }">
            {{ row.nameZh || row.nameEn || row.name || row.code }}
            <span v-if="row.nameZh && row.nameEn" class="region-en"> / {{ row.nameEn }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="handleRegionPick(row)">
              {{ t('pages.rechargeOrderList.selectCurrency') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <!-- Step 2: 选档位建单并打开 HaiPay -->
    <template v-else>
      <div class="step-bar">
        <el-button link type="primary" @click="backToRegion">← {{ t('pages.rechargeOrderList.backToSelectRegion') }}</el-button>
        <span class="step-region">{{ selectedRegion?.code }} · {{ selectedRegion?.currency }} · {{ selectedRegion?.name }}</span>
      </div>
      <el-alert v-if="selectedRegion && selectedRegion.paymentMethods.length === 0"
                :closable="false" show-icon type="warning"
                :title="t('pages.rechargeOrderList.noPaymentMethod')"/>
      <el-form-item v-else :label="t('pages.rechargeOrderList.paymentMethod')">
        <el-select v-model="selectedMethodKey" style="width: 100%" :placeholder="t('pages.rechargeOrderList.selectPaymentMethod')">
          <el-option v-for="method in selectedRegion?.paymentMethods || []"
                     :key="methodKey(method)" :label="methodLabel(method)" :value="methodKey(method)"/>
        </el-select>
      </el-form-item>
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
    </template>
  </el-dialog>

  <!-- 首次（无资料）填写 name / email / phone -->
  <el-dialog
      v-model="payerDialogVisible"
      :close-on-click-modal="false"
      :title="t('pages.rechargeOrderList.payerInfoTitle')"
      width="480px"
      append-to-body
      @closed="pendingCfg = null"
  >
    <p class="payer-tip">{{ t('pages.rechargeOrderList.payerInfoTip') }}</p>
    <el-form ref="payerFormRef" :model="payerForm" :rules="payerRules" label-width="110px">
      <el-form-item :label="t('pages.rechargeOrderList.payName')" prop="payName">
        <el-input
            v-model="payerForm.payName"
            clearable
            :placeholder="t('pages.rechargeOrderList.payNamePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('pages.rechargeOrderList.payEmail')" prop="payEmail">
        <el-input
            v-model="payerForm.payEmail"
            clearable
            :placeholder="t('pages.rechargeOrderList.payEmailPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('pages.rechargeOrderList.payPhone')" prop="payPhone">
        <el-input
            v-model="payerForm.payPhone"
            clearable
            :placeholder="t('pages.rechargeOrderList.payPhonePlaceholder')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="payerDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button :loading="channelTestCreating" type="primary" @click="confirmPayerAndCreate">
        {{ t('pages.rechargeOrderList.payerInfoConfirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {fiatCurrencyApi, rechargeCfgApi, rechargeOrderApi} from '@/api'
import type {HaiPayRegionItem} from '@/api/modules/fiatCurrency'
import type {RechargeCfg} from '@/types/api.ts'
import {formatAmount, formatNumberDisplay} from '@/utils/number-format'

/** 与充值档位页一致：测试档位可能是 0.0001 这类小金额 */
const RECHARGE_PRICE_DECIMALS = 4

type RegionRow = {
  code: string
  currency: string
  name: string
  nameEn: string
  nameZh: string
  paymentMethods: HaiPayRegionItem['paymentMethods']
}

const props = defineProps<{
  modelValue: boolean
  userId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  created: []
}>()

const {t} = useI18n()
const step = ref<'region' | 'cfg'>('region')
const selectedRegion = ref<RegionRow | null>(null)
const selectedMethodKey = ref('')
const regionList = ref<RegionRow[]>([])
const regionLoading = ref(false)
const channelTestCfgLoading = ref(false)
const channelTestCreating = ref(false)
const channelTestCfgList = ref<RechargeCfg[]>([])
const channelTestFormRef = ref<FormInstance>()
const payerFormRef = ref<FormInstance>()
const payerDialogVisible = ref(false)
const pendingCfg = ref<RechargeCfg | null>(null)
const lockedUserId = computed(() => String(props.userId || '').trim())
const channelTestForm = reactive({
  userId: '',
})
const payerForm = reactive({
  payName: '',
  payEmail: '',
  payPhone: '',
})
/** 已加载过的玩家资料缓存（同一次弹窗内） */
const cachedPayerByUserId = ref<Record<string, { name: string; email: string; phone: string }>>({})

const dialogTitle = computed(() => {
  if (step.value === 'cfg' && selectedRegion.value) {
    return t('pages.rechargeOrderList.channelTestCfgTitle', {currency: selectedRegion.value.code})
  }
  return t('pages.rechargeOrderList.channelTestCurrencyTitle')
})

const mapRegionItem = (item: HaiPayRegionItem): RegionRow => {
  const code = String(item.currencyCode || item.symbol || '').trim().toUpperCase()
  const preferredCurrency = String(item.fiatCurrencyCode || '').trim().toUpperCase()
  const currencies = (Array.isArray(item.fiatCurrencyCodes) && item.fiatCurrencyCodes.length > 0
      ? item.fiatCurrencyCodes
      : [item.fiatCurrencyCode])
      .map((currency) => String(currency || '').trim().toUpperCase())
      .filter(Boolean)
  return {
    code,
    currency: currencies.length > 1 ? `${preferredCurrency} / ${currencies.join(', ')}` : preferredCurrency,
    name: String(item.name || item.nameEn || code),
    nameEn: String(item.nameEn || item.name || ''),
    nameZh: String(item.nameZh || ''),
    paymentMethods: Array.isArray(item.paymentMethods) ? item.paymentMethods : [],
  }
}

const methodKey = (method: {payType: string; inBankCode: string}) => `${method.payType}\u0000${method.inBankCode}`
const methodLabel = (method: HaiPayRegionItem['paymentMethods'][number]) => {
  const limit = method.minAmount && method.maxAmount ? `${method.minAmount}-${method.maxAmount}` : ''
  return [method.payType, method.inBankCode, method.description, limit].filter(Boolean).join(' · ')
}
const selectedPaymentMethod = computed(() => selectedRegion.value?.paymentMethods.find(
    (method) => methodKey(method) === selectedMethodKey.value,
) || null)

const loadRegionList = async () => {
  regionLoading.value = true
  try {
    const res = await fiatCurrencyApi.haiPayRegionList()
    regionList.value = (res?.list || []).map(mapRegionItem).filter((r) => !!r.code)
  } catch (error) {
    console.error('load haiPay regions failed:', error)
    ElMessage.error(t('pages.rechargeOrderList.loadRechargeCfgFailed'))
    regionList.value = []
  } finally {
    regionLoading.value = false
  }
}

const channelTestRules = computed<FormRules>(() => ({
  userId: [{required: true, message: t('pages.rechargeOrderList.playerIdRequired'), trigger: 'blur'}],
}))

const payerRules = computed<FormRules>(() => ({
  payName: [{required: true, message: t('pages.rechargeOrderList.payNameRequired'), trigger: 'blur'}],
  payEmail: [
    {required: true, message: t('pages.rechargeOrderList.payEmailRequired'), trigger: 'blur'},
    {
      type: 'email',
      message: t('pages.rechargeOrderList.payEmailInvalid'),
      trigger: ['blur', 'change'],
    },
  ],
  payPhone: [{required: true, message: t('pages.rechargeOrderList.payPhoneRequired'), trigger: 'blur'}],
}))

const cfgGold = (row: RechargeCfg) => {
  const anyRow = row as RechargeCfg & {gold?: number; diamond?: number}
  return anyRow.gold ?? anyRow.diamond ?? 0
}

const resetChannelTest = () => {
  channelTestForm.userId = lockedUserId.value
  channelTestCfgList.value = []
  selectedRegion.value = null
  selectedMethodKey.value = ''
  step.value = 'region'
  payerDialogVisible.value = false
  pendingCfg.value = null
  payerForm.payName = ''
  payerForm.payEmail = ''
  payerForm.payPhone = ''
  cachedPayerByUserId.value = {}
  channelTestFormRef.value?.clearValidate()
}

const backToRegion = () => {
  step.value = 'region'
  selectedRegion.value = null
  selectedMethodKey.value = ''
  channelTestCfgList.value = []
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

const loadPayerProfile = async (userId: string) => {
  if (cachedPayerByUserId.value[userId]) {
    return cachedPayerByUserId.value[userId]
  }
  try {
    const res = await rechargeOrderApi.getChannelPayUserProfile({userId})
    const name = String(res?.name || '').trim()
    const email = String(res?.email || '').trim()
    const phone = String(res?.phone || '').trim()
    const profile = {name, email, phone}
    cachedPayerByUserId.value = {...cachedPayerByUserId.value, [userId]: profile}
    return profile
  } catch {
    return {name: '', email: '', phone: ''}
  }
}

const createOrderWithPayer = async (row: RechargeCfg, payName: string, payEmail: string, payPhone: string) => {
  if (!selectedRegion.value?.code) {
    ElMessage.warning(t('pages.rechargeOrderList.currencyRequired'))
    step.value = 'region'
    return
  }
  if (!selectedPaymentMethod.value) {
    ElMessage.warning(t('pages.rechargeOrderList.selectPaymentMethod'))
    return
  }
  channelTestCreating.value = true
  try {
    const userId = channelTestForm.userId.trim()
    const res = await rechargeOrderApi.createChannelRechargeOrderTest({
      userId,
      cfgId: Number(row.id),
      currencyCode: selectedRegion.value.code,
      payName: payName.trim(),
      payEmail: payEmail.trim(),
      payPhone: payPhone.trim(),
      payType: selectedPaymentMethod.value.payType,
      inBankCode: selectedPaymentMethod.value.inBankCode,
    })
    if (!res?.payUrl) {
      ElMessage.error(t('pages.rechargeOrderList.openPayUrlFailed'))
      return
    }
    cachedPayerByUserId.value = {
      ...cachedPayerByUserId.value,
      [userId]: {name: payName.trim(), email: payEmail.trim(), phone: payPhone.trim()},
    }
    ElMessage.success(t('pages.rechargeOrderList.channelTestCreated', {
      orderId: res.orderId,
      price: formatAmount(res.payAmount ?? res.price),
      currency: res.currency,
    }))
    window.open(res.payUrl, '_blank')
    payerDialogVisible.value = false
    emit('update:modelValue', false)
    emit('created')
  } catch (error) {
    console.error('channel recharge test failed:', error)
  } finally {
    channelTestCreating.value = false
  }
}

const handleRegionPick = async (row: RegionRow) => {
  if (!row?.code || channelTestCfgLoading.value) return
  if (!channelTestFormRef.value) return
  try {
    await channelTestFormRef.value.validate()
  } catch {
    return
  }
  selectedRegion.value = row
  const onlyMethod = row.paymentMethods.length === 1 ? row.paymentMethods[0] : undefined
  selectedMethodKey.value = onlyMethod ? methodKey(onlyMethod) : ''
  step.value = 'cfg'
  await loadChannelTestCfgList()
}

const handleChannelTestCfgPick = async (row: RechargeCfg) => {
  if (!row?.id || channelTestCreating.value) return
  if (!selectedRegion.value?.code) {
    ElMessage.warning(t('pages.rechargeOrderList.currencyRequired'))
    step.value = 'region'
    return
  }
  if (!channelTestFormRef.value) return
  try {
    await channelTestFormRef.value.validate()
  } catch {
    return
  }

  const userId = channelTestForm.userId.trim()
  const profile = await loadPayerProfile(userId)
  if (profile.name && profile.email && profile.phone) {
    await createOrderWithPayer(row, profile.name, profile.email, profile.phone)
    return
  }

  pendingCfg.value = row
  payerForm.payName = profile.name
  payerForm.payEmail = profile.email
  payerForm.payPhone = profile.phone
  payerDialogVisible.value = true
  payerFormRef.value?.clearValidate()
}

const confirmPayerAndCreate = async () => {
  if (!pendingCfg.value || !payerFormRef.value) return
  try {
    await payerFormRef.value.validate()
  } catch {
    return
  }
  await createOrderWithPayer(pendingCfg.value, payerForm.payName, payerForm.payEmail, payerForm.payPhone)
}

watch(
    () => props.modelValue,
    (visible) => {
      if (!visible) return
      channelTestForm.userId = lockedUserId.value
      channelTestCfgList.value = []
      selectedRegion.value = null
      selectedMethodKey.value = ''
      step.value = 'region'
      payerDialogVisible.value = false
      pendingCfg.value = null
      cachedPayerByUserId.value = {}
      void loadRegionList()
    },
)
</script>

<style scoped>
.region-en {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.step-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.step-region {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.payer-tip {
  margin: 0 0 16px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}
</style>
