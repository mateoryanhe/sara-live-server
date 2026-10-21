<template>
  <div v-loading="pageLoading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LuckyGiftMultiplierCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="page-hint"
          show-icon
          :title="t('pages.luckyGiftMultiplierCfg.pageHint')"
          type="info"
      />

      <div class="toolbar">
        <el-button type="primary" @click="handleAddTier">{{ t('pages.luckyGiftMultiplierCfg.addTier') }}</el-button>
        <el-button @click="handleResetDefault">{{ t('pages.luckyGiftMultiplierCfg.resetDefault') }}</el-button>
        <el-button :loading="saveLoading" type="success" @click="handleSave">{{ t('pages.luckyGiftMultiplierCfg.saveConfig') }}</el-button>
      </div>

      <el-table :data="tiers" border style="width: 100%; max-width: 880px">
        <el-table-column :label="t('common.sort')" type="index" width="72"/>
        <el-table-column :label="t('pages.luckyGiftMultiplierCfg.multiplier')" min-width="200">
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
        <el-table-column min-width="200">
          <template #header>
            <span>{{ t('pages.luckyGiftMultiplierCfg.weight') }}</span>
            <el-tooltip :content="t('pages.luckyGiftMultiplierCfg.weightHint')" placement="top">
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
        <el-table-column :label="t('pages.luckyGiftMultiplierCfg.multiplierType')" min-width="200">
          <template #default="{ row }">
            <el-select v-model="row.type" style="width: 100%">
              <el-option :label="t('pages.luckyGiftMultiplierCfg.typeLow')" value="low"/>
              <el-option :label="t('pages.luckyGiftMultiplierCfg.typeMedium')" value="medium"/>
              <el-option :label="t('pages.luckyGiftMultiplierCfg.typeHigh')" value="high"/>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="100">
          <template #default="{ $index }">
            <el-button
                link
                type="danger"
                :disabled="tiers.length <= 1"
                @click="handleRemoveTier($index)"
            >
              {{ t('pages.luckyGiftMultiplierCfg.deleteTier') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {giftApi} from '@/api/modules/gift'
import {
  createDefaultLuckyGiftMultiplierTiers,
  newLuckyGiftMultiplierTierKey,
  type LuckyGiftMultiplierTier,
  type LuckyGiftMultiplierType,
} from '@/utils/lucky-gift-multiplier'

const {t} = useI18n()
const tiers = ref<LuckyGiftMultiplierTier[]>(createDefaultLuckyGiftMultiplierTiers())
const pageLoading = ref(false)
const saveLoading = ref(false)

const isValidType = (value: unknown): value is LuckyGiftMultiplierType => {
  return value === 'low' || value === 'medium' || value === 'high'
}

const normalizeTiers = (raw: unknown): LuckyGiftMultiplierTier[] | null => {
  if (!Array.isArray(raw) || raw.length === 0) {
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

const loadFromServer = async () => {
  pageLoading.value = true
  try {
    const res = await giftApi.getLuckyGiftMultiplierCfg()
    const parsed = normalizeTiers(res.tiers)
    if (parsed) {
      tiers.value = parsed
    } else {
      tiers.value = createDefaultLuckyGiftMultiplierTiers()
    }
  } catch (error) {
    console.error('load lucky gift multiplier cfg failed:', error)
    ElMessage.warning(t('pages.luckyGiftMultiplierCfg.loadFailed'))
  } finally {
    pageLoading.value = false
  }
}

const validateTiers = (): boolean => {
  if (tiers.value.length === 0) {
    ElMessage.warning(t('pages.luckyGiftMultiplierCfg.atLeastOneTier'))
    return false
  }
  for (const row of tiers.value) {
    if (!row.multiplier || row.multiplier <= 0) {
      ElMessage.warning(t('pages.luckyGiftMultiplierCfg.multiplierMin'))
      return false
    }
    if (!isValidType(row.type)) {
      ElMessage.warning(t('pages.luckyGiftMultiplierCfg.multiplierType'))
      return false
    }
    if (!row.weight || row.weight <= 0) {
      ElMessage.warning(t('pages.luckyGiftMultiplierCfg.weightMin'))
      return false
    }
  }
  return true
}

const handleAddTier = () => {
  tiers.value.push({
    key: newLuckyGiftMultiplierTierKey(),
    multiplier: 5,
    weight: 1,
    type: 'low',
  })
}

const handleRemoveTier = (index: number) => {
  if (tiers.value.length <= 1) {
    ElMessage.warning(t('pages.luckyGiftMultiplierCfg.atLeastOneTier'))
    return
  }
  tiers.value.splice(index, 1)
}

const handleResetDefault = () => {
  tiers.value = createDefaultLuckyGiftMultiplierTiers()
}

const handleSave = async () => {
  if (!validateTiers()) {
    return
  }
  saveLoading.value = true
  try {
    await giftApi.saveLuckyGiftMultiplierCfg({
      tiers: tiers.value.map((row) => ({
        multiplier: row.multiplier,
        weight: row.weight,
        type: row.type,
      })),
    })
    ElMessage.success(t('pages.luckyGiftMultiplierCfg.saveSuccess'))
    await loadFromServer()
  } catch (error) {
    console.error('save lucky gift multiplier cfg failed:', error)
    ElMessage.error(t('pages.luckyGiftMultiplierCfg.saveFailed'))
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

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
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
