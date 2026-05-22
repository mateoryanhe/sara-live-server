<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>直播工会管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增工会</el-button>
        </div>

        <!-- 搜索表单 -->
        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="工会名称">
            <el-input v-model="searchForm.name" clearable placeholder="工会名称"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchGuildList">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="工会名称" prop="name"/>
          <el-table-column label="会长ID" prop="leaderId" width="140"/>
          <el-table-column label="联系方式" prop="contact" width="160"/>
          <el-table-column label="简介" prop="description" show-overflow-tooltip/>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column label="操作" width="200">
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

    <!-- 工会编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="工会名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入工会名称"/>
        </el-form-item>
        <el-form-item label="会长ID" prop="leaderId">
          <el-input-number v-model="currentRow.leaderId" :min="0" controls-position="right" placeholder="请输入会长/负责人ID"/>
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="currentRow.contact" placeholder="请输入联系方式"/>
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input v-model="currentRow.description" placeholder="请输入工会简介" type="textarea"/>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
import {guildApi} from '@/api'
import type {Guild} from '@/types/api.ts'

interface SearchForm {
  name: string
}

interface GuildForm {
  id: string
  name: string
  leaderId: number
  contact: string
  description: string
  status: number
}

const loading = ref(false)
const tableData = ref<Guild[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<GuildForm>({
  id: '',
  name: '',
  leaderId: 0,
  contact: '',
  description: '',
  status: 1
})

const formRef = ref<FormInstance>()

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入工会名称', trigger: 'blur'},
    {min: 2, max: 32, message: '工会名称长度在2-32个字符', trigger: 'blur'}
  ],
  description: [
    {max: 200, message: '简介长度不能超过200个字符', trigger: 'blur'}
  ]
}

// 获取工会列表
const fetchGuildList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getGuildList({
      name: searchForm.name,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取工会列表失败:', error)
    ElMessage.error('获取工会列表失败')
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGuildList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGuildList()
}

// 操作处理
const handleAdd = () => {
  dialogTitle.value = '新增工会'
  currentRow.value = {
    id: '',
    name: '',
    leaderId: 0,
    contact: '',
    description: '',
    status: 1
  }
  dialogVisible.value = true
}

const handleEdit = (row: Guild) => {
  dialogTitle.value = '编辑工会'
  currentRow.value = {
    id: row.id,
    name: row.name,
    leaderId: Number(row.leaderId) || 0,
    contact: row.contact,
    description: row.description,
    status: row.status
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Guild) => {
  try {
    await ElMessageBox.confirm(`确定要删除工会 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await guildApi.deleteGuild(row.id)
    ElMessage.success('删除成功')
    fetchGuildList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

// 保存操作
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentRow.value.id) {
          await guildApi.updateGuild(currentRow.value)
        } else {
          const {name, leaderId, contact, description, status} = currentRow.value
          await guildApi.createGuild({name, leaderId, contact, description, status})
        }

        ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchGuildList()
      } catch (error) {
        console.error(currentRow.value.id ? '更新失败:' : '创建失败:', error)
        ElMessage.error(currentRow.value.id ? '更新失败' : '创建失败')
      }
    }
  })
}

// 重置搜索
const resetSearch = () => {
  searchForm.name = ''
  fetchGuildList()
}

onMounted(() => {
  fetchGuildList()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
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
</style>
