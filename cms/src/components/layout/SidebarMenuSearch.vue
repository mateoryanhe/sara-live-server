<template>
  <div class="sidebar-menu-search">
    <el-autocomplete
        v-model="query"
        clearable
        fit-input-width
        :debounce="120"
        :fetch-suggestions="fetchSuggestions"
        :placeholder="t('common.sidebarMenuSearchPlaceholder')"
        :trigger-on-focus="true"
        value-key="path"
        @select="handleSelect"
    >
      <template #prefix>
        <el-icon>
          <Search/>
        </el-icon>
      </template>
      <template #default="{ item }">
        <div class="sidebar-menu-search__option">
          <div class="sidebar-menu-search__title">{{ item.label }}</div>
          <div v-if="item.group" class="sidebar-menu-search__group">{{ item.group }}</div>
        </div>
      </template>
    </el-autocomplete>
  </div>
</template>

<script lang="ts" setup>
import {computed, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {Search} from '@element-plus/icons-vue'
import {
  filterSidebarMenuSearchItems,
  getSidebarMenuSearchItems,
  type SidebarMenuSearchItem,
} from '@/utils/sidebar-menu-search'

type SearchSuggestion = SidebarMenuSearchItem & { label: string }

const emit = defineEmits<{
  navigate: [path: string]
}>()

const {t, te, locale} = useI18n()
const query = ref('')

const menuItems = computed(() => getSidebarMenuSearchItems())

const resolveLabel = (item: SidebarMenuSearchItem) => {
  if (te(`menu.${item.permission}`)) {
    return t(`menu.${item.permission}`)
  }
  return item.metaTitle
}

const fetchSuggestions = (queryString: string, cb: (results: SearchSuggestion[]) => void) => {
  void locale.value
  const filtered = filterSidebarMenuSearchItems(menuItems.value, queryString, resolveLabel)
  cb(
      filtered.map(item => ({
        ...item,
        label: resolveLabel(item),
      })),
  )
}

const handleSelect = (item: SearchSuggestion) => {
  if (!item?.path) {
    return
  }
  query.value = ''
  emit('navigate', item.path)
}
</script>

<style scoped>
.sidebar-menu-search {
  width: min(360px, 42vw);
  margin-left: 4px;
}

.sidebar-menu-search :deep(.el-input__wrapper) {
  border-radius: 20px;
}

.sidebar-menu-search__option {
  line-height: 1.35;
  padding: 2px 0;
}

.sidebar-menu-search__title {
  font-size: 14px;
  color: #303133;
}

.sidebar-menu-search__group {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}
</style>
