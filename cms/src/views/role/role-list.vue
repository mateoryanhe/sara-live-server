<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.RoleManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button v-if="can('create')" type="primary" @click="handleAdd">{{ t('pages.roleList.addRole') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.roleList.roleName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.roleList.roleNamePlaceholder')"/>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="fetchRoleList">{{ t('common.search') }}</el-button>
            <el-button v-if="can('search')" @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.roleList.roleName')" prop="name"/>
          <el-table-column :label="t('pages.roleList.description')" prop="description" show-overflow-tooltip/>
          <el-table-column :label="t('pages.roleList.roleType')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.roleType === ROLE_TYPE_EXTERNAL ? 'warning' : 'info'">
                {{ formatRoleType(row.roleType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column :label="t('common.actions')" width="280">
            <template #default="{ row }">
              <el-button v-if="can('edit')" size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button v-if="can('delete')" size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
              <el-button v-if="can('permission')" size="small" type="primary" @click="handlePermissions(row)">{{ t('pages.roleList.permissions') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.roleList.roleName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.roleList.enterRoleName')"/>
        </el-form-item>
        <el-form-item :label="t('pages.roleList.roleDescription')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="t('pages.roleList.enterRoleDescription')" type="textarea"/>
        </el-form-item>
        <el-form-item :label="t('pages.roleList.roleType')" prop="roleType">
          <el-radio-group v-model="currentRow.roleType">
            <el-radio :label="ROLE_TYPE_INTERNAL">{{ t('pages.roleList.roleTypeInternal') }}</el-radio>
            <el-radio :label="ROLE_TYPE_EXTERNAL">{{ t('pages.roleList.roleTypeExternal') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">{{ t('common.enabled') }}</el-radio>
            <el-radio :label="0">{{ t('common.disabled') }}</el-radio>
          </el-radio-group>
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
import {roleApi} from '@/api'
import type {Role} from '@/types/api'
import {useRouter} from 'vue-router'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('RoleManagement')

const ROLE_TYPE_INTERNAL = 1
const ROLE_TYPE_EXTERNAL = 2

interface SearchForm {
  name: string
}

interface RoleForm {
  id: string
  name: string
  description: string
  roleType: number
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
  roleType: ROLE_TYPE_INTERNAL,
  status: 1
})

const formRef = ref<FormInstance>()
const router = useRouter()

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.roleList.roleNameRequired'), trigger: 'blur'},
    {min: 2, max: 20, message: t('pages.roleList.roleNameLength'), trigger: 'blur'}
  ],
  description: [
    {max: 200, message: t('pages.roleList.descriptionMaxLength'), trigger: 'blur'}
  ],
  roleType: [
    {required: true, message: t('pages.roleList.roleTypeRequired'), trigger: 'change'}
  ]
}))

const formatRoleType = (roleType?: number) => {
  if (roleType === ROLE_TYPE_EXTERNAL) {
    return t('pages.roleList.roleTypeExternal')
  }
  return t('pages.roleList.roleTypeInternal')
}

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
    console.error('Failed to load role list:', error)
    ElMessage.error(t('pages.roleList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchRoleList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchRoleList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.roleList.addDialogTitle')
  currentRow.value = {
    id: '',
    name: '',
    description: '',
    roleType: ROLE_TYPE_INTERNAL,
    status: 1
  }
  dialogVisible.value = true
}

const handleEdit = (row: Role) => {
  dialogTitle.value = t('pages.roleList.editDialogTitle')
  currentRow.value = {
    id: row.id,
    name: row.name,
    description: row.description,
    roleType: row.roleType === ROLE_TYPE_EXTERNAL ? ROLE_TYPE_EXTERNAL : ROLE_TYPE_INTERNAL,
    status: row.status
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(
        t('pages.roleList.deleteConfirm', {name: row.name}),
        t('pages.roleList.deleteTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        })

    await roleApi.deleteRole(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchRoleList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
    }
  }
}

const handlePermissions = (row: Role) => {
  router.push({path: '/role/module-list', query: {roleId: row.id, roleName: row.name}})
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentRow.value.id) {
          await roleApi.updateRole(currentRow.value)
        } else {
          const {name, description, roleType, status} = currentRow.value
          await roleApi.createRole({
            name,
            description,
            roleType,
            status,
            permissions: []
          })
        }

        ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
        dialogVisible.value = false
        fetchRoleList()
      } catch (error) {
        console.error('Save failed:', error)
        ElMessage.error(currentRow.value.id ? t('pages.roleList.updateFailed') : t('pages.roleList.createFailed'))
      }
    }
  })
}

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
