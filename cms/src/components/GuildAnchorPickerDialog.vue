<template>
  <el-dialog
      :model-value="visible"
      :close-on-click-modal="false"
      destroy-on-close
      :title="dialogTitle"
      width="900px"
      @close="handleClose"
  >
    <div v-if="view === 'guild'" class="picker-body">
      <p class="hint">{{ t('pages.liveRecordList.guildPickerHint') }}</p>
      <el-form inline @submit.prevent>
        <el-form-item :label="t('pages.guildList.guildName')">
          <el-input
              v-model="guildKeyword"
              clearable
              :placeholder="t('pages.guildList.guildNameSearchPlaceholder')"
              style="width: 280px"
              @keyup.enter="searchGuilds"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchGuilds">{{ t('common.search') }}</el-button>
          <el-button @click="resetGuildSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="guildLoading"
          :data="guildTableData"
          highlight-current-row
          style="width: 100%"
          @row-click="enterMembers"
      >
        <el-table-column label="ID" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.guildList.guildName')" min-width="160" prop="name" show-overflow-tooltip/>
        <el-table-column :label="t('pages.guildList.leader')" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.leaderName || row.leaderId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildList.description')" min-width="180" prop="description" show-overflow-tooltip/>
        <el-table-column fixed="right" :label="t('common.actions')" width="100">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="enterMembers(row)">
              {{ t('pages.liveRecordList.viewMembers') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="guildPagination.pageIndex"
            v-model:page-size="guildPagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="guildPagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchGuildList"
            @size-change="handleGuildSizeChange"
        />
      </div>
    </div>

    <div v-else class="picker-body">
      <div class="members-toolbar">
        <el-button @click="backToGuildList">{{ t('pages.liveRecordList.backToGuildList') }}</el-button>
        <span v-if="currentGuild" class="current-guild">
          {{ t('pages.liveRecordList.currentGuild', {name: currentGuild.name}) }}
        </span>
      </div>

      <el-form inline @submit.prevent>
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="memberKeyword"
              clearable
              :placeholder="t('pages.liveRecordList.searchGuildAnchor')"
              style="width: 280px"
              @keyup.enter="searchMembers"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchMembers">{{ t('common.search') }}</el-button>
          <el-button @click="resetMemberSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <div class="selected-summary">
        {{ t('pages.liveRecordList.selectedGuildAnchorsCount', {count: selectedCount}) }}
      </div>

      <el-table
          ref="memberTableRef"
          v-loading="memberLoading"
          :data="memberTableData"
          row-key="id"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column reserve-selection type="selection" width="48"/>
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
        <el-table-column :label="t('common.phone')" min-width="130" prop="phone">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveStatus')" width="90">
          <template #default="{ row }">{{ formatLiveStatus(row.liveStatus) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="memberPagination.pageIndex"
            v-model:page-size="memberPagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="memberPagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchMemberList"
            @size-change="handleMemberSizeChange"
        />
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClear">{{ t('pages.liveRecordList.clearGuildAnchor') }}</el-button>
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
import type {AnchorListItem, Guild} from '@/types/api'

const props = withDefaults(defineProps<{
  visible: boolean
  initialAnchors?: AnchorListItem[]
}>(), {
  initialAnchors: () => [],
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: [anchors: AnchorListItem[]]
}>()

type PickerView = 'guild' | 'members'

const {t} = useI18n()
const view = ref<PickerView>('guild')
const currentGuild = ref<Guild | null>(null)

const guildLoading = ref(false)
const guildKeyword = ref('')
const guildTableData = ref<Guild[]>([])
const guildPagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const memberLoading = ref(false)
const memberKeyword = ref('')
const memberTableData = ref<AnchorListItem[]>([])
const memberTableRef = ref<TableInstance>()
const memberPagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const selectedMap = ref(new Map<string, AnchorListItem>())
const selectedCount = computed(() => selectedMap.value.size)

const dialogTitle = computed(() => {
  if (view.value === 'members' && currentGuild.value) {
    return t('pages.liveRecordList.selectGuildMembersTitle', {name: currentGuild.value.name})
  }
  return t('pages.liveRecordList.selectGuildAnchorTitle')
})

const anchorKey = (anchor: AnchorListItem) => String(anchor.id)

const formatLiveStatus = (liveStatus?: number) => {
  return liveStatus === 1 ? t('common.live') : t('common.offline')
}

const resetSelectedMap = (anchors: AnchorListItem[]) => {
  const map = new Map<string, AnchorListItem>()
  anchors.forEach(anchor => {
    map.set(anchorKey(anchor), anchor)
  })
  selectedMap.value = map
}

const syncMemberSelection = async () => {
  if (!memberTableRef.value) {
    return
  }
  await nextTick()
  memberTableRef.value.clearSelection()
  memberTableData.value.forEach(row => {
    if (selectedMap.value.has(anchorKey(row))) {
      memberTableRef.value?.toggleRowSelection(row, true)
    }
  })
}

const fetchGuildList = async () => {
  guildLoading.value = true
  try {
    const response = await guildApi.getGuildList({
      pageIndex: guildPagination.pageIndex,
      pageSize: guildPagination.pageSize,
      name: guildKeyword.value.trim() || undefined,
    })
    guildTableData.value = response.data || []
    guildPagination.total = response.total || 0
  } catch (error) {
    console.error('fetch guild list for picker failed:', error)
    guildTableData.value = []
    guildPagination.total = 0
  } finally {
    guildLoading.value = false
  }
}

const fetchMemberList = async () => {
  if (!currentGuild.value) {
    memberTableData.value = []
    memberPagination.total = 0
    return
  }
  memberLoading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: memberPagination.pageIndex,
      pageSize: memberPagination.pageSize,
      guildId: currentGuild.value.id,
      guildOnly: true,
      key: memberKeyword.value.trim() || undefined,
    })
    memberTableData.value = response.data || []
    memberPagination.total = response.total || 0
    await syncMemberSelection()
  } catch (error) {
    console.error('fetch guild members for picker failed:', error)
    memberTableData.value = []
    memberPagination.total = 0
  } finally {
    memberLoading.value = false
  }
}

