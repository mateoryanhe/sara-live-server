<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CoinMerchantManagement') }}</span>
        </div>
      </template>

      <div class="content">
        <div v-if="can('create')" class="table-header">
          <el-button type="primary" @click="openCreate">{{ t('pages.coinMerchantList.add') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('common.keyword')">
            <el-input
                v-model="searchForm.key"
                clearable
                :placeholder="t('pages.coinMerchantList.keywordPlaceholder')"
                style="width: 220px"
            />
          </el-form-item>
          <el-form-item :label="t('pages.coinMerchantList.cancelStatus')">
            <el-select v-model="searchForm.cancel" clearable :placeholder="t('common.all')" style="width: 120px">
              <el-option :value="0" :label="t('common.normal')"/>
              <el-option :value="1" :label="t('pages.coinMerchantList.canceled')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
            <el-button v-if="can('search')" @click="handleReset">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" min-width="200" prop="id"/>
          <el-table-column :label="t('pages.coinMerchantList.username')" min-width="160" prop="username"/>
          <el-table-column :label="t('pages.coinMerchantList.gold')" align="right" min-width="120" prop="gold">
            <template #default="{row}">
              <span class="currency-gold">{{ formatGold(row.gold) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.coinMerchantList.cancelStatus')" width="120">
            <template #default="{row}">
              <el-tag v-if="row.cancel" type="warning">{{ t('pages.coinMerchantList.canceled') }}</el-tag>
              <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" min-width="180" prop="createdAt">
            <template #default="{row}">{{ formatDate(row.createdAt) }}</template>
          </el-table-column>
          <el-table-column fixed="right" :label="t('common.actions')" width="200">
            <template #default="{row}">
              <el-button
                  v-if="can('resetPassword') && !row.cancel"
                  size="small"
                  type="warning"
                  @click="openResetPassword(row)"
              >
                {{ t('pages.coinMerchantList.resetPassword') }}
              </el-button>
              <el-button
                  v-if="can('cancel') && !row.cancel"
                  size="small"
                  type="danger"
                  @click="handleCancel(row)"
              >
                {{ t('pages.coinMerchantList.cancelAccount') }}
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

    <el-dialog v-model="createVisible" :title="t('pages.coinMerchantList.addDialogTitle')" width="560px" @closed="resetCreateForm">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item :label="t('pages.coinMerchantList.username')" prop="username">
          <el-input v-model="createForm.username" :placeholder="t('pages.coinMerchantList.enterUsername')"/>
        </el-form-item>
        <el-form-item :label="t('pages.coinMerchantList.password')" prop="password">
          <div class="pwd-field">
            <el-input
                v-model="createForm.password"
                :placeholder="t('pages.coinMerchantList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateCreatePassword">{{ t('pages.coinMerchantList.randomGenerate') }}</el-button>
            <el-button :disabled="!createForm.password" @click="copyCreatePassword">
              {{ t('pages.coinMerchantList.copyPassword') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="createSubmitting" type="primary" @click="handleCreateSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="resetVisible"
        :close-on-click-modal="false"
        destroy-on-close
        :title="t('pages.coinMerchantList.resetPasswordTitle')"
        width="560px"
        @closed="resetResetForm"
    >
      <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-width="100px">
        <el-form-item :label="t('pages.coinMerchantList.username')">
          <el-input v-model="resetForm.username" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.coinMerchantList.password')" prop="password">
          <div class="pwd-field">
            <el-input
                v-model="resetForm.password"
                :placeholder="t('pages.coinMerchantList.enterPasswordOrGenerate')"
                show-password
                type="password"
            />
            <el-button @click="generateResetPassword">{{ t('pages.coinMerchantList.randomGenerate') }}</el-button>
            <el-button :disabled="!resetForm.password" @click="copyResetPassword">
              {{ t('pages.coinMerchantList.copyPassword') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="resetSubmitting" type="primary" @click="handleResetSave">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {coinMerchantApi, type CoinMerchantItem} from '@/api/modules/coin-merchant'
import {usePagePermission} from '@/composables/usePagePermission'
import {copyTextToClipboard, createRandomPassword} from '@/utils/random-password'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime} from '@/utils/server-datetime'

const {t} = useI18n()
const {can} = usePagePermission('CoinMerchantManagement')

const loading = ref(false)
const tableData = ref<CoinMerchantItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<{key: string; cancel: number | null}>({
  key: '',
  cancel: null,
})

const createVisible = ref(false)
const createSubmitting = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = ref({username: '', password: ''})

const resetVisible = ref(false)
const resetSubmitting = ref(false)
const resetFormRef = ref<FormInstance>()
const resetForm = ref({accountId: '', username: '', password: ''})

const createRules = computed<FormRules>(() => ({
  username: [
    {required: true, message: t('pages.coinMerchantList.usernameRequired'), trigger: 'blur'},
    {min: 2, max: 32, message: t('pages.coinMerchantList.usernameLength'), trigger: 'blur'},
  ],
  password: [
    {required: true, message: t('pages.coinMerchantList.passwordRequired'), trigger: 'blur'},
    {min: 6, max: 32, message: t('pages.coinMerchantList.passwordLength'), trigger: 'blur'},
  ],
}))

const resetRules = computed<FormRules>(() => ({
  password: [
    {required: true, message: t('pages.coinMerchantList.passwordRequired'), trigger: 'blur'},
    {min: 6, max: 32, message: t('pages.coinMerchantList.passwordLength'), trigger: 'blur'},
  ],
}))

const formatGold = (v: number) => formatWalletBalance(v)
const formatDate = (v?: string) => (v ? formatServerDateTime(v) : '-')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await coinMerchantApi.list({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      key: searchForm.key || undefined,
      cancel: searchForm.cancel,
    })
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.cancel = null
  handleSearch()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

const handleCurrentChange = () => {
  fetchList()
}

const openCreate = () => {
  createForm.value = {username: '', password: createRandomPassword()}
  createVisible.value = true
}

const resetCreateForm = () => {
  createForm.value = {username: '', password: ''}
  createFormRef.value?.clearValidate()
}

const generateCreatePassword = () => {
  createForm.value.password = createRandomPassword()
}

const copyCreatePassword = async () => {
  if (await copyTextToClipboard(createForm.value.password)) {
    ElMessage.success(t('pages.coinMerchantList.copySuccess'))
  }
}

const handleCreateSave = async () => {
  const valid = await createFormRef.value?.validate().catch(() => false)
  if (!valid) return
  createSubmitting.value = true
  try {
    await coinMerchantApi.create({
      username: createForm.value.username.trim(),
      password: createForm.value.password,
    })
    ElMessage.success(t('pages.coinMerchantList.createSuccess'))
    createVisible.value = false
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantList.createFailed'))
  } finally {
    createSubmitting.value = false
  }
}

const openResetPassword = (row: CoinMerchantItem) => {
  resetForm.value = {
    accountId: row.id,
    username: row.username,
    password: createRandomPassword(),
  }
  resetVisible.value = true
}

const resetResetForm = () => {
  resetForm.value = {accountId: '', username: '', password: ''}
  resetFormRef.value?.clearValidate()
}

const generateResetPassword = () => {
  resetForm.value.password = createRandomPassword()
}

const copyResetPassword = async () => {
  if (await copyTextToClipboard(resetForm.value.password)) {
    ElMessage.success(t('pages.coinMerchantList.copySuccess'))
  }
}

const handleResetSave = async () => {
  const valid = await resetFormRef.value?.validate().catch(() => false)
  if (!valid) return
  resetSubmitting.value = true
  try {
    await coinMerchantApi.resetPassword({
      accountId: resetForm.value.accountId,
      password: resetForm.value.password,
    })
    ElMessage.success(t('pages.coinMerchantList.resetPasswordSuccess'))
    resetVisible.value = false
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantList.resetPasswordFailed'))
  } finally {
    resetSubmitting.value = false
  }
}

const handleCancel = async (row: CoinMerchantItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.coinMerchantList.cancelConfirm', {username: row.username}),
        t('pages.coinMerchantList.cancelTitle'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    await coinMerchantApi.cancel({accountId: row.id})
    ElMessage.success(t('pages.coinMerchantList.cancelSuccess'))
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.coinMerchantList.cancelFailed'))
  }
}

onMounted(fetchList)
</script>

<style scoped>
.table-header {
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 12px;
}

.pwd-field {
  display: flex;
  gap: 8px;
  width: 100%;
}

.pwd-field .el-input {
  flex: 1;
}

.currency-gold {
  color: #b8860b;
  font-variant-numeric: tabular-nums;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
