<template>
  <div v-loading="pageLoading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LuckyGiftPoolCfgManagement') }}</span>
          <el-button :loading="pageLoading" @click="handleRefresh">
            {{ t('pages.luckyGiftPoolCfg.refresh') }}
          </el-button>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="page-hint"
          show-icon
          :title="t('pages.luckyGiftPoolCfg.pageHint')"
          type="info"
      />

      <el-form :model="form" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.luckyGiftPoolCfg.poolBalance')">
          <el-input-number
              v-model="form.poolBalance"
              :min="0"
              :precision="0"
              :step="100000"
              controls-position="right"
              style="width: 240px"
          />
          <el-button class="refresh-inline" :loading="pageLoading" @click="handleRefresh">
            {{ t('pages.luckyGiftPoolCfg.refreshPool') }}
          </el-button>
          <span class="field-hint">{{ t('pages.luckyGiftPoolCfg.poolBalanceHint') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.luckyGiftPoolCfg.contributePercent')">
          <el-input-number
              v-model="form.contributePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="0.1"
              controls-position="right"
              style="width: 240px"
          />
          <span class="field-unit">{{ t('pages.luckyGiftPoolCfg.percentUnit') }}</span>
          <span class="field-hint">{{ t('pages.luckyGiftPoolCfg.contributePercentHint') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.luckyGiftPoolCfg.winProbability')">
          <el-input-number
              v-model="form.winProbability"
              :max="100"
              :min="0"
              :precision="2"
              :step="0.1"
              controls-position="right"
              style="width: 240px"
          />
          <span class="field-unit">{{ t('pages.luckyGiftPoolCfg.percentUnit') }}</span>
          <span class="field-hint">{{ t('pages.luckyGiftPoolCfg.winProbabilityHint') }}</span>
        </el-form-item>

        <el-form-item v-if="form.updatedAt" :label="t('pages.luckyGiftPoolCfg.lastUpdated')">
          <span>{{ form.updatedAt }}</span>
        </el-form-item>
      </el-form>

      <div class="tiers-header">
        <span class="tiers-title">{{ t('pages.luckyGiftPoolCfg.tiersTitle') }}</span>
        <div class="toolbar">
          <el-button @click="handleResetDefault">{{ t('pages.luckyGiftPoolCfg.resetDefault') }}</el-button>
          <el-button v-if="canSave" :loading="saveLoading" type="success" @click="handleSave">
            {{ t('pages.luckyGiftPoolCfg.saveConfig') }}
          </el-button>
        </div>
      </div>

      <el-table :data="tiers" border style="width: 100%; max-width: 960px">
        <el-table-column :label="t('common.sort')" type="index" width="72"/>
        <el-table-column :label="t('pages.luckyGiftPoolCfg.multiplier')" min-width="180">
          <template #default="{ row }">
            <el-input-number
                v-model="row.multiplier"
                :min="1"
                :precision="0"
                :step="1"
                controls-position="right"
            />
          </template>
        </el-table-column>
        <el-table-column min-width="180">
          <template #header>
            <span>{{ t('pages.luckyGiftPoolCfg.weight') }}</span>
            <el-tooltip :content="t('pages.luckyGiftPoolCfg.weightHint')" placement="top">
              <span class="weight-hint">?</span>
            </el-tooltip>
          </template>
          <template #default="{ row }">
            <el-input-number
                v-model="row.weight"
                :min="1"
                :precision="0"
                :step="1"
                controls-position="right"
            />
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.luckyGiftPoolCfg.weightShare')" min-width="120">
          <template #default="{ row }">
            {{ luckyGiftTierWeightSharePercent(tiers, row.weight).toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.luckyGiftPoolCfg.multiplierType')" min-width="160">
          <template #default="{ row }">
            <el-select v-model="row.type" style="width: 100%">
              <el-option :label="t('pages.luckyGiftPoolCfg.typeLow')" value="low"/>
              <el-option :label="t('pages.luckyGiftPoolCfg.typeMedium')" value="medium"/>
              <el-option :label="t('pages.luckyGiftPoolCfg.typeHigh')" value="high"/>
            </el-select>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {luckyGiftPoolCfgApi} from '@/api/modules/lucky-gift-pool-cfg'
import {hasButtonPermission} from '@/utils/permission'
import {
  createDefaultLuckyGiftMultiplierTiers,
  luckyGiftTierWeightSharePercent,
  newLuckyGiftMultiplierTierKey,
  type LuckyGiftMultiplierTier,
  type LuckyGiftMultiplierType,
} from '@/utils/lucky-gift-multiplier'

const {t} = useI18n()
const REQUIRED_TIER_COUNT = 5
const DEFAULT_POOL_BALANCE = 10_000_000

const pageLoading = ref(false)
const saveLoading = ref(false)
const canSave = computed(() => hasButtonPermission('LuckyGiftPoolCfgManagement', 'save'))

const form = reactive({
  id: '' as string,
  poolBalance: DEFAULT_POOL_BALANCE,
  contributePercent: 0,
  winProbability: 0,
  updatedAt: '',
})

const tiers = ref<LuckyGiftMultiplierTier[]>(createDefaultLuckyGiftMultiplierTiers())

const isValidType = (value: unknown): value is LuckyGiftMultiplierType => {
  return value === 'low' || value === 'medium' || value === 'high'
}

const normalizeTiers = (raw: unknown): LuckyGiftMultiplierTier[] | null => {
  if (!Array.isArray(raw) || raw.length !== REQUIRED_TIER_COUNT) {
    return null
  }
  const list: LuckyGiftMultiplierTier[] = []
  for (const item of raw) {
    if (!item || typeof item !== 'object') {
      return null
    }
    const row = item as Record<string, unknown>
    const multiplier = Number(row.multiplier)
    if (!Number.isFinite(multiplier) || multiplier <= 0) {
      return null
    }
    if (!isValidType(row.type)) {
      return null
    }
    const weightRaw = row.weight
    const weight = weightRaw === undefined || weightRaw === null ? 1 : Number(weightRaw)
    if (!Number.isFinite(weight) || weight <= 0) {
      return null
    }
    const id = typeof row.id === 'string' && row.id ? row.id : undefined
    list.push({
      key: id || newLuckyGiftMultiplierTierKey(),
      multiplier: Math.floor(multiplier),
      weight: Math.floor(weight),
      type: row.type,
    })
  }
  return list
}

const loadFromServer = async (): Promise<boolean> => {
  pageLoading.value = true
  try {
    const res = await luckyGiftPoolCfgApi.getLuckyGiftPoolCfg()
    const cfg = res.cfg
    if (cfg) {
      form.id = cfg.id || ''
      form.poolBalance = Number.isFinite(cfg.poolBalance) ? cfg.poolBalance : DEFAULT_POOL_BALANCE
      form.contributePercent = Number.isFinite(cfg.contributePercent) ? cfg.contributePercent : 0
      form.winProbability = Number.isFinite(cfg.winProbability) ? cfg.winProbability : 0
      form.updatedAt = cfg.updatedAt || ''
      const parsed = normalizeTiers(cfg.tiers)
      tiers.value = parsed || createDefaultLuckyGiftMultiplierTiers()
    } else {
      form.id = ''
      form.poolBalance = DEFAULT_POOL_BALANCE
      form.contributePercent = 0
      form.winProbability = 0
      form.updatedAt = ''
      tiers.value = createDefaultLuckyGiftMultiplierTiers()
    }
    return true
  } catch (error) {
    console.error('load lucky gift pool cfg failed:', error)
    ElMessage.warning(t('pages.luckyGiftPoolCfg.loadFailed'))
    return false
  } finally {
    pageLoading.value = false
  }
}

const validateForm = (): boolean => {
  if (!Number.isFinite(form.poolBalance) || form.poolBalance < 0) {
    ElMessage.warning(t('pages.luckyGiftPoolCfg.poolBalanceMin'))
    return false
  }
  if (
      !Number.isFinite(form.contributePercent) ||
      form.contributePercent < 0 ||
      form.contributePercent > 100
  ) {
    ElMessage.warning(t('pages.luckyGiftPoolCfg.percentRange'))
    return false
  }
  if (
      !Number.isFinite(form.winProbability) ||
      form.winProbability < 0 ||
      form.winProbability > 100
  ) {
    ElMessage.warning(t('pages.luckyGiftPoolCfg.percentRange'))
    return false
  }
  if (tiers.value.length !== REQUIRED_TIER_COUNT) {
    ElMessage.warning(t('pages.luckyGiftPoolCfg.needFiveTiers'))
    return false
  }
  for (const row of tiers.value) {
    if (!row.multiplier || row.multiplier <= 0) {
      ElMessage.warning(t('pages.luckyGiftPoolCfg.multiplierMin'))
      return false
    }
    if (!row.weight || row.weight <= 0) {
      ElMessage.warning(t('pages.luckyGiftPoolCfg.weightMin'))
      return false
    }
    if (!isValidType(row.type)) {
      ElMessage.warning(t('pages.luckyGiftPoolCfg.multiplierType'))
      return false
    }
  }
  return true
}

const handleRefresh = async () => {
  if (await loadFromServer()) {
    ElMessage.success(t('pages.luckyGiftPoolCfg.refreshSuccess'))
  }
}

const handleResetDefault = () => {
  tiers.value = createDefaultLuckyGiftMultiplierTiers()
}

const handleSave = async () => {
  if (!validateForm()) {
    return
  }
  saveLoading.value = true
  try {
    await luckyGiftPoolCfgApi.saveLuckyGiftPoolCfg({
      // snowflake id 必须用字符串，Number() 会丢精度导致服务端校验失败
      id: form.id || undefined,
      poolBalance: form.poolBalance,
      contributePercent: form.contributePercent,
      winProbability: form.winProbability,
      tiers: tiers.value.map((row) => ({
        multiplier: row.multiplier,
        weight: row.weight,
        type: row.type,
      })),
    })
    ElMessage.success(t('pages.luckyGiftPoolCfg.saveSuccess'))
    await loadFromServer()
  } catch (error) {
    console.error('save lucky gift pool cfg failed:', error)
    ElMessage.error(t('pages.luckyGiftPoolCfg.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => {
  loadFromServer()
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
}

.page-hint {
  margin-bottom: 16px;
}

.cfg-form {
  max-width: 900px;
  margin-bottom: 8px;
}

.field-unit {
  margin-left: 8px;
  color: #606266;
}

.refresh-inline {
  margin-left: 8px;
}

.field-hint {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}

.tiers-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin: 12px 0 16px;
}

.tiers-title {
  font-weight: 600;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.weight-hint {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-left: 4px;
  border-radius: 50%;
  background: #dcdfe6;
  color: #606266;
  font-size: 12px;
  cursor: help;
}
</style>
