<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>主播列表</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item label="关键字">
          <el-input v-model="searchForm.key" clearable placeholder="用户ID/昵称/手机号/分享码"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="用户ID" prop="id" width="180"/>
        <el-table-column label="昵称" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="头像" width="80">
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
        <el-table-column label="手机号" min-width="130" prop="phone">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column label="工会ID" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column label="登录IP" min-width="140" prop="ip">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播间标题" min-width="140" prop="roomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? '直播中' : '未开播' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="直播收益" min-width="110">
          <template #default="{ row }">{{ formatAmount(row.totalIncome) }}</template>
        </el-table-column>
        <el-table-column label="礼物收益" min-width="110">
          <template #default="{ row }">{{ formatAmount(row.totalGiftIncome) }}</template>
        </el-table-column>
        <el-table-column label="付费弹幕收益" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPaidDanmakuIncome) }}</template>
        </el-table-column>
        <el-table-column label="私密门票收益" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPrivateRoomTicketIncome) }}</template>
        </el-table-column>
        <el-table-column label="私密观看收益" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPrivateRoomWatchIncome) }}</template>
        </el-table-column>
        <el-table-column label="封禁状态" prop="ban" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.ban" type="danger">已封禁</el-tag>
            <el-tag v-else type="success">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="封禁截止" prop="banApplyTime" width="170">
          <template #default="{ row }">{{ formatDate(row.banApplyTime) }}</template>
        </el-table-column>
        <el-table-column label="封禁原因" min-width="160" prop="banReason" show-overflow-tooltip>
          <template #default="{ row }">{{ row.banReason || '-' }}</template>
        </el-table-column>
        <el-table-column label="注册时间" prop="registeredAt" width="170">
          <template #default="{ row }">{{ formatDate(row.registeredAt) }}</template>
        </el-table-column>
        <el-table-column label="资料更新时间" prop="createdAt" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="{ row }">
            <el-button
                :type="row.ban ? 'warning' : 'danger'"
                link
                @click="toggleBanStatus(row)"
            >
              {{ row.ban ? '解封' : '封禁' }}
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
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <el-dialog
        v-model="banDialogVisible"
        :close-on-click-modal="false"
        destroy-on-close
        title="封禁主播"
        width="520px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banRules" label-width="100px">
        <el-form-item label="主播ID">
          <el-input v-model="banForm.accountId" disabled/>
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item label="封禁截止" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择封禁截止时间"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="封禁原因" prop="banReason">
          <el-input
              v-model="banForm.banReason"
              :maxlength="512"
              :rows="4"
              placeholder="请输入封禁原因"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">取消</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">确认封禁</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElForm, ElMessage, ElMessageBox} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'

const loading = ref(false)
const tableData = ref<AnchorListItem[]>([])
const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<InstanceType<typeof ElForm>>()

const searchForm = reactive({
  key: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const banForm = reactive({
  accountId: '',
  nickname: '',
  banApplyTime: '',
  banReason: '',
})

const banRules = {
  banApplyTime: [
    {required: true, message: '请选择封禁截止时间', trigger: 'change'},
  ],
  banReason: [
    {required: true, message: '请输入封禁原因', trigger: 'blur'},
    {min: 1, max: 512, message: '封禁原因长度需在1到512之间', trigger: 'blur'},
  ],
}

const defaultBanApplyTime = () => {
  const date = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const disabledDate = (time: Date) => time.getTime() < Date.now()

const fetchList = async () => {
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取主播列表失败:', error)
    ElMessage.error('获取主播列表失败')
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
  pagination.pageIndex = 1
  fetchList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString('zh-CN')
  } catch {
    return '-'
  }
}

const formatAmount = (value: number | null | undefined) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(2)
}

const resetBanForm = () => {
  banForm.accountId = ''
  banForm.nickname = ''
  banForm.banApplyTime = ''
  banForm.banReason = ''
  banFormRef.value?.clearValidate()
}

const openBanDialog = (row: AnchorListItem) => {
  banForm.accountId = row.id
  banForm.nickname = row.nickname || '-'
  banForm.banApplyTime = defaultBanApplyTime()
  banForm.banReason = ''
  banDialogVisible.value = true
}

const submitBan = async () => {
  if (!banFormRef.value) {
    return
  }
  await banFormRef.value.validate(async (valid: boolean) => {
    if (!valid) {
      return
    }
    banSubmitting.value = true
    try {
      const banData: BanAnchorReq = {
        accountId: banForm.accountId,
        banApplyTime: banForm.banApplyTime,
        banReason: banForm.banReason.trim(),
      }
      const response = await accountApi.banAnchor(banData)
      if (response) {
        ElMessage.success('封禁成功，已通知App端')
        banDialogVisible.value = false
        fetchList()
      } else {
        ElMessage.error('封禁失败')
      }
    } catch (error) {
      console.error('封禁主播失败:', error)
      ElMessage.error('封禁请求失败')
    } finally {
      banSubmitting.value = false
    }
  })
}

const toggleBanStatus = async (row: AnchorListItem) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
          `确定要解封主播 ${row.id} 吗？`,
          '确认解封',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
      )
      const unBanData: UnBanAnchorReq = {accountId: row.id}
      const response = await accountApi.unBanAnchor(unBanData)
      if (response) {
        ElMessage.success('解封成功')
        fetchList()
      } else {
        ElMessage.error('解封失败')
      }
    } catch {
      // 用户取消
    }
    return
  }

  openBanDialog(row)
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

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
