<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CMSUserManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.cmsUserList.addUser') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.cmsUserList.username')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.cmsUserList.usernamePlaceholder')"/>
          </el-form-item>
          <el-form-item :label="t('common.status')" width="100">
            <el-select v-model="searchForm.status" clearable :placeholder="t('common.status')">
              <el-option :value="1" :label="t('common.enabled')"/>
              <el-option :value="0" :label="t('common.disabled')"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('pages.cmsUserList.admin')" width="100">
            <el-select v-model="searchForm.admin" clearable :placeholder="t('pages.cmsUserList.isAdminPlaceholder')">
              <el-option :value="true" :label="t('common.yes')"/>
              <el-option :value="false" :label="t('common.no')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchCMSUserList">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.cmsUserList.username')" prop="name"/>
          <el-table-column :label="t('common.status')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.cmsUserList.admin')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.admin ? 'primary' : 'info'">
                {{ row.admin ? t('common.yes') : t('common.no') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.cmsUserList.role')" prop="roleName" width="100"/>
          <el-table-column :label="t('common.remark')" min-width="160" prop="remark" show-overflow-tooltip/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column :label="t('common.actions')" width="300">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button size="small" type="warning" @click="handleResetPassword(row)">
                {{ t('pages.cmsUserList.resetPassword') }}
              </el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.cmsUserList.username')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.cmsUserList.enterUsername')"/>
        </el-form-item>
        <el-form-item v-if="!currentRow.id" :label="t('pages.cmsUserList.password')" prop="pwd">
          <div class="pwd-field">
            <el-input
                v-model="currentRow.pwd"
                :placeholder="t('pages.cmsUserList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateRandomPassword">{{ t('pages.cmsUserList.randomGenerate') }}</el-button>
            <el-button :disabled="!currentRow.pwd" @click="copyPassword">{{ t('pages.cmsUserList.copyPassword') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item v-else :label="t('pages.cmsUserList.password')" prop="pwd">
          <el-input v-model="currentRow.pwd" :placeholder="t('pages.cmsUserList.leavePasswordEmpty')" type="password"/>
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">{{ t('common.enabled') }}</el-radio>
            <el-radio :label="0">{{ t('common.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('pages.cmsUserList.admin')" prop="admin">
          <el-switch
              v-model="currentRow.admin"
              :active-text="t('common.yes')"
              :inactive-text="t('common.no')"
              inline-prompt
          />
        </el-form-item>
        <el-form-item :label="t('pages.cmsUserList.role')" prop="roleId">
          <el-select
              v-model="currentRow.roleId"
              :disabled="currentRow.admin"
              clearable
              filterable
              :placeholder="t('pages.cmsUserList.selectRole')"
              style="width: 100%"
          >
            <el-option
                v-for="role in roleOptions"
                :key="String(role.id)"
                :label="role.name"
                :value="String(role.id)"
            />
          </el-select>
          <div v-if="currentRow.admin" class="field-hint">{{ t('pages.cmsUserList.adminRoleHint') }}</div>
          <div v-else-if="roleOptions.length === 0" class="field-hint warn">{{ t('pages.cmsUserList.noRolesHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.remark')" prop="remark">
          <el-input
              v-model="currentRow.remark"
              :rows="3"
              maxlength="512"
              :placeholder="t('pages.cmsUserList.enterRemark')"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="resetPasswordDialogVisible"
        :close-on-click-modal="false"
        destroy-on-close
        :title="t('pages.cmsUserList.resetPasswordTitle')"
        width="560px"
        @closed="resetResetPasswordForm"
    >
      <el-form ref="resetPasswordFormRef" :model="resetPasswordForm" :rules="resetPasswordRules" label-width="100px">
        <el-form-item :label="t('pages.cmsUserList.username')">
          <el-input v-model="resetPasswordForm.name" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.cmsUserList.password')" prop="pwd">
          <div class="pwd-field">
            <el-input
                v-model="resetPasswordForm.pwd"
                :placeholder="t('pages.cmsUserList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateResetPassword">{{ t('pages.cmsUserList.randomGenerate') }}</el-button>
            <el-button :disabled="!resetPasswordForm.pwd" @click="copyResetPassword">{{ t('pages.cmsUserList.copyPassword') }}</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPasswordDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="resetPasswordSubmitting" type="primary" @click="handleResetPasswordSave">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {cmsUserApi, roleApi} from '@/api'
import type {CMSUser, Role} from '@/types/api'
import {copyTextToClipboard, createRandomPassword} from '@/utils/random-password'

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

const {t} = useI18n()

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
const resetPasswordDialogVisible = ref(false)
const resetPasswordSubmitting = ref(false)
const resetPasswordFormRef = ref<FormInstance>()
const resetPasswordForm = ref({
  id: '',
  name: '',
  pwd: '',
  status: 1,
  admin: false,
  roleId: '',
  remark: '',
})

const resetPasswordRules = computed<FormRules>(() => ({
  pwd: [
    {required: true, message: t('pages.cmsUserList.passwordRequired'), trigger: 'blur'},
    {min: 6, max: 20, message: t('pages.cmsUserList.passwordLength'), trigger: 'blur'},
  ],
}))

const validateRoleId = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (currentRow.value.admin) {
    callback()
    return
  }
  if (!value) {
    callback(new Error(t('pages.cmsUserList.roleRequired')))
    return
  }
  callback()
}

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.cmsUserList.usernameRequired'), trigger: 'blur'},
    {min: 2, max: 20, message: t('pages.cmsUserList.usernameLength'), trigger: 'blur'}
  ],
  pwd: currentRow.value.id
      ? [{min: 6, max: 20, message: t('pages.cmsUserList.passwordLength'), trigger: 'blur'}]
      : [
        {required: true, message: t('pages.cmsUserList.passwordRequired'), trigger: 'blur'},
        {min: 6, max: 20, message: t('pages.cmsUserList.passwordLength'), trigger: 'blur'}
      ],
  roleId: [{validator: validateRoleId, trigger: 'change'}],
}))

watch(
    () => currentRow.value.admin,
    (isAdmin) => {
      if (isAdmin) {
        currentRow.value.roleId = ''
        formRef.value?.clearValidate('roleId')
      }
    },
)

const generateRandomPassword = () => {
  currentRow.value.pwd = createRandomPassword()
  formRef.value?.validateField('pwd')
}

const generateResetPassword = () => {
  resetPasswordForm.value.pwd = createRandomPassword()
  resetPasswordFormRef.value?.validateField('pwd')
}

const copyResetPassword = async () => {
  try {
    const copied = await copyTextToClipboard(resetPasswordForm.value.pwd)
    if (!copied) {
      ElMessage.warning(t('pages.cmsUserList.nothingToCopy'))
      return
    }
    ElMessage.success(t('pages.cmsUserList.passwordCopied'))
  } catch (error) {
    console.error('Copy password failed:', error)
    ElMessage.error(t('pages.cmsUserList.copyFailed'))
  }
}

const resetResetPasswordForm = () => {
  resetPasswordForm.value = {
    id: '',
    name: '',
    pwd: '',
    status: 1,
    admin: false,
    roleId: '',
    remark: '',
  }
  resetPasswordFormRef.value?.clearValidate()
}

const handleResetPassword = (row: CMSUser) => {
  resetPasswordForm.value = {
    id: String(row.id),
    name: row.name,
    pwd: createRandomPassword(),
    status: row.status,
    admin: row.admin,
    roleId: row.roleId ? String(row.roleId) : '',
    remark: row.remark || '',
  }
  resetPasswordDialogVisible.value = true
}

const handleResetPasswordSave = async () => {
  if (!resetPasswordFormRef.value) {
    return
  }
  await resetPasswordFormRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    resetPasswordSubmitting.value = true
    try {
      const {id, name, pwd, status, admin, roleId, remark} = resetPasswordForm.value
      await cmsUserApi.updateCMSUser({
        id,
        name,
        pwd,
        status,
        admin,
        roleId,
        remark,
      })
      ElMessage.success(t('pages.cmsUserList.resetPasswordSuccess'))
      resetPasswordDialogVisible.value = false
    } catch (error) {
      console.error('Reset password failed:', error)
      ElMessage.error(t('pages.cmsUserList.resetPasswordFailed'))
    } finally {
      resetPasswordSubmitting.value = false
    }
  })
}

const copyPassword = async () => {
  try {
    const copied = await copyTextToClipboard(currentRow.value.pwd)
    if (!copied) {
      ElMessage.warning(t('pages.cmsUserList.nothingToCopy'))
      return
    }
    ElMessage.success(t('pages.cmsUserList.passwordCopied'))
  } catch (error) {
    console.error('Copy password failed:', error)
    ElMessage.error(t('pages.cmsUserList.copyFailed'))
  }
}

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
    console.error('Failed to load CMS users:', error)
    ElMessage.error(t('pages.cmsUserList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchCMSUserList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchCMSUserList()
}

const openDialog = async () => {
  await fetchRoleList()
  dialogVisible.value = true
}

const handleAdd = () => {
  dialogTitle.value = t('pages.cmsUserList.addDialogTitle')
  currentRow.value = {
    id: '',
    name: '',
    pwd: createRandomPassword(),
    status: 1,
    admin: false,
    roleId: '',
    remark: ''
  }
  openDialog()
}

const handleEdit = (row: CMSUser) => {
  dialogTitle.value = t('pages.cmsUserList.editDialogTitle')
  currentRow.value = {
    id: String(row.id),
    name: row.name,
    pwd: '',
    status: row.status,
    admin: row.admin,
    roleId: row.roleId ? String(row.roleId) : '',
    remark: row.remark || ''
  }
  openDialog()
}

const handleDelete = async (row: CMSUser) => {
  try {
    await ElMessageBox.confirm(
        t('pages.cmsUserList.deleteConfirm', {name: row.name}),
        t('pages.cmsUserList.deleteTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        })

    await cmsUserApi.deleteCMSUser(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchCMSUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentRow.value.id) {
          const {id, name, pwd, status, admin, roleId, remark} = currentRow.value
          await cmsUserApi.updateCMSUser({
            id: String(id),
            name,
            pwd: pwd || undefined,
            status,
            admin,
            roleId,
            remark,
          })
        } else {
          const {name, pwd, status, admin, roleId, remark} = currentRow.value
          await cmsUserApi.createCMSUser({
            name,
            pwd,
            status,
            admin,
            roleId,
            remark,
          })
        }

        ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
        dialogVisible.value = false
        fetchCMSUserList()
      } catch (error) {
        console.error('Save failed:', error)
        ElMessage.error(currentRow.value.id ? t('pages.cmsUserList.updateFailed') : t('pages.cmsUserList.createFailed'))
      }
    }
  })
}

const fetchRoleList = async () => {
  try {
    const response = await roleApi.getAllRoles()
    const list = response?.data ?? []
    roleOptions.value = list.filter(role => role.status === 1)
  } catch (error) {
    console.error('Failed to load roles:', error)
    roleOptions.value = []
    ElMessage.error(t('pages.cmsUserList.fetchRolesFailed'))
  }
}

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
