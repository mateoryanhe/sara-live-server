<template>
  <div class="page-container">
    <el-card v-loading="serverTimeLoading">
      <template #header>
        <div class="card-header">
          <span>服务器日志</span>
          <span class="server-time-tip">服务器时间: {{ serverTimeDisplay }}</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="访问统计" name="stats">
          <el-form :model="statsForm" class="search-form" inline label-width="90px">
            <el-form-item label="开始日期">
              <el-date-picker
                  v-model="statsForm.startDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="开始日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="结束日期">
              <el-date-picker
                  v-model="statsForm.endDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="结束日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="TopN">
              <el-input-number v-model="statsForm.topN" :max="100" :min="1" controls-position="right" style="width: 120px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchAccessStats">查询</el-button>
            </el-form-item>
          </el-form>

          <el-row v-loading="statsLoading" :element-loading-text="queryStatusTip" :gutter="20">
            <el-col :span="12">
              <el-card shadow="never">
                <template #header>接口访问 TopN</template>
                <el-table :data="urlTopData" size="small">
                  <el-table-column label="#" type="index" width="50"/>
                  <el-table-column label="URL" min-width="220" prop="key" show-overflow-tooltip/>
                  <el-table-column label="访问数" prop="count" width="100"/>
                </el-table>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never">
                <template #header>IP访问 TopN</template>
                <el-table :data="ipTopData" size="small">
                  <el-table-column label="#" type="index" width="50"/>
                  <el-table-column label="IP" min-width="180" prop="key" show-overflow-tooltip/>
                  <el-table-column label="访问数" prop="count" width="100"/>
                </el-table>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <el-tab-pane label="Access日志" name="access">
          <el-form :model="accessForm" class="search-form" inline label-width="100px">
            <el-form-item label="开始日期">
              <el-date-picker
                  v-model="accessForm.startDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="开始日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="结束日期">
              <el-date-picker
                  v-model="accessForm.endDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="结束日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="TraceId">
              <el-input v-model="accessForm.traceId" clearable placeholder="支持模糊匹配" style="width: 280px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="accessForm.url" clearable placeholder="支持模糊匹配" style="width: 220px"/>
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="accessForm.ip" clearable placeholder="支持模糊匹配" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="状态码">
              <el-input-number v-model="accessForm.statusCode" :controls="false" :min="0" placeholder="全部" style="width: 120px"/>
            </el-form-item>
            <el-form-item label="处理耗时ms">
              <el-input-number v-model="accessForm.minHandlerMs" :controls="false" :min="0" placeholder="最小" style="width: 100px"/>
              <span class="range-sep">-</span>
              <el-input-number v-model="accessForm.maxHandlerMs" :controls="false" :min="0" placeholder="最大" style="width: 100px"/>
            </el-form-item>
            <el-form-item label="聚合粒度">
              <el-select v-model="accessForm.intervalMinutes" placeholder="自动" style="width: 120px">
                <el-option :value="0" label="自动"/>
                <el-option :value="1" label="1分钟"/>
                <el-option :value="5" label="5分钟"/>
                <el-option :value="15" label="15分钟"/>
                <el-option :value="60" label="1小时"/>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleAccessSearch">查询</el-button>
              <el-button @click="resetAccessForm">重置</el-button>
            </el-form-item>
          </el-form>

          <div v-loading="accessTrendLoading" :element-loading-text="queryStatusTip" class="access-trend-panel">
            <AccessTrendChart ref="accessTrendChartRef" :data="accessTrendData"/>
            <div v-if="accessTrendData?.peakTime" class="trend-summary">
              总访问量 {{ accessTrendData.totalCount }}，
              峰值 {{ accessTrendData.peakCount }} 次/每{{ accessTrendData.intervalMinutes }}分钟
              ({{ accessTrendData.peakTime }})
            </div>
          </div>

          <el-table v-loading="accessLoading" :data="accessTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column label="时间" prop="time" width="210"/>
            <el-table-column label="TraceId" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column label="状态码" prop="statusCode" width="90"/>
            <el-table-column label="方法" prop="method" width="90"/>
            <el-table-column label="URL" min-width="220" prop="url" show-overflow-tooltip/>
            <el-table-column label="耗时(ms)" width="100">
              <template #default="{ row }">{{ formatHandlerMs(row.handlerMs) }}</template>
            </el-table-column>
            <el-table-column label="IP" min-width="160" prop="ip" show-overflow-tooltip/>
            <el-table-column label="UserAgent" min-width="180" prop="userAgent" show-overflow-tooltip/>
          </el-table>

          <div class="pagination-wrap">
            <el-pagination
                v-model:current-page="accessPageIndex"
                v-model:page-size="accessPageSize"
                :page-sizes="[20, 50, 100, 200]"
                :total="accessTotal"
                background
                layout="total, sizes, prev, pager, next"
                @current-change="fetchAccessLogs"
                @size-change="handleAccessSearch"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="Error日志" name="error">
          <el-form :model="errorForm" class="search-form" inline label-width="100px">
            <el-form-item label="开始日期">
              <el-date-picker
                  v-model="errorForm.startDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="开始日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="结束日期">
              <el-date-picker
                  v-model="errorForm.endDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="结束日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="TraceId">
              <el-input v-model="errorForm.traceId" clearable placeholder="支持模糊匹配" style="width: 280px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="errorForm.url" clearable placeholder="支持模糊匹配" style="width: 220px"/>
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="errorForm.ip" clearable placeholder="支持模糊匹配" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="状态码">
              <el-input-number v-model="errorForm.statusCode" :controls="false" :min="0" placeholder="全部" style="width: 120px"/>
            </el-form-item>
            <el-form-item label="关键词">
              <el-input v-model="errorForm.keyword" clearable placeholder="ErrorLog/堆栈/错误信息" style="width: 200px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleErrorSearch">查询</el-button>
              <el-button @click="resetErrorForm">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table v-loading="errorLoading" :data="errorTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column label="时间" prop="time" width="210"/>
            <el-table-column label="TraceId" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column label="状态码" prop="statusCode" width="90"/>
            <el-table-column label="方法" prop="method" width="90"/>
            <el-table-column label="URL" min-width="200" prop="url" show-overflow-tooltip/>
            <el-table-column label="IP" min-width="140" prop="ip" show-overflow-tooltip/>
            <el-table-column label="错误码" prop="errorCode" width="90"/>
            <el-table-column label="错误信息" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">{{ formatErrorSummary(row) }}</template>
            </el-table-column>
            <el-table-column label="详情" min-width="180" prop="detail" show-overflow-tooltip/>
            <el-table-column label="堆栈" min-width="220" show-overflow-tooltip>
              <template #default="{ row }">{{ formatErrorStack(row) }}</template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrap">
            <el-pagination
                v-model:current-page="errorPageIndex"
                v-model:page-size="errorPageSize"
                :page-sizes="[20, 50, 100, 200]"
                :total="errorTotal"
                background
                layout="total, sizes, prev, pager, next"
                @current-change="fetchErrorLogs"
                @size-change="handleErrorSearch"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="详情日志" lazy name="detail">
          <div class="detail-query-tip">查询已改为后台排队执行（单线程串行），通常 1 分钟内返回；建议填写 TraceId / URL 等条件</div>
          <el-form :model="detailForm" class="search-form" inline label-width="90px">
            <el-form-item label="开始日期">
              <el-date-picker
                  v-model="detailForm.startDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="开始日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="结束日期">
              <el-date-picker
                  v-model="detailForm.endDate"
                  clearable
                  format="YYYY-MM-DD"
                  placeholder="结束日期"
                  style="width: 140px"
                  teleported
                  type="date"
                  value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="TraceId">
              <el-input v-model="detailForm.traceId" clearable placeholder="支持模糊匹配" style="width: 280px"/>
            </el-form-item>
            <el-form-item label="ReqId">
              <el-input v-model="detailForm.reqId" clearable placeholder="支持模糊匹配" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="AuthId">
              <el-input v-model="detailForm.authId" clearable placeholder="支持模糊匹配(含推送)" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="detailForm.url" clearable placeholder="支持模糊匹配" style="width: 220px"/>
            </el-form-item>
            <el-form-item label="关键词">
              <el-input v-model="detailForm.keyword" clearable placeholder="ErrorLog/内容模糊匹配" style="width: 180px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleDetailSearch">查询</el-button>
              <el-button @click="resetDetailForm">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table v-loading="detailLoading" :data="detailTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column label="时间" prop="time" width="210"/>
            <el-table-column label="级别" prop="level" width="80"/>
            <el-table-column label="TraceId" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column label="ReqId" prop="reqId" width="160"/>
            <el-table-column label="AuthId" prop="authId" width="180"/>
            <el-table-column label="URL" min-width="180" prop="url" show-overflow-tooltip/>
            <el-table-column label="耗时(ms)" width="100">
              <template #default="{ row }">{{ formatHandlerMs(row.elapsedMs) }}</template>
            </el-table-column>
            <el-table-column label="内容" min-width="320" prop="message" show-overflow-tooltip/>
          </el-table>

          <div class="pagination-wrap">
            <el-pagination
                v-model:current-page="detailPageIndex"
                v-model:page-size="detailPageSize"
                :page-sizes="[20, 50, 100, 200]"
                :total="detailTotal"
                background
                layout="total, sizes, prev, pager, next"
                @current-change="fetchDetailLogs"
                @size-change="handleDetailSearch"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-drawer v-model="traceDrawerVisible" :title="`TraceId: ${traceDetail.traceId}`" size="60%">
      <div class="trace-drawer">
        <el-form class="trace-search-form" inline label-width="90px">
          <el-form-item label="开始日期">
            <el-date-picker
                v-model="traceStartDate"
                clearable
                format="YYYY-MM-DD"
                placeholder="开始日期"
                style="width: 140px"
                teleported
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="结束日期">
            <el-date-picker
                v-model="traceEndDate"
                clearable
                format="YYYY-MM-DD"
                placeholder="结束日期"
                style="width: 140px"
                teleported
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item>
            <el-button :loading="traceLoading" type="primary" @click="fetchTraceDetail">查询</el-button>
          </el-form-item>
        </el-form>
        <p v-if="resolveTraceAuthId()" class="trace-meta">AuthId: {{ resolveTraceAuthId() }}</p>

        <div v-loading="traceLoading" :element-loading-text="queryStatusTip">
        <h4>Error日志</h4>
        <div v-for="(item, index) in traceDetail.errorLogs" :key="'error-' + index" class="trace-log-line">
          <div class="trace-log-meta">
            <div class="trace-log-meta-main">
              <span>{{ item.time }}</span>
              <span v-if="hasElapsedMs(item.handlerMs)" :class="elapsedMsClass(item.handlerMs)" class="trace-elapsed">
                {{ formatElapsedMs(item.handlerMs) }}
              </span>
              <span v-if="item.authId" class="trace-auth-id">AuthId: {{ item.authId }}</span>
              <span>[{{ item.level }}] {{ formatErrorSummary(item) }}</span>
            </div>
          </div>
          <pre v-if="formatErrorStack(item)" class="trace-log-content">{{ formatErrorStack(item) }}</pre>
        </div>

        <h4>Access日志</h4>
        <el-table :data="traceDetail.accessLogs" size="small" style="margin-bottom: 20px">
          <el-table-column label="时间" prop="time" width="210"/>
          <el-table-column label="耗时(ms)" width="100">
            <template #default="{ row }">
              <span v-if="hasElapsedMs(row.handlerMs)" :class="elapsedMsClass(row.handlerMs)" class="trace-elapsed">
                {{ formatElapsedMs(row.handlerMs) }}
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="状态码" prop="statusCode" width="80"/>
          <el-table-column label="URL" min-width="180" prop="url" show-overflow-tooltip/>
          <el-table-column label="IP" prop="ip" width="140"/>
        </el-table>

        <h4>详情日志</h4>
        <div v-for="(item, index) in traceDetail.detailLogs" :key="index" class="trace-log-line">
          <div class="trace-log-meta">
            <div class="trace-log-meta-main">
              <span>{{ item.time }}</span>
              <span v-if="hasElapsedMs(item.elapsedMs)" :class="elapsedMsClass(item.elapsedMs)" class="trace-elapsed">
                {{ formatElapsedMs(item.elapsedMs) }}
              </span>
              <span v-if="item.authId" class="trace-auth-id">AuthId: {{ item.authId }}</span>
              <span>[{{ item.level }}]</span>
            </div>
            <el-button link size="small" type="primary" @click="copyTraceLog(item.raw || item.message)">
              <el-icon>
                <CopyDocument/>
              </el-icon>
              复制
            </el-button>
          </div>
          <pre class="trace-log-content">{{ item.raw }}</pre>
        </div>
        <el-empty v-if="!traceLoading && traceDetail.detailLogs.length === 0 && traceDetail.accessLogs.length === 0 && traceDetail.errorLogs.length === 0" description="未找到日志"/>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import {nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {logQueryApi} from '@/api'
import AccessTrendChart from './components/access-trend-chart.vue'
import {formatErrorStack, formatErrorSummary} from '@/utils/logParsers'
import type {AccessLogItem, AccessTrendData, DetailLogItem, ErrorLogItem, LogQueryJobResult, TopStatItem, TraceLogDetail} from '@/types/api'

const activeTab = ref('stats')
const queryStatusTip = ref('查询中...')
const serverTimeLoading = ref(false)
const serverTimeDisplay = ref('-')
let serverTimeBaseMs = 0
let serverTimeSyncClientMs = 0
let serverTimeTickTimer: ReturnType<typeof setInterval> | null = null

const formatLocalDate = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const formatServerDateTime = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  const millisecond = String(date.getMilliseconds()).padStart(3, '0')
  return `${year}-${month}-${day} ${hour}:${minute}:${second}.${millisecond}`
}

