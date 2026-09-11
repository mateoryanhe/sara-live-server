<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.SimulatorDeviceWhitelistManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-input
              v-model="searchKey"
              clearable
              class="search-input"
              :placeholder="t('pages.simulatorDeviceWhitelistList.searchPlaceholder')"
              @keyup.enter="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button type="primary" @click="handleAdd">{{ t('pages.simulatorDeviceWhitelistList.addDevice') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.simulatorDeviceWhitelistList.deviceId')" prop="deviceId" min-width="220"/>
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
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.simulatorDeviceWhitelistList.deviceId')" prop="deviceId">
          <el-input v-model="currentRow.deviceId" maxlength="128" show-word-limit/>
          <div class="form-tip">{{ t('pages.simulatorDeviceWhitelistList.deviceIdTip') }}</div>
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
import {simulatorDeviceWhitelistApi} from '@/api/modules/simulator-device-whitelist'
import type {SimulatorDeviceWhitelist} from '@/types/api'

interface DeviceForm {
  id: string
  deviceId: string
}

const {t} = useI18n()

const loading = ref(false)
const tableData = ref<SimulatorDeviceWhitelist[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKey = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): DeviceForm => ({
  id: '',
  deviceId: '',
})
const currentRow = ref<DeviceForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  deviceId: [
    {required: true, message: t('pages.simulatorDeviceWhitelistList.deviceIdRequired'), trigger: 'blur'},
  ],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await simulatorDeviceWhitelistApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      key: searchKey.value.trim(),
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch simulator device whitelist failed:', error)
    ElMessage.error(t('pages.simulatorDeviceWhitelistList.fetchFailed'))
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
  dialogTitle.value = t('pages.simulatorDeviceWhitelistList.addDevice')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: SimulatorDeviceWhitelist) => {
  dialogTitle.value = t('pages.simulatorDeviceWhitelistList.editDevice')
  currentRow.value = {
    id: row.id,
    deviceId: row.deviceId || '',
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const deviceId = currentRow.value.deviceId.trim()
    if (currentRow.value.id) {
      await simulatorDeviceWhitelistApi.update({
        id: currentRow.value.id,
        deviceId,
      })
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await simulatorDeviceWhitelistApi.create({deviceId})
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save simulator device whitelist failed:', error)
  }
}

const handleDelete = async (row: SimulatorDeviceWhitelist) => {
  try {
    await ElMessageBox.confirm(
        t('pages.simulatorDeviceWhitelistList.deleteConfirm', {deviceId: row.deviceId}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await simulatorDeviceWhitelistApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete simulator device whitelist failed:', error)
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
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
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
  line-height: 1.4;
}
</style>
