<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>VIP配置管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增VIP等级</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="等级名称">
            <el-input v-model="searchForm.levelName" clearable placeholder="等级名称(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.statusFilter" placeholder="全部" style="width: 140px">
              <el-option :value="0" label="全部"/>
              <el-option :value="2" label="只看开启"/>
              <el-option :value="1" label="只看关闭"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="等级" prop="level" width="80"/>
          <el-table-column label="等级名称" prop="levelName" min-width="120"/>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '开启' : '关闭' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="升级充值上限" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.upgradeRechargeLimit) }}
            </template>
          </el-table-column>
          <el-table-column label="最低提现金额" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.minWithdrawAmount) }}
            </template>
          </el-table-column>
          <el-table-column label="最高提现金额" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.maxWithdrawAmount) }}
            </template>
          </el-table-column>
          <el-table-column label="手续费" width="100">
            <template #default="{ row }">
              {{ formatAmount(row.fee) }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="180">
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="130px">
        <el-form-item label="等级" prop="level">
          <el-input-number v-model="currentRow.level" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item label="等级名称" prop="levelName">
          <el-input v-model="currentRow.levelName" placeholder="请输入等级名称"/>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">开启</el-radio>
            <el-radio :label="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="升级充值上限" prop="upgradeRechargeLimit">
          <el-input-number
              v-model="currentRow.upgradeRechargeLimit"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
          />
          <div class="form-tip">保留4位小数，例如 100.0000</div>
        </el-form-item>
        <el-form-item label="最低提现金额" prop="minWithdrawAmount">
          <el-input-number
              v-model="currentRow.minWithdrawAmount"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
          />
          <div class="form-tip">保留4位小数</div>
        </el-form-item>
        <el-form-item label="最高提现金额" prop="maxWithdrawAmount">
          <el-input-number
              v-model="currentRow.maxWithdrawAmount"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
          />
          <div class="form-tip">保留4位小数，0 表示不限制</div>
        </el-form-item>
        <el-form-item label="手续费" prop="fee">
          <el-input-number
              v-model="currentRow.fee"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
          />
          <div class="form-tip">保留4位小数</div>
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {vipCfgApi} from '@/api'
import type {VipCfg} from '@/types/api.ts'

interface SearchForm {
  levelName: string
  statusFilter: number
}

interface VipCfgForm {
  id: string
  level: number
  levelName: string
  status: number
  upgradeRechargeLimit: number
  minWithdrawAmount: number
  maxWithdrawAmount: number
  fee: number
}

const loading = ref(false)
const tableData = ref<VipCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  levelName: '',
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): VipCfgForm => ({
  id: '',
  level: 1,
  levelName: '',
  status: 1,
  upgradeRechargeLimit: 0,
  minWithdrawAmount: 0,
  maxWithdrawAmount: 0,
  fee: 0
})
const currentRow = ref<VipCfgForm>(defaultForm())
const formRef = ref<FormInstance>()

const formatAmount = (val: number | null | undefined) => {
  if (val === null || val === undefined) {
    return '-'
  }
  const n = Number(val)
  if (Number.isNaN(n)) {
    return '-'
  }
  return n.toLocaleString('zh-CN', {minimumFractionDigits: 0, maximumFractionDigits: 4})
}

const validateWithdrawRange = (_rule: unknown, _value: unknown, callback: (error?: Error) => void) => {
  const {minWithdrawAmount, maxWithdrawAmount} = currentRow.value
  if (maxWithdrawAmount > 0 && minWithdrawAmount > maxWithdrawAmount) {
    callback(new Error('最低提现金额不能大于最高提现金额'))
    return
  }
  callback()
}

const formRules: FormRules = {
  level: [{required: true, message: '请输入等级', trigger: 'change'}],
  levelName: [
    {required: true, message: '请输入等级名称', trigger: 'blur'},
    {min: 1, max: 64, message: '等级名称长度在1-64个字符', trigger: 'blur'}
  ],
  status: [{required: true, message: '请选择状态', trigger: 'change'}],
  minWithdrawAmount: [{validator: validateWithdrawRange, trigger: 'change'}],
  maxWithdrawAmount: [{validator: validateWithdrawRange, trigger: 'change'}]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await vipCfgApi.getVipCfgList({
      levelName: searchForm.levelName,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取VIP配置列表失败:', error)
    ElMessage.error('获取VIP配置列表失败')
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
  dialogTitle.value = '新增VIP等级'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: VipCfg) => {
  dialogTitle.value = '编辑VIP等级'
  currentRow.value = {
    id: row.id,
    level: Number(row.level) || 1,
    levelName: row.levelName,
    status: Number(row.status) || 0,
    upgradeRechargeLimit: Number(row.upgradeRechargeLimit) || 0,
    minWithdrawAmount: Number(row.minWithdrawAmount) || 0,
    maxWithdrawAmount: Number(row.maxWithdrawAmount) || 0,
    fee: Number(row.fee) || 0
  }
  dialogVisible.value = true
}

const handleDelete = async (row: VipCfg) => {
  try {
    await ElMessageBox.confirm(`确定要删除 VIP等级 "${row.levelName}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await vipCfgApi.deleteVipCfg(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        level: currentRow.value.level,
        levelName: currentRow.value.levelName,
        status: currentRow.value.status,
        upgradeRechargeLimit: currentRow.value.upgradeRechargeLimit,
        minWithdrawAmount: currentRow.value.minWithdrawAmount,
        maxWithdrawAmount: currentRow.value.maxWithdrawAmount,
        fee: currentRow.value.fee
      }
      if (currentRow.value.id) {
        await vipCfgApi.updateVipCfg({id: currentRow.value.id, ...payload})
      } else {
        await vipCfgApi.createVipCfg(payload)
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    }
  })
}

const resetSearch = () => {
  searchForm.levelName = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
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
  font-size: 16px;
  font-weight: bold;
}

.table-header {
  margin-bottom: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
