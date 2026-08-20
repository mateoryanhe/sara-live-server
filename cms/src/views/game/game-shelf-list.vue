<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GameShelfListManagement') }}</span>
          <div class="header-actions">
            <el-button
                v-if="isPickMode"
                @click="cancelPickUser"
            >
              {{ t('pages.gameShelfList.cancelPickUser') }}
            </el-button>
            <el-button
                v-if="!isPickMode"
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
          v-if="isPickMode"
          :closable="false"
          class="tip-alert"
          show-icon
          type="warning"
      >
        <p>{{ pickUserHintText }}</p>
      </el-alert>

      <el-alert
          v-else
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
        <el-table-column v-if="!isPickMode" type="selection" width="48"/>
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
        <el-table-column fixed="right" :label="t('common.actions')" :width="isPickMode ? 120 : 220">
          <template #default="{ row }">
            <el-button
                v-if="isPickMode && can('startGame')"
                :loading="startingGameCode === row.gameCode"
                link
                type="success"
                @click="handleStartGame(row)"
            >
              {{ t('pages.gameShelfList.startGame') }}
            </el-button>
            <template v-else-if="!isPickMode">
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import type {GameShelfItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {can} = usePagePermission('GameShelfListManagement')

const pickUserId = ref('')
const pickUserNickname = ref('')
const isPickMode = computed(() => !!pickUserId.value)
const pickUserHintText = computed(() => {
  if (!pickUserId.value) {
    return ''
  }
  if (pickUserNickname.value) {
    return t('pages.gameShelfList.pickUserHint', {
      name: pickUserNickname.value,
      id: pickUserId.value,
    })
  }
  return t('pages.gameShelfList.pickUserHintNoName', {id: pickUserId.value})
})

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
const startingGameCode = ref('')
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

const syncPickUserFromRoute = () => {
  pickUserId.value = typeof route.query.pickUserId === 'string' ? route.query.pickUserId.trim() : ''
  pickUserNickname.value = typeof route.query.pickUserNickname === 'string' ? route.query.pickUserNickname.trim() : ''
}

watch(() => route.query, () => {
  syncPickUserFromRoute()
}, {deep: true})

const cancelPickUser = () => {
  router.replace({path: '/game/game-shelf-list'})
}

const handleStartGame = async (row: GameShelfItem) => {
  if (!pickUserId.value) {
    return
  }
  startingGameCode.value = row.gameCode
  try {
    const response = await gamePlatformApi.getCMSGameStartLink({
      userId: pickUserId.value,
      gameCode: row.gameCode,
      platform: row.platform,
    })
    const link = response?.link?.trim()
    if (!link) {
      ElMessage.error(t('pages.gameShelfList.startGameEmpty'))
      return
    }
    window.open(link, '_blank', 'noopener,noreferrer')
  } catch (error) {
    console.error('get cms game start link failed:', error)
    ElMessage.error(t('pages.gameShelfList.startGameFailed'))
  } finally {
    startingGameCode.value = ''
  }
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

const handleEditSave = async () => {
  if (!editForm.gameCode) {
    return
  }
  editSaving.value = true
  try {
    const response = await gamePlatformApi.updateGameShelf({
      gameCode: editForm.gameCode,
      liveGameName: editForm.liveGameName.trim(),
      liveGameCover: editForm.liveGameCover.trim(),
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

const handleOffShelf = async (row: GameShelfItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.gameList.offShelfConfirm', {name: row.name || row.nameEn || row.gameCode}),
        t('common.confirmOffShelf'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.deleteGameShelf({gameCode: row.gameCode})
    if (response?.success) {
      ElMessage.success(t('pages.gameList.offShelfSuccess'))
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameList.offShelfFailed'))
    }
  } catch (error) {
    console.error('off shelf failed:', error)
    ElMessage.error(t('pages.gameList.offShelfFailed'))
  } finally {
    shelfOperating.value = false
  }
}

const handleBatchOffShelf = async () => {
  const gameCodes = selectedRows.value.map(row => row.gameCode)
  if (!gameCodes.length) {
    ElMessage.warning(t('pages.gameList.selectPublished'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.gameList.batchOffShelfConfirm', {count: gameCodes.length}),
        t('pages.gameList.batchOffShelfTitle'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.batchDeleteGameShelf({gameCodes})
    if (response?.success) {
      ElMessage.success(t('pages.gameList.batchOffShelfSuccess', {count: response.successCount}))
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameList.batchOffShelfFailed'))
    }
  } catch (error) {
    console.error('batch off shelf failed:', error)
    ElMessage.error(t('pages.gameList.batchOffShelfFailed'))
  } finally {
    shelfOperating.value = false
  }
}

onMounted(() => {
  syncPickUserFromRoute()
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
  gap: 12px;
  font-size: 16px;
  font-weight: bold;
}

.header-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.selection-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>
