import {trackingEventApi} from '@/api/modules/tracking-event'
import type {HotLiveRoomJoinTrendRes} from '@/types/api'
import type {DashboardCsvBundleFile} from '@/utils/dashboard-export'
import {
  buildTrackingCountTrendCsvContent,
  buildTrackingExportFilename,
  pickTrackingTrendPeriodPoints,
  type TrackingTrendExportSummaryLabels,
} from '@/utils/tracking-trend-export'

type TranslateFn = (key: string) => string

type TrackingTrendPeriod = 'daily' | 'weekly' | 'monthly'

type DashboardTrackingExportSpec = {
  fileSlug: string
  eventSlug: string
  fetch: () => Promise<HotLiveRoomJoinTrendRes>
  summaryLabels: (t: TranslateFn) => TrackingTrendExportSummaryLabels
  seriesName: (t: TranslateFn) => string
}

const TRACKING_PERIODS: TrackingTrendPeriod[] = ['daily', 'weekly', 'monthly']

const normalizeTrend = (data: HotLiveRoomJoinTrendRes): HotLiveRoomJoinTrendRes => ({
  todayCount: data.todayCount ?? 0,
  weekCount: data.weekCount ?? 0,
  monthCount: data.monthCount ?? 0,
  daily: data.daily || [],
  weekly: data.weekly || [],
  monthly: data.monthly || [],
})

const trendKey = (fileSlug: string, eventSlug: string) => `${fileSlug}:${eventSlug}`