const searchGuilds = () => {
  guildPagination.pageIndex = 1
  fetchGuildList()
}

const resetGuildSearch = () => {
  guildKeyword.value = ''
  guildPagination.pageIndex = 1
  fetchGuildList()
}

const handleGuildSizeChange = (size: number) => {
  guildPagination.pageSize = size
  guildPagination.pageIndex = 1
  fetchGuildList()
}

const enterMembers = (row: Guild) => {
  currentGuild.value = row
  view.value = 'members'
  memberKeyword.value = ''
  memberPagination.pageIndex = 1
  fetchMemberList()
}

const backToGuildList = () => {
  view.value = 'guild'
  currentGuild.value = null
}

const searchMembers = () => {
  memberPagination.pageIndex = 1
  fetchMemberList()
}

const resetMemberSearch = () => {
  memberKeyword.value = ''
  memberPagination.pageIndex = 1
  fetchMemberList()
}

const handleMemberSizeChange = (size: number) => {
  memberPagination.pageSize = size
  memberPagination.pageIndex = 1
  fetchMemberList()
}

const handleSelectionChange = (rows: AnchorListItem[]) => {
  const currentPageIds = new Set(memberTableData.value.map(row => anchorKey(row)))
  currentPageIds.forEach(id => {
    if (!rows.some(row => anchorKey(row) === id)) {
      selectedMap.value.delete(id)
    }
  })
  rows.forEach(row => {
    const item = currentGuild.value && !row.guildName
        ? {...row, guildName: currentGuild.value.name}
        : row
    selectedMap.value.set(anchorKey(item), item)
  })
}

const resetDialogState = () => {
  view.value = 'guild'
  currentGuild.value = null
  guildKeyword.value = ''
  memberKeyword.value = ''
  guildPagination.pageIndex = 1
  memberPagination.pageIndex = 1
  resetSelectedMap(props.initialAnchors || [])
}

const handleClear = () => {
  selectedMap.value = new Map()
  emit('confirm', [])
  emit('update:visible', false)
}

const handleConfirm = () => {
  emit('confirm', [...selectedMap.value.values()])
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
    resetDialogState()
    fetchGuildList()
  },
)
</script>

<style scoped>
.picker-body {
  min-height: 360px;
}

.hint {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.members-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.current-guild {
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.selected-summary {
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

:deep(.el-table__row) {
  cursor: pointer;
}
</style>
