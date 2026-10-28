<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t(menuTitleKey) }}</span>
          <el-button :loading="loading" @click="fetchList">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t(pageDescKey) }}</p>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('pages.userAppSubmissionRecord.userId')">
          <el-input
              v-model="searchForm.userId"
              clearable
              :placeholder="t('pages.userAppSubmissionRecord.userIdPlaceholder')"
              style="width: 180px"
              @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('common.startTime')">
          <el-date-picker
              v-model="searchForm.startTime"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('common.startTime')"
              style="width: 200px"
              teleported
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item :label="t('common.endTime')">
          <el-date-picker
              v-model="searchForm.endTime"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('common.endTime')"
              style="width: 200px"
              teleported
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
        <el-table-column :label="t('pages.userAppSubmissionRecord.userId')" prop="userId" width="120"/>
        <el-table-column :label="t('pages.userAppSubmissionRecord.content')" min-width="280">
          <template #default="{ row }">
            <el-button class="content-link" link type="primary" @click="openPreview(row.content)">
              {{ summarizeContent(row.content) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.userAppSubmissionRecord.createdAt')" prop="createdAt" width="170"/>
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

    <el-dialog
        v-model="previewVisible"
        :title="t('pages.userAppSubmissionRecord.contentPreview')"
        width="720px"
        destroy-on-close
    >
      <pre class="content-preview">{{ previewText }}</pre>
      <template #footer>
        <el-button @click="previewVisible = false">{{ t('tabs.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import type {CMSUserAppSubmissionRecordItem, PageResponse} from '@/types/api'

const props = defineProps<{
  menuTitleKey: string
  pageDescKey: string
  fetchListApi: (params: Record<string, unknown>) => Promise<PageResponse<CMSUserAppSubmissionRecordItem>>
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<CMSUserAppSubmissionRecordItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const previewVisible = ref(false)
const previewText = ref('')

const searchForm = reactive({
  userId: '',
  startTime: '',
  endTime: '',
})

const toServerDateTimeUnix = (value: string) => {
  if (!value) {
    return 0
  }
  return Math.floor(Date.parse(`${value.replace(' ', 'T')}.000Z`) / 1000)
}

const buildQueryParams = () => {
  const userIdText = searchForm.userId.trim()
  const userId = userIdText ? Number(userIdText) : 0
  return {
    userId: userIdText && Number.isFinite(userId) ? userId : undefined,
    startTime: toServerDateTimeUnix(searchForm.startTime),
    endTime: toServerDateTimeUnix(searchForm.endTime),
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await props.fetchListApi(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load user app submission records:', error)
    ElMessage.error(t('pages.userAppSubmissionRecord.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  void fetchList()
}

const handleReset = () => {
  searchForm.userId = ''
  searchForm.startTime = ''
  searchForm.endTime = ''
  currentPage.value = 1
  void fetchList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  void fetchList()
}

const handleCurrentChange = () => {
  void fetchList()
}

const formatRowIndex = (index: number) => (currentPage.value - 1) * pageSize.value + index + 1

const summarizeContent = (content: string) => {
  const text = (content || '').trim()
  if (!text) {
    return t('pages.userAppSubmissionRecord.contentEmpty')
  }
  if (text.length <= 120) {
    return text
  }
  return `${text.slice(0, 120)}…`
}

const openPreview = (content: string) => {
  previewText.value = content || ''
  previewVisible.value = true
}

onMounted(() => {
  void fetchList()
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

.search-form {
  margin-bottom: 12px;
}

.page-desc {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.content-link {
  max-width: 100%;
  white-space: normal;
  text-align: left;
  line-height: 1.4;
  height: auto;
}

.content-preview {
  margin: 0;
  max-height: 60vh;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 13px;
  line-height: 1.5;
}
</style>