const DASHBOARD_TRACKING_EXPORT_SPECS: DashboardTrackingExportSpec[] = [
  {fileSlug: 'download', eventSlug: 'h5', fetch: () => trackingEventApi.getH5DownloadTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingDownload.todayTotal'),
    week: t('pages.trackingDownload.weekTotal'),
    month: t('pages.trackingDownload.monthTotal'),
  }), seriesName: (t) => t('pages.trackingDownload.h5SeriesName')},
  {fileSlug: 'download', eventSlug: 'official', fetch: () => trackingEventApi.getOfficialSiteDownloadTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingDownload.todayTotal'),
    week: t('pages.trackingDownload.weekTotal'),
    month: t('pages.trackingDownload.monthTotal'),
  }), seriesName: (t) => t('pages.trackingDownload.officialSeriesName')},
  {fileSlug: 'download', eventSlug: 'shareShort', fetch: () => trackingEventApi.getShareShortLinkClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingDownload.todayTotal'),
    week: t('pages.trackingDownload.weekTotal'),
    month: t('pages.trackingDownload.monthTotal'),
  }), seriesName: (t) => t('pages.trackingDownload.shareShortSeriesName')},

  {fileSlug: 'hot-live-room-join', eventSlug: 'showcaseJoin', fetch: () => trackingEventApi.getHotLiveRoomJoinTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingHotLiveRoomJoin.todayTotal'),
    week: t('pages.trackingHotLiveRoomJoin.weekTotal'),
    month: t('pages.trackingHotLiveRoomJoin.monthTotal'),
  }), seriesName: (t) => t('pages.trackingHotLiveRoomJoin.seriesName')},
  {fileSlug: 'hot-live-room-join', eventSlug: 'firstFrame', fetch: () => trackingEventApi.getLiveFirstFrameRenderTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingHotLiveRoomJoin.todayTotal'),
    week: t('pages.trackingHotLiveRoomJoin.weekTotal'),
    month: t('pages.trackingHotLiveRoomJoin.monthTotal'),
  }), seriesName: (t) => t('pages.trackingHotLiveRoomJoin.firstFrameSeriesName')},
  {fileSlug: 'hot-live-room-join', eventSlug: 'showcaseLeave', fetch: () => trackingEventApi.getHotLiveRoomLeaveTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingHotLiveRoomJoin.todayTotal'),
    week: t('pages.trackingHotLiveRoomJoin.weekTotal'),
    month: t('pages.trackingHotLiveRoomJoin.monthTotal'),
  }), seriesName: (t) => t('pages.trackingHotLiveRoomJoin.leaveSeriesName')},
  {fileSlug: 'hot-live-room-join', eventSlug: 'giftClick', fetch: () => trackingEventApi.getHotLiveRoomGiftClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingHotLiveRoomJoin.todayTotal'),
    week: t('pages.trackingHotLiveRoomJoin.weekTotal'),
    month: t('pages.trackingHotLiveRoomJoin.monthTotal'),
  }), seriesName: (t) => t('pages.trackingHotLiveRoomJoin.giftClickSeriesName')},

  {fileSlug: 'call-1v1-initiate', eventSlug: 'liveRoom', fetch: () => trackingEventApi.getCall1v1InitiateTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingCall1v1Initiate.todayTotal'),
    week: t('pages.trackingCall1v1Initiate.weekTotal'),
    month: t('pages.trackingCall1v1Initiate.monthTotal'),
  }), seriesName: (t) => t('pages.trackingCall1v1Initiate.seriesName')},
  {fileSlug: 'call-1v1-initiate', eventSlug: 'oneToOneRoom', fetch: () => trackingEventApi.getCall1v1RoomCallTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingCall1v1Initiate.todayTotal'),
    week: t('pages.trackingCall1v1Initiate.weekTotal'),
    month: t('pages.trackingCall1v1Initiate.monthTotal'),
  }), seriesName: (t) => t('pages.trackingCall1v1Initiate.oneToOneRoomSeriesName')},
  {fileSlug: 'call-1v1-initiate', eventSlug: 'connectSuccess', fetch: () => trackingEventApi.getCall1v1ConnectSuccessTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingCall1v1Initiate.todayTotal'),
    week: t('pages.trackingCall1v1Initiate.weekTotal'),
    month: t('pages.trackingCall1v1Initiate.monthTotal'),
  }), seriesName: (t) => t('pages.trackingCall1v1Initiate.connectSuccessSeriesName')},
  {fileSlug: 'call-1v1-initiate', eventSlug: 'oneToOneRoomList', fetch: () => trackingEventApi.getOneToOneRoomListTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingCall1v1Initiate.todayTotal'),
    week: t('pages.trackingCall1v1Initiate.weekTotal'),
    month: t('pages.trackingCall1v1Initiate.monthTotal'),
  }), seriesName: (t) => t('pages.trackingCall1v1Initiate.oneToOneRoomListSeriesName')},

  {fileSlug: 'mini-game-round-start', eventSlug: 'spinClick', fetch: () => trackingEventApi.getMiniGameSpinClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingMiniGameRoundStart.todayTotal'),
    week: t('pages.trackingMiniGameRoundStart.weekTotal'),
    month: t('pages.trackingMiniGameRoundStart.monthTotal'),
  }), seriesName: (t) => t('pages.trackingMiniGameRoundStart.spinClickSeriesName')},
  {fileSlug: 'mini-game-round-start', eventSlug: 'plusClick', fetch: () => trackingEventApi.getMiniGamePlusClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingMiniGameRoundStart.todayTotal'),
    week: t('pages.trackingMiniGameRoundStart.weekTotal'),
    month: t('pages.trackingMiniGameRoundStart.monthTotal'),
  }), seriesName: (t) => t('pages.trackingMiniGameRoundStart.plusClickSeriesName')},

  {fileSlug: 'short-video', eventSlug: 'playStart', fetch: () => trackingEventApi.getShortVideoPlayStartTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingShortVideo.todayTotal'),
    week: t('pages.trackingShortVideo.weekTotal'),
    month: t('pages.trackingShortVideo.monthTotal'),
  }), seriesName: (t) => t('pages.trackingShortVideo.playStartSeriesName')},
  {fileSlug: 'short-video', eventSlug: 'paywallClick', fetch: () => trackingEventApi.getShortVideoPaywallClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingShortVideo.todayTotal'),
    week: t('pages.trackingShortVideo.weekTotal'),
    month: t('pages.trackingShortVideo.monthTotal'),
  }), seriesName: (t) => t('pages.trackingShortVideo.paywallClickSeriesName')},
  {fileSlug: 'short-video', eventSlug: 'paySuccess', fetch: () => trackingEventApi.getShortVideoPaySuccessTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingShortVideo.todayTotal'),
    week: t('pages.trackingShortVideo.weekTotal'),
    month: t('pages.trackingShortVideo.monthTotal'),
  }), seriesName: (t) => t('pages.trackingShortVideo.paySuccessSeriesName')},
  {fileSlug: 'short-video', eventSlug: 'watchEnd', fetch: () => trackingEventApi.getShortVideoWatchEndTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingShortVideo.todayTotal'),
    week: t('pages.trackingShortVideo.weekTotal'),
    month: t('pages.trackingShortVideo.monthTotal'),
  }), seriesName: (t) => t('pages.trackingShortVideo.watchEndSeriesName')},

  {fileSlug: 'voice-chat-room', eventSlug: 'enter', fetch: () => trackingEventApi.getVoiceRoomEnterTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingVoiceChatRoom.todayTotal'),
    week: t('pages.trackingVoiceChatRoom.weekTotal'),
    month: t('pages.trackingVoiceChatRoom.monthTotal'),
  }), seriesName: (t) => t('pages.trackingVoiceChatRoom.enterSeriesName')},
  {fileSlug: 'voice-chat-room', eventSlug: 'micRequest', fetch: () => trackingEventApi.getVoiceRoomMicRequestTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingVoiceChatRoom.todayTotal'),
    week: t('pages.trackingVoiceChatRoom.weekTotal'),
    month: t('pages.trackingVoiceChatRoom.monthTotal'),
  }), seriesName: (t) => t('pages.trackingVoiceChatRoom.micRequestSeriesName')},
  {fileSlug: 'voice-chat-room', eventSlug: 'giftClick', fetch: () => trackingEventApi.getVoiceRoomGiftClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingVoiceChatRoom.todayTotal'),
    week: t('pages.trackingVoiceChatRoom.weekTotal'),
    month: t('pages.trackingVoiceChatRoom.monthTotal'),
  }), seriesName: (t) => t('pages.trackingVoiceChatRoom.giftClickSeriesName')},

  {fileSlug: 'recharge', eventSlug: 'panelEnter', fetch: () => trackingEventApi.getChargePanelEnterTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingRecharge.todayTotal'),
    week: t('pages.trackingRecharge.weekTotal'),
    month: t('pages.trackingRecharge.monthTotal'),
  }), seriesName: (t) => t('pages.trackingRecharge.panelEnterSeriesName')},
  {fileSlug: 'recharge', eventSlug: 'chargeClick', fetch: () => trackingEventApi.getChargeClickTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingRecharge.todayTotal'),
    week: t('pages.trackingRecharge.weekTotal'),
    month: t('pages.trackingRecharge.monthTotal'),
  }), seriesName: (t) => t('pages.trackingRecharge.chargeClickSeriesName')},
  {fileSlug: 'recharge', eventSlug: 'chargeResult', fetch: () => trackingEventApi.getChargeResultServerTrend(), summaryLabels: (t) => ({
    today: t('pages.trackingRecharge.todayTotal'),
    week: t('pages.trackingRecharge.weekTotal'),
    month: t('pages.trackingRecharge.monthTotal'),
  }), seriesName: (t) => t('pages.trackingRecharge.chargeResultSeriesName')},
]

