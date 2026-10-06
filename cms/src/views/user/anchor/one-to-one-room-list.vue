<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.OneToOneRoomManagement') }}</span>
          <el-button v-if="can('create')" type="primary" @click="openCreate">
            {{ t('pages.oneToOneRoomList.addRoom') }}
          </el-button>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.key"
              clearable
              :placeholder="t('pages.oneToOneRoomList.keywordPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.oneToOneRoomList.oneToOneStatus')">
          <el-select v-model="searchForm.status" clearable style="width: 140px">
            <el-option :label="t('common.onShelf')" :value="1"/>
            <el-option :label="t('common.offShelf')" :value="0"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.oneToOneRoomList.userId')" prop="userId" width="180"/>
        <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.avatar"
                :preview-src-list="[row.avatar]"
                :src="row.avatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.phone')" min-width="130" prop="phone">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRoomRecycleBin.guildId')" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.oneToOneRoomList.oneToOneTitle')" min-width="140" prop="title">
          <template #default="{ row }">{{ row.title || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.oneToOneRoomList.liveRoomStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveRoomStatus === 1 ? 'success' : 'info'">
              {{ row.liveRoomStatus === 1 ? t('common.onShelf') : t('common.offShelf') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.oneToOneRoomList.oneToOneStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.oneToOneRoomList.billing')" prop="billing" width="120"/>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170">
          <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="180">
          <template #default="{ row }">
            <el-button
                v-if="can('edit')"
                link
                type="primary"
                @click="openEdit(row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button
                v-if="row.status !== 1 && can('onShelf')"
                link
                type="success"
                @click="handleSetStatus(row, 1)"
            >
              {{ t('common.onShelf') }}
            </el-button>
            <el-button
                v-if="row.status === 1 && can('offShelf')"
                link
                type="warning"
                @click="handleSetStatus(row, 0)"
            >
              {{ t('common.offShelf') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
            v-model:current-page="pagination.pageIndex"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="formVisible" :title="dialogTitle" width="480px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.oneToOneRoomList.userId')" prop="userId">
          <el-input v-model="form.userId" :disabled="isEdit" :placeholder="t('pages.oneToOneRoomList.userIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.oneToOneRoomList.oneToOneTitle')">
          <el-input v-model="form.title" clearable/>
        </el-form-item>
        <el-form-item :label="t('pages.oneToOneRoomList.tagId')">
          <el-input v-model="form.tagId" clearable placeholder="0"/>
        </el-form-item>
        <el-form-item :label="t('pages.oneToOneRoomList.coverObject')">
          <el-input v-model="form.cover" clearable :placeholder="t('pages.oneToOneRoomList.coverObjectTip')"/>
        </el-form-item>
        <el-form-item :label="t('pages.oneToOneRoomList.billing')" prop="billing">
          <el-input-number v-model="form.billing" :min="0" :precision="4" :step="1" style="width: 100%"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submitForm">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {oneToOneRoomApi} from '@/api'
import type {OneToOneRoomItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const {can} = usePagePermission('OneToOneRoomManagement')

const loading = ref(false)
const saving = ref(false)
const formVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const tableData = ref<OneToOneRoomItem[]>([])
const searchForm = reactive({key: '', status: undefined as number | undefined})
const form = reactive({
  userId: '',
  title: '',
  cover: '',
  tagId: '',
  billing: 0,
})
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const dialogTitle = computed(() => isEdit.value ? t('pages.oneToOneRoomList.editRoom') : t('pages.oneToOneRoomList.addRoom'))

const formRules: FormRules = {
  userId: [{required: true, message: t('pages.oneToOneRoomList.userIdRequired'), trigger: 'blur'}],
  billing: [
    {required: true, message: t('pages.oneToOneRoomList.billingRequired'), trigger: 'change'},
    {type: 'number', min: 0, message: t('pages.oneToOneRoomList.billingMin'), trigger: 'change'},
  ],
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await oneToOneRoomApi.list({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
      status: searchForm.status,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch {
    ElMessage.error(t('pages.oneToOneRoomList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.status = undefined
  handleSearch()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const openCreate = () => {
  isEdit.value = false
  form.userId = ''
  form.title = ''
  form.cover = ''
  form.tagId = ''
  form.billing = 0
  formVisible.value = true
}

const openEdit = (row: OneToOneRoomItem) => {
  isEdit.value = true
  form.userId = row.userId
  form.title = row.title || ''
  form.cover = ''
  form.tagId = row.tagId || ''
  form.billing = Number(row.billing || 0)
  formVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate()
  saving.value = true
  try {
    const payload = {
      userId: form.userId.trim(),
      billing: form.billing,
      title: form.title.trim(),
      cover: form.cover.trim(),
      tagId: form.tagId.trim() || '0',
    }
    if (isEdit.value) {
      await oneToOneRoomApi.update(payload)
    } else {
      await oneToOneRoomApi.create(payload)
    }
    ElMessage.success(t('pages.oneToOneRoomList.saveSuccess'))
    formVisible.value = false
    fetchList()
  } catch {
    ElMessage.error(t('pages.oneToOneRoomList.saveFailed'))
  } finally {
    saving.value = false
  }
}

const handleSetStatus = async (row: OneToOneRoomItem, status: number) => {
  const confirmKey = status === 1 ? 'onShelfConfirm' : 'offShelfConfirm'
  const successKey = status === 1 ? 'onShelfSuccess' : 'offShelfSuccess'
  try {
    await ElMessageBox.confirm(t(`pages.oneToOneRoomList.${confirmKey}`), t('common.confirm'), {type: 'warning'})
    await oneToOneRoomApi.setStatus(row.userId, status)
    ElMessage.success(t(`pages.oneToOneRoomList.${successKey}`))
    fetchList()
  } catch (e) {
    if (e === 'cancel' || e === 'close') {
      return
    }
    ElMessage.error(t('pages.oneToOneRoomList.statusFailed'))
  }
}

onMounted(fetchList)
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-form {
  margin-bottom: 12px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
