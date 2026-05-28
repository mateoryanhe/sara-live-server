<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>App Token</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="用户ID">
          <el-input
              v-model="searchForm.userId"
              clearable
              placeholder="请输入用户ID"
              style="width: 220px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <div class="table-header">
        <el-button type="primary" @click="handleAdd">新增 Token</el-button>
      </div>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="用户ID" prop="id" width="220"/>
        <el-table-column label="Token" min-width="320" prop="token" show-overflow-tooltip/>
        <el-table-column label="过期时间" prop="expireAt" width="200">
          <template #default="scope">
            {{ formatDate(scope.row.expireAt) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="expired" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.expired" type="danger">已过期</el-tag>
            <el-tag v-else type="success">有效</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
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
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="用户ID" prop="id">
          <el-input
              v-model="currentRow.id"
              :disabled="isEdit"
              placeholder="请输入用户ID"
          />
        </el-form-item>
        <el-form-item label="Token" prop="token">
          <el-input
              v-model="currentRow.token"
              placeholder="新增留空自动生成,编辑留空则保持不变"
              type="textarea"
          />
        </el-form-item>
        <el-form-item label="过期时间" prop="expireAt">
          <el-date-picker
              v-model="currentRow.expireAt"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="留空则默认7天或保持原值"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import appTokenApi from '@/api/modules/appToken'
import type {AppToken, SaveAppTokenReq} from '@/types/api'

const tableData = ref<AppToken[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  userId: '',
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增 Token')
const isEdit = ref(false)
const currentRow = ref<SaveAppTokenReq>({
  id: '',
  token: '',
  expireAt: '',
})
const formRef = ref()

const formRules = reactive({
  id: [{required: true, message: '请输入用户ID', trigger: 'blur'}],
})

const formatDate = (value?: string | null) => {
  if (!value) {
    return '-'
  }
  return value.replace('T', ' ').slice(0, 19)
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await appTokenApi.getAppToken({
      userId: searchForm.userId || undefined,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取App Token列表失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleReset = () => {
  searchForm.userId = ''
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
  dialogTitle.value = '新增 Token'
  isEdit.value = false
  currentRow.value = {
    id: searchForm.userId || '',
    token: '',
    expireAt: '',
  }
  dialogVisible.value = true
}

const handleEdit = (row: AppToken) => {
  dialogTitle.value = '编辑 Token'
  isEdit.value = true
  currentRow.value = {
    id: row.id,
    token: row.token,
    expireAt: formatDate(row.expireAt) === '-' ? '' : formatDate(row.expireAt),
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    const payload: SaveAppTokenReq = {
      id: currentRow.value.id,
      token: currentRow.value.token || undefined,
      expireAt: currentRow.value.expireAt || undefined,
    }
    const response = await appTokenApi.saveAppToken(payload)
    if (response) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      await fetchList()
    } else {
      ElMessage.error('保存失败')
    }
  } catch (error) {
    console.error('保存App Token失败:', error)
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
}

.search-form {
  margin-bottom: 15px;
}

.table-header {
  margin-bottom: 15px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
