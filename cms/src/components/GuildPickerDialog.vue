<template>
  <el-dialog
      :model-value="visible"
      :close-on-click-modal="false"
      destroy-on-close
      :title="dialogTitle"
      width="1100px"
      @close="handleClose"
  >
    <p v-if="tip" class="picker-tip">{{ tip }}</p>

    <el-form inline @submit.prevent>
      <el-form-item :label="t('pages.guildList.guildName')">
        <el-input
            v-model="keyword"
            clearable
            :placeholder="t('pages.guildList.guildNameSearchPlaceholder')"
            style="width: 280px"
            @keyup.enter="handleSearch"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <div class="selected-summary">
      {{ t('common.selectedCount', {count: selectedCount}) }}
    </div>

    <el-table
        ref="tableRef"
        v-loading="loading"
        :data="tableData"
        row-key="id"
        style="width: 100%"
        @selection-change="handleSelectionChange"
    >
      <el-table-column reserve-selection type="selection" width="48"/>
      <el-table-column label="ID" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.guildList.guildName')" min-width="140" prop="name" show-overflow-tooltip/>
      <el-table-column :label="t('pages.guildList.guildType')" min-width="120">
        <template #default="{ row }">{{ formatGuildType(row.guildType) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildList.sharePercent')" width="110" align="right">
        <template #default="{ row }">
          {{ row.sharePercent == null ? '-' : row.sharePercent }}
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.guildList.unsettledTotalIncome')" min-width="140" align="right">
        <template #default="{ row }">
          <span class="money-amount">{{ formatWalletBalance(row.unsettledTotalIncome) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.guildList.leader')" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">{{ formatPerson(row.leaderName, row.leaderId) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildList.creator')" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">{{ formatPerson(row.creatorName, row.creatorId) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildList.description')" min-width="180" prop="description" show-overflow-tooltip>
        <template #default="{ row }">{{ row.description || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.updatedAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
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
      <el-button @click="handleClear">{{ t('common.reset') }}</el-button>
      <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleConfirm">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {computed, nextTick, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type TableInstance} from 'element-plus'
import {guildApi} from '@/api'
import type {Guild} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const GUILD_TYPE_COIN_MERCHANT = 1

const props = withDefaults(defineProps<{
  visible: boolean
  title?: string
  tip?: string
  initialGuilds?: Guild[]
}>(), {
  initialGuilds: () => [],
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: [guilds: Guild[]]
}>()

const {t} = useI18n()
const loading = ref(false)
const keyword = ref('')
const tableData = ref<Guild[]>([])
const tableRef = ref<TableInstance>()
const selectedMap = ref(new Map<string, Guild>())
const selectedCount = computed(() => selectedMap.value.size)
const syncingSelection = ref(false)

const dialogTitle = computed(() => props.title || t('pages.liveDailyEffectiveLiveList.selectGuild'))

const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const guildKey = (guild: Guild) => String(guild.id)

const formatPerson = (name?: string, id?: string) => {
  if (name && id) return `${name} (${id})`
  if (name) return name
  if (id) return String(id)
  return '-'
}

const formatGuildType = (guildType?: number) => {
  if (guildType === GUILD_TYPE_COIN_MERCHANT) {
    return t('pages.guildList.guildTypeCoinMerchant')
  }
  return t('pages.guildList.guildTypeNormal')
}

const resetSelectedMap = (guilds: Guild[]) => {
  const map = new Map<string, Guild>()
  guilds.forEach(guild => {
    map.set(guildKey(guild), guild)
  })
  selectedMap.value = map
}

const syncTableSelection = async () => {
  if (!tableRef.value) {
    return
  }
  syncingSelection.value = true
  await nextTick()
  tableRef.value.clearSelection()
  tableData.value.forEach(row => {
    if (selectedMap.value.has(guildKey(row))) {
      tableRef.value?.toggleRowSelection(row, true)
    }
  })
  await nextTick()
  syncingSelection.value = false
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getGuildList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      name: keyword.value.trim(),
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
    await syncTableSelection()
  } catch (error) {
    console.error('Failed to load guild picker list:', error)
    ElMessage.error(t('pages.liveDailyEffectiveLiveList.fetchGuildFailed'))
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

const handleSelectionChange = (rows: Guild[]) => {
  if (syncingSelection.value) {
    return
  }
  const pageIds = new Set(tableData.value.map(guildKey))
  const next = new Map(selectedMap.value)
  pageIds.forEach(id => {
    if (!rows.some(row => guildKey(row) === id)) {
      next.delete(id)
    }
  })
  rows.forEach(row => {
    next.set(guildKey(row), row)
  })
  selectedMap.value = next
}

const handleClear = () => {
  selectedMap.value = new Map()
  tableRef.value?.clearSelection()
}

const handleClose = () => {
  emit('update:visible', false)
}

const handleConfirm = () => {
  const guilds = Array.from(selectedMap.value.values())
  if (guilds.length === 0) {
    ElMessage.warning(t('pages.liveDailyEffectiveLiveList.selectGuildRequired'))
    return
  }
  emit('confirm', guilds)
  handleClose()
}

watch(
    () => props.visible,
    (visible) => {
      if (!visible) {
        return
      }
      keyword.value = ''
      pagination.pageIndex = 1
      resetSelectedMap(props.initialGuilds || [])
      fetchList()
    },
)
</script>

<style scoped>
.picker-tip {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
  font-size: 13px;
}

.selected-summary {
  margin-bottom: 8px;
  color: var(--el-text-color-regular);
  font-size: 13px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.money-amount {
  font-variant-numeric: tabular-nums;
}
</style>
