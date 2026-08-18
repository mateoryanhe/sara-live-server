<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorList') }}</span>
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
        <el-table-column fixed="right" :label="t('common.actions')" width="220">
          <template #default="{ row }">
            <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
              {{ t('common.detail') }}
            </el-button>
            <el-button v-if="can('offShelf')" type="warning" link @click="handleOffShelf(row)">
              {{ t('common.offShelf') }}
            </el-button>
            <el-button
                v-if="row.ban ? can('unban') : can('ban')"
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
import {useRouter} from 'vue-router'
import {ElForm, ElMessage, ElMessageBox, type FormRules} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('PlatformAnchorList')
const canViewDetail = computed(() => can('viewDetail'))

const loading = ref(false)
const tableData = ref<AnchorListItem[]>([])
const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<InstanceType<typeof ElForm>>()

const searchForm = reactive({key: ''})
const pagination = reactive({pageIndex: 1, pageSize: 10, total: 0})
const banForm = reactive({accountId: '', nickname: '', banApplyTime: '', banReason: ''})

const banRules = computed<FormRules>(() => ({
  banApplyTime: [{required: true, message: t('pages.anchorList.banApplyTimeRequired'), trigger: 'change'}],
  banReason: [
    {required: true, message: t('pages.anchorList.banReasonRequired'), trigger: 'blur'},
    {min: 1, max: 512, message: t('pages.anchorList.banReasonLength'), trigger: 'blur'},
  ],
}))

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

const categoryLabel = (category?: number) => {
  if (category === 1) return t('pages.anchorList.categoryHot')
  if (category === 2) return t('pages.anchorList.categoryGame')
  if (category === 3) return t('pages.anchorList.categoryPrivate')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === 3) return 'warning'
  if (category === 2) return 'success'
  if (category === 1) return 'danger'
  return 'info'
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
      platformOnly: true,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load platform anchor list:', error)
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
      const response = await accountApi.banAnchor({
        accountId: banForm.accountId,
        banApplyTime: banForm.banApplyTime,
        banReason: banForm.banReason.trim(),
      } as BanAnchorReq)
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
      const response = await accountApi.unBanAnchor({accountId: row.id} as UnBanAnchorReq)
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
    await accountApi.setLiveRoomStatus({anchorId: row.id, status: 0})
    ElMessage.success(t('pages.anchorList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') return
    console.error('off shelf live room failed:', error)
    ElMessage.error(t('pages.anchorList.offShelfFailed'))
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

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
