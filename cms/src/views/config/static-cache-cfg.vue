<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.StaticCacheCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :title="t('pages.staticCacheCfg.noticeTitle')"
          :description="t('pages.staticCacheCfg.noticeBody')"
          type="info"
          show-icon
          :closable="false"
          class="notice"
      />

      <div class="table-header">
        <el-input
            v-model="searchKey"
            clearable
            class="search-input"
            :placeholder="t('pages.staticCacheCfg.searchPlaceholder')"
            @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button type="primary" @click="handleAdd">{{ t('pages.staticCacheCfg.addRule') }}</el-button>
      </div>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="ID" prop="id" width="100"/>
        <el-table-column :label="t('pages.staticCacheCfg.fileName')" prop="fileName" min-width="220"/>
        <el-table-column :label="t('pages.staticCacheCfg.cachePolicy')" min-width="240">
          <template #default>
            <el-tag type="warning">no-cache, no-store, must-revalidate</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.staticCacheCfg.remark')" prop="remark" min-width="180"/>
        <el-table-column :label="t('common.createdAt')" prop="createdAt" width="170"/>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170"/>
        <el-table-column fixed="right" :label="t('common.actions')" width="160">
          <template #default="{row}">
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
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
        <el-form-item :label="t('pages.staticCacheCfg.fileName')" prop="fileName">
          <el-input v-model="currentRow.fileName" maxlength="255" show-word-limit/>
          <div class="form-tip">{{ t('pages.staticCacheCfg.fileNameTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.staticCacheCfg.remark')" prop="remark">
          <el-input v-model="currentRow.remark" maxlength="255" show-word-limit/>
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
import {staticCacheCfgApi} from '@/api/modules/static-cache-cfg'
import type {StaticCacheRule} from '@/types/api'

interface RuleForm {
  id: string
  fileName: string
  remark: string
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<StaticCacheRule[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKey = ref('')
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()

const defaultForm = (): RuleForm => ({id: '', fileName: '', remark: ''})
const currentRow = ref<RuleForm>(defaultForm())
const formRules = computed<FormRules>(() => ({
  fileName: [
    {required: true, message: t('pages.staticCacheCfg.fileNameRequired'), trigger: 'blur'},
    {pattern: /^[^/\\]+$/, message: t('pages.staticCacheCfg.fileNameInvalid'), trigger: 'blur'},
  ],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await staticCacheCfgApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      key: searchKey.value.trim(),
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch static cache rule list failed:', error)
    ElMessage.error(t('pages.staticCacheCfg.fetchFailed'))
  } finally {
    loading.value = false
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
  dialogTitle.value = t('pages.staticCacheCfg.addRule')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: StaticCacheRule) => {
  dialogTitle.value = t('pages.staticCacheCfg.editRule')
  currentRow.value = {id: row.id, fileName: row.fileName || '', remark: row.remark || ''}
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const payload = {
      fileName: currentRow.value.fileName.trim(),
      remark: currentRow.value.remark.trim(),
    }
    if (currentRow.value.id) {
      await staticCacheCfgApi.update({id: currentRow.value.id, ...payload})
      ElMessage.success(t('pages.staticCacheCfg.updateSuccess'))
    } else {
      await staticCacheCfgApi.create(payload)
      ElMessage.success(t('pages.staticCacheCfg.createSuccess'))
    }
    dialogVisible.value = false
    await fetchList()
  } catch (error) {
    if (error !== false) {
      console.error('save static cache rule failed:', error)
    }
  }
}

const handleDelete = async (row: StaticCacheRule) => {
  try {
    await ElMessageBox.confirm(
        t('pages.staticCacheCfg.deleteConfirm', {fileName: row.fileName}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await staticCacheCfgApi.remove(row.id)
    ElMessage.success(t('pages.staticCacheCfg.deleteSuccess'))
    await fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete static cache rule failed:', error)
    }
  }
}

onMounted(fetchList)
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

.notice {
  margin-bottom: 16px;
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
</style>
