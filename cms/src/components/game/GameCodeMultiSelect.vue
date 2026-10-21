<template>
  <div class="game-code-multi-select">
    <div class="game-code-multi-select__tags">
      <el-tag
          v-for="code in modelValue"
          :key="code"
          closable
          size="small"
          @close="removeCode(code)"
      >
        {{ displayLabel(code) }}
      </el-tag>
      <span v-if="!modelValue.length" class="game-code-multi-select__empty">{{ t('pages.gameList.gameCodeSelectEmpty') }}</span>
    </div>
    <el-button type="primary" link @click="openDialog">{{ t('pages.gameList.selectGames') }}</el-button>

    <el-dialog
        v-model="dialogVisible"
        :title="t('pages.gameList.selectGamesDialogTitle')"
        destroy-on-close
        width="920px"
        @closed="handleDialogClosed"
        @open="handleDialogOpen"
    >
      <el-form :model="dialogSearch" class="search-form" inline>
        <el-form-item :label="t('pages.gameList.gameName')">
          <el-input v-model="dialogSearch.name" clearable :placeholder="t('pages.gameList.gameNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameList.platform')">
          <el-input v-model="dialogSearch.platform" clearable :placeholder="t('pages.gameList.platformPlaceholder')"/>
        </el-form-item>
        <el-form-item v-if="listSource === 'vendor'" :label="t('common.category')">
          <el-input v-model="dialogSearch.category" clearable :placeholder="t('pages.gameList.categoryPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleDialogSearch">{{ t('common.search') }}</el-button>
          <el-button @click="resetDialogSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          ref="tableRef"
          v-loading="dialogLoading"
          :data="dialogTableData"
          :row-key="rowKey"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column reserve-selection type="selection" width="48"/>
        <el-table-column :label="t('pages.gameList.gameCode')" min-width="120" prop="gameCode"/>
        <el-table-column :label="t('common.name')" min-width="140">
          <template #default="{ row }">{{ row.name || row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.nameEn')" min-width="120" prop="nameEn">
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.platform')" prop="platform" width="100"/>
      </el-table>

      <div class="pagination">
        <el-pagination
            v-model:current-page="dialogPagination.pageIndex"
            v-model:page-size="dialogPagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="dialogPagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchDialogList"
            @size-change="handleDialogSizeChange"
        />
      </div>

      <template #footer>
        <span class="dialog-selected-tip">{{ t('common.selectedCount', {count: dialogSelection.length}) }}</span>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmSelection">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {nextTick, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import type {TableInstance} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import type {GameShelfItem, VendorGame} from '@/types/api'

type ListSource = 'shelf' | 'vendor'
type PickerRow = GameShelfItem | VendorGame

const props = withDefaults(
    defineProps<{
      modelValue: string[]
      listSource?: ListSource
      nameMap?: Record<string, string>
    }>(),
    {
      listSource: 'shelf',
      nameMap: () => ({}),
    },
)

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const {t} = useI18n()
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const dialogTableData = ref<PickerRow[]>([])
const dialogSelection = ref<PickerRow[]>([])
const tableRef = ref<TableInstance>()
const rowLabelMap = ref<Record<string, string>>({})

const dialogSearch = reactive({
  name: '',
  platform: '',
  category: '',
})

const dialogPagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const rowKey = (row: PickerRow) => row.gameCode

const displayLabel = (code: string) => {
  const fromProp = props.nameMap[code]
  if (fromProp) {
    return `${fromProp} (${code})`
  }
  const fromDialog = rowLabelMap.value[code]
  if (fromDialog) {
    return `${fromDialog} (${code})`
  }
  return code
}

const removeCode = (code: string) => {
  emit('update:modelValue', props.modelValue.filter(item => item !== code))
}

const openDialog = () => {
  dialogVisible.value = true
}

const buildListParams = () => {
  const base = {
    name: dialogSearch.name.trim(),
    platform: dialogSearch.platform.trim(),
    pageIndex: dialogPagination.pageIndex,
    pageSize: dialogPagination.pageSize,
  }
  if (props.listSource === 'vendor') {
    return {
      ...base,
      category: dialogSearch.category.trim(),
    }
  }
  return base
}

const fetchDialogList = async () => {
  dialogLoading.value = true
  try {
    const params = buildListParams()
    const response =
        props.listSource === 'vendor'
            ? await gamePlatformApi.getVendorGameList(params)
            : await gamePlatformApi.getGameShelfList(params)
    dialogTableData.value = response.data || []
    dialogPagination.total = response.total || 0
    await nextTick()
    syncTableSelection()
  } catch (error) {
    console.error('fetch game picker list failed:', error)
    dialogTableData.value = []
    dialogPagination.total = 0
  } finally {
    dialogLoading.value = false
  }
}

const syncTableSelection = () => {
  const table = tableRef.value
  if (!table) {
    return
  }
  table.clearSelection()
  const selectedSet = new Set(props.modelValue)
  for (const row of dialogTableData.value) {
    if (selectedSet.has(row.gameCode)) {
      table.toggleRowSelection(row, true)
    }
  }
}

const handleSelectionChange = (rows: PickerRow[]) => {
  dialogSelection.value = rows
  for (const row of rows) {
    const label = row.name || row.nameEn || ''
    if (label) {
      rowLabelMap.value[row.gameCode] = label
    }
  }
}

const handleDialogSearch = () => {
  dialogPagination.pageIndex = 1
  fetchDialogList()
}

const resetDialogSearch = () => {
  dialogSearch.name = ''
  dialogSearch.platform = ''
  dialogSearch.category = ''
  dialogPagination.pageIndex = 1
  fetchDialogList()
}

const handleDialogSizeChange = () => {
  dialogPagination.pageIndex = 1
  fetchDialogList()
}

const handleDialogOpen = () => {
  dialogSelection.value = []
  dialogPagination.pageIndex = 1
  fetchDialogList()
}

const handleDialogClosed = () => {
  dialogSearch.name = ''
  dialogSearch.platform = ''
  dialogSearch.category = ''
  dialogSelection.value = []
}

const confirmSelection = () => {
  const codes = new Set(props.modelValue)
  for (const row of dialogSelection.value) {
    if (row.gameCode) {
      codes.add(row.gameCode)
      const label = row.name || row.nameEn || ''
      if (label) {
        rowLabelMap.value[row.gameCode] = label
      }
    }
  }
  emit('update:modelValue', Array.from(codes))
  dialogVisible.value = false
}
</script>

<style scoped>
.game-code-multi-select {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 8px;
  max-width: 520px;
}

.game-code-multi-select__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
  min-width: 200px;
}

.game-code-multi-select__empty {
  color: var(--el-text-color-placeholder);
  font-size: 13px;
  line-height: 24px;
}

.search-form {
  margin-bottom: 12px;
}

.pagination {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

.dialog-selected-tip {
  margin-right: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
