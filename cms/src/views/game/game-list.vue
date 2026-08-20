<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GameVendorGameListManagement') }}</span>
          <div class="header-actions">
            <el-button
                v-if="can('sync')"
                :loading="syncing"
                type="primary"
                @click="handleSyncVendorLibrary"
            >
              {{ t('common.syncData') }}
            </el-button>
            <el-button
                :disabled="!canBatchOnShelf"
                :loading="shelfOperating"
                type="success"
                @click="handleBatchOnShelf"
            >
              {{ t('common.batchOnShelf') }}
            </el-button>
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
          :title="t('pages.gameList.noteTitle')"
          class="tip-alert"
          show-icon
          type="info"
      >
        <p>{{ t('pages.gameList.tipLine1') }}</p>
        <p>{{ t('pages.gameList.tipLine2') }}</p>
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
        <el-form-item :label="t('common.category')">
          <el-input v-model="searchForm.category" clearable :placeholder="t('pages.gameList.categoryPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="loading"
          :data="tableData"
          :row-key="vendorGameRowKey"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48"/>
        <el-table-column :label="t('pages.gameList.gameCode')" min-width="140" prop="gameCode"/>
        <el-table-column :label="t('common.name')" min-width="140" prop="name"/>
        <el-table-column :label="t('pages.gameList.nameEn')" min-width="140" prop="nameEn"/>
        <el-table-column :label="t('common.category')" prop="category" width="100"/>
        <el-table-column :label="t('pages.gameList.platform')" prop="platform" width="100"/>
        <el-table-column :label="t('pages.gameList.shelfStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.onShelf ? 'success' : 'info'">{{ row.onShelf ? t('pages.gameList.onShelfStatus') : t('pages.gameList.offShelfStatus') }}</el-tag>
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
        <el-table-column fixed="right" :label="t('common.actions')" width="160">
          <template #default="{ row }">
            <el-button
                v-if="!row.onShelf"
                link
                type="success"
                @click="handleOnShelf(row)"
            >
              {{ t('common.onShelf') }}
            </el-button>
            <el-button
                v-else
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
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import type {VendorGame} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('GameVendorGameListManagement')

interface SearchForm {
  gameCode: string
  name: string
  platform: string
  category: string
}

const loading = ref(false)
const syncing = ref(false)
const shelfOperating = ref(false)
const tableData = ref<VendorGame[]>([])
const selectedRows = ref<VendorGame[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive<SearchForm>({
  gameCode: '',
  name: '',
  platform: '',
  category: '',
})

const canBatchOnShelf = computed(() => selectedRows.value.some(row => !row.onShelf))
const canBatchOffShelf = computed(() => selectedRows.value.some(row => row.onShelf))

const vendorGameRowKey = (row: VendorGame) => `${row.gameCode}@${row.platform}`

const fetchList = async () => {
  loading.value = true
  try {
    const response = await gamePlatformApi.getVendorGameList({
      gameCode: searchForm.gameCode,
      name: searchForm.name,
      platform: searchForm.platform,
      category: searchForm.category,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('fetch game list failed:', error)
    ElMessage.error(t('pages.gameList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: VendorGame[]) => {
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
  searchForm.category = ''
  currentPage.value = 1
  fetchList()
}

const handleSyncVendorLibrary = async () => {
  try {
    await ElMessageBox.confirm(
        t('pages.gameList.syncConfirm'),
        t('common.syncData'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  syncing.value = true
  try {
    const response = await gamePlatformApi.reloadVendorGameCache()
    if (response?.success) {
      ElMessage.success(t('pages.gameList.syncSuccess', {count: response.count || 0}))
      currentPage.value = 1
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameList.syncFailed'))
    }
  } catch (error) {
    console.error('sync vendor game library failed:', error)
    ElMessage.error(t('pages.gameList.syncFailed'))
  } finally {
    syncing.value = false
  }
}

onMounted(() => {
  fetchList()
})

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchList()
}

const handleOnShelf = async (row: VendorGame) => {
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.addGameShelf({gameCode: row.gameCode, platform: row.platform})
    if (response?.success) {
      ElMessage.success(t('pages.gameList.onShelfSuccess'))
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameList.onShelfFailed'))
    }
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.gameList.onShelfFailed'))
  } finally {
    shelfOperating.value = false
  }
}

const handleOffShelf = async (row: VendorGame) => {
  try {
    await ElMessageBox.confirm(
        t('pages.gameList.offShelfConfirm', {name: row.name || row.gameCode}),
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

const handleBatchOnShelf = async () => {
  const items = selectedRows.value.filter(row => !row.onShelf).map(row => ({
    gameCode: row.gameCode,
    platform: row.platform,
  }))
  if (!items.length) {
    ElMessage.warning(t('pages.gameList.selectUnpublished'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.gameList.batchOnShelfConfirm', {count: items.length}),
        t('pages.gameList.batchOnShelfTitle'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.batchAddGameShelf({items})
    if (response?.success) {
      if (response.skipCount > 0) {
        ElMessage.success(t('pages.gameList.batchOnShelfDoneWithSkip', {
          success: response.successCount,
          skip: response.skipCount,
        }))
      } else {
        ElMessage.success(t('pages.gameList.batchOnShelfSuccess', {count: response.successCount}))
      }
      await fetchList()
    } else {
      ElMessage.error(t('pages.gameList.batchOnShelfFailed'))
    }
  } catch (error) {
    console.error('batch on shelf failed:', error)
    ElMessage.error(t('pages.gameList.batchOnShelfFailed'))
  } finally {
    shelfOperating.value = false
  }
}

const handleBatchOffShelf = async () => {
  const gameCodes = selectedRows.value.filter(row => row.onShelf).map(row => row.gameCode)
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
