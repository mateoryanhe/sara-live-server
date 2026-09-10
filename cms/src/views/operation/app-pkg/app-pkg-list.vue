<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AppPkgManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.appPkgList.addAppPkg') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.appPkgList.packageName')">
            <el-input v-model="searchForm.packageName" clearable :placeholder="t('pages.appPkgList.packageNameFuzzy')"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.appPkgList.packageName')" min-width="220" prop="packageName" show-overflow-tooltip/>
          <el-table-column :label="t('pages.appPkgList.attributionEnabled')" prop="attributionEnabled" width="110">
            <template #default="{ row }">
              <el-tag :type="row.attributionEnabled ? 'success' : 'info'" size="small">
                {{ row.attributionEnabled ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.appPkgList.attributionProvider')" min-width="120" prop="attributionProvider" show-overflow-tooltip/>
          <el-table-column :label="t('common.remark')" min-width="160" prop="remark" show-overflow-tooltip/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="680px" @closed="activeTab = 'basic'">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="140px">
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="t('pages.appPkgList.tabBasic')" name="basic">
            <el-form-item :label="t('pages.appPkgList.packageName')" prop="packageName">
              <el-input v-model="currentRow.packageName" :placeholder="t('pages.appPkgList.packageNamePlaceholder')"/>
            </el-form-item>
            <el-form-item :label="t('common.remark')" prop="remark">
              <el-input v-model="currentRow.remark" :rows="3" :placeholder="t('pages.appPkgList.remarkOptional')" type="textarea"/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane :label="t('pages.appPkgList.tabAppsFlyer')" name="appsFlyer">
            <el-form-item :label="t('pages.appPkgList.attributionEnabled')" prop="attributionEnabled">
              <el-switch v-model="currentRow.attributionEnabled"/>
            </el-form-item>
            <el-form-item :label="t('pages.appPkgList.attributionProvider')" prop="attributionProvider">
              <el-select
                  v-model="currentRow.attributionProvider"
                  clearable
                  allow-create
                  filterable
                  :placeholder="t('pages.appPkgList.attributionProviderPlaceholder')"
                  style="width: 100%"
              >
                <el-option label="appsFlyer" value="appsFlyer"/>
              </el-select>
            </el-form-item>
            <el-form-item :label="t('pages.appPkgList.appsFlyerDevKey')" prop="appsFlyerDevKey">
              <el-input v-model="currentRow.appsFlyerDevKey" clearable :placeholder="t('pages.appPkgList.appsFlyerDevKeyPlaceholder')"/>
            </el-form-item>
            <el-form-item :label="t('pages.appPkgList.appsFlyerAppId')" prop="appsFlyerAppId">
              <el-input v-model="currentRow.appsFlyerAppId" clearable :placeholder="t('pages.appPkgList.appsFlyerAppIdPlaceholder')"/>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
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
import {appPkgApi} from '@/api'
import type {AppPkg} from '@/types/api.ts'

interface SearchForm {
  packageName: string
}

interface AppPkgForm {
  id: string
  packageName: string
  remark: string
  attributionEnabled: boolean
  attributionProvider: string
  appsFlyerDevKey: string
  appsFlyerAppId: string
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<AppPkg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const activeTab = ref('basic')

const searchForm = reactive<SearchForm>({
  packageName: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): AppPkgForm => ({
  id: '',
  packageName: '',
  remark: '',
  attributionEnabled: false,
  attributionProvider: '',
  appsFlyerDevKey: '',
  appsFlyerAppId: ''
})
const currentRow = ref<AppPkgForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  packageName: [
    {required: true, message: t('pages.appPkgList.packageNameRequired'), trigger: 'blur'},
    {min: 1, max: 128, message: t('pages.appPkgList.packageNameLength'), trigger: 'blur'}
  ]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await appPkgApi.getAppPkgList({
      packageName: searchForm.packageName.trim(),
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch app pkg list failed:', error)
    ElMessage.error(t('pages.appPkgList.fetchFailed'))
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

const resetSearch = () => {
  searchForm.packageName = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.appPkgList.addAppPkg')
  currentRow.value = defaultForm()
  activeTab.value = 'basic'
  dialogVisible.value = true
}

const handleEdit = (row: AppPkg) => {
  dialogTitle.value = t('pages.appPkgList.editAppPkg')
  currentRow.value = {
    id: row.id,
    packageName: row.packageName,
    remark: row.remark || '',
    attributionEnabled: !!row.attributionEnabled,
    attributionProvider: row.attributionProvider || '',
    appsFlyerDevKey: row.appsFlyerDevKey || '',
    appsFlyerAppId: row.appsFlyerAppId || ''
  }
  activeTab.value = 'basic'
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      activeTab.value = 'basic'
      return
    }
    try {
      const payload = {
        packageName: currentRow.value.packageName.trim(),
        remark: currentRow.value.remark.trim(),
        attributionEnabled: currentRow.value.attributionEnabled,
        attributionProvider: currentRow.value.attributionProvider.trim(),
        appsFlyerDevKey: currentRow.value.appsFlyerDevKey.trim(),
        appsFlyerAppId: currentRow.value.appsFlyerAppId.trim()
      }
      if (currentRow.value.id) {
        await appPkgApi.updateAppPkg({
          id: currentRow.value.id,
          ...payload
        })
        ElMessage.success(t('common.updateSuccess'))
      } else {
        await appPkgApi.createAppPkg(payload)
        ElMessage.success(t('common.createSuccess'))
      }
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save app pkg failed:', error)
      ElMessage.error(t('pages.appPkgList.saveFailed'))
    }
  })
}

const handleDelete = async (row: AppPkg) => {
  try {
    await ElMessageBox.confirm(t('pages.appPkgList.deleteConfirm', {name: row.packageName}), t('pages.appPkgList.promptTitle'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await appPkgApi.deleteAppPkg(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete app pkg failed:', error)
    }
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
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
