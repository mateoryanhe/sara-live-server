<template>
  <div class="lucky-gift-tier-editor">
    <div class="toolbar">
      <el-button size="small" type="primary" @click="addTier">{{ t('pages.giftList.luckyAddTier') }}</el-button>
      <el-button size="small" @click="resetDefault">{{ t('pages.giftList.luckyResetDefaultTiers') }}</el-button>
    </div>
    <el-table :data="modelValue" border size="small" style="width: 100%">
      <el-table-column :label="t('common.sort')" type="index" width="56"/>
      <el-table-column :label="t('pages.luckyGiftMultiplierCfg.multiplier')" min-width="120">
        <template #default="{ row }">
          <el-input-number v-model="row.multiplier" :min="1" :precision="0" :step="1" controls-position="right"/>
        </template>
      </el-table-column>
      <el-table-column min-width="120">
        <template #header>
          <span>{{ t('pages.luckyGiftMultiplierCfg.weight') }}</span>
          <el-tooltip :content="t('pages.luckyGiftMultiplierCfg.weightHint')" placement="top">
            <span class="weight-hint">?</span>
          </el-tooltip>
        </template>
        <template #default="{ row }">
          <el-input-number v-model="row.weight" :min="1" :precision="0" :step="1" controls-position="right"/>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.luckyGiftMultiplierCfg.multiplierType')" min-width="120">
        <template #default="{ row }">
          <el-select v-model="row.type" style="width: 100%">
            <el-option :label="t('pages.luckyGiftMultiplierCfg.typeLow')" value="low"/>
            <el-option :label="t('pages.luckyGiftMultiplierCfg.typeMedium')" value="medium"/>
            <el-option :label="t('pages.luckyGiftMultiplierCfg.typeHigh')" value="high"/>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column fixed="right" :label="t('common.actions')" width="72">
        <template #default="{ $index }">
          <el-button
              link
              type="danger"
              :disabled="modelValue.length <= 1"
              @click="removeTier($index)"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {
  createDefaultLuckyGiftMultiplierTiers,
  newLuckyGiftMultiplierTierKey,
  type LuckyGiftMultiplierTier,
} from '@/utils/lucky-gift-multiplier'

const props = defineProps<{
  modelValue: LuckyGiftMultiplierTier[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: LuckyGiftMultiplierTier[]]
}>()

const {t} = useI18n()

const patch = (next: LuckyGiftMultiplierTier[]) => {
  emit('update:modelValue', next)
}

const addTier = () => {
  patch([
    ...props.modelValue,
    {key: newLuckyGiftMultiplierTierKey(), multiplier: 5, weight: 1, type: 'low'},
  ])
}

const removeTier = (index: number) => {
  if (props.modelValue.length <= 1) {
    return
  }
  const next = [...props.modelValue]
  next.splice(index, 1)
  patch(next)
}

const resetDefault = () => {
  patch(createDefaultLuckyGiftMultiplierTiers())
}
</script>

<style scoped>
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
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
