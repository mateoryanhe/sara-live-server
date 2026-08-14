<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.guildMembers.back') }}</el-button>
        </div>
      </template>

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
        <el-table-column :label="t('pages.anchorList.guildId')" prop="guildId" width="190">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.anchorType')" width="110">
          <template #default="{ row }">
            <el-tag :type="anchorTypeTagType(row.userType)">{{ anchorTypeLabel(row.userType) }}</el-tag>
          </template>
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
        <el-table-column fixed="right" :label="t('common.actions')" :width="readonly ? 90 : 220">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">
              {{ t('common.detail') }}
            </el-button>
            <template v-if="!readonly">
              <el-button
                  v-if="row.ban ? can('unban') : can('ban')"
                  :type="row.ban ? 'warning' : 'danger'"
                  link
                  @click="toggleBanStatus(row)"
              >
                {{ row.ban ? t('pages.anchorList.unban') : t('pages.anchorList.ban') }}
              </el-button>
              <el-button
                  v-if="can('exitGuild')"
                  link
                  type="danger"
                  @click="handleExitGuild(row)"
              >
                {{ t('pages.anchorList.exitGuild') }}
              </el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="pagination.pageIndex"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>

      <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildMembers.noMembers')"/>
    </el-card>

    <el-dialog
        v-if="!readonly"
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
import {computed, onActivated, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElForm, ElMessage, ElMessageBox, type FormRules} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const readonly = computed(() => route.name === 'GuildProfileMembers')
const {can} = usePagePermission(
  route.name === 'GuildProfileMembers' ? 'GuildProfileManagement' : 'GuildManagement',
)

const loading = ref(false)
const tableData = ref<AnchorListItem[]>([])
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<InstanceType<typeof ElForm>>()
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

const guildId = computed(() => {
  const value = route.query.guildId
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
})

const guildName = computed(() => {
  const value = route.query.guildName
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
})

const pageTitle = computed(() => {
  if (guildName.value) {
    return t('pages.guildMembers.titleWithName', {name: guildName.value})
  }
  return t('pages.guildMembers.title')
})

const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_PRIVATE = 3
const LIVE_ROOM_PRIVATE_INVITE_ALL = 1
const LIVE_ROOM_PRIVATE_INVITE_VIP = 2
const LIVE_ROOM_PRIVATE_INVITE_REJECT = 3
const USER_TYPE_ANCHOR = 1
const USER_TYPE_SENIOR_ANCHOR = 7

const isPrivateRoom = (category?: number) => category === LIVE_ROOM_CATEGORY_PRIVATE

const anchorTypeLabel = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return t('pages.anchorList.anchorTypeSenior')
  if (userType === USER_TYPE_ANCHOR) return t('pages.anchorList.anchorTypeNormal')
  return '-'
}

const anchorTypeTagType = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return 'warning'
  if (userType === USER_TYPE_ANCHOR) return 'success'
  return 'info'
}

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

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const defaultBanApplyTime = () => {
  const date = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const disabledDate = (time: Date) => time.getTime() < Date.now()

const fetchList = async () => {
  if (!guildId.value) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      guildId: guildId.value,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load guild members:', error)
    ElMessage.error(t('pages.guildMembers.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
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
          },
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

const handleExitGuild = async (row: AnchorListItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.anchorList.exitGuildConfirm', {id: row.id}),
        t('pages.anchorList.exitGuildTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await accountApi.exitGuild({anchorId: row.id})
    ElMessage.success(t('pages.anchorList.exitGuildSuccess'))
    fetchList()
  } catch {
    // cancelled
  }
}

const goBack = () => {
  router.push({name: readonly.value ? 'GuildProfileManagement' : 'GuildManagement'})
}

const openDetail = (row: AnchorListItem) => {
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(row.id)},
  })
}

// keep-alive 按 path 缓存，query 变化不会重新挂载；激活或 guildId 变化时都要拉列表
watch(guildId, (_id, prev) => {
  if (prev !== undefined) {
    pagination.pageIndex = 1
  }
  fetchList()
})

onActivated(() => {
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
  font-size: 16px;
  font-weight: bold;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