const parseServerTime = (value?: string) => {
  if (!value) {
    return null
  }
  const match = value.trim().match(/^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}):(\d{2}):(\d{2})(?:\.(\d{1,3}))?$/)
  if (!match) {
    return null
  }
  const [, year, month, day, hour, minute, second, millisecond = '0'] = match
  return new Date(
      Number(year),
      Number(month) - 1,
      Number(day),
      Number(hour),
      Number(minute),
      Number(second),
      Number(millisecond.padEnd(3, '0')),
  )
}

const updateServerTimeDisplay = () => {
  if (!serverTimeBaseMs) {
    serverTimeDisplay.value = '-'
    return
  }
  serverTimeDisplay.value = formatServerDateTime(new Date(serverTimeBaseMs + Date.now() - serverTimeSyncClientMs))
}

const syncServerTime = async () => {
  serverTimeLoading.value = true
  try {
    const data = await logQueryApi.getLogPaths()
    const parsed = parseServerTime(data.serverTime)
    if (!parsed) {
      throw new Error('invalid server time')
    }
    serverTimeBaseMs = parsed.getTime()
    serverTimeSyncClientMs = Date.now()
    updateServerTimeDisplay()
  } catch (error) {
    console.error('获取服务器时间失败:', error)
    ElMessage.error('获取服务器时间失败')
  } finally {
    serverTimeLoading.value = false
  }
}

