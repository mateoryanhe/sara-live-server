<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AnchorSalaryCfgManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.anchorSalaryCfgList.addTier') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.anchorSalaryCfgList.weeklyWorkDays')" prop="weeklyWorkDays" min-width="140"/>
          <el-table-column :label="t('pages.anchorSalaryCfgList.dailyLiveDurationMinutes')" prop="dailyLiveDurationMinutes" min-width="160"/>
          <el-table-column :label="t('pages.anchorSalaryCfgList.salaryAmount')" prop="salaryAmount" min-width="120"/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="180px">
        <el-form-item :label="t('pages.anchorSalaryCfgList.weeklyWorkDays')" prop="weeklyWorkDays">
          <el-input-number v-model="currentRow.weeklyWorkDays" :max="7" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorSalaryCfgList.dailyLiveDurationMinutes')" prop="dailyLiveDurationMinutes">
          <el-input-number v-model="currentRow.dailyLiveDurationMinutes" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorSalaryCfgList.salaryAmount')" prop="salaryAmount">
          <el-input-number v-model="currentRow.salaryAmount" :min="0" :precision="4" :step="1" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.anchorSalaryCfgList.sortHigherFirst') }}</div>
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
import {computed, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {anchorSalaryCfgApi} from '@/api/modules/anchor-salary-cfg'
import type {AnchorSalaryCfg} from '@/types/api'

interface TierForm {
  id: string
  weeklyWorkDays: number
  dailyLiveDurationMinutes: number
  salaryAmount: number
  sort: number
}

const {t} = useI18n()

const loading = ref(false)
const tableData = ref<AnchorSalaryCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): TierForm => ({
  id: '',
  weeklyWorkDays: 0,
  dailyLiveDurationMinutes: 0,
  salaryAmount: 0,
  sort: 0,
})
const currentRow = ref<TierForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  weeklyWorkDays: [
    {required: true, message: t('pages.anchorSalaryCfgList.weeklyWorkDaysRequired'), trigger: 'change'},
  ],
  dailyLiveDurationMinutes: [
    {required: true, message: t('pages.anchorSalaryCfgList.dailyLiveDurationRequired'), trigger: 'change'},
  ],
  salaryAmount: [
    {required: true, message: t('pages.anchorSalaryCfgList.salaryRequired'), trigger: 'change'},
  ],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await anchorSalaryCfgApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch anchor salary cfg list failed:', error)
    ElMessage.error(t('pages.anchorSalaryCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
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
  dialogTitle.value = t('pages.anchorSalaryCfgList.addTier')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: AnchorSalaryCfg) => {
  dialogTitle.value = t('pages.anchorSalaryCfgList.editTier')
  currentRow.value = {
    id: row.id,
    weeklyWorkDays: Number(row.weeklyWorkDays) || 0,
    dailyLiveDurationMinutes: Number(row.dailyLiveDurationMinutes) || 0,
    salaryAmount: Number(row.salaryAmount) || 0,
    sort: Number(row.sort) || 0,
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const payload = {
      weeklyWorkDays: currentRow.value.weeklyWorkDays,
      dailyLiveDurationMinutes: currentRow.value.dailyLiveDurationMinutes,
      salaryAmount: currentRow.value.salaryAmount,
      sort: currentRow.value.sort,
    }
    if (currentRow.value.id) {
      await anchorSalaryCfgApi.update({
        id: currentRow.value.id,
        ...payload,
      })
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await anchorSalaryCfgApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save anchor salary cfg failed:', error)
  }
}

const handleDelete = async (row: AnchorSalaryCfg) => {
  try {
    await ElMessageBox.confirm(
        t('pages.anchorSalaryCfgList.deleteConfirm', {id: row.id}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await anchorSalaryCfgApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete anchor salary cfg failed:', error)
    }
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

.table-header {
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
}
</style>
