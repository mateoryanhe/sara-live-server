<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增角色</el-button>
        </div>

        <!-- 搜索表单 -->
        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="角色名称">
            <el-input v-model="searchForm.name" clearable placeholder="角色名称"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchRoleList">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="角色名称" prop="name"/>
          <el-table-column label="描述" prop="description" show-overflow-tooltip/>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column label="操作" width="280">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              <el-button size="small" type="primary" @click="handlePermissions(row)">权限</el-button>
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

    <!-- 角色编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入角色名称"/>
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input v-model="currentRow.description" placeholder="请输入角色描述" type="textarea"/>
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
import {roleApi} from '@/api'
import type {Role} from '@/types/api'
import {useRouter} from 'vue-router'

interface SearchForm {
  name: string
}

interface RoleForm {
  id: string
  name: string
  description: string
  status: number
  permissions?: string[]
}

const loading = ref(false)
const tableData = ref<Role[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<RoleForm>({
  id: '',
  name: '',
  description: '',
  status: 1
})

const formRef = ref<FormInstance>()
const router = useRouter()

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入角色名称', trigger: 'blur'},
    {min: 2, max: 20, message: '角色名称长度在2-20个字符', trigger: 'blur'}
  ],
  description: [
    {max: 200, message: '描述长度不能超过200个字符', trigger: 'blur'}
  ]
}


// 获取角色列表
const fetchRoleList = async () => {
  loading.value = true
  try {
    const response = await roleApi.getRoleList({
      name: searchForm.name,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}


// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchRoleList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchRoleList()
}

// 操作处理
const handleAdd = () => {
  dialogTitle.value = '新增角色'
  currentRow.value = {
    id: '',
    name: '',
    description: '',
    status: 1
  }
  dialogVisible.value = true
}

const handleEdit = (row: Role) => {
  dialogTitle.value = '编辑角色'
  currentRow.value = {
    id: row.id,
    name: row.name,
    description: row.description,
    status: row.status
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await roleApi.deleteRole(row.id)
    ElMessage.success('删除成功')
    fetchRoleList() // 重新获取数据
  } catch (error) {
    console.error('删除失败:', error)
  }
}

// 处理权限设置
const handlePermissions = (row: Role) => {
  // 跳转到模块管理页面，传递角色ID
  router.push({path: '/role/module-list', query: {roleId: row.id}})
}


// 保存操作
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let response
        if (currentRow.value.id) {
          // 更新角色
          response = await roleApi.updateRole(currentRow.value)
        } else {
          // 创建角色
          const {id, name, description, status} = currentRow.value
          response = await roleApi.createRole({
            name,
            description,
            status,
            permissions: [] // 权限单独设置
          })
        }

        ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchRoleList() // 重新获取数据
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
  fetchRoleList()
}

onMounted(() => {
  fetchRoleList()
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