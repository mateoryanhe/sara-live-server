<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>基础数据</span>
          <el-button :loading="loading" @click="fetchSysStat">刷新</el-button>
        </div>
      </template>

      <div class="basic-stat-rows">
        <div
            v-for="(row, rowIndex) in basicStatCardRows"
            :key="rowIndex"
            class="basic-stat-row"
        >
          <div
              v-for="card in row"
              :key="card.key"
              class="stat-card"
              :class="card.theme"
          >
            <div class="stat-label">{{ card.label }}</div>
            <div class="stat-value">
              {{ card.format === 'count' ? formatCount(card.value) : formatAmount(card.value) }}
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="user-stat-card">
      <template #header>
        <div class="card-header">
          <span>用户数据</span>
          <el-button :loading="trendLoading" @click="fetchUserStatTrend">刷新</el-button>
        </div>
      </template>

      <el-tabs v-model="activePeriod" @tab-change="handleTabChange">
        <el-tab-pane label="日" name="daily">
          <UserStatChart ref="dailyLineChartRef" :data="userStatTrend.daily" title="活跃用户 / 新注册 (最近30天)"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="dailyBarChartRef"
                :data="userStatTrend.daily"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
        <el-tab-pane label="周" name="weekly">
          <UserStatChart ref="weeklyLineChartRef" :data="userStatTrend.weekly" title="活跃用户 / 新注册 (最近12周)"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="weeklyBarChartRef"
                :data="userStatTrend.weekly"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
        <el-tab-pane label="月" name="monthly">
          <UserStatChart ref="monthlyLineChartRef" :data="userStatTrend.monthly" title="活跃用户 / 新注册 (最近12月)"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="monthlyBarChartRef"
                :data="userStatTrend.monthly"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {sysStatApi} from '@/api'
import type {SysStat, UserStatTrend} from '@/types/api'
import UserStatChart from './components/user-stat-chart.vue'
import UserStatBarChart from './components/user-stat-bar-chart.vue'
import BarMetricSection from './components/bar-metric-section.vue'
import {getUserStatBarMetricTabs, USER_STAT_BAR_SERIES} from './user-stat-bar-series'

const loading = ref(false)
const trendLoading = ref(false)
const activePeriod = ref('daily')
const barMetricTabs = getUserStatBarMetricTabs()
const activeBarMetric = ref(barMetricTabs[0]?.key ?? 'rechargeUser')

const barChartPeriodSuffix = computed(() => {
  if (activePeriod.value === 'weekly') {
    return '最近12周'
  }
  if (activePeriod.value === 'monthly') {
    return '最近12月'
  }
  return '最近30天'
})

const barChartTitle = computed(() => {
  const metric = USER_STAT_BAR_SERIES.find((item) => item.key === activeBarMetric.value)
  return `${metric?.label ?? ''} (${barChartPeriodSuffix.value})`
})

const dailyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const weeklyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const monthlyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const dailyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const weeklyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const monthlyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()

const sysStat = reactive<SysStat>({
  totalGold: 0,
  totalGoldConsume: 0,
  totalDiamondConsume: 0,
  totalRecharge: 0,
  totalWithdraw: 0,
  totalRegisterUser: 0,
  todayRecharge: 0,
  todayGoldConsume: 0,
  todayDiamondConsume: 0,
  todayRegisterUser: 0,
})

const userStatTrend = reactive<UserStatTrend>({
  daily: [],
  weekly: [],
  monthly: [],
})

type BasicStatCardKey = keyof Pick<
    SysStat,
    | 'totalGold'
    | 'todayRecharge'
    | 'todayGoldConsume'
    | 'todayDiamondConsume'
    | 'todayRegisterUser'
    | 'totalRecharge'
    | 'totalWithdraw'
    | 'totalGoldConsume'
    | 'totalDiamondConsume'
    | 'totalRegisterUser'
>

type BasicStatCardConfig = {
  key: BasicStatCardKey
  label: string
  theme: string
  format: 'amount' | 'count'
}

type BasicStatCard = BasicStatCardConfig & {
  value: number | string
}

/** 基础数据展示顺序(自上而下、从左到右); 每 5 项一行 */
const BASIC_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'totalGold', label: '金币钱包现存总金额', theme: 'stat-card-gold', format: 'amount'},
  {key: 'todayRecharge', label: '今日充值金额', theme: 'stat-card-today-recharge', format: 'amount'},
  {key: 'todayGoldConsume', label: '今日消费金额-金币', theme: 'stat-card-today-gold-consume', format: 'amount'},
  {key: 'todayDiamondConsume', label: '今日消费金额-钻石', theme: 'stat-card-today-diamond-consume', format: 'amount'},
  {key: 'todayRegisterUser', label: '今日注册用户数', theme: 'stat-card-today-register', format: 'count'},
  {key: 'totalRecharge', label: '总充值金额', theme: 'stat-card-recharge', format: 'amount'},
  {key: 'totalWithdraw', label: '总提现金额', theme: 'stat-card-withdraw', format: 'amount'},
  {key: 'totalGoldConsume', label: '总消费金额-金币', theme: 'stat-card-total-gold-consume', format: 'amount'},
  {key: 'totalDiamondConsume', label: '总消费金额-钻石', theme: 'stat-card-total-diamond-consume', format: 'amount'},
  {key: 'totalRegisterUser', label: '总注册用户数', theme: 'stat-card-register', format: 'count'},
]

