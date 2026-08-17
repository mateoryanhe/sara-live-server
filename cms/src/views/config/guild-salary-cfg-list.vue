<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildSalaryCfgManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.guildSalaryCfgList.addTier') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.guildSalaryCfgList.dailyEffectiveLiveCount')" prop="dailyEffectiveLiveCount" min-width="160"/>
          <el-table-column :label="t('pages.guildSalaryCfgList.weeklyEffectiveLiveCount')" prop="weeklyEffectiveLiveCount" min-width="160"/>
          <el-table-column :label="t('pages.guildSalaryCfgList.salaryAmount')" prop="salaryAmount" min-width="120"/>
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
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="160px">
        <el-form-item :label="t('pages.guildSalaryCfgList.dailyEffectiveLiveCount')" prop="dailyEffectiveLiveCount">
          <el-input-number v-model="currentRow.dailyEffectiveLiveCount" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildSalaryCfgList.weeklyEffectiveLiveCount')" prop="weeklyEffectiveLiveCount">
          <el-input-number v-model="currentRow.weeklyEffectiveLiveCount" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildSalaryCfgList.salaryAmount')" prop="salaryAmount">
          <el-input-number v-model="currentRow.salaryAmount" :min="0" :precision="4" :step="1" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.guildSalaryCfgList.sortHigherFirst') }}</div>
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
import {guildSalaryCfgApi} from '@/api/modules/guild-salary-cfg'
import type {GuildSalaryCfg} from '@/types/api'

interface TierForm {
  id: string
  dailyEffectiveLiveCount: number
  weeklyEffectiveLiveCount: number
  salaryAmount: number
  sort: number
}

const {t} = useI18n()

const loading = ref(false)
const tableData = ref<GuildSalaryCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): TierForm => ({
  id: '',
  dailyEffectiveLiveCount: 0,
  weeklyEffectiveLiveCount: 0,
  salaryAmount: 0,
  sort: 0,
})
const currentRow = ref<TierForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  dailyEffectiveLiveCount: [
    {required: true, message: t('pages.guildSalaryCfgList.dailyRequired'), trigger: 'change'},
  ],
  weeklyEffectiveLiveCount: [
    {required: true, message: t('pages.guildSalaryCfgList.weeklyRequired'), trigger: 'change'},
  ],
  salaryAmount: [
    {required: true, message: t('pages.guildSalaryCfgList.salaryRequired'), trigger: 'change'},
  ],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildSalaryCfgApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch guild salary cfg list failed:', error)
    ElMessage.error(t('pages.guildSalaryCfgList.fetchFailed'))
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
  dialogTitle.value = t('pages.guildSalaryCfgList.addTier')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: GuildSalaryCfg) => {
  dialogTitle.value = t('pages.guildSalaryCfgList.editTier')
  currentRow.value = {
    id: row.id,
    dailyEffectiveLiveCount: Number(row.dailyEffectiveLiveCount) || 0,
    weeklyEffectiveLiveCount: Number(row.weeklyEffectiveLiveCount) || 0,
    salaryAmount: Number(row.salaryAmount) || 0,
    sort: Number(row.sort) || 0,
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const payload = {
      dailyEffectiveLiveCount: currentRow.value.dailyEffectiveLiveCount,
      weeklyEffectiveLiveCount: currentRow.value.weeklyEffectiveLiveCount,
      salaryAmount: currentRow.value.salaryAmount,
      sort: currentRow.value.sort,
    }
    if (currentRow.value.id) {
      await guildSalaryCfgApi.update({
        id: currentRow.value.id,
        ...payload,
      })
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await guildSalaryCfgApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save guild salary cfg failed:', error)
  }
}

const handleDelete = async (row: GuildSalaryCfg) => {
  try {
    await ElMessageBox.confirm(
        t('pages.guildSalaryCfgList.deleteConfirm', {id: row.id}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await guildSalaryCfgApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete guild salary cfg failed:', error)
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
