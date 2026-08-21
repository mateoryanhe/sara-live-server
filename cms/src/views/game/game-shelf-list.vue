<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GameShelfListManagement') }}</span>
          <div class="header-actions">
            <el-button
                :disabled="!canBatchOffShelf"
                :loading="shelfOperating"
                type="warning"
                @click="handleBatchOffShelf"
            >
              {{ t('common.batchOffShelf') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.gameShelfList.noteTitle')"
          class="tip-alert"
          show-icon
          type="info"
      >
        <p>{{ t('pages.gameShelfList.tipLine1') }}</p>
      </el-alert>

      <div v-if="selectedRows.length" class="selection-tip">{{ t('common.selectedCount', {count: selectedRows.length}) }}</div>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item :label="t('pages.gameList.gameCode')">
          <el-input v-model="searchForm.gameCode" clearable :placeholder="t('pages.gameList.gameCodePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameList.gameName')">
          <el-input v-model="searchForm.name" clearable :placeholder="t('pages.gameList.gameNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameList.platform')">
          <el-input v-model="searchForm.platform" clearable :placeholder="t('pages.gameList.platformPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="loading"
          :data="tableData"
          row-key="gameCode"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48"/>
        <el-table-column :label="t('pages.gameList.gameCode')" min-width="140" prop="gameCode"/>
        <el-table-column :label="t('common.name')" min-width="120" prop="name">
          <template #default="{ row }">{{ row.name || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameShelfList.liveGameName')" min-width="120">
          <template #default="{ row }">{{ row.liveGameName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.nameEn')" min-width="120" prop="nameEn">
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.platform')" prop="platform" width="100"/>
        <el-table-column :label="t('pages.gameShelfList.liveGameCover')" min-width="120">
          <template #default="{ row }">
            <el-image
                v-if="row.liveGameCoverUrl"
                :preview-src-list="[row.liveGameCoverUrl]"
                :src="row.liveGameCoverUrl"
                fit="cover"
                preview-teleported
                style="width: 72px; height: 48px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.cover')" min-width="120">
          <template #default="{ row }">
            <el-image
                v-if="row.cover"
                :preview-src-list="[row.cover]"
                :src="row.cover"
                fit="cover"
                preview-teleported
                style="width: 72px; height: 48px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="220">
          <template #default="{ row }">
            <el-button
                v-if="can('vendorConfig')"
                link
                type="success"
                @click="handleOpenVendorConfig(row)"
            >
              {{ t('pages.gameShelfList.vendorConfig') }}
            </el-button>
            <el-button
                v-if="can('edit')"
                link
                type="primary"
                @click="openEditDialog(row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button
                v-if="can('shelf')"
                link
                type="warning"
                @click="handleOffShelf(row)"
            >
              {{ t('common.offShelf') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog
        v-model="editDialogVisible"
        :title="t('pages.gameShelfList.editTitle')"
        destroy-on-close
        width="520px"
        @closed="resetEditForm"
    >
      <el-form ref="editFormRef" :model="editForm" label-width="120px">
        <el-form-item :label="t('pages.gameList.gameCode')">
          <el-input v-model="editForm.gameCode" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.gameShelfList.defaultNameEn')">
          <span>{{ editForm.defaultNameEn || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.gameShelfList.defaultCover')">
          <el-image
              v-if="editForm.defaultCover"
              :preview-src-list="[editForm.defaultCover]"
              :src="editForm.defaultCover"
              fit="cover"
              preview-teleported
              style="width: 72px; height: 48px"
          />
          <span v-else>-</span>
        </el-form-item>
        <el-form-item :label="t('pages.gameShelfList.liveGameName')">
          <el-input
              v-model="editForm.liveGameName"
              clearable
              :placeholder="t('pages.gameShelfList.liveGameNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.gameShelfList.liveGameCover')">
          <el-input
              v-model="editForm.liveGameCover"
              clearable
              :placeholder="t('pages.gameShelfList.liveGameCoverPlaceholder')"
          />
        </el-form-item>
        <el-form-item v-if="editCoverPreview" :label="t('pages.gameShelfList.liveCoverPreview')">
          <el-image
              :preview-src-list="[editCoverPreview]"
              :src="editCoverPreview"
              fit="cover"
              preview-teleported
              style="width: 120px; height: 72px"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="editSaving" type="primary" @click="handleEditSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import type {GameShelfItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('GameShelfListManagement')

interface SearchForm {
  gameCode: string
  name: string
  platform: string
}

interface EditForm {
  gameCode: string
  defaultNameEn: string
  defaultCover: string
  liveGameName: string
  liveGameCover: string
  savedLiveCoverUrl: string
}

const loading = ref(false)
const shelfOperating = ref(false)
const editSaving = ref(false)
const editDialogVisible = ref(false)
const editFormRef = ref<FormInstance>()
const tableData = ref<GameShelfItem[]>([])
const selectedRows = ref<GameShelfItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive<SearchForm>({
  gameCode: '',
  name: '',
  platform: '',
})

const editForm = reactive<EditForm>({
  gameCode: '',
  defaultNameEn: '',
  defaultCover: '',
  liveGameName: '',
  liveGameCover: '',
  savedLiveCoverUrl: '',
})

const canBatchOffShelf = computed(() => selectedRows.value.length > 0)

const editCoverPreview = computed(() => {
  const cover = editForm.liveGameCover.trim()
  if (cover.startsWith('http://') || cover.startsWith('https://')) {
    return cover
  }
  if (!cover && editForm.savedLiveCoverUrl) {
    return editForm.savedLiveCoverUrl
  }
  return ''
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await gamePlatformApi.getGameShelfList({
      gameCode: searchForm.gameCode,
      name: searchForm.name,
      platform: searchForm.platform,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('fetch game shelf list failed:', error)
    ElMessage.error(t('pages.gameShelfList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: GameShelfItem[]) => {
  selectedRows.value = rows
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.gameCode = ''
  searchForm.name = ''
  searchForm.platform = ''
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchList()
}

const openEditDialog = (row: GameShelfItem) => {
  editForm.gameCode = row.gameCode
  editForm.defaultNameEn = row.nameEn || row.name || ''
  editForm.defaultCover = row.cover || ''
  editForm.liveGameName = row.liveGameName || ''
  editForm.liveGameCover = row.liveGameCover || ''
  editForm.savedLiveCoverUrl = row.liveGameCoverUrl || ''
  editDialogVisible.value = true
}

const resetEditForm = () => {
  editForm.gameCode = ''
  editForm.defaultNameEn = ''
  editForm.defaultCover = ''
  editForm.liveGameName = ''
  editForm.liveGameCover = ''
  editForm.savedLiveCoverUrl = ''
  editFormRef.value?.clearValidate()
}

const handleOpenVendorConfig = (row: GameShelfItem) => {
  router.push({
    name: 'GameVendorConfig',
    query: {
      gameCode: row.gameCode,
      platform: row.platform,
      name: row.name || row.nameEn || row.gameCode,
    },
  })
}

const handleOffShelf = async (row: GameShelfItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.gameShelfList.offShelfConfirm', {name: row.name || row.nameEn || row.gameCode}),
        t('common.confirm'),
        {type: 'warning'},
    )
    const response = await gamePlatformApi.deleteGameShelf({gameCode: row.gameCode})
    if (response?.success) {
      ElMessage.success(t('pages.gameShelfList.offShelfSuccess'))
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameShelfList.offShelfFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('off shelf failed:', error)
      ElMessage.error(t('pages.gameShelfList.offShelfFailed'))
    }
  }
}

const handleBatchOffShelf = async () => {
  const gameCodes = selectedRows.value.map(row => row.gameCode).filter(Boolean)
  if (!gameCodes.length) {
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.gameShelfList.batchOffShelfConfirm', {count: gameCodes.length}),
        t('pages.gameShelfList.batchOffShelfTitle'),
        {type: 'warning'},
    )
    shelfOperating.value = true
    const response = await gamePlatformApi.batchDeleteGameShelf({gameCodes})
    if (response?.success) {
      ElMessage.success(t('pages.gameShelfList.batchOffShelfSuccess', {count: response.successCount || gameCodes.length}))
      selectedRows.value = []
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameShelfList.batchOffShelfFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('batch off shelf failed:', error)
      ElMessage.error(t('pages.gameShelfList.batchOffShelfFailed'))
    }
  } finally {
    shelfOperating.value = false
  }
}

const handleEditSave = async () => {
  if (!editForm.gameCode) {
    return
  }
  editSaving.value = true
  try {
    const response = await gamePlatformApi.updateGameShelf({
      gameCode: editForm.gameCode,
      liveGameName: editForm.liveGameName,
      liveGameCover: editForm.liveGameCover,
    })
    if (response?.success) {
      ElMessage.success(t('pages.gameShelfList.editSuccess'))
      editDialogVisible.value = false
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameShelfList.editFailed'))
    }
  } catch (error) {
    console.error('update game shelf failed:', error)
    ElMessage.error(t('pages.gameShelfList.editFailed'))
  } finally {
    editSaving.value = false
  }
}

onMounted(() => {
  fetchList()
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

.tip-alert {
  margin-bottom: 16px;
}

.selection-tip {
  margin-bottom: 12px;
  color: var(--el-color-primary);
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
