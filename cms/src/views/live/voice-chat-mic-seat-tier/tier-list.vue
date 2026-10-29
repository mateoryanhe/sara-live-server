<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="card-header-main">
            <span>{{ t('menu.VoiceChatMicSeatTierManagement') }}</span>
            <span class="card-tip">{{ t('pages.voiceChatMicSeatTierList.seatCountRangeHint') }}</span>
          </div>
          <el-button
              v-if="can('sync')"
              :loading="syncingToTarget"
              type="warning"
              @click="handleSyncToTarget"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.voiceChatMicSeatTierList.addTier') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.voiceChatMicSeatTierList.seatCount')">
            <el-input-number
                v-model="searchForm.seatCount"
                :min="0"
                :max="25"
                controls-position="right"
                :placeholder="t('common.all')"
            />
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

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.voiceChatMicSeatTierList.seatCount')" prop="seatCount" width="160"/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="260">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="row.status !== 1"
                  size="small"
                  type="success"
                  @click="handleOnShelf(row)"
              >
                {{ t('common.onShelf') }}
              </el-button>
              <el-button
                  v-else
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                {{ t('common.offShelf') }}
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="140px">
        <el-form-item :label="t('pages.voiceChatMicSeatTierList.seatCount')" prop="seatCount">
          <el-input-number
              v-model="currentRow.seatCount"
              :min="2"
              :max="25"
              controls-position="right"
          />
          <div class="form-tip">{{ t('pages.voiceChatMicSeatTierList.seatCountFormHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.voiceChatMicSeatTierList.sortHigherFirst') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {voiceChatMicSeatTierApi, type VoiceChatMicSeatTier} from '@/api/modules/voice-chat-mic-seat-tier'
import {dataSyncApi} from '@/api/modules/data-sync'
import {confirmDataSync} from '@/utils/confirm-data-sync'
import {usePagePermission} from '@/composables/usePagePermission'

const {can} = usePagePermission('VoiceChatMicSeatTierManagement')

interface SearchForm {
  seatCount: number
  statusFilter: number
}

interface TierForm {
  id: string
  seatCount: number
  sort: number
}

const {t} = useI18n()
const loading = ref(false)
const syncingToTarget = ref(false)
const tableData = ref<VoiceChatMicSeatTier[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  seatCount: 0,
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): TierForm => ({
  id: '',
  seatCount: 2,
  sort: 0
})
const currentRow = ref<TierForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  seatCount: [{required: true, message: t('pages.voiceChatMicSeatTierList.seatCountRequired'), trigger: 'blur'}]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await voiceChatMicSeatTierApi.list({
      seatCount: searchForm.seatCount > 0 ? searchForm.seatCount : undefined,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch tier list failed:', error)
    ElMessage.error(t('pages.voiceChatMicSeatTierList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
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

const handleAdd = () => {
  dialogTitle.value = t('pages.voiceChatMicSeatTierList.addTier')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: VoiceChatMicSeatTier) => {
  dialogTitle.value = t('pages.voiceChatMicSeatTierList.editTier')
  currentRow.value = {
    id: row.id,
    seatCount: row.seatCount,
    sort: Number(row.sort) || 0
  }
  dialogVisible.value = true
}

const handleDelete = async (row: VoiceChatMicSeatTier) => {
  try {
    await ElMessageBox.confirm(
        t('pages.voiceChatMicSeatTierList.deleteConfirm', {seatCount: row.seatCount}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        })
    await voiceChatMicSeatTierApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    console.error('delete failed:', error)
  }
}

const handleOnShelf = async (row: VoiceChatMicSeatTier) => {
  try {
    await voiceChatMicSeatTierApi.onShelf(row.id)
    ElMessage.success(t('pages.voiceChatMicSeatTierList.onShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.voiceChatMicSeatTierList.operationFailed'))
  }
}

const handleOffShelf = async (row: VoiceChatMicSeatTier) => {
  try {
    await voiceChatMicSeatTierApi.offShelf(row.id)
    ElMessage.success(t('pages.voiceChatMicSeatTierList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('off shelf failed:', error)
    ElMessage.error(t('pages.voiceChatMicSeatTierList.operationFailed'))
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        seatCount: currentRow.value.seatCount,
        sort: currentRow.value.sort
      }
      if (currentRow.value.id) {
        await voiceChatMicSeatTierApi.update({id: currentRow.value.id, ...payload})
      } else {
        await voiceChatMicSeatTierApi.create(payload)
      }
      ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.voiceChatMicSeatTierList.saveFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.seatCount = 0
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
}

const handleSyncToTarget = async () => {
  try {
    await confirmDataSync({
      detail: t('pages.voiceChatMicSeatTierList.syncToTargetConfirm'),
      title: t('common.syncData'),
    })
  } catch {
    return
  }
  syncingToTarget.value = true
  try {
    const response = await dataSyncApi.syncVoiceChatMicSeatTier()
    if (response?.success) {
      ElMessage.success(response.message || t('pages.voiceChatMicSeatTierList.syncToTargetSuccess'))
    } else {
      ElMessage.error(t('pages.voiceChatMicSeatTierList.syncToTargetFailed'))
    }
  } catch (error) {
    console.error('sync voice chat mic seat tier failed:', error)
    ElMessage.error(t('pages.voiceChatMicSeatTierList.syncToTargetFailedCheckConfig'))
  } finally {
    syncingToTarget.value = false
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
  font-size: 16px;
  font-weight: bold;
}

.card-header-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.card-tip {
  font-size: 13px;
  font-weight: normal;
  color: var(--el-text-color-secondary);
}

.table-header {
  margin-bottom: 20px;
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

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
