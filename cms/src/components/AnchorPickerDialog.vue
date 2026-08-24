<template>
  <el-dialog
      :model-value="visible"
      :close-on-click-modal="false"
      destroy-on-close
      :title="t('pages.guildAnchorDailyLiveList.selectAnchorTitle')"
      width="860px"
      @close="handleClose"
  >
    <el-form inline @submit.prevent>
      <el-form-item :label="t('common.keyword')">
        <el-input
            v-model="keyword"
            clearable
            :placeholder="t('pages.guildAnchorDailyLiveList.anchorSearchPlaceholder')"
            style="width: 280px"
            @keyup.enter="handleSearch"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table
        v-loading="loading"
        :data="tableData"
        highlight-current-row
        style="width: 100%"
        @current-change="handleCurrentChange"
        @row-dblclick="handleConfirm"
    >
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
import {reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {accountApi, guildApi} from '@/api'
import type {AnchorListItem} from '@/types/api'

const props = defineProps<{
  visible: boolean
  guildId?: string
  ownedGuildOnly?: boolean
  initialAnchor?: AnchorListItem | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: [anchor: AnchorListItem | null]
}>()

const {t} = useI18n()
const loading = ref(false)
const keyword = ref('')
const tableData = ref<AnchorListItem[]>([])
const selectedAnchor = ref<AnchorListItem | null>(null)
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const fetchList = async () => {
  if (!props.ownedGuildOnly && !props.guildId) {
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
        })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch guild anchors for picker failed:', error)
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
  selectedAnchor.value = row ?? null
}

const handleClear = () => {
  selectedAnchor.value = null
  emit('confirm', null)
  emit('update:visible', false)
}

const handleConfirm = () => {
  emit('confirm', selectedAnchor.value)
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
    selectedAnchor.value = props.initialAnchor ?? null
    keyword.value = ''
    pagination.pageIndex = 1
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
</style>