const BASIC_STAT_ROW_SIZE = 5

const buildBasicStatCards = (): BasicStatCard[] =>
    BASIC_STAT_CARD_CONFIG.map((cfg) => ({
      ...cfg,
      value: sysStat[cfg.key] ?? 0,
    }))

/** 按配置顺序拆成 2 行 × 5 列 */
const basicStatCardRows = computed<BasicStatCard[][]>(() => {
  const cards = buildBasicStatCards()
  const rows: BasicStatCard[][] = []
  for (let i = 0; i < cards.length; i += BASIC_STAT_ROW_SIZE) {
    rows.push(cards.slice(i, i + BASIC_STAT_ROW_SIZE))
  }
  return rows
})

const fetchSysStat = async () => {
  loading.value = true
  try {
    const data = await sysStatApi.getSysStat()
    sysStat.totalGold = data.totalGold ?? 0
    sysStat.totalGoldConsume = data.totalGoldConsume ?? 0
    sysStat.totalDiamondConsume = data.totalDiamondConsume ?? 0
    sysStat.totalRecharge = data.totalRecharge ?? 0
    sysStat.totalWithdraw = data.totalWithdraw ?? 0
    sysStat.totalRegisterUser = data.totalRegisterUser ?? 0
    sysStat.todayRecharge = data.todayRecharge ?? 0
    sysStat.todayGoldConsume = data.todayGoldConsume ?? 0
    sysStat.todayDiamondConsume = data.todayDiamondConsume ?? 0
    sysStat.todayRegisterUser = data.todayRegisterUser ?? 0
  } catch (error) {
    console.error('获取系统总数据失败:', error)
    ElMessage.error('获取系统总数据失败')
  } finally {
    loading.value = false
  }
}

const fetchUserStatTrend = async () => {
  trendLoading.value = true
  try {
    const data = await sysStatApi.getUserStatTrend()
    userStatTrend.daily = data.daily || []
    userStatTrend.weekly = data.weekly || []
    userStatTrend.monthly = data.monthly || []
    await nextTick()
    setTimeout(() => {
      resizeActiveChart()
    }, 0)
  } catch (error) {
    console.error('获取用户数据趋势失败:', error)
    ElMessage.error('获取用户数据趋势失败')
  } finally {
    trendLoading.value = false
  }
}

const resizeActiveChart = () => {
  if (activePeriod.value === 'daily') {
    dailyLineChartRef.value?.resize()
    dailyBarChartRef.value?.resize()
    return
  }
  if (activePeriod.value === 'weekly') {
    weeklyLineChartRef.value?.resize()
    weeklyBarChartRef.value?.resize()
    return
  }
  monthlyLineChartRef.value?.resize()
  monthlyBarChartRef.value?.resize()
}

const handleTabChange = async () => {
  await nextTick()
  resizeActiveChart()
}

const handleBarMetricChange = async () => {
  await nextTick()
  resizeActiveChart()
}

const formatAmount = (value: number | null | undefined) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

const formatCount = (value: string | number | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchSysStat()
  fetchUserStatTrend()
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

.user-stat-card {
  margin-top: 20px;
}

.basic-stat-rows {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.basic-stat-row {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 20px;
}

.stat-card {
  border-radius: 12px;
  padding: 20px 16px;
  color: #fff;
  min-height: 108px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 12px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  word-break: break-all;
}

@media (max-width: 1200px) {
  .basic-stat-row {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .basic-stat-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .stat-value {
    font-size: 20px;
  }
}

.stat-card-gold {
  background: linear-gradient(135deg, #f6ad55, #ed8936);
}

.stat-card-recharge {
  background: linear-gradient(135deg, #63b3ed, #3182ce);
}

.stat-card-withdraw {
  background: linear-gradient(135deg, #68d391, #38a169);
}

.stat-card-register {
  background: linear-gradient(135deg, #b794f4, #805ad5);
}

.stat-card-today-recharge {
  background: linear-gradient(135deg, #4fd1c5, #319795);
}

.stat-card-today-register {
  background: linear-gradient(135deg, #f687b3, #d53f8c);
}

.stat-card-total-gold-consume {
  background: linear-gradient(135deg, #fbd38d, #dd6b20);
}

.stat-card-today-gold-consume {
  background: linear-gradient(135deg, #fc8181, #c53030);
}

.stat-card-total-diamond-consume {
  background: linear-gradient(135deg, #9f7aea, #6b46c1);
}

.stat-card-today-diamond-consume {
  background: linear-gradient(135deg, #76e4f7, #3182ce);
}
</style>