const startServerTimeClock = () => {
  serverTimeTickTimer = setInterval(updateServerTimeDisplay, 1000)
}

const stopServerTimeClock = () => {
  if (serverTimeTickTimer) {
    clearInterval(serverTimeTickTimer)
    serverTimeTickTimer = null
  }
}

type DateRangeForm = {
  startDate: string
  endDate: string
}

const normalizeDateRange = (startDate: string, endDate: string): string[] | null => {
  if (!startDate || !endDate) {
    return null
  }
  if (startDate > endDate) {
    return [endDate, startDate]
  }
  return [startDate, endDate]
}

const extractLogDate = (time?: string) => {
  if (!time || time.length < 10) {
    return ''
  }
  return time.slice(0, 10)
}

const expandRangeToIncludeDate = (range: string[], logDate: string) => {
  const normalized = normalizeDateRange(range[0], range[1])
  if (!normalized) {
    return null
  }
  if (!logDate) {
    return normalized
  }
  let [start, end] = normalized
  if (logDate < start) {
    start = logDate
  }
  if (logDate > end) {
    end = logDate
  }
  return normalizeDateRange(start, end)
}

const getActiveTabDateRange = (): string[] => {
  let form: DateRangeForm
  if (activeTab.value === 'access') {
    form = accessForm
  } else if (activeTab.value === 'error') {
    form = errorForm
  } else if (activeTab.value === 'stats') {
    form = statsForm
  } else {
    form = detailForm
  }
  return [form.startDate, form.endDate]
}

