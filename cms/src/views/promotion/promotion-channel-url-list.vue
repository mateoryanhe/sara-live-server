<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionChannelUrlListManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <el-alert
            :closable="false"
            show-icon
            type="info"
            class="domain-hint"
            :title="shareBaseHost ? t('pages.promotionChannelUrlList.shareBaseHostCurrent', {host: shareBaseHost}) : t('pages.promotionChannelUrlList.shareBaseHostEmpty')"
        >
          <template #default>
            <p class="hint-text">{{ t('pages.promotionChannelUrlList.shareBaseHostHint') }}</p>
          </template>
        </el-alert>

        <div class="table-header">
          <el-input
              v-model="searchKey"
              clearable
              class="search-input"
              :placeholder="t('pages.promotionChannelUrlList.searchPlaceholder')"
              @keyup.enter="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button type="primary" :disabled="!shareBaseHost" @click="handleAdd">{{ t('pages.promotionChannelUrlList.add') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.promotionChannelUrlList.channelName')" prop="channelName" min-width="140"/>
          <el-table-column :label="t('pages.promotionChannelUrlList.channelCode')" prop="channelCode" width="120"/>
          <el-table-column :label="t('pages.promotionChannelUrlList.shareUrl')" min-width="300">
            <template #default="{ row }">
              <div class="share-url-cell">
                <span class="share-url-text" :title="row.shareUrl || row.url">{{ row.shareUrl || row.url || '-' }}</span>
                <el-button
                    v-if="row.shareUrl || row.url"
                    link
                    type="primary"
                    @click="() => copyShareUrl(row.shareUrl || row.url)"
                >
                  {{ t('pages.promotionChannelUrlList.copyShareUrl') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="120px">
        <el-form-item :label="t('pages.promotionChannelUrlList.channelName')" prop="channelName">
          <el-input v-model="currentRow.channelName" maxlength="128" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('pages.promotionChannelUrlList.channelCode')" prop="channelCode">
          <el-input-number v-model="currentRow.channelCode" :min="1" :max="2147483647" controls-position="right" style="width: 100%"/>
          <div class="form-tip">{{ t('pages.promotionChannelUrlList.channelCodeTip') }}</div>
        </el-form-item>
        <el-form-item v-if="previewShareUrl" :label="t('pages.promotionChannelUrlList.shareUrl')">
          <el-input :model-value="previewShareUrl" readonly/>
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
import {computed, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {promotionChannelUrlApi} from '@/api/modules/promotion-channel-url'
import type {PromotionChannelUrl} from '@/types/api'

interface ChannelForm {
  id: string
  channelName: string
  channelCode: number | undefined
}

const {t} = useI18n()

const loading = ref(false)
const tableData = ref<PromotionChannelUrl[]>([])
const total = ref(0)
const shareBaseHost = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const searchKey = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): ChannelForm => ({
  id: '',
  channelName: '',
  channelCode: undefined,
})
const currentRow = ref<ChannelForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  channelName: [
    {required: true, message: t('pages.promotionChannelUrlList.channelNameRequired'), trigger: 'blur'},
  ],
  channelCode: [
    {required: true, message: t('pages.promotionChannelUrlList.channelCodeRequired'), trigger: 'change'},
    {type: 'number', min: 1, message: t('pages.promotionChannelUrlList.channelCodeMin'), trigger: 'change'},
  ],
}))

const previewShareUrl = computed(() => {
  const code = currentRow.value.channelCode
  const host = shareBaseHost.value.trim()
  if (!host || code == null || code < 1) {
    return ''
  }
  return `https://${host}?ch=${code}`
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await promotionChannelUrlApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      key: searchKey.value.trim(),
    })
    tableData.value = response.data
    total.value = response.total
    shareBaseHost.value = (response.shareBaseHost || '').trim()
  } catch (error) {
    console.error('fetch promotion channel url list failed:', error)
    ElMessage.error(t('pages.promotionChannelUrlList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const copyShareUrl = async (link: string) => {
  if (!link) {
    return
  }
  try {
    await navigator.clipboard.writeText(link)
    ElMessage.success(t('pages.promotionChannelUrlList.copyShareUrlSuccess'))
  } catch (error) {
    console.error('copy share url failed:', error)
    ElMessage.error(t('pages.promotionChannelUrlList.copyShareUrlFailed'))
  }
}

const handleSearch = () => {
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
  dialogTitle.value = t('pages.promotionChannelUrlList.add')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: PromotionChannelUrl) => {
  dialogTitle.value = t('pages.promotionChannelUrlList.edit')
  currentRow.value = {
    id: row.id,
    channelName: row.channelName || '',
    channelCode: row.channelCode,
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const code = currentRow.value.channelCode
    if (code == null || code < 1) {
      return
    }
    const payload = {
      channelName: currentRow.value.channelName.trim(),
      channelCode: code,
    }
    if (currentRow.value.id) {
      await promotionChannelUrlApi.update({
        id: currentRow.value.id,
        ...payload,
      })
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await promotionChannelUrlApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save promotion channel url failed:', error)
  }
}

const handleDelete = async (row: PromotionChannelUrl) => {
  try {
    await ElMessageBox.confirm(
        t('pages.promotionChannelUrlList.deleteConfirm', {name: row.channelName}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await promotionChannelUrlApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete promotion channel url failed:', error)
    }
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
  justify-content: space-between;
}

.domain-hint {
  margin-bottom: 16px;
}

.hint-text {
  margin: 4px 0 0;
  font-size: 13px;
  line-height: 1.5;
}

.table-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
}

.share-url-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.share-url-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
