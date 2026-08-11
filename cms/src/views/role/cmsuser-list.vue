<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>CMS用户管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增用户</el-button>
        </div>

        <!-- 搜索表单 -->
        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="用户名">
            <el-input v-model="searchForm.name" clearable placeholder="用户名"/>
          </el-form-item>
          <el-form-item label="状态" width="100">
            <el-select v-model="searchForm.status" clearable placeholder="状态">
              <el-option :value="1" label="启用"/>
              <el-option :value="0" label="禁用"/>
            </el-select>
          </el-form-item>
          <el-form-item label="管理员" width="100">
            <el-select v-model="searchForm.admin" clearable placeholder="是否管理员">
              <el-option :value="true" label="是"/>
              <el-option :value="false" label="否"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchCMSUserList">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="用户名" prop="name"/>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="管理员" width="100">
            <template #default="{ row }">
              <el-tag :type="row.admin ? 'primary' : 'info'">
                {{ row.admin ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="角色" prop="roleName" width="100"/>
          <el-table-column label="备注" min-width="160" prop="remark" show-overflow-tooltip/>
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

    <!-- CMS用户编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="用户名" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入用户名"/>
        </el-form-item>
        <el-form-item v-if="!currentRow.id" label="密码" prop="pwd">
          <div class="pwd-field">
            <el-input
                v-model="currentRow.pwd"
                placeholder="请输入密码或点击随机生成"
                show-password
                type="password"
            />
            <el-button @click="generateRandomPassword">随机生成</el-button>
            <el-button :disabled="!currentRow.pwd" @click="copyPassword">复制</el-button>
          </div>
        </el-form-item>
        <el-form-item v-else label="密码" prop="pwd">
          <el-input v-model="currentRow.pwd" placeholder="留空则不修改密码" type="password"/>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="管理员" prop="admin">
          <el-switch
              v-model="currentRow.admin"
              active-text="是"
              inactive-text="否"
              inline-prompt
          />
        </el-form-item>
        <el-form-item label="角色" prop="roleId">
          <el-select
              v-model="currentRow.roleId"
              :disabled="currentRow.admin"
              clearable
              filterable
              placeholder="请选择角色"
              style="width: 100%"
          >
            <el-option
                v-for="role in roleOptions"
                :key="String(role.id)"
                :label="role.name"
                :value="String(role.id)"
            />
          </el-select>
          <div v-if="currentRow.admin" class="field-hint">管理员拥有全部权限，无需分配角色</div>
          <div v-else-if="roleOptions.length === 0" class="field-hint warn">暂无可用角色，请先在「角色权限管理」中创建</div>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
              v-model="currentRow.remark"
              :rows="3"
              maxlength="512"
              placeholder="请输入备注"
              show-word-limit
              type="textarea"
          />
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
import {onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {cmsUserApi, roleApi} from '@/api'
import type {CMSUser, Role} from '@/types/api'

interface SearchForm {
  name: string
  status: number | null
  admin: boolean | null
}

interface CMSUserForm {
  id: string
  name: string
  pwd: string
  status: number
  admin: boolean
  roleId: string
  remark: string
}

const loading = ref(false)
const tableData = ref<CMSUser[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const roleOptions = ref<Role[]>([])

const searchForm = reactive<SearchForm>({
  name: '',
  status: null,
  admin: null
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<CMSUserForm>({
  id: '',
  name: '',
  pwd: '',
  status: 1,
  admin: false,
  roleId: '',
  remark: ''
})

const formRef = ref<FormInstance>()

const validateRoleId = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (currentRow.value.admin) {
    callback()
    return
  }
  if (!value) {
    callback(new Error('请选择角色'))
    return
  }
  callback()
}

const formRules = ref<FormRules>({
  name: [
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 2, max: 20, message: '用户名长度在2-20个字符', trigger: 'blur'}
  ],
  pwd: [],
  roleId: [{validator: validateRoleId, trigger: 'change'}],
})

watch(
    () => currentRow.value.admin,
    (isAdmin) => {
      if (isAdmin) {
        currentRow.value.roleId = ''
        formRef.value?.clearValidate('roleId')
      }
    },
)

const PASSWORD_CHARS = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789'

const createRandomPassword = (length = 20) => {
  const values = crypto.getRandomValues(new Uint32Array(length))
  return Array.from(values, (value) => PASSWORD_CHARS[value % PASSWORD_CHARS.length]).join('')
}

const generateRandomPassword = () => {
  currentRow.value.pwd = createRandomPassword()
  formRef.value?.validateField('pwd')
}

const copyPassword = async () => {
  const value = currentRow.value.pwd?.trim()
  if (!value) {
    ElMessage.warning('无可复制内容')
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success('密码已复制')
  } catch (error) {
    console.error('复制密码失败:', error)
    ElMessage.error('复制失败')
  }
}


// 获取CMS用户列表
const fetchCMSUserList = async () => {
  loading.value = true
  try {
    const response = await cmsUserApi.getCMSUserList({
      name: searchForm.name,
      status: searchForm.status || undefined,
      admin: searchForm.admin ?? undefined,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取CMS用户列表失败:', error)
    ElMessage.error('获取CMS用户列表失败')
  } finally {
    loading.value = false
  }
}


// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchCMSUserList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchCMSUserList()
}

// 操作处理
const openDialog = async () => {
  await fetchRoleList()
  dialogVisible.value = true
}

const handleAdd = () => {
  dialogTitle.value = '新增用户'
  currentRow.value = {
    id: '',
    name: '',
    pwd: createRandomPassword(),
    status: 1,
    admin: false,
    roleId: '',
    remark: ''
  }
  // 新增模式下密码为必填项
  formRules.value.pwd = [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 6, max: 20, message: '密码长度在6-20个字符', trigger: 'blur'}
  ]
  openDialog()
}

const handleEdit = (row: CMSUser) => {
  dialogTitle.value = '编辑用户'
  currentRow.value = {
    id: String(row.id),
    name: row.name,
    pwd: '',
    status: row.status,
    admin: row.admin,
    roleId: row.roleId ? String(row.roleId) : '',
    remark: row.remark || ''
  }
  // 编辑模式下密码非必填
  formRules.value.pwd = [
    {min: 6, max: 20, message: '密码长度在6-20个字符', trigger: 'blur'}
  ]
  openDialog()
}

const handleDelete = async (row: CMSUser) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await cmsUserApi.deleteCMSUser(row.id)
    ElMessage.success('删除成功')
    fetchCMSUserList() // 重新获取数据
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
        let response
        if (currentRow.value.id) {
          // 更新用户
          const {id, name, pwd, status, admin, roleId, remark} = currentRow.value
          response = await cmsUserApi.updateCMSUser({
            id: String(id),
            name,
            pwd: pwd || undefined, // 如果密码为空则不更新
            status,
            admin,
            roleId,
            remark,
          })
        } else {
          // 创建用户
          const {name, pwd, status, admin, roleId, remark} = currentRow.value
          response = await cmsUserApi.createCMSUser({
            name,
            pwd,
            status,
            admin,
            roleId,
            remark,
          })
        }

        ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchCMSUserList() // 重新获取数据
      } catch (error) {
        console.error(currentRow.value.id ? '更新失败:' : '创建失败:', error)
        ElMessage.error(currentRow.value.id ? '更新失败' : '创建失败')
      }
    }
  })
}


// 获取角色列表（仅启用状态，供下拉选择）
const fetchRoleList = async () => {
  try {
    const response = await roleApi.getAllRoles()
    const list = response?.data ?? []
    roleOptions.value = list.filter(role => role.status === 1)
  } catch (error) {
    console.error('获取角色列表失败:', error)
    roleOptions.value = []
    ElMessage.error('获取角色列表失败')
  }
}

// 重置搜索
const resetSearch = () => {
  searchForm.name = ''
  searchForm.status = null
  searchForm.admin = null
  fetchCMSUserList()
}

onMounted(() => {
  fetchCMSUserList()
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

.pwd-field {
  display: flex;
  gap: 8px;
  width: 100%;
}

.pwd-field .el-input {
  flex: 1;
}

.field-hint {
  margin-top: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.field-hint.warn {
  color: var(--el-color-warning);
}
</style>