/** 默认8天; 结束日期+2天,兼容不同时区日志时间 */
const buildDefaultDateRange = (baseDate = new Date()) => {
  const end = new Date(baseDate)
  end.setDate(end.getDate() + 2)
  const start = new Date(end)
  start.setDate(end.getDate() - 7)
  return [formatLocalDate(start), formatLocalDate(end)]
}

const defaultDateRange = () => buildDefaultDateRange(serverTimeBaseMs ? new Date(serverTimeBaseMs) : new Date())

const createDefaultDateRangeForm = (): DateRangeForm => {
  const [startDate, endDate] = defaultDateRange()
  return {startDate, endDate}
}

const resetFormDateRange = (form: DateRangeForm) => {
  const [startDate, endDate] = defaultDateRange()
  form.startDate = startDate
  form.endDate = endDate
}

const applyDefaultDateRanges = () => {
  resetFormDateRange(statsForm)
  resetFormDateRange(accessForm)
  resetFormDateRange(errorForm)
  resetFormDateRange(detailForm)
  const [startDate, endDate] = defaultDateRange()
  traceStartDate.value = startDate
  traceEndDate.value = endDate
}

const detailForm = reactive({
  ...createDefaultDateRangeForm(),
  traceId: '',
  reqId: '',
  authId: '',
  url: '',
  keyword: '',
})
const detailLoading = ref(false)
const detailTableData = ref<DetailLogItem[]>([])
const detailTotal = ref(0)
const detailPageIndex = ref(1)
const detailPageSize = ref(50)

