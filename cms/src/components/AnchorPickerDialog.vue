<template>
  <el-dialog
      :model-value="visible"
      :close-on-click-modal="false"
      destroy-on-close
      :title="dialogTitle"
      width="860px"
      @close="handleClose"
  >
    <el-form inline @submit.prevent>
      <el-form-item :label="t('common.keyword')">
        <el-input
            v-model="keyword"
            clearable
            :placeholder="searchPlaceholder"
            style="width: 280px"
            @keyup.enter="handleSearch"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <div v-if="multiple" class="selected-summary">
      {{ selectedCountText }}
    </div>

    <el-table
        ref="tableRef"
        v-loading="loading"
        :data="tableData"
        :highlight-current-row="!multiple"
        row-key="id"
        style="width: 100%"
        @current-change="handleCurrentChange"
        @row-dblclick="handleRowDblClick"
        @selection-change="handleSelectionChange"
    >
      <el-table-column v-if="multiple" reserve-selection type="selection" width="48"/>
      <el-table-column :label="t('common.userId')" min-width="170" prop="id"/>
      <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname" show-overflow-tooltip>
        <template #default="{ row }">{{ row.nickname || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.avatar')" width="80">
        <template #default="{ row }">
          <el-image
              v-if="row.avatar"
              :src="row.avatar"
              fit="cover"
              style="width:40px;height:40px;border-radius:50%"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column v-if="platformOnly" :label="t('pages.anchorList.anchorType')" min-width="100">
        <template #default="{ row }">{{ formatUserType(row.userType) }}</template>
      </el-table-column>
      <el-table-column v-if="platformOnly" :label="t('pages.anchorList.liveStatus')" width="90">
        <template #default="{ row }">{{ formatLiveStatus(row.liveStatus) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.phone')" min-width="130" prop="phone">
        <template #default="{ row }">{{ row.phone || '-' }}</template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
          v-model:current-page="pagination.pageIndex"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchList"
          @size-change="handleSizeChange"
      />
    </div>

    <template #footer>
      <el-button @click="handleClear">{{ t('pages.guildAnchorDailyLiveList.clearAnchor') }}</el-button>
      <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleConfirm">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {computed, nextTick, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import type {TableInstance} from 'element-plus'
import {accountApi, guildApi} from '@/api'
import type {AnchorListItem} from '@/types/api'

const props = withDefaults(defineProps<{
  visible: boolean
  guildId?: string
  ownedGuildOnly?: boolean
  platformOnly?: boolean
  multiple?: boolean
  initialAnchor?: AnchorListItem | null
  initialAnchors?: AnchorListItem[]
}>(), {
  multiple: false,
  initialAnchors: () => [],
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: [anchor: AnchorListItem | null]
  confirmMultiple: [anchors: AnchorListItem[]]
}>()

const {t} = useI18n()
const loading = ref(false)
const keyword = ref('')
const tableData = ref<AnchorListItem[]>([])
const tableRef = ref<TableInstance>()
const selectedAnchor = ref<AnchorListItem | null>(null)
const selectedMap = ref(new Map<string, AnchorListItem>())
const selectedCount = computed(() => selectedMap.value.size)

const dialogTitle = computed(() => {
  if (props.platformOnly) {
    return props.multiple
        ? t('pages.liveRecordList.selectPlatformAnchorsTitle')
        : t('pages.liveRecordList.selectPlatformAnchorTitle')
  }
  return props.multiple
      ? t('pages.guildAnchorDailyLiveList.selectAnchorsTitle')
      : t('pages.guildAnchorDailyLiveList.selectAnchorTitle')
})

const searchPlaceholder = computed(() => {
  if (props.platformOnly) {
    return t('pages.liveRecordList.searchPlatformAnchor')
  }
  return t('pages.guildAnchorDailyLiveList.anchorSearchPlaceholder')
})

const selectedCountText = computed(() => {
  if (props.platformOnly) {
    return t('pages.liveRecordList.selectedPlatformAnchorsCount', {count: selectedCount.value})
  }
  return t('pages.guildAnchorDailyLiveList.selectedAnchorsCount', {count: selectedCount.value})
})

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
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const anchorKey = (anchor: AnchorListItem) => String(anchor.id)

const resetSelectedMap = (anchors: AnchorListItem[]) => {
  const map = new Map<string, AnchorListItem>()
  anchors.forEach(anchor => {
    map.set(anchorKey(anchor), anchor)
  })
  selectedMap.value = map
}

const syncTableSelection = async () => {
  if (!props.multiple || !tableRef.value) {
    return
  }
  await nextTick()
  tableRef.value.clearSelection()
  tableData.value.forEach(row => {
    if (selectedMap.value.has(anchorKey(row))) {
      tableRef.value?.toggleRowSelection(row, true)
    }
  })
}

const fetchList = async () => {
  if (!props.platformOnly && !props.ownedGuildOnly && !props.guildId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = props.ownedGuildOnly
        ? await guildApi.getMyOwnedGuildAnchorList({
          pageIndex: pagination.pageIndex,
          pageSize: pagination.pageSize,
          key: keyword.value.trim() || undefined,
        })
        : await accountApi.getAnchorList({
          pageIndex: pagination.pageIndex,
          pageSize: pagination.pageSize,
          guildId: props.guildId,
          key: keyword.value.trim() || undefined,
          ...(props.platformOnly ? {platformOnly: true} : {}),
        })
    tableData.value = response.data || []
    pagination.total = response.total || 0
    await syncTableSelection()
  } catch (error) {
    console.error('fetch anchors for picker failed:', error)
    tableData.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  keyword.value = ''
  pagination.pageIndex = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const handleCurrentChange = (row: AnchorListItem | null | undefined) => {
  if (props.multiple) {
    return
  }
  selectedAnchor.value = row ?? null
}

const handleSelectionChange = (rows: AnchorListItem[]) => {
  if (!props.multiple) {
    return
  }
  const currentPageIds = new Set(tableData.value.map(row => anchorKey(row)))
  currentPageIds.forEach(id => {
    if (!rows.some(row => anchorKey(row) === id)) {
      selectedMap.value.delete(id)
    }
  })
  rows.forEach(row => {
    selectedMap.value.set(anchorKey(row), row)
  })
}

const handleRowDblClick = (row: AnchorListItem) => {
  if (props.multiple) {
    return
  }
  selectedAnchor.value = row
  handleConfirm()
}

const handleClear = () => {
  if (props.multiple) {
    selectedMap.value = new Map()
    tableRef.value?.clearSelection()
    emit('confirmMultiple', [])
  } else {
    selectedAnchor.value = null
    emit('confirm', null)
  }
  emit('update:visible', false)
}

const handleConfirm = () => {
  if (props.multiple) {
    emit('confirmMultiple', [...selectedMap.value.values()])
  } else {
    emit('confirm', selectedAnchor.value)
  }
  emit('update:visible', false)
}

const handleClose = () => {
  emit('update:visible', false)
}

watch(
  () => props.visible,
  (visible) => {
    if (!visible) {
      return
    }
    keyword.value = ''
    pagination.pageIndex = 1
    if (props.multiple) {
      resetSelectedMap(props.initialAnchors || [])
      selectedAnchor.value = null
    } else {
      selectedAnchor.value = props.initialAnchor ?? null
      selectedMap.value = new Map()
    }
    fetchList()
  },
)
</script>

<style scoped>
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.selected-summary {
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
