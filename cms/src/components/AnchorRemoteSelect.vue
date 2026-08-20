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
      remote
      :remote-method="searchAnchors"
      reserve-keyword
      style="width: 420px"
      @update:model-value="handleChange"
      @focus="loadInitial"
  >
    <el-option
        v-for="item in displayOptions"
        :key="item.id"
        :label="formatAnchorOption(item)"
        :value="item.id"
    />
  </el-select>
</template>

<script lang="ts" setup>
import {computed, ref, watch} from 'vue'
import {accountApi} from '@/api'
import type {AnchorListItem} from '@/types/api'

const props = withDefaults(defineProps<{
  modelValue: string[]
  placeholder?: string
}>(), {
  modelValue: () => [],
  placeholder: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const loading = ref(false)
const searchResults = ref<AnchorListItem[]>([])
const optionCache = ref<Map<string, AnchorListItem>>(new Map())

const hasGuild = (item: AnchorListItem) => {
  const guildId = item.guildId
  return guildId !== undefined && guildId !== null && String(guildId) !== '0' && String(guildId) !== ''
}

const formatAnchorOption = (item: AnchorListItem) => {
  const nickname = item.nickname || '-'
  const phone = item.phone ? ` · ${item.phone}` : ''
  if (hasGuild(item)) {
    const guildName = item.guildName || '-'
    return `[${guildName}] ${nickname} (${item.id})${phone}`
  }
  return `${nickname} (${item.id})${phone}`
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
