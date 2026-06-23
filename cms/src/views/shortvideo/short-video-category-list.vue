<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>短视频分类</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增分类</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="分类名称" prop="name" min-width="160"/>
          <el-table-column label="排序" prop="sort" width="80"/>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="currentRow.name" maxlength="64" placeholder="请输入分类名称" show-word-limit/>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">数值越大越靠前</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {shortVideoApi} from '@/api/modules/shortVideo'
import type {ShortVideoCategory} from '@/types/api'

interface CategoryForm {
  id: string
  name: string
  sort: number
}

const loading = ref(false)
const tableData = ref<ShortVideoCategory[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): CategoryForm => ({
  id: '',
  name: '',
  sort: 0
})
const currentRow = ref<CategoryForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入分类名称', trigger: 'blur'},
    {max: 64, message: '分类名称最长64字符', trigger: 'blur'},
  ]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoCategoryList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取短视频分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
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
  dialogTitle.value = '新增分类'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: ShortVideoCategory) => {
  dialogTitle.value = '编辑分类'
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
      await shortVideoApi.updateShortVideoCategory({
        id: currentRow.value.id,
        name: currentRow.value.name,
        sort: currentRow.value.sort
      })
      ElMessage.success('更新成功')
    } else {
      await shortVideoApi.createShortVideoCategory({
        name: currentRow.value.name,
        sort: currentRow.value.sort
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('保存短视频分类失败:', error)
  }
}

const handleDelete = async (row: ShortVideoCategory) => {
  try {
    await ElMessageBox.confirm(`确定要删除分类「${row.name}」吗？`, '确认删除', {
      type: 'warning'
    })
    await shortVideoApi.deleteShortVideoCategory(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除短视频分类失败:', error)
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
