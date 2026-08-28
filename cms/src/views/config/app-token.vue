<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AppTokenConfig') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item :label="t('common.userId')">
          <el-input
              v-model="searchForm.userId"
              clearable
              :placeholder="t('pages.appToken.userIdPlaceholder')"
              style="width: 220px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <div class="table-header">
        <el-button type="primary" @click="handleAdd">{{ t('pages.appToken.addToken') }}</el-button>
      </div>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('common.userId')" prop="id" width="220"/>
        <el-table-column :label="t('pages.appToken.token')" min-width="320" prop="token" show-overflow-tooltip/>
        <el-table-column :label="t('pages.appToken.expireAt')" prop="expireAt" width="200">
          <template #default="scope">
            {{ formatDate(scope.row.expireAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')" prop="expired" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.expired" type="danger">{{ t('pages.appToken.expired') }}</el-tag>
            <el-tag v-else type="success">{{ t('pages.appToken.valid') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">{{ t('common.edit') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" destroy-on-close width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('common.userId')" prop="id">
          <el-input
              v-model="currentRow.id"
              :disabled="isEdit"
              :placeholder="t('pages.appToken.userIdPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.appToken.token')" prop="token">
          <el-input
              v-model="currentRow.token"
              :placeholder="t('pages.appToken.tokenPlaceholder')"
              type="textarea"
          />
        </el-form-item>
        <el-form-item :label="t('pages.appToken.expireAt')" prop="expireAt">
          <el-date-picker
              v-model="currentRow.expireAt"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.appToken.selectExpireAt')"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
          <div v-if="!isEdit" class="form-tip">{{ t('pages.appToken.defaultExpireTip', {time: currentRow.expireAt || '-'}) }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import appTokenApi from '@/api/modules/appToken'
import type {AppToken, SaveAppTokenReq} from '@/types/api'
import {formatServerDateTime as formatDate, formatServerNowPlusDays} from '@/utils/server-datetime'

const {t} = useI18n()
const tableData = ref<AppToken[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  userId: '',
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentRow = ref<SaveAppTokenReq>({
  id: '',
  token: '',
  expireAt: '',
})
const formRef = ref()

const formRules = computed(() => ({
  id: [{required: true, message: t('pages.appToken.userIdRequired'), trigger: 'blur'}],
}))

const defaultExpireAt = () => formatServerNowPlusDays(30)

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
    console.error('fetch app token list failed:', error)
    ElMessage.error(t('pages.appToken.fetchListFailed'))
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

const handleAdd = async () => {
  dialogTitle.value = t('pages.appToken.addTokenTitle')
  isEdit.value = false
  currentRow.value = {
    id: searchForm.userId || '',
    token: '',
    expireAt: '',
  }
  dialogVisible.value = true
  await nextTick()
  currentRow.value.expireAt = defaultExpireAt()
}

const handleEdit = (row: AppToken) => {
  dialogTitle.value = t('pages.appToken.editTokenTitle')
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
      ElMessage.success(t('pages.appToken.saveSuccess'))
      dialogVisible.value = false
      await fetchList()
    } else {
      ElMessage.error(t('pages.appToken.saveFailed'))
    }
  } catch (error) {
    console.error('save app token failed:', error)
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

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
