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

          <el-input
              v-model="searchForm.key"
              clearable
              :placeholder="t('pages.anchorList.keywordPlaceholder')"
              style="width: 200px"
          />

        </el-form-item>

        <el-form-item>

          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>

          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>

        </el-form-item>

      </el-form>

      <div class="table-scroll">
      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
        <el-table-column :label="t('menu.UserDetail')" width="110">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail" link type="primary" @click="openUserDetail(row.id)">
              {{ t('pages.userList.viewDetail') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.userId')" prop="id" width="180">
          <template #default="{ row }">
            <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
              {{ row.id }}
            </el-button>
            <span v-else>{{ row.id }}</span>
          </template>
        </el-table-column>
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
        <el-table-column :label="t('pages.anchorList.guildName')" min-width="140" prop="guildName" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button
                v-if="hasGuild(row) && canViewGuildDetail"
                link
                type="primary"
                @click="openGuildDetail(row.guildId)"
            >
              {{ row.guildName || '-' }}
            </el-button>
            <span v-else>{{ row.guildName || '-' }}</span>
          </template>
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
        <el-table-column :label="t('pages.anchorList.ticketPrice')" align="right" min-width="110">
          <template #default="{ row }">
            <span v-if="isPrivateRoom(row.category)" class="money-amount">{{ formatWalletBalance(row.ticket) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.billingPricePerMinute')" align="right" min-width="120">
          <template #default="{ row }">
            <span v-if="isPrivateRoom(row.category)" class="money-amount">{{ formatWalletBalance(row.billing) }}</span>
            <span v-else>-</span>
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
        <el-table-column :label="t('pages.anchorList.liveIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.giftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoTicketIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoBillingIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.banStatus')" prop="ban" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.ban" type="danger">{{ t('common.banned') }}</el-tag>
            <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.shelfStatus')" prop="status" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">{{ t('common.onShelf') }}</el-tag>
            <el-tag v-else type="info">{{ t('common.offShelf') }}</el-tag>
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
        <el-table-column fixed="right" :label="t('common.actions')" width="280">
          <template #default="{ row }">
            <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
              {{ t('common.detail') }}
            </el-button>
            <el-button
                v-if="can('offShelf')"
                type="warning"
                link
                @click="handleOffShelf(row)"
            >
              {{ t('common.offShelf') }}
            </el-button>
            <el-button
                :type="row.ban ? 'warning' : 'danger'"
                link
                @click="toggleBanStatus(row)"
            >
              {{ row.ban ? t('pages.anchorList.unban') : t('pages.anchorList.ban') }}
            </el-button>
            <el-button
                v-if="Number(row.guildId) !== 0"
                link
                type="danger"
                @click="handleExitGuild(row)"
            >
              {{ t('pages.anchorList.exitGuild') }}
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

import {useRouter} from 'vue-router'

import {ElForm, ElMessage, ElMessageBox, type FormRules} from 'element-plus'

import {accountApi} from '@/api'

import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'

import {formatWalletBalance} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'

const {t} = useI18n()

const router = useRouter()
const {can} = usePagePermission('AnchorListManagement')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('AnchorListManagement')
const canViewDetail = computed(() => can('viewDetail'))
const canViewGuildDetail = computed(() => can('viewGuildDetail'))

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
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return t('pages.anchorList.privateInviteRejectAll')
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return t('pages.anchorList.privateInviteAcceptAll')
  return '-'
}

const privateInviteTagType = (type?: number) => {
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

const formatRowIndex = (index: number) =>
    (pagination.pageIndex - 1) * pagination.pageSize + index + 1

const resetBanForm = () => {

  banForm.accountId = ''

  banForm.nickname = ''

  banForm.banApplyTime = ''

  banForm.banReason = ''

  banFormRef.value?.clearValidate()

}

const openDetail = (row: AnchorListItem) => {
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(row.id)},
  })
}

const hasGuild = (row: AnchorListItem) => {
  const guildId = row.guildId
  if (guildId == null || guildId === '') {
    return false
  }
  return Number(guildId) !== 0
}

const openGuildDetail = (guildId: string | number | undefined) => {
  if (!guildId || Number(guildId) === 0) {
    return
  }
  router.push({
    name: 'GuildDetail',
    query: {id: String(guildId)},
  })
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

const handleOffShelf = async (row: AnchorListItem) => {
  try {
    await ElMessageBox.confirm(
      t('pages.anchorList.offShelfConfirm', {id: row.id}),
      t('common.confirmOffShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await accountApi.setLiveRoomStatus({
      anchorId: row.id,
      status: 0,
    })
    ElMessage.success(t('pages.anchorList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('off shelf live room failed:', error)
    ElMessage.error(t('pages.anchorList.offShelfFailed'))
  }
}

const handleExitGuild = async (row: AnchorListItem) => {
  try {
    await ElMessageBox.confirm(
      t('pages.anchorList.exitGuildConfirm', {id: row.id}),
      t('pages.anchorList.exitGuildTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await accountApi.exitGuild({anchorId: row.id})
    ElMessage.success(t('pages.anchorList.exitGuildSuccess'))
    fetchList()
  } catch {
    // cancelled
  }
}

</script>



<style scoped>

.page-container {

  padding: 20px;

  max-width: 100%;

  min-width: 0;

}

.page-container :deep(.el-card__body) {

  max-width: 100%;

  overflow-x: hidden;

}



.card-header {

  display: flex;

  align-items: center;

  justify-content: space-between;

}



.search-form {

  margin-bottom: 20px;

  width: 100%;

  min-width: 0;

}

.search-form :deep(.el-form-item__label) {

  white-space: nowrap;

}

.search-form :deep(.el-form--inline .el-form-item) {

  margin-right: 12px;

  margin-bottom: 8px;

}

.table-scroll {

  width: 100%;

  max-width: 100%;

  overflow-x: auto;

}



.pagination {

  margin-top: 20px;

  display: flex;

  justify-content: flex-end;

}

.table-header {

  margin-bottom: 12px;

  display: flex;

  gap: 12px;

  flex-wrap: wrap;

}

</style>

