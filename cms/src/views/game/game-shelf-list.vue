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
        <el-table-column :label="t('common.name')" min-width="140" prop="name">
          <template #default="{ row }">{{ row.name || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.nameEn')" min-width="140" prop="nameEn">
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.platform')" prop="platform" width="100"/>
        <el-table-column :label="t('pages.gameList.shelfStatus')" width="100">
          <template #default>
            <el-tag type="success">{{ t('pages.gameList.onShelfStatus') }}</el-tag>
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
        <el-table-column fixed="right" :label="t('common.actions')" width="100">
          <template #default="{ row }">
            <el-button
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
import type {GameShelfItem} from '@/types/api'

const {t} = useI18n()

interface SearchForm {
  gameCode: string
  name: string
  platform: string
}

const loading = ref(false)
const shelfOperating = ref(false)
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

const canBatchOffShelf = computed(() => selectedRows.value.length > 0)

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
