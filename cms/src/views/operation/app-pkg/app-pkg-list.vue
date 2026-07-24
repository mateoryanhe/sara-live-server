<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>App包管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增App包</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="包名">
            <el-input v-model="searchForm.packageName" clearable placeholder="包名(模糊匹配)"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="包名" min-width="220" prop="packageName" show-overflow-tooltip/>
          <el-table-column label="密钥" min-width="300">
            <template #default="{ row }">
              <div class="secret-key-cell">
                <span
                    :title="isSecretKeyVisible(row.id) ? '点击隐藏' : '点击显示'"
                    class="secret-key-text"
                    @click="toggleSecretKeyVisible(row.id)"
                >
                  {{ formatSecretKeyDisplay(row) }}
                </span>
                <el-button link type="primary" @click="copySecretKey(row.secretKey)">复制</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="180" prop="remark" show-overflow-tooltip/>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
        <el-form-item label="包名" prop="packageName">
          <el-input v-model="currentRow.packageName" placeholder="例如 com.example.app"/>
        </el-form-item>
        <el-form-item label="密钥" prop="secretKey">
          <div class="secret-key-input">
            <el-input v-model="currentRow.secretKey" placeholder="请输入密钥或自动生成" show-password type="password"/>
            <el-button @click="generateSecretKeyForForm">自动生成</el-button>
          </div>
        </el-form-item>
        <el-form-item label="隐私政策 URL" prop="privacyPolicyUrl">
          <el-input v-model="currentRow.privacyPolicyUrl" clearable placeholder="留空则使用全局法律合规配置"/>
        </el-form-item>
        <el-form-item label="用户协议 URL" prop="termsOfServiceUrl">
          <el-input v-model="currentRow.termsOfServiceUrl" clearable placeholder="留空则使用全局法律合规配置"/>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="currentRow.remark" :rows="3" placeholder="可选" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
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
    ElMessage.warning('无可复制内容')
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success('已复制')
  } catch (error) {
    console.error('复制密钥失败:', error)
    ElMessage.error('复制失败')
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
    callback(new Error('URL 长度不能超过 512'))
    return
  }
  if (!/^https?:\/\//i.test(url)) {
    callback(new Error('URL 需以 http:// 或 https:// 开头'))
    return
  }
  callback()
}

const formRules: FormRules = {
  packageName: [
    {required: true, message: '请输入包名', trigger: 'blur'},
    {min: 1, max: 128, message: '包名长度在1-128个字符', trigger: 'blur'}
  ],
  secretKey: [
    {required: true, message: '请输入密钥', trigger: 'blur'},
    {min: 1, max: 256, message: '密钥长度在1-256个字符', trigger: 'blur'}
  ],
  privacyPolicyUrl: [{validator: validateOptionalUrl, trigger: 'blur'}],
  termsOfServiceUrl: [{validator: validateOptionalUrl, trigger: 'blur'}]
}

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
    console.error('获取App包列表失败:', error)
    ElMessage.error('获取App包列表失败')
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
  dialogTitle.value = '新增App包'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: AppPkg) => {
  dialogTitle.value = '编辑App包'
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
        ElMessage.success('更新成功')
      } else {
        await appPkgApi.createAppPkg(payload)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存App包失败:', error)
    }
  })
}

const handleDelete = async (row: AppPkg) => {
  try {
    await ElMessageBox.confirm(`确定删除App包「${row.packageName}」吗？`, '提示', {
      type: 'warning'
    })
    await appPkgApi.deleteAppPkg(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除App包失败:', error)
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
