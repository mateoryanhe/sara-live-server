<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CoinMerchantRechargeCfgManagement') }}</span>
        </div>
      </template>

      <div class="content">
        <div v-if="can('create')" class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.coinMerchantRechargeCfgList.add') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.coinMerchantRechargeCfgList.name')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.coinMerchantRechargeCfgList.nameFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('common.status')">
            <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option :value="2" :label="t('common.onlyOnShelf')"/>
              <el-option :value="1" :label="t('common.onlyOffShelf')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button v-if="can('search')" @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.coinMerchantRechargeCfgList.name')" min-width="140" prop="name"/>
          <el-table-column :label="t('pages.coinMerchantRechargeCfgList.usdPrice')" width="140">
            <template #default="{row}">
              {{ formatNumberDisplay(row.price, '-', PRICE_DECIMALS) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.coinMerchantRechargeCfgList.gold')" prop="gold" width="120"/>
          <el-table-column :label="t('common.status')" width="100">
            <template #default="{row}">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="170"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="260">
            <template #default="{row}">
              <el-button v-if="can('edit')" size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="can('edit') && row.status !== 1"
                  size="small"
                  type="success"
                  @click="handleOnShelf(row)"
              >
                {{ t('common.onShelf') }}
              </el-button>
              <el-button
                  v-else-if="can('edit') && row.status === 1"
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                {{ t('common.offShelf') }}
              </el-button>
              <el-button v-if="can('delete')" size="small" type="danger" @click="handleDelete(row)">
                {{ t('common.delete') }}
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
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="110px">
        <el-form-item :label="t('pages.coinMerchantRechargeCfgList.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('pages.coinMerchantRechargeCfgList.namePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.coinMerchantRechargeCfgList.usdPrice')" prop="price">
          <el-input-number v-model="form.price" :min="0.0001" :precision="PRICE_DECIMALS" :step="0.01" style="width: 100%"/>
        </el-form-item>
        <el-form-item :label="t('pages.coinMerchantRechargeCfgList.gold')" prop="gold">
          <el-input-number v-model="form.gold" :min="1" :precision="0" :step="1" style="width: 100%"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="submitting" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {
  coinMerchantRechargeCfgApi,
  type CoinMerchantRechargeCfg,
} from '@/api/modules/coin-merchant-recharge-cfg'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatNumberDisplay, truncateNumber} from '@/utils/number-format'

const PRICE_DECIMALS = 4
const {t} = useI18n()
const {can} = usePagePermission('CoinMerchantRechargeCfgManagement')

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<CoinMerchantRechargeCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const dialogVisible = ref(false)
const editingId = ref('')
const formRef = ref<FormInstance>()

const searchForm = reactive({name: '', statusFilter: 0})
const form = ref({name: '', price: 0.99, gold: 100})

const dialogTitle = computed(() =>
    editingId.value
        ? t('pages.coinMerchantRechargeCfgList.editTitle')
        : t('pages.coinMerchantRechargeCfgList.addTitle'),
)

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.coinMerchantRechargeCfgList.nameRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.coinMerchantRechargeCfgList.nameLength'), trigger: 'blur'},
  ],
  price: [{required: true, message: t('pages.coinMerchantRechargeCfgList.priceRequired'), trigger: 'change'}],
  gold: [{required: true, message: t('pages.coinMerchantRechargeCfgList.goldRequired'), trigger: 'change'}],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const res = await coinMerchantRechargeCfgApi.list({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      name: searchForm.name || undefined,
      statusFilter: searchForm.statusFilter,
    })
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantRechargeCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.statusFilter = 0
  handleSearch()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

const handleCurrentChange = () => fetchList()

const handleAdd = () => {
  editingId.value = ''
  form.value = {name: '', price: 0.99, gold: 100}
  dialogVisible.value = true
}

const handleEdit = (row: CoinMerchantRechargeCfg) => {
  editingId.value = row.id
  form.value = {
    name: row.name,
    price: truncateNumber(row.price, PRICE_DECIMALS) || 0.99,
    gold: row.gold || 1,
  }
  dialogVisible.value = true
}

const resetForm = () => {
  formRef.value?.clearValidate()
}

const handleSave = async () => {
  const ok = await formRef.value?.validate().catch(() => false)
  if (!ok) return
  submitting.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      price: Number(form.value.price),
      gold: Number(form.value.gold),
    }
    if (editingId.value) {
      await coinMerchantRechargeCfgApi.update({id: editingId.value, ...payload})
      ElMessage.success(t('pages.coinMerchantRechargeCfgList.updateSuccess'))
    } else {
      await coinMerchantRechargeCfgApi.create(payload)
      ElMessage.success(t('pages.coinMerchantRechargeCfgList.createSuccess'))
    }
    dialogVisible.value = false
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(
        editingId.value
            ? t('pages.coinMerchantRechargeCfgList.updateFailed')
            : t('pages.coinMerchantRechargeCfgList.createFailed'),
    )
  } finally {
    submitting.value = false
  }
}

const handleOnShelf = async (row: CoinMerchantRechargeCfg) => {
  try {
    await coinMerchantRechargeCfgApi.onShelf(row.id)
    ElMessage.success(t('common.operationSuccess'))
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantRechargeCfgList.operationFailed'))
  }
}

const handleOffShelf = async (row: CoinMerchantRechargeCfg) => {
  try {
    await coinMerchantRechargeCfgApi.offShelf(row.id)
    ElMessage.success(t('common.operationSuccess'))
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantRechargeCfgList.operationFailed'))
  }
}

const handleDelete = async (row: CoinMerchantRechargeCfg) => {
  try {
    await ElMessageBox.confirm(
        t('pages.coinMerchantRechargeCfgList.deleteConfirm', {name: row.name}),
        t('common.confirm'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    await coinMerchantRechargeCfgApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantRechargeCfgList.deleteFailed'))
  }
}

onMounted(fetchList)
</script>

<style scoped>
.table-header {
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
