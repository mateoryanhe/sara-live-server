<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GiftManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.giftList.addGift') }}</el-button>
          <el-button
              v-if="hasButtonPermission('GiftManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
          <el-button
              v-if="hasButtonPermission('GiftManagement', 'syncAssets')"
              :disabled="selectedRows.length === 0"
              :loading="syncingAssets"
              type="warning"
              plain
              @click="handleSyncAssets"
          >
            {{ t('pages.giftList.syncAssets') }}
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">
          {{ t('common.selectedCount', {count: selectedRows.length}) }}
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.giftList.giftName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.giftList.giftNameFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('common.category')">
            <el-select v-model="searchForm.categoryType" clearable :placeholder="t('common.all')" style="width: 160px">
              <el-option :label="t('pages.giftList.categoryNormal')" :value="0"/>
              <el-option :label="t('pages.giftList.categoryLucky')" :value="1"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('common.status')">
            <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option :value="2" :label="t('common.onlyOnShelf')"/>
              <el-option :value="1" :label="t('common.onlyOffShelf')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table
            v-loading="loading"
            :data="tableData"
            row-key="id"
            style="width: 100%"
            @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="48"/>
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.giftList.giftName')" prop="name" min-width="120"/>
          <el-table-column :label="t('common.nameEn')" prop="nameEn" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameEs')" prop="nameEs" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.namePt')" prop="namePt" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameHi')" prop="nameHi" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameId')" prop="nameId" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.icon')" width="90">
            <template #default="{ row }">
              <el-image
                  v-if="row.icon"
                  :preview-src-list="[row.icon]"
                  :src="row.icon"
                  fit="cover"
                  preview-teleported
                  style="width: 48px; height: 48px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.animationResource')" min-width="220">
            <template #default="{ row }">
              <video
                  v-if="isVideoUrl(row.animation)"
                  :src="row.animation"
                  class="table-media-preview"
                  controls
                  preload="metadata"
              />
              <el-image
                  v-else-if="isImageUrl(row.animation)"
                  :preview-src-list="[row.animation]"
                  :src="row.animation"
                  fit="cover"
                  preview-teleported
                  style="width: 48px; height: 48px"
              />
              <span v-else class="media-url-text">{{ row.animationName || row.animation || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.price')" prop="price" width="110">
            <template #default="{ row }">{{ formatPrice(row.price) }}</template>
          </el-table-column>
          <el-table-column :label="t('common.category')" prop="categoryType" width="120">
            <template #default="{ row }">{{ formatGiftCategoryLabel(row.categoryType) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.giftList.luckyGiftRtp')" width="100">
            <template #default="{ row }">
              {{ row.categoryType === GIFT_CATEGORY_LUCKY ? formatPrice(row.luckyGiftRtp ?? 0) : '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.description')" prop="description" min-width="160" show-overflow-tooltip/>
          <el-table-column :label="t('common.publishedAt')" prop="publishedAt" width="160">
            <template #default="{ row }">
              {{ row.publishedAt || '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="100">
            <template #default="{ row }">
              <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(row, cmd)">
                <el-button size="small" type="primary">
                  {{ t('common.actions') }}
                  <el-icon class="el-icon--right">
                    <ArrowDown/>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="canEditGift" command="edit">{{ t('common.edit') }}</el-dropdown-item>
                    <el-dropdown-item
                        v-if="canEditGift && row.status !== 1"
                        :divided="canEditGift"
                        command="onShelf"
                    >
                      {{ t('common.onShelf') }}
                    </el-dropdown-item>
                    <el-dropdown-item
                        v-if="canEditGift && row.status === 1"
                        :divided="canEditGift"
                        command="offShelf"
                    >
                      {{ t('common.offShelf') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="canDeleteGift" :divided="canEditGift" command="delete">{{ t('common.delete') }}</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ArrowDown} from '@element-plus/icons-vue'
import {dataSyncApi, giftApi} from '@/api'
import {confirmDataSync} from '@/utils/confirm-data-sync'
import type {Gift} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'
import {isImageUrl, isVideoUrl} from '@/utils/media-preview'
import {formatPrice} from '@/utils/number-format'
import {GIFT_CATEGORY_LUCKY, GIFT_CATEGORY_NORMAL} from '@/utils/gift-category'

interface SearchForm {
  name: string
  categoryType: number | undefined
  statusFilter: number
}

const {t} = useI18n()
const router = useRouter()

const canEditGift = computed(() => hasButtonPermission('GiftManagement', 'edit'))
const canDeleteGift = computed(() => hasButtonPermission('GiftManagement', 'delete'))
const hasRowActions = computed(() => canEditGift.value || canDeleteGift.value)

const handleRowCommand = (row: Gift, command: string) => {
  switch (command) {
    case 'edit':
      handleEdit(row)
      break
    case 'onShelf':
      handleOnShelf(row)
      break
    case 'offShelf':
      handleOffShelf(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const formatGiftCategoryLabel = (categoryType: number | undefined) => {
  if (categoryType === GIFT_CATEGORY_NORMAL) {
    return t('pages.giftList.categoryNormal')
  }
  if (categoryType === GIFT_CATEGORY_LUCKY) {
    return t('pages.giftList.categoryLucky')
  }
  return '-'
}

const loading = ref(false)
const syncing = ref(false)
const syncingAssets = ref(false)
const selectedRows = ref<Gift[]>([])
const tableData = ref<Gift[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  categoryType: undefined,
  statusFilter: 0
})

const fetchGiftList = async () => {
  loading.value = true
  try {
    const response = await giftApi.getGiftList({
      name: searchForm.name,
      categoryType: searchForm.categoryType,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch gift list failed:', error)
    ElMessage.error(t('pages.giftList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchGiftList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGiftList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGiftList()
}

const handleAdd = () => {
  router.push({name: 'GiftManagementEdit', params: {id: 'new'}})
}

const handleEdit = (row: Gift) => {
  router.push({
    name: 'GiftManagementEdit',
    params: {id: row.id},
    state: {giftRow: row},
  })
}

const handleDelete = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.deleteConfirm', {name: row.name}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )

    await giftApi.deleteGift(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchGiftList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete gift failed:', error)
    }
  }
}

const handleOnShelf = async (row: Gift) => {
  try {
    await giftApi.onShelfGift(row.id)
    ElMessage.success(t('pages.giftList.onShelfSuccess'))
    fetchGiftList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.giftList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.offShelfConfirm', {name: row.name}),
        t('common.confirmOffShelf'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await giftApi.offShelfGift(row.id)
    ElMessage.success(t('pages.giftList.offShelfSuccess'))
    fetchGiftList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('off shelf failed:', error)
    }
  }
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.categoryType = undefined
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchGiftList()
}

const handleSelectionChange = (rows: Gift[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.giftList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.giftList.invalidSelection'))
    return
  }
  try {
    await confirmDataSync({
      detail: t('pages.giftList.syncConfirm', {count: ids.length}),
      title: t('common.syncData'),
    })
    syncing.value = true
    const response = await dataSyncApi.syncGift({ids})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.giftList.syncSuccessDetail', {
        rows: response.rowCount,
        files: response.fileCount
      }))
    } else {
      ElMessage.error(t('pages.giftList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync failed:', error)
      ElMessage.error(t('pages.giftList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

const handleSyncAssets = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.giftList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.giftList.invalidSelection'))
    return
  }
  try {
    await confirmDataSync({
      detail: t('pages.giftList.syncAssetsConfirm', {count: ids.length}),
      title: t('pages.giftList.syncAssets'),
    })
    syncingAssets.value = true
    const response = await dataSyncApi.syncGiftAssets({ids})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.giftList.syncAssetsSuccessDetail', {
        rows: response.rowCount,
        files: response.fileCount
      }))
    } else {
      ElMessage.error(t('pages.giftList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync gift assets failed:', error)
      ElMessage.error(t('pages.giftList.syncFailedCheckConfig'))
    }
  } finally {
    syncingAssets.value = false
  }
}

onMounted(() => {
  fetchGiftList()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.table-header {
  margin-bottom: 20px;
}

.selection-tip {
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--el-color-primary);
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

.table-media-preview {
  width: 160px;
  max-height: 90px;
  display: block;
  background: #000;
  border-radius: 4px;
}

.media-url-text {
  word-break: break-all;
  line-height: 1.4;
}
</style>
