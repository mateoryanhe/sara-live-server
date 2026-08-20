<template>
  <el-select
      :model-value="modelValue"
      clearable
      collapse-tags
      collapse-tags-tooltip
      filterable
      :loading="loading"
      multiple
      :placeholder="placeholder"
      popper-class="anchor-remote-select-popper"
      remote
      :remote-method="searchAnchors"
      reserve-keyword
      :style="{width: selectWidth}"
      @update:model-value="handleChange"
      @focus="loadInitial"
  >
    <template #header>
      <div class="anchor-option-row anchor-option-header" :style="rowGridStyle">
        <span v-for="column in columns" :key="column.key" class="anchor-option-cell">{{ column.label }}</span>
      </div>
    </template>
    <el-option
        v-for="item in displayOptions"
        :key="item.id"
        :label="formatSelectedLabel(item)"
        :value="item.id"
    >
      <div class="anchor-option-row" :style="rowGridStyle">
        <template v-for="column in columns" :key="column.key">
          <span v-if="column.key === 'check'" class="anchor-option-cell anchor-option-check">
            <el-checkbox
                :model-value="isSelected(item.id)"
                @change="(checked: boolean) => setSelected(item.id, checked)"
                @click.stop
            />
          </span>
          <span v-else-if="column.key === 'avatar'" class="anchor-option-cell anchor-option-avatar">
            <el-avatar :size="28" :src="item.avatar || undefined">
              {{ avatarFallback(item.nickname) }}
            </el-avatar>
          </span>
          <span
              v-else
              class="anchor-option-cell"
              :title="formatCellValue(item, column.key)"
          >
            {{ formatCellValue(item, column.key) }}
          </span>
        </template>
      </div>
    </el-option>
  </el-select>
</template>

<script lang="ts" setup>
import {computed, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {accountApi} from '@/api'
import type {AnchorListItem} from '@/types/api'

type ColumnKey = 'check' | 'avatar' | 'guildName' | 'nickname' | 'id' | 'phone' | 'userType' | 'liveStatus'

const props = withDefaults(defineProps<{
  modelValue: string[]
  mode: 'platform' | 'guild'
  placeholder?: string
}>(), {
  modelValue: () => [],
  placeholder: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const {t} = useI18n()
const loading = ref(false)
const searchResults = ref<AnchorListItem[]>([])
const optionCache = ref<Map<string, AnchorListItem>>(new Map())

const selectWidth = computed(() => props.mode === 'guild' ? '360px' : '340px')

const columns = computed<Array<{ key: ColumnKey, label: string, width: string }>>(() => {
  if (props.mode === 'guild') {
    return [
      {key: 'check', label: '', width: '32px'},
      {key: 'avatar', label: t('common.avatar'), width: '40px'},
      {key: 'guildName', label: t('pages.guildList.guildName'), width: '1.3fr'},
      {key: 'nickname', label: t('common.nickname'), width: '1fr'},
      {key: 'id', label: t('common.userId'), width: '1fr'},
      {key: 'phone', label: t('common.phone'), width: '1fr'},
      {key: 'liveStatus', label: t('pages.anchorList.liveStatus'), width: '0.8fr'},
    ]
  }
  return [
    {key: 'check', label: '', width: '32px'},
    {key: 'avatar', label: t('common.avatar'), width: '40px'},
    {key: 'id', label: t('common.userId'), width: '1fr'},
    {key: 'nickname', label: t('common.nickname'), width: '1fr'},
    {key: 'phone', label: t('common.phone'), width: '1fr'},
    {key: 'userType', label: t('pages.anchorList.anchorType'), width: '0.9fr'},
    {key: 'liveStatus', label: t('pages.anchorList.liveStatus'), width: '0.8fr'},
  ]
})

const rowGridStyle = computed(() => ({
  gridTemplateColumns: columns.value.map(column => column.width).join(' '),
}))

const formatSelectedLabel = (item: AnchorListItem) => {
  const nickname = item.nickname || '-'
  if (props.mode === 'guild' && item.guildName) {
    return `[${item.guildName}] ${nickname} (${item.id})`
  }
  return `${nickname} (${item.id})`
}

const isSelected = (id: string) => props.modelValue.includes(id)

const setSelected = (id: string, checked: boolean) => {
  if (checked) {
    if (isSelected(id)) {
      return
    }
    emit('update:modelValue', [...props.modelValue, id])
    return
  }
  emit('update:modelValue', props.modelValue.filter(itemId => itemId !== id))
}

const avatarFallback = (nickname?: string) => {
  const text = (nickname || '?').trim()
  return text ? text.slice(0, 1).toUpperCase() : '?'
}

const formatUserType = (userType?: number) => {
  if (userType === 7) {
    return t('pages.anchorList.anchorTypeSenior')
  }
  if (userType === 1) {
    return t('pages.anchorList.anchorTypeNormal')
  }
  return '-'
}

const formatLiveStatus = (liveStatus?: number) => {
  return liveStatus === 1 ? t('common.live') : t('common.offline')
}

const formatCellValue = (item: AnchorListItem, key: ColumnKey) => {
  switch (key) {
    case 'guildName':
      return item.guildName || '-'
    case 'nickname':
      return item.nickname || '-'
    case 'id':
      return item.id || '-'
    case 'phone':
      return item.phone || '-'
    case 'userType':
      return formatUserType(item.userType)
    case 'liveStatus':
      return formatLiveStatus(item.liveStatus)
    default:
      return '-'
  }
}

const cacheAnchorItems = (items: AnchorListItem[]) => {
  for (const item of items) {
    optionCache.value.set(item.id, item)
  }
}

const displayOptions = computed(() => {
  const merged = new Map<string, AnchorListItem>()
  for (const id of props.modelValue) {
    const cached = optionCache.value.get(id)
    if (cached) {
      merged.set(id, cached)
    }
  }
  for (const item of searchResults.value) {
    merged.set(item.id, item)
  }
  return Array.from(merged.values())
})

const searchAnchors = async (query: string) => {
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: 1,
      pageSize: 50,
      key: query.trim(),
      ...(props.mode === 'platform'
        ? {platformOnly: true}
        : {guildOnly: true}),
    })
    searchResults.value = response.data || []
    cacheAnchorItems(searchResults.value)
  } catch (error) {
    console.error('Failed to search anchors:', error)
    searchResults.value = []
  } finally {
    loading.value = false
  }
}