const accessForm = reactive({
  ...createDefaultDateRangeForm(),
  traceId: '',
  url: '',
  ip: '',
  statusCode: undefined as number | undefined,
  minHandlerMs: undefined as number | undefined,
  maxHandlerMs: undefined as number | undefined,
  intervalMinutes: 0,
})
const accessLoading = ref(false)
const accessTrendLoading = ref(false)
const accessTrendData = ref<AccessTrendData | null>(null)
const accessTrendChartRef = ref<InstanceType<typeof AccessTrendChart>>()
const accessTableData = ref<AccessLogItem[]>([])
const accessTotal = ref(0)
const accessPageIndex = ref(1)
const accessPageSize = ref(50)

const errorForm = reactive({
  ...createDefaultDateRangeForm(),
  traceId: '',
  url: '',
  ip: '',
  statusCode: undefined as number | undefined,
  keyword: '',
})
const errorLoading = ref(false)
const errorTableData = ref<ErrorLogItem[]>([])
const errorTotal = ref(0)
const errorPageIndex = ref(1)
const errorPageSize = ref(50)

const statsForm = reactive({
  ...createDefaultDateRangeForm(),
  topN: 20,
})
const statsLoading = ref(false)
const urlTopData = ref<TopStatItem[]>([])
const ipTopData = ref<TopStatItem[]>([])

