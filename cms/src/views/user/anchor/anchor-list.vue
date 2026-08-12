<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AnchorListManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.keyword')">
          <el-input v-model="searchForm.key" clearable :placeholder="t('pages.anchorList.keywordPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('common.userId')" prop="id" width="180"/>
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
        <el-table-column :label="t('pages.anchorList.guildId')" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.loginIp')" min-width="140" prop="ip">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveRoom')" prop="roomId" width="180">
          <template #default="{ row }">{{ row.roomId || row.id || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.roomType')" width="100">
          <template #default="{ row }">
            <el-tag :type="categoryTagType(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.privateInviteType')" min-width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category === LIVE_ROOM_CATEGORY_HOT" :type="privateInviteTagType(row.privateInviteType)">
              {{ privateInviteLabel(row.privateInviteType) }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.ticketPrice')" min-width="100">
          <template #default="{ row }">
            {{ isPrivateRoom(row.category) ? formatAmount(row.ticket) : '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.billingPricePerMinute')" min-width="110">
          <template #default="{ row }">
            {{ isPrivateRoom(row.category) ? formatAmount(row.billing) : '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.roomTitle')" min-width="140" prop="roomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? t('common.live') : t('common.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveIncome')" min-width="110">
          <template #default="{ row }">{{ formatAmount(row.totalIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.giftIncome')" min-width="110">
          <template #default="{ row }">{{ formatAmount(row.totalGiftIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPaidDanmakuIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoTicketIncome')" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalVideoCallTicketIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoBillingIncome')" min-width="140">
          <template #default="{ row }">{{ formatAmount(row.totalVideoCallBillingIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoCallIncome')" min-width="120">
          <template #default="{ row }">{{ formatAmount(row.totalVideoCallIncome) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.banStatus')" prop="ban" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.ban" type="danger">{{ t('common.banned') }}</el-tag>
            <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.banUntil')" prop="banApplyTime" width="170">
          <template #default="{ row }">{{ formatDate(row.banApplyTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.banReason')" min-width="160" prop="banReason" show-overflow-tooltip>
          <template #default="{ row }">{{ row.banReason || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.registeredAt')" prop="registeredAt" width="170">
          <template #default="{ row }">{{ formatDate(row.registeredAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.profileUpdatedAt')" prop="createdAt" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button
                :type="row.ban ? 'warning' : 'danger'"
                link
                @click="toggleBanStatus(row)"
            >
              {{ row.ban ? t('pages.anchorList.unban') : t('pages.anchorList.ban') }}
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
        :title="t('pages.anchorList.banDialogTitle')"
        width="520px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banRules" label-width="100px">
        <el-form-item :label="t('pages.anchorList.anchorId')">
          <el-input v-model="banForm.accountId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.banUntil')" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.anchorList.selectBanUntil')"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.banReason')" prop="banReason">
          <el-input
              v-model="banForm.banReason"
              :maxlength="512"
              :rows="4"
              :placeholder="t('pages.anchorList.enterBanReason')"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">{{ t('pages.anchorList.confirmBan') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElForm, ElMessage, ElMessageBox, type FormRules} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const {t} = useI18n()

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

const banRules = computed<FormRules>(() => ({
  banApplyTime: [
    {required: true, message: t('pages.anchorList.banApplyTimeRequired'), trigger: 'change'},
  ],
  banReason: [
    {required: true, message: t('pages.anchorList.banReasonRequired'), trigger: 'blur'},
    {min: 1, max: 512, message: t('pages.anchorList.banReasonLength'), trigger: 'blur'},
  ],
}))

const defaultBanApplyTime = () => {
  const date = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const disabledDate = (time: Date) => time.getTime() < Date.now()

const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_PRIVATE = 3

const LIVE_ROOM_PRIVATE_INVITE_ALL = 1
const LIVE_ROOM_PRIVATE_INVITE_VIP = 2
const LIVE_ROOM_PRIVATE_INVITE_REJECT = 3

const isPrivateRoom = (category?: number) => category === LIVE_ROOM_CATEGORY_PRIVATE

const privateInviteLabel = (type?: number) => {
  if (type === LIVE_ROOM_PRIVATE_INVITE_VIP) return t('pages.anchorList.privateInviteVipOnly')
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return t('pages.anchorList.privateInviteRejectAll')
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return t('pages.anchorList.privateInviteAcceptAll')
  return '-'
}

const privateInviteTagType = (type?: number) => {
  if (type === LIVE_ROOM_PRIVATE_INVITE_VIP) return 'warning'
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return 'danger'
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return 'success'
  return 'info'
}

const categoryLabel = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_HOT) return t('pages.anchorList.categoryHot')
  if (category === LIVE_ROOM_CATEGORY_GAME) return t('pages.anchorList.categoryGame')
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return t('pages.anchorList.categoryPrivate')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return 'warning'
  if (category === LIVE_ROOM_CATEGORY_GAME) return 'success'
  if (category === LIVE_ROOM_CATEGORY_HOT) return 'danger'
  return 'info'
}

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
    console.error('Failed to load anchor list:', error)
    ElMessage.error(t('pages.anchorList.fetchFailed'))
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
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
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
  if (!banFormRef.value) return
  await banFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    banSubmitting.value = true
    try {
      const banData: BanAnchorReq = {
        accountId: banForm.accountId,
        banApplyTime: banForm.banApplyTime,
        banReason: banForm.banReason.trim(),
      }
      const response = await accountApi.banAnchor(banData)
      if (response) {
        ElMessage.success(t('pages.anchorList.banSuccessNotify'))
        banDialogVisible.value = false
        fetchList()
      } else {
        ElMessage.error(t('pages.anchorList.banFailed'))
      }
    } catch (error) {
      console.error('Ban anchor failed:', error)
      ElMessage.error(t('pages.anchorList.banRequestFailed'))
    } finally {
      banSubmitting.value = false
    }
  })
}

const toggleBanStatus = async (row: AnchorListItem) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
          t('pages.anchorList.unbanConfirm', {id: row.id}),
          t('pages.anchorList.unbanTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning',
          }
      )
      const unBanData: UnBanAnchorReq = {accountId: row.id}
      const response = await accountApi.unBanAnchor(unBanData)
      if (response) {
        ElMessage.success(t('pages.anchorList.unbanSuccess'))
        fetchList()
      } else {
        ElMessage.error(t('pages.anchorList.unbanFailed'))
      }
    } catch {
      // cancelled
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