const loadInitial = () => {
  if (searchResults.value.length === 0) {
    void searchAnchors('')
  }
}

const handleChange = (value: string[] | null | undefined) => {
  emit('update:modelValue', value ? [...value] : [])
}

watch(
  () => props.modelValue,
  (ids) => {
    const missingIds = ids.filter(id => !optionCache.value.has(id))
    if (missingIds.length === 0) {
      return
    }
    void (async () => {
      for (const id of missingIds) {
        try {
          const response = await accountApi.getAnchorList({
            pageIndex: 1,
            pageSize: 1,
            key: id,
            ...(props.mode === 'platform'
              ? {platformOnly: true}
              : {guildOnly: true}),
          })
          const item = (response.data || []).find(row => row.id === id)
          if (item) {
            optionCache.value.set(item.id, item)
          }
        } catch (error) {
          console.error('Failed to load selected anchor:', error)
        }
      }
    })()
  },
  {immediate: true},
)
</script>

<style>
.anchor-remote-select-popper.el-popper {
  min-width: 800px !important;
}

.anchor-remote-select-popper.el-select-dropdown.is-multiple .el-select-dropdown__item {
  padding-right: 12px;
}

.anchor-remote-select-popper.el-select-dropdown.is-multiple .el-select-dropdown__item::after {
  display: none;
}

.anchor-remote-select-popper .el-select-dropdown__header {
  padding: 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.anchor-remote-select-popper .el-select-dropdown__item {
  height: auto;
  line-height: 1.4;
  padding: 8px 12px;
}

.anchor-remote-select-popper .anchor-option-row {
  display: grid;
  gap: 8px;
  align-items: center;
  width: 100%;
}

.anchor-remote-select-popper .anchor-option-header {
  padding: 8px 12px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  font-weight: 600;
  background: var(--el-fill-color-light);
}

.anchor-remote-select-popper .anchor-option-cell {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.anchor-remote-select-popper .anchor-option-avatar {
  display: flex;
  justify-content: center;
  overflow: visible;
}

.anchor-remote-select-popper .anchor-option-check {
  display: flex;
  justify-content: center;
  overflow: visible;
}

.anchor-remote-select-popper .anchor-option-check .el-checkbox {
  height: auto;
  margin-right: 0;
}
</style>