const traceDrawerVisible = ref(false)
const traceLoading = ref(false)
const defaultTraceDates = createDefaultDateRangeForm()
const traceStartDate = ref(defaultTraceDates.startDate)
const traceEndDate = ref(defaultTraceDates.endDate)
const traceDetail = reactive<TraceLogDetail>({
  traceId: '',
  startDate: '',
  endDate: '',
  detailLogs: [],
  accessLogs: [],
  errorLogs: [],
})

const ensureFormDateRange = (form: DateRangeForm): string[] | null => {
  if (!form.startDate || !form.endDate) {
    ElMessage.warning('请选择开始和结束日期')
    return null
  }
  const normalized = normalizeDateRange(form.startDate, form.endDate)
  if (!normalized) {
    return null
  }
  form.startDate = normalized[0]
  form.endDate = normalized[1]
  return normalized
}

const handleLogQueryStatus = (job: LogQueryJobResult) => {
  if (job.status === 'pending') {
    queryStatusTip.value = job.queuePosition > 1
        ? `排队中，前面还有 ${job.queuePosition - 1} 个查询...`
        : '排队中，即将开始查询...'
    return
  }
  if (job.status === 'running') {
    queryStatusTip.value = '正在查询日志，请稍候...'
  }
}

const fetchDetailLogs = async () => {
  const range = ensureFormDateRange(detailForm)
  if (!range) {
    return
  }
  detailLoading.value = true
  try {
    const response = await logQueryApi.queryDetailLogs({
      pageIndex: detailPageIndex.value,
      pageSize: detailPageSize.value,
      startDate: range[0],
      endDate: range[1],
      traceId: detailForm.traceId.trim(),
      reqId: detailForm.reqId.trim(),
      authId: detailForm.authId.trim(),
      url: detailForm.url.trim(),
      keyword: detailForm.keyword.trim(),
    }, handleLogQueryStatus)
    detailTableData.value = response.data || []
    detailTotal.value = response.total || 0
  } catch (error) {
    console.error('查询详情日志失败:', error)
    ElMessage.error('查询详情日志失败')
  } finally {
    detailLoading.value = false
  }
}

const buildAccessQueryParams = (range: string[]) => ({
  startDate: range[0],
  endDate: range[1],
  traceId: accessForm.traceId.trim(),
  url: accessForm.url.trim(),
  ip: accessForm.ip.trim(),
  statusCode: accessForm.statusCode || undefined,
  minHandlerMs: accessForm.minHandlerMs,
  maxHandlerMs: accessForm.maxHandlerMs,
  intervalMinutes: accessForm.intervalMinutes || undefined,
})

const fetchAccessTrend = async () => {
  const range = ensureFormDateRange(accessForm)
  if (!range) {
    return
  }
  accessTrendLoading.value = true
  try {
    accessTrendData.value = await logQueryApi.getAccessTrend(buildAccessQueryParams(range), handleLogQueryStatus)
    await nextTick()
    accessTrendChartRef.value?.resize()
  } catch (error) {
    console.error('查询Access趋势失败:', error)
    ElMessage.error('查询Access趋势失败')
  } finally {
    accessTrendLoading.value = false
  }
}

const fetchAccessLogs = async () => {
  const range = ensureFormDateRange(accessForm)
  if (!range) {
    return
  }
  accessLoading.value = true
  try {
    const response = await logQueryApi.queryAccessLogs({
      pageIndex: accessPageIndex.value,
      pageSize: accessPageSize.value,
      ...buildAccessQueryParams(range),
    }, handleLogQueryStatus)
    accessTableData.value = response.data || []
    accessTotal.value = response.total || 0
  } catch (error) {
    console.error('查询Access日志失败:', error)
    ElMessage.error('查询Access日志失败')
  } finally {
    accessLoading.value = false
  }
}

