<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>全局配置</span>
        </div>
      </template>
      <div class="content">
        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="模块编码">
            <el-select
                v-model="searchForm.module"
                allow-create
                default-first-option
                filterable
                placeholder="请选择模块编码"
                style="width: 200px"
            >
              <el-option label="全部" :value="ALL_MODULE"/>
              <el-option
                  v-for="item in moduleOptionList"
                  :key="item.module"
                  :label="formatModuleLabel(item)"
                  :value="item.module"
              >
                <span>{{ item.module }}</span>
                <span v-if="item.moduleName" class="option-desc">{{ item.moduleName }}</span>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="模块名称">
            <el-input
                :model-value="displaySearchModuleName"
                disabled
                placeholder="选择模块编码后自动显示"
                style="width: 200px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>

        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增配置</el-button>
        </div>
        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="序号" type="index" width="80"/>
          <el-table-column label="ID" prop="id" width="200"/>
          <el-table-column label="模块编码" prop="module" width="120"/>
          <el-table-column label="模块名称" prop="moduleName" width="150"/>
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
        <el-form-item label="模块编码" prop="module">
          <el-select
              v-model="currentRow.module"
              allow-create
              default-first-option
              filterable
              placeholder="请选择或输入模块编码"
              style="width: 100%"
              @change="handleFormModuleChange"
          >
            <el-option
                v-for="item in moduleOptionList"
                :key="item.module"
                :label="formatModuleLabel(item)"
                :value="item.module"
            >
              <span>{{ item.module }}</span>
              <span v-if="item.moduleName" class="option-desc">{{ item.moduleName }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模块名称" prop="moduleName">
          <el-input
              v-model="currentRow.moduleName"
              :disabled="isExistingModuleWithName(currentRow.module)"
              placeholder="新模块请填写名称，已有模块自动带出"
          />
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
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import globalCfgApi from '@/api/modules/globalCfg'
import type {GetGlobalCfgReq, GlobalCfg} from '@/types/api'

interface ModuleOption {
  module: string
  moduleName: string
}

const ALL_MODULE = ''

// 表格数据
const tableData = ref<GlobalCfg[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const moduleOptionList = ref<ModuleOption[]>([])

const searchForm = reactive<GetGlobalCfgReq>({
  module: ALL_MODULE,
})

const getModuleNameByCode = (module: string) => {
  return moduleOptionList.value.find(item => item.module === module)?.moduleName || ''
}

const displaySearchModuleName = computed(() => {
  if (!searchForm.module) {
    return '全部'
  }
  return getModuleNameByCode(searchForm.module) || '-'
})

// 弹窗相关
const dialogVisible = ref(false)
const dialogTitle = ref('编辑配置')
const currentRow = ref<GlobalCfg>({
  id: '0',
  module: '',
  moduleName: '',
  key: '',
  value: '',
  desc: ''
})
const formRef = ref()

// 表单验证规则
const formRules = reactive({
  module: [
    {required: true, message: '请输入模块编码', trigger: 'blur'}
  ],
  moduleName: [
    {required: true, message: '请输入模块名称', trigger: 'blur'}
  ],
  key: [
    {required: true, message: '请输入配置键', trigger: 'blur'}
  ],
  value: [
    {required: true, message: '请输入配置值', trigger: 'blur'}
  ]
})

// 从已有配置中提取模块编码与名称映射（module 字段已在库中）
const fetchModuleOptions = async () => {
  try {
    const response = await globalCfgApi.getGlobalCfg()
    const moduleMap = new Map<string, string>()
    ;(response.data || []).forEach(item => {
      if (!item.module) {
        return
      }
      const currentName = moduleMap.get(item.module) || ''
      if (!moduleMap.has(item.module) || (!currentName && item.moduleName)) {
        moduleMap.set(item.module, item.moduleName || currentName)
      }
    })
    moduleOptionList.value = Array.from(moduleMap.entries())
        .map(([module, moduleName]) => ({module, moduleName}))
        .sort((a, b) => a.module.localeCompare(b.module))
  } catch (error) {
    console.error('获取模块列表失败:', error)
  }
}

const formatModuleLabel = (item: ModuleOption) => {
  return item.moduleName ? `${item.module}（${item.moduleName}）` : item.module
}

const isExistingModuleWithName = (module: string) => {
  const item = moduleOptionList.value.find(option => option.module === module)
  return !!(item && item.moduleName)
}

const handleFormModuleChange = (module: string) => {
  const moduleName = getModuleNameByCode(module)
  if (moduleName) {
    currentRow.value.moduleName = moduleName
    return
  }
  if (!module) {
    currentRow.value.moduleName = ''
  }
}

// 获取全局配置列表
const fetchGlobalCfgList = async () => {
  loading.value = true
  try {
    const response = await globalCfgApi.getGlobalCfg({
      module: searchForm.module || undefined,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取全局配置列表失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchGlobalCfgList()
}

const handleReset = () => {
  searchForm.module = ALL_MODULE
  currentPage.value = 1
  fetchGlobalCfgList()
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
    id: '0',
    module: '',
    moduleName: '',
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
      await fetchModuleOptions()
      await fetchGlobalCfgList()
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
      await fetchModuleOptions()
      await fetchGlobalCfgList()
    } else {
      ElMessage.error('保存失败')
    }
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  }
}

// 初始化数据
onMounted(async () => {
  await fetchModuleOptions()
  await fetchGlobalCfgList()
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

.search-form {
  margin-bottom: 15px;
}

.option-desc {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}
</style>