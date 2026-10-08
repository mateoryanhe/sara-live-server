import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

export type GuildTransferListScope = 'normal' | 'coinMerchant'

const scopeConfig = {
  normal: {
    pagePermission: 'GuildTransferManagement',
    detailRouteName: 'GuildTransferDetail',
    pageTitleKey: 'menu.GuildTransferManagement',
    guildType: 0 as const,
    normalGuildOnly: true,
    coinMerchantGuildOnly: false,
    transferInfoFrom: 'guildTransfer',
  },
  coinMerchant: {
    pagePermission: 'CoinMerchantGuildTransferManagement',
    detailRouteName: 'CoinMerchantGuildTransferDetail',
    pageTitleKey: 'menu.CoinMerchantGuildTransferManagement',
    guildType: 1 as const,
    normalGuildOnly: false,
    coinMerchantGuildOnly: true,
    transferInfoFrom: 'coinMerchantGuildTransfer',
  },
} as const

export function useGuildTransferList(scope: GuildTransferListScope) {
  const {t} = useI18n()
  const router = useRouter()
  const cfg = scopeConfig[scope]
  const {can} = usePagePermission(cfg.pagePermission)

  const pageTitle = computed(() => t(cfg.pageTitleKey))
  const listHint = computed(() => t('pages.guildTransferList.historyHint'))

  const loading = ref(false)
  const tableData = ref<GuildIncomeSettlementLogItem[]>([])
  const selectedRows = ref<GuildIncomeSettlementLogItem[]>([])
  const reopeningId = ref('')
  const copyingId = ref('')
  const editReceivableVisible = ref(false)
  const editReceivableSaving = ref(false)
  const editReceivableRow = ref<GuildIncomeSettlementLogItem | null>(null)
  const editReceivableAmount = ref(0)

  const settlementStatus = (row: GuildIncomeSettlementLogItem) => Number(row.status)
  const selectedPendingRows = computed(() => selectedRows.value.filter(row => settlementStatus(row) === 0))
  const selectedApprovedRows = computed(() => selectedRows.value.filter(row => settlementStatus(row) === 1))
  const selectedPendingCount = computed(() => selectedPendingRows.value.length)
  const selectedApprovedCount = computed(() => selectedApprovedRows.value.length)

  const searchForm = reactive({
    guildId: '',
    status: -1 as number | undefined,
  })

  const pagination = reactive({
    pageIndex: 1,
    pageSize: 20,
    total: 0,
  })

  const buildQueryParams = () => ({
    guildId: searchForm.guildId.trim(),
    guildType: cfg.guildType,
    normalGuildOnly: cfg.normalGuildOnly,
    coinMerchantGuildOnly: cfg.coinMerchantGuildOnly,
    status: searchForm.status !== undefined && searchForm.status !== null && searchForm.status >= 0
        ? searchForm.status
        : undefined,
    hideHistoricalTransferred: true,
    orderByReceivableUsdDesc: true,
    includeDetail: false,
    includeTransfer: true,
    includeTransferInfo: false,
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
  })

  const fetchList = async () => {
    loading.value = true
    try {
      const response = await guildIncomeSettlementLogApi.getList(buildQueryParams())
      tableData.value = response.data || []
      pagination.total = response.total || 0
    } catch (error) {
      console.error('Failed to load guild transfer list:', error)
      ElMessage.error(t('pages.guildTransferList.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  const handleSearch = () => {
    pagination.pageIndex = 1
    void fetchList()
  }

  const handleReset = () => {
    searchForm.guildId = ''
    searchForm.status = -1
    pagination.pageIndex = 1
    void fetchList()
  }

  const handlePageChange = (page: number) => {
    pagination.pageIndex = page
    void fetchList()
  }

  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.pageIndex = 1
    void fetchList()
  }

  const handleSelectionChange = (rows: GuildIncomeSettlementLogItem[]) => {
    selectedRows.value = rows || []
  }

  const hasRowAction = (row: GuildIncomeSettlementLogItem) => (
      can('viewDetail')
      || (settlementStatus(row) === 1 && can('reopenApproval'))
      || (settlementStatus(row) === 2 && can('copyPayout'))
      || (settlementStatus(row) === 0 && can('editReceivable'))
  )

  const handleReopenApproval = async (row: GuildIncomeSettlementLogItem) => {
    const id = String(row.id || '').trim()
    if (!id) return
    try {
      await ElMessageBox.confirm(
          t('pages.guildTransferList.reopenConfirm', {guild: row.guildName || row.guildId || '-'}),
          t('pages.guildTransferList.reopenApproval'),
          {type: 'warning'},
      )
    } catch {
      return
    }
    reopeningId.value = id
    try {
      await guildIncomeSettlementLogApi.reopenApproval({id})
      ElMessage.success(t('pages.guildTransferList.reopenSuccess'))
      await fetchList()
    } catch (error) {
      console.error('Reopen guild settlement approval failed:', error)
      ElMessage.error(t('pages.guildTransferList.reopenFailed'))
    } finally {
      reopeningId.value = ''
    }
  }

  const handleCopyPayout = async (row: GuildIncomeSettlementLogItem) => {
    const id = String(row.id || '').trim()
    if (!id) return
    try {
      await ElMessageBox.confirm(
          t('pages.guildTransferList.copyPayoutConfirm', {guild: row.guildName || row.guildId || '-'}),
          t('pages.guildTransferList.copyPayout'),
          {type: 'warning'},
      )
    } catch {
      return
    }
    copyingId.value = id
    try {
      const res = await guildIncomeSettlementLogApi.copyPayout({id})
      ElMessage.success(t('pages.guildTransferList.copyPayoutSuccess', {id: res.id || '-'}))
      searchForm.status = -1
      pagination.pageIndex = 1
      selectedRows.value = []
      await fetchList()
      const copiedId = String(res.id || '').trim()
      if (copiedId && !tableData.value.some(item => String(item.id) === copiedId)) {
        tableData.value.unshift({
          ...row,
          id: copiedId,
          status: 0,
          transferAt: undefined,
          transferOrderId: '',
          transferPlatformNo: '',
          transferLocalAmount: 0,
          transferCurrency: '',
          transferFailMsg: '',
          createdAt: new Date().toISOString(),
        })
        pagination.total += 1
      }
    } catch (error) {
      console.error('Copy guild settlement payout failed:', error)
      ElMessage.error(t('pages.guildTransferList.copyPayoutFailed'))
    } finally {
      copyingId.value = ''
    }
  }

  const openGuildDetail = (row: GuildIncomeSettlementLogItem) => {
    const guildId = String(row.guildId || '').trim()
    if (!guildId) return
    router.push({
      name: 'GuildDetail',
      query: {
        id: guildId,
        name: row.guildName || '',
      },
    })
  }

  const openSettlementDetail = (row: GuildIncomeSettlementLogItem) => {
    const id = String(row.id || '').trim()
    if (!id) return
    router.push({name: cfg.detailRouteName, params: {id}})
  }

  const openEditReceivable = (row: GuildIncomeSettlementLogItem) => {
    editReceivableRow.value = row
    editReceivableAmount.value = Number(row.settlementReceivableUsd || 0)
    editReceivableVisible.value = true
  }

  const handleRowCommand = (row: GuildIncomeSettlementLogItem, command: string) => {
    switch (command) {
      case 'viewDetail':
        openSettlementDetail(row)
        break
      case 'reopenApproval':
        void handleReopenApproval(row)
        break
      case 'copyPayout':
        void handleCopyPayout(row)
        break
      case 'editReceivable':
        openEditReceivable(row)
        break
    }
  }

  const handleSaveReceivable = async () => {
    const row = editReceivableRow.value
    const id = String(row?.id || '').trim()
    const amount = Number(editReceivableAmount.value)
    if (!id || !Number.isFinite(amount) || amount <= 0) {
      ElMessage.warning(t('pages.guildTransferList.receivableInvalid'))
      return
    }
    editReceivableSaving.value = true
    try {
      await guildIncomeSettlementLogApi.updateReceivableUsd({
        id,
        settlementReceivableUsd: amount,
      })
      ElMessage.success(t('pages.guildTransferList.updateReceivableSuccess'))
      editReceivableVisible.value = false
      await fetchList()
    } catch (error) {
      console.error('Update guild settlement receivable failed:', error)
      ElMessage.error(t('pages.guildTransferList.updateReceivableFailed'))
    } finally {
      editReceivableSaving.value = false
    }
  }

  const handleBatchApprove = async () => {
    if (!selectedRows.value.length) {
      ElMessage.warning(t('pages.guildTransferList.selectRows'))
      return
    }
    const pendingIds = selectedPendingRows.value
        .map(row => String(row.id || '').trim())
        .filter(Boolean)
    if (!pendingIds.length) {
      ElMessage.warning(t('pages.guildTransferList.selectPendingRows'))
      return
    }
    try {
      await ElMessageBox.confirm(
          t('pages.guildTransferList.approveConfirm', {count: pendingIds.length}),
          t('pages.guildTransferList.batchApprove'),
          {type: 'warning'},
      )
    } catch {
      return
    }
    try {
      const res = await guildIncomeSettlementLogApi.batchApprove({ids: pendingIds})
      ElMessage.success(t('pages.guildTransferList.approveSuccess', {
        success: res.successCount || 0,
        fail: res.failCount || 0,
      }))
      await fetchList()
    } catch (error) {
      console.error('Batch approve failed:', error)
      ElMessage.error(t('pages.guildTransferList.approveFailed'))
    }
  }

  const handleBatchTransfer = async () => {
    if (!selectedRows.value.length) {
      ElMessage.warning(t('pages.guildTransferList.selectRows'))
      return
    }
    const approvedIds = selectedApprovedRows.value
        .map(row => String(row.id || '').trim())
        .filter(Boolean)
    if (!approvedIds.length) {
      ElMessage.warning(t('pages.guildTransferList.selectApprovedRows'))
      return
    }
    try {
      await ElMessageBox.confirm(
          t('pages.guildTransferList.transferConfirm', {count: approvedIds.length}),
          t('pages.guildTransferList.batchTransfer'),
          {type: 'warning'},
      )
    } catch {
      return
    }
    try {
      const res = await guildIncomeSettlementLogApi.batchTransfer({ids: approvedIds})
      const success = res.successCount || 0
      const fail = res.failCount || 0
      if (success > 0) {
        ElMessage.success(res.message || t('pages.guildTransferList.transferSuccess', {success, fail}))
      } else {
        ElMessage.warning(res.message || t('pages.guildTransferList.transferFailed'))
      }
      await fetchList()
    } catch (error) {
      console.error('Batch transfer failed:', error)
      ElMessage.error(t('pages.guildTransferList.transferFailed'))
    }
  }

  const statusLabel = (status: number | undefined) => {
    const normalizedStatus = Number(status)
    if (normalizedStatus === 1) return t('pages.guildTransferList.statusApproved')
    if (normalizedStatus === 2) return t('pages.guildTransferList.statusTransferred')
    if (normalizedStatus === 3) return t('pages.guildTransferList.statusTransferring')
    return t('pages.guildTransferList.statusPending')
  }

  const statusTagType = (status: number | undefined) => {
    const normalizedStatus = Number(status)
    if (normalizedStatus === 1) return 'success'
    if (normalizedStatus === 2) return 'info'
    if (normalizedStatus === 3) return ''
    return 'warning'
  }

  onMounted(() => {
    void fetchList()
  })

  return {
    t,
    can,
    pageTitle,
    listHint,
    loading,
    tableData,
    selectedRows,
    reopeningId,
    copyingId,
    editReceivableVisible,
    editReceivableSaving,
    editReceivableRow,
    editReceivableAmount,
    searchForm,
    pagination,
    selectedPendingCount,
    selectedApprovedCount,
    settlementStatus,
    hasRowAction,
    handleSearch,
    handleReset,
    handlePageChange,
    handleSizeChange,
    handleSelectionChange,
    handleRowCommand,
    handleSaveReceivable,
    handleBatchApprove,
    handleBatchTransfer,
    statusLabel,
    statusTagType,
    openGuildDetail,
    openSettlementDetail,
  }
}
