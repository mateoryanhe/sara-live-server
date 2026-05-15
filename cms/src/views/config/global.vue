<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>全局配置</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增配置</el-button>
        </div>
        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="序号" type="index" width="80"/>
          <el-table-column label="ID" prop="id" width="200"/>
          <el-table-column label="模块" prop="module" width="150"/>
          <el-table-column label="配置键" prop="key" width="200"/>
          <el-table-column label="配置值" prop="value" show-overflow-tooltip/>
          <el-table-column label="描述" prop="desc" show-overflow-tooltip width="200"/>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
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

    <!-- 编辑/新增弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="模块" prop="module">
          <el-input v-model="currentRow.module" placeholder="请输入模块名称"/>
        </el-form-item>
        <el-form-item label="配置键" prop="key">
          <el-input v-model="currentRow.key" placeholder="请输入配置键"/>
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="currentRow.value" placeholder="请输入配置值" type="textarea"/>
        </el-form-item>
        <el-form-item label="描述" prop="desc">
          <el-input v-model="currentRow.desc" placeholder="请输入描述" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import globalCfgApi from '@/api/modules/globalCfg'
import type {GlobalCfg} from '@/types/api'

// 表格数据
const tableData = ref<GlobalCfg[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 弹窗相关
const dialogVisible = ref(false)
const dialogTitle = ref('编辑配置')
const currentRow = ref<GlobalCfg>({
  id: '0',
  module: '',
  key: '',
  value: '',
  desc: ''
})
const formRef = ref()

// 表单验证规则
const formRules = reactive({
  module: [
    {required: true, message: '请输入模块名称', trigger: 'blur'}
  ],
  key: [
    {required: true, message: '请输入配置键', trigger: 'blur'}
  ],
  value: [
    {required: true, message: '请输入配置值', trigger: 'blur'}
  ]
})

// 获取全局配置列表
const fetchGlobalCfgList = async () => {
  loading.value = true
  try {
    // 使用POST请求获取数据
    const response = await globalCfgApi.getGlobalCfg()
    // 直接处理PageResponse<GlobalCfg>格式的响应
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取全局配置列表失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGlobalCfgList()
}

// 当前页变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGlobalCfgList()
}

// 新增配置
const handleAdd = () => {
  dialogTitle.value = '新增配置'
  currentRow.value = {
    id: '0',  // 设置默认ID为'0'字符串
    module: '',
    key: '',
    value: '',
    desc: ''
  }
  dialogVisible.value = true
}

// 编辑操作
const handleEdit = (row: GlobalCfg) => {
  dialogTitle.value = '编辑配置'
  currentRow.value = {...row}
  dialogVisible.value = true
}

// 删除操作
const handleDelete = async (row: GlobalCfg) => {
  try {
    await ElMessageBox.confirm(`确定要删除配置 "${row.key}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await globalCfgApi.delGlobalCfg(row)
    if (response) {
      ElMessage.success('删除成功')
      await fetchGlobalCfgList() // 重新获取数据
    } else {
      ElMessage.error('删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
  }
}

// 保存操作
const handleSave = async () => {
  try {
    await formRef.value.validate()

    let response
    if (currentRow.value.id && currentRow.value.id !== '0') {
      // 更新现有配置
      response = await globalCfgApi.saveGlobalCfg(currentRow.value)
    } else {
      // 新增配置
      response = await globalCfgApi.saveGlobalCfg(currentRow.value)
    }

    if (response) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      await fetchGlobalCfgList() // 重新获取数据
    } else {
      ElMessage.error('保存失败')
    }
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  }
}

// 初始化数据
onMounted(() => {
  fetchGlobalCfgList()
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

.content {
  min-height: 400px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.table-header {
  margin-bottom: 15px;
}
</style>