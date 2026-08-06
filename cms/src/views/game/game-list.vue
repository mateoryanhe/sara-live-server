<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>游戏列表</span>
          <div class="header-actions">
            <el-button
                :disabled="!canBatchOnShelf"
                :loading="shelfOperating"
                type="success"
                @click="handleBatchOnShelf"
            >
              批量上架
            </el-button>
            <el-button
                :disabled="!canBatchOffShelf"
                :loading="shelfOperating"
                type="warning"
                @click="handleBatchOffShelf"
            >
              批量下架
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>点击「搜索」从第三方全量拉取游戏列表(缓存 30 分钟,每次搜索都会重新拉取覆盖)。已上架游戏写入 game_cfgs 并永久缓存,App 端只展示上架游戏。</p>
        <p>勾选游戏后可批量上架/下架；表格右侧「操作」列可单条上架/下架。</p>
      </el-alert>

      <div v-if="selectedRows.length" class="selection-tip">已选 {{ selectedRows.length }} 项</div>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="游戏编码">
          <el-input v-model="searchForm.gameCode" clearable placeholder="游戏编码(模糊匹配)"/>
        </el-form-item>
        <el-form-item label="游戏名称">
          <el-input v-model="searchForm.name" clearable placeholder="中/英文名称(模糊匹配)"/>
        </el-form-item>
        <el-form-item label="平台">
          <el-input v-model="searchForm.platform" clearable placeholder="如 pg / pp"/>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="searchForm.category" clearable placeholder="如 slot"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
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
        <el-table-column label="游戏编码" min-width="140" prop="gameCode"/>
        <el-table-column label="名称" min-width="140" prop="name"/>
        <el-table-column label="英文名称" min-width="140" prop="nameEn"/>
        <el-table-column label="分类" prop="category" width="100"/>
        <el-table-column label="平台" prop="platform" width="100"/>
        <el-table-column label="上架状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.onShelf ? 'success' : 'info'">{{ row.onShelf ? '已上架' : '未上架' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="封面" min-width="120">
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
        <el-table-column fixed="right" label="操作" width="160">
          <template #default="{ row }">
            <el-button
                v-if="!row.onShelf"
                link
                type="success"
                @click="handleOnShelf(row)"
            >
              上架
            </el-button>
            <el-button
                v-else
                link
                type="warning"
                @click="handleOffShelf(row)"
            >
              下架
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
import {computed, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import type {VendorGame} from '@/types/api'

interface SearchForm {
  gameCode: string
  name: string
  platform: string
  category: string
}

const loading = ref(false)
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

const fetchList = async (refreshFromVendor = false) => {
  loading.value = true
  try {
    const response = await gamePlatformApi.getVendorGameList({
      gameCode: searchForm.gameCode,
      name: searchForm.name,
      platform: searchForm.platform,
      category: searchForm.category,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      refreshFromVendor,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取游戏列表失败:', error)
    ElMessage.error('获取游戏列表失败')
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: VendorGame[]) => {
  selectedRows.value = rows
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList(true)
}

const resetSearch = () => {
  searchForm.gameCode = ''
  searchForm.name = ''
  searchForm.platform = ''
  searchForm.category = ''
  currentPage.value = 1
  fetchList(true)
}

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
    const response = await gamePlatformApi.addGameShelf({gameCode: row.gameCode})
    if (response?.success) {
      ElMessage.success('上架成功')
      await fetchList()
    } else {
      ElMessage.error('上架失败')
    }
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  } finally {
    shelfOperating.value = false
  }
}

const handleOffShelf = async (row: VendorGame) => {
  try {
    await ElMessageBox.confirm(`确定下架游戏「${row.name || row.gameCode}」吗？`, '确认下架', {type: 'warning'})
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.deleteGameShelf({gameCode: row.gameCode})
    if (response?.success) {
      ElMessage.success('下架成功')
      await fetchList()
    } else {
      ElMessage.error('下架失败')
    }
  } catch (error) {
    console.error('下架失败:', error)
    ElMessage.error('下架失败')
  } finally {
    shelfOperating.value = false
  }
}

const handleBatchOnShelf = async () => {
  const gameCodes = selectedRows.value.filter(row => !row.onShelf).map(row => row.gameCode)
  if (!gameCodes.length) {
    ElMessage.warning('请选择未上架的游戏')
    return
  }
  try {
    await ElMessageBox.confirm(`确定批量上架 ${gameCodes.length} 个游戏吗？`, '确认批量上架', {type: 'warning'})
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.batchAddGameShelf({gameCodes})
    if (response?.success) {
      if (response.skipCount > 0) {
        ElMessage.success(`批量上架完成：成功 ${response.successCount} 个，跳过 ${response.skipCount} 个`)
      } else {
        ElMessage.success(`批量上架成功，共 ${response.successCount} 个`)
      }
      await fetchList()
    } else {
      ElMessage.error('批量上架失败')
    }
  } catch (error) {
    console.error('批量上架失败:', error)
    ElMessage.error('批量上架失败')
  } finally {
    shelfOperating.value = false
  }
}

const handleBatchOffShelf = async () => {
  const gameCodes = selectedRows.value.filter(row => row.onShelf).map(row => row.gameCode)
  if (!gameCodes.length) {
    ElMessage.warning('请选择已上架的游戏')
    return
  }
  try {
    await ElMessageBox.confirm(`确定批量下架 ${gameCodes.length} 个游戏吗？`, '确认批量下架', {type: 'warning'})
  } catch {
    return
  }
  shelfOperating.value = true
  try {
    const response = await gamePlatformApi.batchDeleteGameShelf({gameCodes})
    if (response?.success) {
      ElMessage.success(`批量下架成功，共 ${response.successCount} 个`)
      await fetchList()
    } else {
      ElMessage.error('批量下架失败')
    }
  } catch (error) {
    console.error('批量下架失败:', error)
    ElMessage.error('批量下架失败')
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
