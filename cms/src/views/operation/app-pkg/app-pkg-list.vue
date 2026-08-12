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
          <el-table-column :label="t('pages.appPkgList.secretKey')" min-width="300">
            <template #default="{ row }">
              <div class="secret-key-cell">
                <span
                    :title="isSecretKeyVisible(row.id) ? t('pages.appPkgList.clickToHide') : t('pages.appPkgList.clickToShow')"
                    class="secret-key-text"
                    @click="toggleSecretKeyVisible(row.id)"
                >
                  {{ formatSecretKeyDisplay(row) }}
                </span>
                <el-button link type="primary" @click="copySecretKey(row.secretKey)">{{ t('common.copy') }}</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.remark')" min-width="180" prop="remark" show-overflow-tooltip/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="120px">
        <el-form-item :label="t('pages.appPkgList.packageName')" prop="packageName">
          <el-input v-model="currentRow.packageName" :placeholder="t('pages.appPkgList.packageNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.appPkgList.secretKey')" prop="secretKey">
          <div class="secret-key-input">
            <el-input v-model="currentRow.secretKey" :placeholder="t('pages.appPkgList.secretKeyPlaceholder')" show-password type="password"/>
            <el-button @click="generateSecretKeyForForm">{{ t('pages.appPkgList.autoGenerate') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.appPkgList.privacyPolicyUrl')" prop="privacyPolicyUrl">
          <el-input v-model="currentRow.privacyPolicyUrl" clearable :placeholder="t('pages.appPkgList.privacyPolicyPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.appPkgList.termsOfServiceUrl')" prop="termsOfServiceUrl">
          <el-input v-model="currentRow.termsOfServiceUrl" clearable :placeholder="t('pages.appPkgList.termsPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.remark')" prop="remark">
          <el-input v-model="currentRow.remark" :rows="3" :placeholder="t('pages.appPkgList.remarkOptional')" type="textarea"/>
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
import {appPkgApi} from '@/api'
import type {AppPkg} from '@/types/api.ts'

interface SearchForm {
  packageName: string
}

interface AppPkgForm {
  id: string
  packageName: string
  secretKey: string
  privacyPolicyUrl: string
  termsOfServiceUrl: string
  remark: string
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<AppPkg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  packageName: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): AppPkgForm => ({
  id: '',
  packageName: '',
  secretKey: '',
  privacyPolicyUrl: '',
  termsOfServiceUrl: '',
  remark: ''
})
const currentRow = ref<AppPkgForm>(defaultForm())
const formRef = ref<FormInstance>()
const visibleSecretKeyIds = ref<Set<string>>(new Set())

const generateSecretKey = () => crypto.randomUUID().replace(/-/g, '')

const generateSecretKeyForForm = () => {
  currentRow.value.secretKey = generateSecretKey()
}

const isSecretKeyVisible = (id: string) => visibleSecretKeyIds.value.has(id)

const toggleSecretKeyVisible = (id: string) => {
  const next = new Set(visibleSecretKeyIds.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  visibleSecretKeyIds.value = next
}

const maskSecretKey = (value: string) => {
  if (!value) {
    return '-'
  }
  return '•'.repeat(16)
}

const formatSecretKeyDisplay = (row: AppPkg) => {
  if (!row.secretKey) {
    return '-'
  }
  return isSecretKeyVisible(row.id) ? row.secretKey : maskSecretKey(row.secretKey)
}

const copySecretKey = async (value?: string) => {
  if (!value) {
    ElMessage.warning(t('pages.appPkgList.nothingToCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(t('pages.appPkgList.copied'))
  } catch (error) {
    console.error('copy secret key failed:', error)
    ElMessage.error(t('pages.appPkgList.copyFailed'))
  }
}

const clearVisibleSecretKeys = () => {
  visibleSecretKeyIds.value = new Set()
}

const validateOptionalUrl = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const url = value?.trim()
  if (!url) {
    callback()
    return
  }
  if (url.length > 512) {
    callback(new Error(t('pages.appPkgList.urlMaxLength')))
    return
  }
  if (!/^https?:\/\//i.test(url)) {
    callback(new Error(t('pages.appPkgList.urlMustHttp')))
    return
  }
  callback()
}

const formRules = computed<FormRules>(() => ({
  packageName: [
    {required: true, message: t('pages.appPkgList.packageNameRequired'), trigger: 'blur'},
    {min: 1, max: 128, message: t('pages.appPkgList.packageNameLength'), trigger: 'blur'}
  ],
  secretKey: [
    {required: true, message: t('pages.appPkgList.secretKeyRequired'), trigger: 'blur'},
    {min: 1, max: 256, message: t('pages.appPkgList.secretKeyLength'), trigger: 'blur'}
  ],
  privacyPolicyUrl: [{validator: validateOptionalUrl, trigger: 'blur'}],
  termsOfServiceUrl: [{validator: validateOptionalUrl, trigger: 'blur'}]
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
    clearVisibleSecretKeys()
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
  dialogVisible.value = true
}

const handleEdit = (row: AppPkg) => {
  dialogTitle.value = t('pages.appPkgList.editAppPkg')
  currentRow.value = {
    id: row.id,
    packageName: row.packageName,
    secretKey: row.secretKey,
    privacyPolicyUrl: row.privacyPolicyUrl || '',
    termsOfServiceUrl: row.termsOfServiceUrl || '',
    remark: row.remark || ''
  }
  dialogVisible.value = true
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
      const payload = {
        packageName: currentRow.value.packageName.trim(),
        secretKey: currentRow.value.secretKey.trim(),
        privacyPolicyUrl: currentRow.value.privacyPolicyUrl.trim(),
        termsOfServiceUrl: currentRow.value.termsOfServiceUrl.trim(),
        remark: currentRow.value.remark.trim()
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

.secret-key-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.secret-key-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  font-family: Consolas, Monaco, monospace;
  user-select: none;
}

.secret-key-input {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.secret-key-input .el-input {
  flex: 1;
}
</style>
