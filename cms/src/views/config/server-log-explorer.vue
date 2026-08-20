<template>
  <div class="page-container">
    <el-card v-loading="serverTimeLoading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ServerLogExplorer') }}</span>
          <span class="server-time-tip">
            {{ t('pages.serverLogExplorer.serverTime') }}: {{ serverTimeDisplay }}
            <template v-if="exportPathTip"> · {{ t('pages.serverLogExplorer.exportPath') }}: {{ exportPathTip }}</template>
            · {{ t('pages.serverLogExplorer.queryMaxWait') }}
          </span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="t('pages.serverLogExplorer.tabAccessStats')" name="stats">
          <el-form :model="statsForm" class="search-form" inline label-width="90px">
            <el-form-item :label="t('common.startTime')">
              <el-date-picker
                  v-model="statsForm.startDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.startTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('common.endTime')">
              <el-date-picker
                  v-model="statsForm.endDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.endTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item label="TopN">
              <el-input-number v-model="statsForm.topN" :max="100" :min="1" controls-position="right" style="width: 120px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchAccessStats">{{ t('common.query') }}</el-button>
            </el-form-item>
          </el-form>

          <el-row v-loading="statsLoading" :element-loading-text="queryStatusTip" :gutter="20">
            <el-col :span="12">
              <el-card shadow="never">
                <template #header>{{ t('pages.serverLogExplorer.urlTopN') }}</template>
                <el-table :data="urlTopData" size="small">
                  <el-table-column label="#" type="index" width="50"/>
                  <el-table-column label="URL" min-width="220" prop="key" show-overflow-tooltip/>
                  <el-table-column :label="t('pages.serverLogExplorer.visitCount')" prop="count" width="100"/>
                </el-table>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never">
                <template #header>{{ t('pages.serverLogExplorer.ipTopN') }}</template>
                <el-table :data="ipTopData" size="small">
                  <el-table-column label="#" type="index" width="50"/>
                  <el-table-column label="IP" min-width="180" prop="key" show-overflow-tooltip/>
                  <el-table-column :label="t('pages.serverLogExplorer.visitCount')" prop="count" width="100"/>
                </el-table>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.serverLogExplorer.tabAccessLog')" name="access">
          <el-form :model="accessForm" class="search-form" inline label-width="100px">
            <el-form-item :label="t('common.startTime')">
              <el-date-picker
                  v-model="accessForm.startDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.startTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('common.endTime')">
              <el-date-picker
                  v-model="accessForm.endDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.endTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.traceId')">
              <el-input v-model="accessForm.traceId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 280px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.authId')">
              <el-input v-model="accessForm.authId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="accessForm.url" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 220px"/>
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="accessForm.ip" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.statusCode')">
              <el-input-number v-model="accessForm.statusCode" :controls="false" :min="0" :placeholder="t('common.all')" style="width: 120px"/>
            </el-form-item>
            <el-form-item>
              <template #label>
                <span class="clickable-label" :title="t('pages.serverLogExplorer.minHandlerMsTip')" @click="setAccessMinHandlerMsDefault">{{ t('pages.serverLogExplorer.minHandlerMs') }}</span>
              </template>
              <el-input-number v-model="accessForm.minHandlerMs" :controls="false" :min="0" :placeholder="t('common.all')" style="width: 120px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.intervalGranularity')">
              <el-select v-model="accessForm.intervalMinutes" :placeholder="t('pages.serverLogExplorer.intervalAuto')" style="width: 120px">
                <el-option :value="0" :label="t('pages.serverLogExplorer.intervalAuto')"/>
                <el-option :value="1" :label="t('pages.serverLogExplorer.interval1Min')"/>
                <el-option :value="5" :label="t('pages.serverLogExplorer.interval5Min')"/>
                <el-option :value="15" :label="t('pages.serverLogExplorer.interval15Min')"/>
                <el-option :value="60" :label="t('pages.serverLogExplorer.interval1Hour')"/>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleAccessSearch">{{ t('common.query') }}</el-button>
              <el-button @click="resetAccessForm">{{ t('common.reset') }}</el-button>
            </el-form-item>
          </el-form>

          <div v-loading="accessTrendLoading" :element-loading-text="queryStatusTip" class="access-trend-panel">
            <AccessTrendChart ref="accessTrendChartRef" :data="accessTrendData"/>
            <div v-if="accessTrendData?.peakTime" class="trend-summary">
              {{ t('pages.serverLogExplorer.trendSummary', {
                total: accessTrendData.totalCount,
                peak: accessTrendData.peakCount,
                interval: accessTrendData.intervalMinutes,
                peakTime: accessTrendData.peakTime,
              }) }}
            </div>
          </div>

          <el-table v-loading="accessLoading" :data="accessTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column :label="t('pages.serverLogExplorer.time')" prop="time" width="210"/>
            <el-table-column :label="t('pages.serverLogExplorer.traceId')" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.statusCode')" width="90">
              <template #default="{ row }">
                <span :class="{ 'log-alert': isAbnormalStatusCode(row.statusCode) }">{{ row.statusCode }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.method')" prop="method" width="90"/>
            <el-table-column label="URL" min-width="220" prop="url" show-overflow-tooltip/>
            <el-table-column :label="t('pages.serverLogExplorer.authId')" prop="authId" width="180"/>
            <el-table-column :label="t('pages.serverLogExplorer.handlerMs')" width="100">
              <template #default="{ row }">
                <span :class="{ 'log-alert': isSlowHandlerMs(row.handlerMs) }">{{ formatHandlerMs(row.handlerMs) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="IP" min-width="160" prop="ip" show-overflow-tooltip/>
            <el-table-column :label="t('pages.serverLogExplorer.userAgent')" min-width="180" prop="userAgent" show-overflow-tooltip/>
          </el-table>

          <div class="pagination-wrap">
            <el-pagination
                v-model:current-page="accessPageIndex"
                v-model:page-size="accessPageSize"
                :page-sizes="[20, 50, 100, 200]"
                :total="accessTotal"
                background
                layout="sizes, prev, pager, next"
                @current-change="fetchAccessLogs"
                @size-change="handleAccessSearch"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.serverLogExplorer.tabErrorLog')" name="error">
          <el-form :model="errorForm" class="search-form" inline label-width="100px">
            <el-form-item :label="t('common.startTime')">
              <el-date-picker
                  v-model="errorForm.startDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.startTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('common.endTime')">
              <el-date-picker
                  v-model="errorForm.endDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.endTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.traceId')">
              <el-input v-model="errorForm.traceId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 280px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="errorForm.url" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 220px"/>
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="errorForm.ip" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.statusCode')">
              <el-input-number v-model="errorForm.statusCode" :controls="false" :min="0" :placeholder="t('common.all')" style="width: 120px"/>
            </el-form-item>
            <el-form-item :label="t('common.keyword')">
              <el-input v-model="errorForm.keyword" clearable :placeholder="t('pages.serverLogExplorer.errorKeywordPlaceholder')" style="width: 200px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleErrorSearch">{{ t('common.query') }}</el-button>
              <el-button @click="resetErrorForm">{{ t('common.reset') }}</el-button>
            </el-form-item>
          </el-form>

          <el-table v-loading="errorLoading" :data="errorTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column :label="t('pages.serverLogExplorer.time')" prop="time" width="210"/>
            <el-table-column :label="t('pages.serverLogExplorer.traceId')" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.statusCode')" prop="statusCode" width="90"/>
            <el-table-column :label="t('pages.serverLogExplorer.method')" prop="method" width="90"/>
            <el-table-column label="URL" min-width="200" prop="url" show-overflow-tooltip/>
            <el-table-column label="IP" min-width="140" prop="ip" show-overflow-tooltip/>
            <el-table-column :label="t('pages.serverLogExplorer.errorCode')" prop="errorCode" width="90"/>
            <el-table-column :label="t('pages.serverLogExplorer.errorMessage')" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">{{ formatErrorSummary(row) }}</template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.detail')" min-width="180" prop="detail" show-overflow-tooltip/>
            <el-table-column :label="t('pages.serverLogExplorer.stack')" min-width="220" show-overflow-tooltip>
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
                layout="sizes, prev, pager, next"
                @current-change="fetchErrorLogs"
                @size-change="handleErrorSearch"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.serverLogExplorer.tabDetailLog')" lazy name="detail">
          <div class="detail-query-tip">{{ t('pages.serverLogExplorer.detailQueryTip') }}</div>
          <el-form :model="detailForm" class="search-form" inline label-width="90px">
            <el-form-item :label="t('common.startTime')">
              <el-date-picker
                  v-model="detailForm.startDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.startTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('common.endTime')">
              <el-date-picker
                  v-model="detailForm.endDate"
                  clearable
                  format="YYYY-MM-DD HH:mm:ss"
                  :placeholder="t('common.endTime')"
                  style="width: 190px"
                  teleported
                  type="datetime"
                  value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.traceId')">
              <el-input v-model="detailForm.traceId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 280px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.reqId')">
              <el-input v-model="detailForm.reqId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item :label="t('pages.serverLogExplorer.authId')">
              <el-input v-model="detailForm.authId" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPushPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item label="URL">
              <el-input v-model="detailForm.url" clearable :placeholder="t('pages.serverLogExplorer.fuzzyMatchPlaceholder')" style="width: 220px"/>
            </el-form-item>
            <el-form-item :label="t('common.keyword')">
              <el-input v-model="detailForm.keyword" clearable :placeholder="t('pages.serverLogExplorer.detailKeywordPlaceholder')" style="width: 180px"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleDetailSearch">{{ t('common.query') }}</el-button>
              <el-button @click="resetDetailForm">{{ t('common.reset') }}</el-button>
            </el-form-item>
          </el-form>

          <el-table v-loading="detailLoading" :data="detailTableData" :element-loading-text="queryStatusTip" style="width: 100%">
            <el-table-column :label="t('pages.serverLogExplorer.time')" prop="time" width="210"/>
            <el-table-column :label="t('pages.serverLogExplorer.level')" prop="level" width="80"/>
            <el-table-column :label="t('pages.serverLogExplorer.traceId')" min-width="280">
              <template #default="{ row }">
                <el-link type="primary" @click="openTraceDetail(row.traceId, row.time)">{{ row.traceId }}</el-link>
              </template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.reqId')" prop="reqId" width="160"/>
            <el-table-column :label="t('pages.serverLogExplorer.authId')" prop="authId" width="180"/>
            <el-table-column label="URL" min-width="180" prop="url" show-overflow-tooltip/>
            <el-table-column :label="t('pages.serverLogExplorer.handlerMs')" width="100">
              <template #default="{ row }">{{ formatHandlerMs(row.elapsedMs) }}</template>
            </el-table-column>
            <el-table-column :label="t('pages.serverLogExplorer.content')" min-width="320">
              <template #default="{ row }">
                <div v-if="row.syndbFlush" class="syndb-flush-cell">
                  <SyndbFlushLogView :flush="row.syndbFlush" compact/>
                  <el-button link size="small" type="primary" @click="openSyndbFlushDialog(row.syndbFlush)">{{ t('pages.serverLogExplorer.detailBtn') }}</el-button>
                </div>
                <span v-else class="log-message-text">{{ row.message }}</span>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrap">
            <el-pagination
                v-model:current-page="detailPageIndex"
                v-model:page-size="detailPageSize"
                :page-sizes="[20, 50, 100, 200]"
                :total="detailTotal"
                background
                layout="sizes, prev, pager, next"
                @current-change="fetchDetailLogs"
                @size-change="handleDetailSearch"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-drawer v-model="traceDrawerVisible" :title="`TraceId: ${traceDetail.traceId}`" size="60%">
      <div class="trace-drawer">
        <p v-if="traceStartDate && traceEndDate" class="trace-meta">
          {{ t('pages.serverLogExplorer.traceQueryRange', {start: traceStartDate, end: traceEndDate}) }}
        </p>
        <p v-if="traceAnchorTime" class="trace-meta">{{ t('pages.serverLogExplorer.traceAnchorLog', {time: traceAnchorTime}) }}</p>
        <p v-if="resolveTraceAuthId()" class="trace-meta">{{ t('pages.serverLogExplorer.traceAuthId', {authId: resolveTraceAuthId()}) }}</p>

        <div v-loading="traceLoading" :element-loading-text="queryStatusTip">
        <h4>{{ t('pages.serverLogExplorer.errorLogSection') }}</h4>
        <div v-for="(item, index) in traceDetail.errorLogs" :key="'error-' + index" class="trace-log-line">
          <div class="trace-log-meta">
            <div class="trace-log-meta-main">
              <span>{{ item.time }}</span>
              <span v-if="hasElapsedMs(item.handlerMs)" :class="elapsedMsClass(item.handlerMs)" class="trace-elapsed">
                {{ formatElapsedMs(item.handlerMs) }}
              </span>
              <span v-if="item.authId" class="trace-auth-id">{{ t('pages.serverLogExplorer.traceAuthId', {authId: item.authId}) }}</span>
              <span>[{{ item.level }}] {{ formatErrorSummary(item) }}</span>
            </div>
          </div>
          <pre v-if="formatErrorStack(item)" class="trace-log-content">{{ formatErrorStack(item) }}</pre>
        </div>

        <h4>{{ t('pages.serverLogExplorer.accessLogSection') }}</h4>
        <el-table :data="traceDetail.accessLogs" size="small" style="margin-bottom: 20px">
          <el-table-column :label="t('pages.serverLogExplorer.time')" prop="time" width="210"/>
          <el-table-column :label="t('pages.serverLogExplorer.handlerMs')" width="100">
            <template #default="{ row }">
              <span v-if="hasElapsedMs(row.handlerMs)" :class="elapsedMsClass(row.handlerMs)" class="trace-elapsed">
                {{ formatElapsedMs(row.handlerMs) }}
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.serverLogExplorer.statusCode')" prop="statusCode" width="80"/>
          <el-table-column label="URL" min-width="180" prop="url" show-overflow-tooltip/>
          <el-table-column :label="t('pages.serverLogExplorer.authId')" prop="authId" width="140"/>
          <el-table-column label="IP" prop="ip" width="140"/>
        </el-table>

        <h4>{{ t('pages.serverLogExplorer.detailLogSection') }}</h4>
        <div v-for="(item, index) in traceDetail.detailLogs" :key="index" class="trace-log-line">
            <div class="trace-log-meta">
            <div class="trace-log-meta-main">
              <span>{{ item.time }}</span>
              <span v-if="hasElapsedMs(item.elapsedMs)" :class="elapsedMsClass(item.elapsedMs)" class="trace-elapsed">
                {{ formatElapsedMs(item.elapsedMs) }}
              </span>
              <span v-if="item.authId" class="trace-auth-id">{{ t('pages.serverLogExplorer.traceAuthId', {authId: item.authId}) }}</span>
              <span>[{{ item.level }}]</span>
            </div>
            <div class="trace-log-actions">
              <el-button
                  v-for="field in listLogJsonFields(item.raw || item.message)"
                  :key="field.key"
                  link
                  size="small"
                  type="primary"
                  @click="openLogJsonDialog(item.raw || item.message, field.key, translateLogFieldLabel(field.key, field.label))"
              >
                {{ translateLogFieldLabel(field.key, field.label) }}
              </el-button>
              <el-button link size="small" type="primary" @click="copyTraceLog(item.raw || item.message)">
                <el-icon>
                  <CopyDocument/>
                </el-icon>
                {{ t('common.copy') }}
              </el-button>
            </div>
          </div>
          <pre v-if="!item.syndbFlush" class="trace-log-content">{{ item.raw }}</pre>
          <div v-else class="trace-log-content syndb-flush-trace-body">
            <SyndbFlushLogView :flush="item.syndbFlush"/>
          </div>
        </div>
        <el-empty v-if="!traceLoading && traceDetail.detailLogs.length === 0 && traceDetail.accessLogs.length === 0 && traceDetail.errorLogs.length === 0" :description="t('pages.serverLogExplorer.noLogsFound')"/>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="syndbFlushDialogVisible" destroy-on-close :title="t('pages.serverLogExplorer.syndbFlushDialogTitle')" width="78%">
      <SyndbFlushLogView v-if="syndbFlushDialogData" :flush="syndbFlushDialogData"/>
      <template #footer>
        <el-button type="primary" @click="syndbFlushDialogVisible = false">{{ t('common.close') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="jsonDialogVisible" :title="jsonDialogTitle" destroy-on-close width="72%">
      <div v-if="jsonDialogParseFailed" class="json-dialog-tip">{{ t('pages.serverLogExplorer.jsonParseFailed') }}</div>
      <pre class="json-dialog-content">{{ jsonDialogContent }}</pre>
      <template #footer>
        <el-button @click="copyJsonDialogContent">{{ t('common.copy') }}</el-button>
        <el-button type="primary" @click="jsonDialogVisible = false">{{ t('common.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {logQueryApi} from '@/api'
import AccessTrendChart from './components/access-trend-chart.vue'
import SyndbFlushLogView from './components/syndb-flush-log.vue'
import {formatErrorStack, formatErrorSummary, extractAndFormatLogJsonField, filterTraceLogsToRound, listLogJsonFields, parseLogTimeMs} from '@/utils/logParsers'
import type {LogJsonFieldKey} from '@/utils/logParsers'
import type {AccessLogItem, AccessTrendData, DetailLogItem, ErrorLogItem, LogQueryJobResult, SyndbFlushLog, TopStatItem, TraceLogDetail} from '@/types/api'

const {t} = useI18n()
const activeTab = ref('stats')
const queryStatusTip = ref('')
const serverTimeLoading = ref(false)
const serverTimeDisplay = ref('-')
const exportPathTip = ref('')
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
    exportPathTip.value = data.fileExportUrlPrefix || ''
  } catch (error) {
    console.error('获取服务器时间失败:', error)
    ElMessage.error(t('pages.serverLogExplorer.fetchServerTimeFailed'))
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

const formatLocalDateTime = (date: Date) => {
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  return `${formatLocalDate(date)} ${hour}:${minute}:${second}`
}

const buildDefaultLogQueryDateRange = (baseDate = new Date()) => {
  const end = new Date(baseDate)
  end.setDate(end.getDate() + 2)
  end.setHours(23, 59, 59, 0)
  const start = new Date(end)
  start.setDate(end.getDate() - 7)
  start.setHours(0, 0, 0, 0)
  return [formatLocalDateTime(start), formatLocalDateTime(end)]
}

const defaultLogQueryDateRange = () => buildDefaultLogQueryDateRange(serverTimeBaseMs ? new Date(serverTimeBaseMs) : new Date())

const createDefaultDateRangeForm = (): DateRangeForm => {
  const [startDate, endDate] = defaultLogQueryDateRange()
  return {startDate, endDate}
}

const resetFormDateRange = (form: DateRangeForm) => {
  const [startDate, endDate] = defaultLogQueryDateRange()
  form.startDate = startDate
  form.endDate = endDate
}

const applyDefaultDateRanges = () => {
  resetFormDateRange(statsForm)
  resetFormDateRange(accessForm)
  resetFormDateRange(errorForm)
  resetFormDateRange(detailForm)
  const [startDate, endDate] = defaultLogQueryDateRange()
  traceStartDate.value = startDate.slice(0, 10)
  traceEndDate.value = endDate.slice(0, 10)
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
  authId: '',
  url: '',
  ip: '',
  statusCode: undefined as number | undefined,
  minHandlerMs: undefined as number | undefined,
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
const traceAnchorTime = ref('')
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

const jsonDialogVisible = ref(false)
const jsonDialogTitle = ref('')
const jsonDialogContent = ref('')
const jsonDialogParseFailed = ref(false)
const syndbFlushDialogVisible = ref(false)
const syndbFlushDialogData = ref<SyndbFlushLog | null>(null)

const translateLogFieldLabel = (field: LogJsonFieldKey, label: string) => {
  if (field === 'respContent') {
    return t('pages.serverLogExplorer.logFieldResponse')
  }
  return label
}

const ensureFormDateRange = (form: DateRangeForm): string[] | null => {
  if (!form.startDate || !form.endDate) {
    ElMessage.warning(t('pages.serverLogExplorer.selectDateRange'))
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
        ? t('pages.serverLogExplorer.queuePending', {count: job.queuePosition - 1})
        : t('pages.serverLogExplorer.queueStarting')
    return
  }
  if (job.status === 'running') {
    queryStatusTip.value = t('pages.serverLogExplorer.queryRunning')
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
    ElMessage.error(t('pages.serverLogExplorer.queryDetailFailed'))
  } finally {
    detailLoading.value = false
  }
}

const buildAccessQueryParams = (range: string[]) => ({
  startDate: range[0],
  endDate: range[1],
  traceId: accessForm.traceId.trim(),
  authId: accessForm.authId.trim(),
  url: accessForm.url.trim(),
  ip: accessForm.ip.trim(),
  statusCode: accessForm.statusCode || undefined,
  minHandlerMs: accessForm.minHandlerMs != null && accessForm.minHandlerMs > 0 ? accessForm.minHandlerMs : undefined,
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
    ElMessage.error(t('pages.serverLogExplorer.queryAccessTrendFailed'))
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
    ElMessage.error(t('pages.serverLogExplorer.queryAccessFailed'))
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
    const message = error instanceof Error ? error.message : t('pages.serverLogExplorer.queryErrorFailed')
    ElMessage.error(message.includes('log') || message.includes('Log') || message.includes('日志') ? message : `${t('pages.serverLogExplorer.queryErrorFailed')}: ${message}`)
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
    ElMessage.error(t('pages.serverLogExplorer.queryStatsFailed'))
  } finally {
    statsLoading.value = false
  }
}

const TRACE_QUERY_WINDOW_MS = 120_000

const resolveTraceQueryRange = (logTime?: string, tabRange?: string[]): string[] => {
  const anchorMs = parseLogTimeMs(logTime)
  if (anchorMs) {
    return [
      formatLocalDateTime(new Date(anchorMs - TRACE_QUERY_WINDOW_MS)),
      formatLocalDateTime(new Date(anchorMs + TRACE_QUERY_WINDOW_MS)),
    ]
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
    const filtered = filterTraceLogsToRound({
      traceId: data.traceId || traceDetail.traceId,
      startDate: data.startDate || range[0],
      endDate: data.endDate || range[1],
      detailLogs: data.detailLogs || [],
      accessLogs: data.accessLogs || [],
      errorLogs: data.errorLogs || [],
    }, traceAnchorTime.value)
    traceDetail.traceId = filtered.traceId
    traceDetail.startDate = filtered.startDate
    traceDetail.endDate = filtered.endDate
    traceDetail.detailLogs = filtered.detailLogs
    traceDetail.accessLogs = filtered.accessLogs
    traceDetail.errorLogs = filtered.errorLogs
  } catch (error) {
    console.error('查询Trace日志失败:', error)
    ElMessage.error(t('pages.serverLogExplorer.queryTraceFailed'))
  } finally {
    traceLoading.value = false
  }
}

const openTraceDetail = async (traceId: string, logTime?: string) => {
  if (!traceId) {
    return
  }
  traceDetail.traceId = traceId
  traceAnchorTime.value = logTime || ''
  const [startDate, endDate] = resolveTraceQueryRange(logTime, getActiveTabDateRange())
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
  accessForm.authId = ''
  accessForm.url = ''
  accessForm.ip = ''
  accessForm.statusCode = undefined
  accessForm.minHandlerMs = undefined
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

const ACCESS_SLOW_MS_THRESHOLD = 400
const TRACE_ELAPSED_THRESHOLD_MS = ACCESS_SLOW_MS_THRESHOLD

const setAccessMinHandlerMsDefault = () => {
  accessForm.minHandlerMs = ACCESS_SLOW_MS_THRESHOLD
}

const hasElapsedMs = (value?: number | null) => value !== null && value !== undefined

const isSlowHandlerMs = (value?: number | null) => hasElapsedMs(value) && Number(value) > ACCESS_SLOW_MS_THRESHOLD

const isAbnormalStatusCode = (statusCode?: number) => statusCode !== 200

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
    ElMessage.warning(t('pages.serverLogExplorer.nothingToCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success(t('pages.serverLogExplorer.copied'))
  } catch (error) {
    console.error('复制日志失败:', error)
    ElMessage.error(t('pages.serverLogExplorer.copyFailed'))
  }
}

const openSyndbFlushDialog = (flush: SyndbFlushLog) => {
  syndbFlushDialogData.value = flush
  syndbFlushDialogVisible.value = true
}

const openLogJsonDialog = (content: string, field: LogJsonFieldKey, label: string) => {
  const result = extractAndFormatLogJsonField(content, field)
  if (!result) {
    ElMessage.warning(t('pages.serverLogExplorer.contentNotFound', {label}))
    return
  }
  jsonDialogTitle.value = `${label}${t('pages.serverLogExplorer.jsonDialogSuffix')}`
  jsonDialogContent.value = result.formatted
  jsonDialogParseFailed.value = !result.parsed
  jsonDialogVisible.value = true
}

const copyJsonDialogContent = async () => {
  if (!jsonDialogContent.value) {
    ElMessage.warning(t('pages.serverLogExplorer.nothingToCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(jsonDialogContent.value)
    ElMessage.success(t('pages.serverLogExplorer.copied'))
  } catch (error) {
    console.error('复制 JSON 失败:', error)
    ElMessage.error(t('pages.serverLogExplorer.copyFailed'))
  }
}

const resolveTraceAuthId = () => {
  for (const item of traceDetail.detailLogs) {
    if (item.authId) {
      return item.authId
    }
  }
  for (const item of traceDetail.accessLogs) {
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
  queryStatusTip.value = t('pages.serverLogExplorer.querying')
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

.clickable-label {
  cursor: pointer;
  color: #409eff;
}

.clickable-label:hover {
  text-decoration: underline;
}

.log-alert {
  color: #f56c6c;
  font-weight: 600;
}

.trace-drawer h4 {
  margin: 0 0 12px;
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

.trace-log-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 4px;
  flex-shrink: 0;
}

.json-dialog-tip {
  margin-bottom: 8px;
  color: #e6a23c;
  font-size: 12px;
}

.json-dialog-content {
  margin: 0;
  max-height: 65vh;
  overflow: auto;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #1e1e1e;
  color: #dcdcdc;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  line-height: 1.6;
  font-family: Consolas, Monaco, 'Courier New', monospace;
}

.trace-log-content {
  margin: 0;
  padding: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  line-height: 1.6;
}

.syndb-flush-trace-body {
  white-space: normal;
}

.syndb-flush-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.log-message-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