const fetchErrorLogs = async () => {
  const range = ensureFormDateRange(errorForm)
  if (!range) {
    return
  }
  errorLoading.value = true
  try {
    const response = await logQueryApi.queryErrorLogs({
      pageIndex: errorPageIndex.value,
      pageSize: errorPageSize.value,
      startDate: range[0],
      endDate: range[1],
      traceId: errorForm.traceId.trim(),
      url: errorForm.url.trim(),
      ip: errorForm.ip.trim(),
      statusCode: errorForm.statusCode || undefined,
      keyword: errorForm.keyword.trim(),
    }, handleLogQueryStatus)
    errorTableData.value = response.data || []
    errorTotal.value = response.total || 0
  } catch (error) {
    console.error('查询Error日志失败:', error)
    ElMessage.error('查询Error日志失败')
  } finally {
    errorLoading.value = false
  }
}

const fetchAccessStats = async () => {
  const range = ensureFormDateRange(statsForm)
  if (!range) {
    return
  }
  statsLoading.value = true
  try {
    const response = await logQueryApi.getAccessStats({
      startDate: range[0],
      endDate: range[1],
      topN: statsForm.topN,
    }, handleLogQueryStatus)
    urlTopData.value = response.urlTop || []
    ipTopData.value = response.ipTop || []
  } catch (error) {
    console.error('查询访问统计失败:', error)
    ElMessage.error('查询访问统计失败')
  } finally {
    statsLoading.value = false
  }
}

const resolveInitialTraceDateRange = (logTime?: string, tabRange?: string[]) => {
  const logDate = extractLogDate(logTime)
  if (logDate) {
    if (tabRange?.length === 2) {
      const expanded = expandRangeToIncludeDate(tabRange, logDate)
      if (expanded) {
        return expanded
      }
    }
    return [logDate, logDate]
  }
  if (tabRange?.length === 2) {
    const normalized = normalizeDateRange(tabRange[0], tabRange[1])
    if (normalized) {
      return normalized
    }
  }
  return defaultDateRange()
}

const fetchTraceDetail = async () => {
  if (!traceDetail.traceId) {
    return
  }
  const range = ensureFormDateRange({
    startDate: traceStartDate.value,
    endDate: traceEndDate.value,
  })
  if (!range) {
    return
  }
  traceStartDate.value = range[0]
  traceEndDate.value = range[1]
  traceLoading.value = true
  traceDetail.startDate = range[0]
  traceDetail.endDate = range[1]
  traceDetail.detailLogs = []
  traceDetail.accessLogs = []
  traceDetail.errorLogs = []
  try {
    const data = await logQueryApi.getTraceLogs(traceDetail.traceId, range[0], range[1], handleLogQueryStatus)
    traceDetail.traceId = data.traceId || traceDetail.traceId
    traceDetail.startDate = data.startDate || range[0]
    traceDetail.endDate = data.endDate || range[1]
    traceDetail.detailLogs = data.detailLogs || []
    traceDetail.accessLogs = data.accessLogs || []
    traceDetail.errorLogs = data.errorLogs || []
  } catch (error) {
    console.error('查询Trace日志失败:', error)
    ElMessage.error('查询Trace日志失败')
  } finally {
    traceLoading.value = false
  }
}

const openTraceDetail = async (traceId: string, logTime?: string) => {
  if (!traceId) {
    return
  }
  traceDetail.traceId = traceId
  const [startDate, endDate] = resolveInitialTraceDateRange(logTime, getActiveTabDateRange())
  traceStartDate.value = startDate
  traceEndDate.value = endDate
  traceDrawerVisible.value = true
  await fetchTraceDetail()
}

const handleDetailSearch = () => {
  detailPageIndex.value = 1
  fetchDetailLogs()
}

const handleAccessSearch = () => {
  accessPageIndex.value = 1
  fetchAccessTrend()
  fetchAccessLogs()
}

