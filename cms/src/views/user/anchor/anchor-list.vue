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
        <el-table-column
            :label="t('pages.anchorList.liveRoomCover')"
            align="center"
            class-name="cover-col"
            label-class-name="cover-col"
            width="100"
        >
          <template #default="{ row }">
            <div class="cover-cell">
              <el-image
                  v-if="listCoverUrl(row)"
                  :preview-src-list="[listCoverUrl(row)]"
                  :src="listCoverUrl(row)"
                  fit="cover"
                  hide-on-click-modal
                  preview-teleported
                  :class="isAvatarCoverFallback(row) ? 'cover-cell-avatar' : 'cover-cell-image'"
              />
              <span v-else>-</span>
            </div>
          </template>
        </el-table-column>
        
         <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>


 <el-table-column :label="t('common.userId')" prop="id" width="180">
          <template #default="{ row }">
            <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
              {{ row.id }}
            </el-button>
            <span v-else>{{ row.id }}</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('menu.UserDetail')" width="110">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail" link type="primary" @click="openUserDetail(row.id)">
              {{ t('pages.userList.viewDetail') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
       
       
       
        <el-table-column :label="t('pages.anchorList.unsettledTotalIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
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
       
        <el-table-column :label="t('pages.anchorList.liveStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? t('common.live') : t('common.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.unsettledGiftIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.unsettledPaidDanmakuIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.unsettledVideoTicketIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.unsettledVideoBillingIncome')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.unsettledVideoCallIncome')" align="right" min-width="140">
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
         <el-table-column :label="t('pages.anchorList.anchorType')" width="110">
          <template #default="{ row }">
            <el-tag :type="anchorTypeTagType(row.userType)">{{ anchorTypeLabel(row.userType) }}</el-tag>
          </template>
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
            <span class="money-amount">{{ formatWalletBalance(row.ticket) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.billingPricePerMinute')" align="right" min-width="120">
          <template #default="{ row }">
            <span class="money-amount">{{ formatWalletBalance(row.billing) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.roomTitle')" min-width="140" prop="roomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.registeredAt')" prop="registeredAt" width="170">
          <template #default="{ row }">{{ formatDate(row.registeredAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.profileUpdatedAt')" prop="createdAt" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="320">
          <template #default="{ row }">
            <el-button
                v-if="can('uploadRoomCover')"
                link
                type="primary"
                @click="openRoomCoverDialog(row)"
            >
              {{ t('pages.anchorList.uploadRoomCover') }}
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

    <el-dialog
        v-model="roomCoverDialogVisible"
        :title="t('pages.anchorList.uploadRoomCoverTitle')"
        width="440px"
        @closed="resetRoomCoverForm"
    >
      <el-form label-width="100px">
        <el-form-item :label="t('pages.anchorList.anchorId')">
          <el-input v-model="roomCoverForm.anchorId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="roomCoverForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.liveRoomCover')">
          <div class="room-cover-upload-wrap">
            <el-upload
                :before-upload="beforeRoomCoverUpload"
                :disabled="roomCoverUploading"
                :http-request="doRoomCoverUpload"
                :show-file-list="false"
                accept="image/*"
                class="room-cover-uploader"
            >
              <el-image
                  v-if="roomCoverPreviewUrl"
                  :src="roomCoverPreviewUrl"
                  fit="cover"
                  style="width:120px;height:120px;border-radius:4px"
              />
              <div v-else class="room-cover-uploader-placeholder">
                <el-icon class="room-cover-uploader-icon">
                  <Plus/>
                </el-icon>
              </div>
            </el-upload>
            <el-button v-if="roomCoverForm.cover || roomCoverPreviewUrl" link type="danger" @click="clearRoomCover">
              {{ t('pages.anchorList.clearRoomCover') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roomCoverDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="roomCoverSubmitting" type="primary" @click="submitRoomCover">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

  </div>

</template>



<script lang="ts" setup>

import {computed, onMounted, reactive, ref} from 'vue'

import {useI18n} from 'vue-i18n'

import {useRouter} from 'vue-router'

import {ElForm, ElMessage, ElMessageBox, type FormRules, type UploadRequestOptions} from 'element-plus'

import {accountApi, uploadApi} from '@/api'

import {Plus} from '@element-plus/icons-vue'

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

const roomCoverDialogVisible = ref(false)
const roomCoverSubmitting = ref(false)
const roomCoverUploading = ref(false)
const roomCoverPreviewUrl = ref('')
const roomCoverChanged = ref(false)
const roomCoverForm = reactive({
  anchorId: '',
  nickname: '',
  cover: '',
})



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

const listCoverUrl = (row: AnchorListItem) => row.roomCover || row.avatar || ''

const isAvatarCoverFallback = (row: AnchorListItem) => !row.roomCover && !!row.avatar

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

const resetRoomCoverForm = () => {
  roomCoverForm.anchorId = ''
  roomCoverForm.nickname = ''
  roomCoverForm.cover = ''
  roomCoverPreviewUrl.value = ''
  roomCoverChanged.value = false
}

const openRoomCoverDialog = (row: AnchorListItem) => {
  roomCoverForm.anchorId = String(row.id)
  roomCoverForm.nickname = row.nickname || '-'
  roomCoverForm.cover = ''
  roomCoverPreviewUrl.value = row.roomCover || row.avatar || ''
  roomCoverChanged.value = false
  roomCoverDialogVisible.value = true
}

const beforeRoomCoverUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.userList.imageOnly'))
    return false
  }
  return true
}

const doRoomCoverUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  roomCoverUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    roomCoverForm.cover = res.fileName
    roomCoverChanged.value = true
    roomCoverPreviewUrl.value = URL.createObjectURL(file)
    ElMessage.success(t('pages.userList.uploadSuccess'))
  } catch (error) {
    console.error('upload room cover failed:', error)
    ElMessage.error(t('pages.userList.uploadFailed'))
  } finally {
    roomCoverUploading.value = false
  }
}

const clearRoomCover = () => {
  roomCoverForm.cover = ''
  roomCoverPreviewUrl.value = ''
  roomCoverChanged.value = true
}

const patchAnchorRow = (anchorId: string, patch: Partial<AnchorListItem>) => {
  const idx = tableData.value.findIndex(item => String(item.id) === anchorId)
  if (idx >= 0) {
    tableData.value[idx] = {...tableData.value[idx], ...patch}
  }
}

const submitRoomCover = async () => {
  if (!roomCoverChanged.value) {
    roomCoverDialogVisible.value = false
    return
  }
  roomCoverSubmitting.value = true
  try {
    const res = await accountApi.setLiveRoomCover({
      anchorId: roomCoverForm.anchorId,
      cover: roomCoverForm.cover,
    })
    if (res?.success) {
      patchAnchorRow(roomCoverForm.anchorId, {roomCover: res.cover || ''})
      roomCoverDialogVisible.value = false
      ElMessage.success(t('pages.anchorList.uploadRoomCoverSuccess'))
    } else {
      ElMessage.error(t('pages.userList.uploadFailed'))
    }
  } catch (error) {
    console.error('setLiveRoomCover failed:', error)
  } finally {
    roomCoverSubmitting.value = false
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

.search-form :deep(.el-form-item) {
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

.room-cover-upload-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.room-cover-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 4px;
  cursor: pointer;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.room-cover-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.room-cover-uploader-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.room-cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

:deep(.cover-col .cell) {
  white-space: nowrap;
}

.cover-cell {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  vertical-align: middle;
  line-height: 1;
}

.cover-cell-image {
  width: 56px;
  height: 56px;
  border-radius: 4px;
  display: block;
  flex-shrink: 0;
}

.cover-cell-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: block;
  flex-shrink: 0;
}

</style>

