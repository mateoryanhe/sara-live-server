<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PrivateRoomBillingManagement') }}</span>
          <span class="card-tip">{{ t('pages.billingList.billingByMinute') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.billingList.addBilling') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
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
          <el-table-column :label="t('pages.billingList.pricePerMinute')" width="160">
            <template #default="{ row }">
              {{ formatPrice(row.pricePerMinute) }}
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
        <el-form-item :label="t('pages.billingList.pricePerMinute')" prop="pricePerMinute">
          <el-input-number
              v-model="currentRow.pricePerMinute"
              :min="0"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
          />
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.billingList.sortHigherFirst') }}</div>
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
import {privateRoomBillingApi} from '@/api'
import type {PrivateRoomBilling} from '@/types/api.ts'
import {formatPrice, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  statusFilter: number
}

interface BillingForm {
  id: string
  pricePerMinute: number
  sort: number
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<PrivateRoomBilling[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): BillingForm => ({
  id: '',
  pricePerMinute: 0,
  sort: 0
})
const currentRow = ref<BillingForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  pricePerMinute: [{required: true, message: t('pages.billingList.priceRequired'), trigger: 'blur'}]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await privateRoomBillingApi.getBillingList({
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch billing list failed:', error)
    ElMessage.error(t('pages.billingList.fetchFailed'))
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
  dialogTitle.value = t('pages.billingList.addBilling')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: PrivateRoomBilling) => {
  dialogTitle.value = t('pages.billingList.editBilling')
  currentRow.value = {
    id: row.id,
    pricePerMinute: truncateNumber(row.pricePerMinute),
    sort: Number(row.sort) || 0
  }
  dialogVisible.value = true
}

const handleDelete = async (row: PrivateRoomBilling) => {
  try {
    await ElMessageBox.confirm(t('pages.billingList.deleteConfirm', {id: row.id}), t('common.confirmDelete'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await privateRoomBillingApi.deleteBilling(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    console.error('delete failed:', error)
  }
}

const handleOnShelf = async (row: PrivateRoomBilling) => {
  try {
    await privateRoomBillingApi.onShelfBilling(row.id)
    ElMessage.success(t('pages.billingList.onShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.billingList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: PrivateRoomBilling) => {
  try {
    await privateRoomBillingApi.offShelfBilling(row.id)
    ElMessage.success(t('pages.billingList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('off shelf failed:', error)
    ElMessage.error(t('pages.billingList.offShelfFailed'))
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        pricePerMinute: currentRow.value.pricePerMinute,
        sort: currentRow.value.sort
      }
      if (currentRow.value.id) {
        await privateRoomBillingApi.updateBilling({id: currentRow.value.id, ...payload})
      } else {
        await privateRoomBillingApi.createBilling(payload)
      }
      ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.billingList.saveFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
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