/** 并行拉取埋点目录全部事件趋势(单事件失败不影响其它) */
export async function fetchAllTrackingTrendsForDashboardExport(): Promise<Map<string, HotLiveRoomJoinTrendRes>> {
  const map = new Map<string, HotLiveRoomJoinTrendRes>()
  const results = await Promise.allSettled(
      DASHBOARD_TRACKING_EXPORT_SPECS.map(async (spec) => {
        const data = await spec.fetch()
        return {key: trendKey(spec.fileSlug, spec.eventSlug), data: normalizeTrend(data)}
      }),
  )
  for (const result of results) {
    if (result.status === 'fulfilled') {
      map.set(result.value.key, result.value.data)
    }
  }
  return map
}

export function buildDashboardTrackingExportCsvFiles(
  t: TranslateFn,
  stamp: string,
  timeHeader: string,
  trends: Map<string, HotLiveRoomJoinTrendRes>,
): DashboardCsvBundleFile[] {
  const files: DashboardCsvBundleFile[] = []
  for (const spec of DASHBOARD_TRACKING_EXPORT_SPECS) {
    const trend = trends.get(trendKey(spec.fileSlug, spec.eventSlug))
    if (!trend) {
      continue
    }
    const summaryLabels = spec.summaryLabels(t)
    const countColumnHeader = spec.seriesName(t)
    for (const period of TRACKING_PERIODS) {
      files.push({
        filename: buildTrackingExportFilename(spec.fileSlug, spec.eventSlug, period, stamp),
        content: buildTrackingCountTrendCsvContent({
          todayCount: Number(trend.todayCount) || 0,
          weekCount: Number(trend.weekCount) || 0,
          monthCount: Number(trend.monthCount) || 0,
          points: pickTrackingTrendPeriodPoints(trend, period),
          countColumnHeader,
          summaryLabels,
          timeHeader,
        }),
      })
    }
  }
  return files
}
