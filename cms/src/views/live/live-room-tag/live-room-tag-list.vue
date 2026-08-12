<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRoomTagManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.liveRoomTagList.addTag') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.liveRoomTagList.tagName')" prop="name" min-width="160"/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.liveRoomTagList.tagName')" prop="name">
          <el-input v-model="currentRow.name" maxlength="64" :placeholder="t('pages.liveRoomTagList.tagNamePlaceholder')" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.liveRoomTagList.sortHigherFirst') }}</div>
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
import {liveRoomTagApi} from '@/api/modules/liveRoomTag'
import type {LiveRoomTag} from '@/types/api'

interface TagForm {
  id: string
  name: string
  sort: number
}

const {t} = useI18n()

const loading = ref(false)
const tableData = ref<LiveRoomTag[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): TagForm => ({
  id: '',
  name: '',
  sort: 0
})
const currentRow = ref<TagForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.liveRoomTagList.nameRequired'), trigger: 'blur'},
    {max: 64, message: t('pages.liveRoomTagList.nameMaxLength'), trigger: 'blur'},
  ]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRoomTagApi.getLiveRoomTagList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch live room tag list failed:', error)
    ElMessage.error(t('pages.liveRoomTagList.fetchFailed'))
  } finally {
    loading.value = false
  }
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
  dialogTitle.value = t('pages.liveRoomTagList.addTag')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: LiveRoomTag) => {
  dialogTitle.value = t('pages.liveRoomTagList.editTag')
  currentRow.value = {
    id: row.id,
    name: row.name,
    sort: Number(row.sort) || 0
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    if (currentRow.value.id) {
      await liveRoomTagApi.updateLiveRoomTag({
        id: currentRow.value.id,
        name: currentRow.value.name,
        sort: currentRow.value.sort
      })
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await liveRoomTagApi.createLiveRoomTag({
        name: currentRow.value.name,
        sort: currentRow.value.sort
      })
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save live room tag failed:', error)
  }
}

const handleDelete = async (row: LiveRoomTag) => {
  try {
    await ElMessageBox.confirm(t('pages.liveRoomTagList.deleteConfirm', {name: row.name}), t('common.confirmDelete'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await liveRoomTagApi.deleteLiveRoomTag(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete live room tag failed:', error)
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

.table-header {
  margin-bottom: 16px;
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
}
</style>
