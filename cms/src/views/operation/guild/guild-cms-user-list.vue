<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildCMSUserManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <el-alert
            :closable="false"
            class="tip-alert"
            show-icon
            :title="t('pages.guildCmsUserList.nonAdminHint')"
            type="info"
        />

        <div v-if="can('create')" class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.guildCmsUserList.addUser') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.guildCmsUserList.username')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.guildCmsUserList.usernamePlaceholder')"/>
          </el-form-item>
          <el-form-item :label="t('common.status')">
            <el-select v-model="searchForm.status" clearable :placeholder="t('common.status')">
              <el-option :value="1" :label="t('common.enabled')"/>
              <el-option :value="0" :label="t('common.disabled')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchList">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.guildCmsUserList.username')" prop="name"/>
          <el-table-column :label="t('common.status')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildCmsUserList.role')" prop="roleName" width="140"/>
          <el-table-column :label="t('common.remark')" min-width="160" prop="remark" show-overflow-tooltip/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column v-if="can('resetPassword')" :label="t('common.actions')" width="120">
            <template #default="{ row }">
              <el-button size="small" type="warning" @click="handleResetPassword(row)">
                {{ t('pages.guildCmsUserList.resetPassword') }}
              </el-button>
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

    <el-dialog v-model="dialogVisible" :title="t('pages.guildCmsUserList.addDialogTitle')" width="600px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.guildCmsUserList.username')" prop="name">
          <el-input v-model="form.name" :placeholder="t('pages.guildCmsUserList.enterUsername')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildCmsUserList.password')" prop="pwd">
          <div class="pwd-field">
            <el-input
                v-model="form.pwd"
                :placeholder="t('pages.guildCmsUserList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateRandomPassword">{{ t('pages.guildCmsUserList.randomGenerate') }}</el-button>
            <el-button :disabled="!form.pwd" @click="copyPassword">{{ t('pages.guildCmsUserList.copyPassword') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">{{ t('common.enabled') }}</el-radio>
            <el-radio :label="0">{{ t('common.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('pages.guildCmsUserList.role')" prop="roleId">
          <el-select
              v-model="form.roleId"
              clearable
              filterable
              :placeholder="t('pages.guildCmsUserList.selectRole')"
              style="width: 100%"
          >
            <el-option
                v-for="role in roleOptions"
                :key="String(role.id)"
                :label="role.name"
                :value="String(role.id)"
            />
          </el-select>
          <div v-if="roleOptions.length === 0" class="field-hint warn">{{ t('pages.guildCmsUserList.noRolesHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.remark')" prop="remark">
          <el-input
              v-model="form.remark"
              :rows="3"
              maxlength="512"
              :placeholder="t('pages.guildCmsUserList.enterRemark')"
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
        :title="t('pages.guildCmsUserList.resetPasswordTitle')"
        width="560px"
        @closed="resetResetPasswordForm"
    >
      <el-form ref="resetPasswordFormRef" :model="resetPasswordForm" :rules="resetPasswordRules" label-width="100px">
        <el-form-item :label="t('pages.guildCmsUserList.username')">
          <el-input v-model="resetPasswordForm.name" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.guildCmsUserList.password')" prop="pwd">
          <div class="pwd-field">
            <el-input
                v-model="resetPasswordForm.pwd"
                :placeholder="t('pages.guildCmsUserList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateResetPassword">{{ t('pages.guildCmsUserList.randomGenerate') }}</el-button>
            <el-button :disabled="!resetPasswordForm.pwd" @click="copyResetPassword">{{ t('pages.guildCmsUserList.copyPassword') }}</el-button>
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
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {cmsUserApi, roleApi} from '@/api'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Role} from '@/api/modules/role'
import {usePagePermission} from '@/composables/usePagePermission'
import {copyTextToClipboard, createRandomPassword} from '@/utils/random-password'

const EXTERNAL_ROLE_TYPE = 2

interface SearchForm {
  name: string
  status: number | null
}

interface CMSUserForm {
  name: string
  pwd: string
  status: number
  roleId: string
  remark: string
}

interface ResetPasswordForm {
  id: string
  name: string
  pwd: string
  status: number
  roleId: string
  remark: string
}

const {t} = useI18n()
const {can} = usePagePermission('GuildCMSUserManagement')

const loading = ref(false)
const tableData = ref<CMSUser[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const roleOptions = ref<Role[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const resetPasswordDialogVisible = ref(false)
const resetPasswordSubmitting = ref(false)
const resetPasswordFormRef = ref<FormInstance>()
const resetPasswordForm = ref<ResetPasswordForm>({
  id: '',
  name: '',
  pwd: '',
  status: 1,
  roleId: '',
  remark: '',
})

const searchForm = reactive<SearchForm>({
  name: '',
  status: null,
})

const form = ref<CMSUserForm>({
  name: '',
  pwd: '',
  status: 1,
  roleId: '',
  remark: '',
})

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.guildCmsUserList.usernameRequired'), trigger: 'blur'},
    {min: 2, max: 20, message: t('pages.guildCmsUserList.usernameLength'), trigger: 'blur'},
  ],
  pwd: [
    {required: true, message: t('pages.guildCmsUserList.passwordRequired'), trigger: 'blur'},
    {min: 6, max: 20, message: t('pages.guildCmsUserList.passwordLength'), trigger: 'blur'},
  ],
  roleId: [{required: true, message: t('pages.guildCmsUserList.roleRequired'), trigger: 'change'}],
}))

const resetPasswordRules = computed<FormRules>(() => ({
  pwd: [
    {required: true, message: t('pages.guildCmsUserList.passwordRequired'), trigger: 'blur'},
    {min: 6, max: 20, message: t('pages.guildCmsUserList.passwordLength'), trigger: 'blur'},
  ],
}))

const generateRandomPassword = () => {
  form.value.pwd = createRandomPassword()
  formRef.value?.validateField('pwd')
}

const copyPassword = async () => {
  try {
    const ok = await copyTextToClipboard(form.value.pwd)
    if (!ok) {
      ElMessage.warning(t('pages.guildCmsUserList.nothingToCopy'))
      return
    }
    ElMessage.success(t('pages.guildCmsUserList.passwordCopied'))
  } catch (error) {
    console.error('Copy password failed:', error)
    ElMessage.error(t('pages.guildCmsUserList.copyFailed'))
  }
}

const generateResetPassword = () => {
  resetPasswordForm.value.pwd = createRandomPassword()
  resetPasswordFormRef.value?.validateField('pwd')
}

const copyResetPassword = async () => {
  try {
    const ok = await copyTextToClipboard(resetPasswordForm.value.pwd)
    if (!ok) {
      ElMessage.warning(t('pages.guildCmsUserList.nothingToCopy'))
      return
    }
    ElMessage.success(t('pages.guildCmsUserList.passwordCopied'))
  } catch (error) {
    console.error('Copy password failed:', error)
    ElMessage.error(t('pages.guildCmsUserList.copyFailed'))
  }
}

const resetResetPasswordForm = () => {
  resetPasswordForm.value = {
    id: '',
    name: '',
    pwd: '',
    status: 1,
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
      const {id, name, pwd, status, roleId, remark} = resetPasswordForm.value
      await cmsUserApi.updateCMSUser({
        id,
        name,
        pwd,
        status,
        admin: false,
        roleId,
        remark,
      })
      ElMessage.success(t('pages.guildCmsUserList.resetPasswordSuccess'))
      resetPasswordDialogVisible.value = false
    } catch (error) {
      console.error('Reset password failed:', error)
      ElMessage.error(t('pages.guildCmsUserList.resetPasswordFailed'))
    } finally {
      resetPasswordSubmitting.value = false
    }
  })
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await cmsUserApi.getCMSUserList({
      name: searchForm.name,
      status: searchForm.status ?? undefined,
      roleType: EXTERNAL_ROLE_TYPE,
      nonAdmin: true,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data ?? []
    total.value = response.total ?? 0
  } catch (error) {
    console.error('Failed to load guild CMS users:', error)
    ElMessage.error(t('pages.guildCmsUserList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchRoleList = async () => {
  try {
    const response = await roleApi.getRoleList({
      pageIndex: 1,
      pageSize: 9999,
      roleType: EXTERNAL_ROLE_TYPE,
    })
    const list = response?.data ?? []
    roleOptions.value = list.filter(role => role.status === 1 && role.roleType === EXTERNAL_ROLE_TYPE)
  } catch (error) {
    console.error('Failed to load roles:', error)
    roleOptions.value = []
    ElMessage.error(t('pages.guildCmsUserList.fetchRolesFailed'))
  }
}

const handleAdd = async () => {
  form.value = {
    name: '',
    pwd: createRandomPassword(),
    status: 1,
    roleId: '',
    remark: '',
  }
  await fetchRoleList()
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
      await cmsUserApi.createGuildCMSUser({
        name: form.value.name,
        pwd: form.value.pwd,
        status: form.value.status,
        roleId: form.value.roleId,
        remark: form.value.remark,
      })
      ElMessage.success(t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('Create guild CMS user failed:', error)
      ElMessage.error(t('pages.guildCmsUserList.createFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.status = null
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

onMounted(() => {
  fetchList()
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

.tip-alert {
  margin-bottom: 16px;
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
