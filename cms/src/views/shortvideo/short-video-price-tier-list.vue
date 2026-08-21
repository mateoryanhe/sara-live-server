<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ShortVideoPriceTierManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.shortVideoPriceTierList.addPriceTier') }}</el-button>
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
          <el-table-column :label="t('pages.shortVideoPriceTierList.diamondPrice')" width="140">
            <template #default="{ row }">
              {{ formatPrice(row.price) }}
            </template>
          </el-table-column>
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
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.shortVideoPriceTierList.diamondPrice')" prop="price">
          <el-input-number
              v-model="currentRow.price"
              :min="0"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
          />
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
import {shortVideoApi} from '@/api/modules/shortVideo'
import type {ShortVideoPriceTier} from '@/types/api'
import {formatPrice, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  statusFilter: number
}

interface PriceTierForm {
  id: string
  price: number
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<ShortVideoPriceTier[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  statusFilter: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): PriceTierForm => ({
  id: '',
  price: 0,
})
const currentRow = ref<PriceTierForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  price: [{required: true, message: t('pages.shortVideoPriceTierList.priceRequired'), trigger: 'blur'}],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoPriceTierList({
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch short video price tier list failed:', error)
    ElMessage.error(t('pages.shortVideoPriceTierList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.statusFilter = 0
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
  dialogTitle.value = t('pages.shortVideoPriceTierList.addTitle')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: ShortVideoPriceTier) => {
  dialogTitle.value = t('pages.shortVideoPriceTierList.editTitle')
  currentRow.value = {
    id: row.id,
    price: truncateNumber(row.price),
  }
  dialogVisible.value = true
}

const handleDelete = async (row: ShortVideoPriceTier) => {
  try {
    await ElMessageBox.confirm(
      t('pages.shortVideoPriceTierList.deleteConfirm'),
      t('common.confirm'),
      {type: 'warning'},
    )
    await shortVideoApi.deleteShortVideoPriceTier(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete short video price tier failed:', error)
    }
  }
}

const handleOnShelf = async (row: ShortVideoPriceTier) => {
  try {
    await shortVideoApi.onShelfShortVideoPriceTier(row.id)
    ElMessage.success(t('pages.shortVideoPriceTierList.onShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('on shelf short video price tier failed:', error)
  }
}

const handleOffShelf = async (row: ShortVideoPriceTier) => {
  try {
    await shortVideoApi.offShelfShortVideoPriceTier(row.id)
    ElMessage.success(t('pages.shortVideoPriceTierList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('off shelf short video price tier failed:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    try {
      if (currentRow.value.id) {
        await shortVideoApi.updateShortVideoPriceTier({
          id: currentRow.value.id,
          price: currentRow.value.price,
        })
      } else {
        await shortVideoApi.createShortVideoPriceTier({
          price: currentRow.value.price,
        })
      }
      ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save short video price tier failed:', error)
    }
  })
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

.table-header {
  margin-bottom: 16px;
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