const handleErrorSearch = () => {
  errorPageIndex.value = 1
  fetchErrorLogs()
}

const resetDetailForm = () => {
  resetFormDateRange(detailForm)
  detailForm.traceId = ''
  detailForm.reqId = ''
  detailForm.authId = ''
  detailForm.url = ''
  detailForm.keyword = ''
  handleDetailSearch()
}

const resetAccessForm = () => {
  resetFormDateRange(accessForm)
  accessForm.traceId = ''
  accessForm.url = ''
  accessForm.ip = ''
  accessForm.statusCode = undefined
  accessForm.minHandlerMs = undefined
  accessForm.maxHandlerMs = undefined
  accessForm.intervalMinutes = 0
  handleAccessSearch()
}

const resetErrorForm = () => {
  resetFormDateRange(errorForm)
  errorForm.traceId = ''
  errorForm.url = ''
  errorForm.ip = ''
  errorForm.statusCode = undefined
  errorForm.keyword = ''
  handleErrorSearch()
}

const formatHandlerMs = (value: number) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(1)
}

const TRACE_ELAPSED_THRESHOLD_MS = 200

const hasElapsedMs = (value?: number | null) => value !== null && value !== undefined

const formatElapsedMs = (value?: number | null) => {
  if (!hasElapsedMs(value)) {
    return ''
  }
  return `${Number(value).toFixed(1)}ms`
}

const elapsedMsClass = (value?: number | null) => {
  if (!hasElapsedMs(value)) {
    return ''
  }
  return Number(value) > TRACE_ELAPSED_THRESHOLD_MS ? 'elapsed-slow' : 'elapsed-fast'
}

const copyTraceLog = async (content?: string) => {
  if (!content) {
    ElMessage.warning('无可复制内容')
    return
  }
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success('已复制')
  } catch (error) {
    console.error('复制日志失败:', error)
    ElMessage.error('复制失败')
  }
}

const resolveTraceAuthId = () => {
  for (const item of traceDetail.detailLogs) {
    if (item.authId) {
      return item.authId
    }
  }
  for (const item of traceDetail.errorLogs) {
    if (item.authId) {
      return item.authId
    }
  }
  return ''
}

onMounted(async () => {
  await syncServerTime()
  applyDefaultDateRanges()
  startServerTimeClock()
  await fetchAccessStats()
})

onUnmounted(() => {
  stopServerTimeClock()
})

watch(activeTab, async (tab) => {
  if (tab === 'access' && !accessTrendData.value) {
    await fetchAccessTrend()
  }
  if (tab === 'error' && errorTableData.value.length === 0) {
    await fetchErrorLogs()
  }
  if (tab === 'detail' && detailTableData.value.length === 0) {
    await fetchDetailLogs()
  }
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
  gap: 12px;
}

.path-tip,
.server-time-tip,
.detail-query-tip {
  color: #909399;
  font-size: 12px;
  word-break: break-all;
}

.detail-query-tip {
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 16px;
}

.access-trend-panel {
  margin-bottom: 20px;
}

.trend-summary {
  margin-top: 8px;
  color: #606266;
  font-size: 13px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.range-sep {
  margin: 0 8px;
  color: #909399;
}

.trace-drawer h4 {
  margin: 0 0 12px;
}

.trace-search-form {
  margin-bottom: 12px;
}

.trace-meta {
  margin: 0 0 16px;
  color: #606266;
}

.trace-log-line {
  margin-bottom: 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  overflow: hidden;
}

.trace-log-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  background: #f5f7fa;
  color: #606266;
  font-size: 12px;
}

.trace-log-meta-main {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
}

.trace-elapsed {
  font-weight: 600;
}

.trace-auth-id {
  color: #409eff;
  font-weight: 500;
}

.elapsed-fast {
  color: #67c23a;
}

.elapsed-slow {
  color: #f56c6c;
}

.trace-log-content {
  margin: 0;
  padding: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  line-height: 1.6;
}
</style>
