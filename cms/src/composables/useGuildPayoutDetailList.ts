import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

export type GuildPayoutDetailListScope = 'normal' | 'coinMerchant'

const scopeConfig = {
  normal: {
    pagePermission: 'GuildPayoutDetailList',
    pageTitleKey: 'menu.GuildPayoutDetailList',
    guildType: 0 as const,
    normalGuildOnly: true,
    coinMerchantGuildOnly: false,
  },
  coinMerchant: {
    pagePermission: 'CoinMerchantPayoutDetailList',
    pageTitleKey: 'menu.CoinMerchantPayoutDetailList',
    guildType: 1 as const,
    normalGuildOnly: false,
    coinMerchantGuildOnly: true,
  },
} as const

export function useGuildPayoutDetailList(scope: GuildPayoutDetailListScope) {
  const {t} = useI18n()
  const cfg = scopeConfig[scope]
  const {can} = usePagePermission(cfg.pagePermission)

  const pageTitle = computed(() => t(cfg.pageTitleKey))
  const loading = ref(false)
  const tableData = ref<GuildIncomeSettlementLogItem[]>([])

  const searchForm = reactive({
    dateRange: [] as string[],
  })

  const pagination = reactive({
    pageIndex: 1,
    pageSize: 20,
    total: 0,
  })

  const buildQueryParams = () => {
    const [startDate, endDate] = searchForm.dateRange || []
    return {
      guildType: cfg.guildType,
      normalGuildOnly: cfg.normalGuildOnly,
      coinMerchantGuildOnly: cfg.coinMerchantGuildOnly,
      payoutOnly: true,
      transferStartTime: startDate ? toServerDayStartUnix(startDate) : 0,
      transferEndTime: endDate ? toServerDayEndUnix(endDate) : 0,
      includeTransfer: true,
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
    }
  }

  const fetchList = async () => {
    loading.value = true
    try {
      const response = await guildIncomeSettlementLogApi.getList(buildQueryParams())
      tableData.value = response.data || []
      pagination.total = response.total || 0
    } catch (error) {
      console.error('Failed to load guild payout details:', error)
      ElMessage.error(t('pages.guildTransferList.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  const payoutTime = (row: GuildIncomeSettlementLogItem) => row.transferAt || row.updatedAt || null

  const statusLabel = (status: number | undefined) => {
    const value = Number(status)
    if (value === 1) return t('pages.guildTransferList.statusApproved')
    if (value === 2) return t('pages.guildTransferList.statusTransferred')
    if (value === 3) return t('pages.guildTransferList.statusTransferring')
    return t('pages.guildTransferList.statusPending')
  }

  const statusTagType = (status: number | undefined) => {
    const value = Number(status)
    if (value === 1) return 'success'
    if (value === 2) return 'info'
    if (value === 3) return ''
    return 'warning'
  }

  const handleSearch = () => {
    pagination.pageIndex = 1
    void fetchList()
  }

  const handleReset = () => {
    searchForm.dateRange = []
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

  onMounted(() => {
    void fetchList()
  })

  return {
    t,
    can,
    pageTitle,
    loading,
    tableData,
    searchForm,
    pagination,
    payoutTime,
    statusLabel,
    statusTagType,
    handleSearch,
    handleReset,
    handlePageChange,
    handleSizeChange,
  }
}
