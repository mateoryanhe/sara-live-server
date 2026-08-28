import type {RouteRecordRaw} from 'vue-router'

/** views/log */
export const logRoutes: RouteRecordRaw = {
    path: '/log',
    meta: {title: '日志', icon: 'Document'},
    redirect: '/log/live/live-record-list',
    children: [
        {
            path: 'live/live-record-list',
            name: 'LiveRecordList',
            component: () => import('@/views/log/live/live-record-list.vue'),
            meta: {title: '直播记录', parentTitle: '直播日志'},
        },
        {
            path: 'live/revenue-log-list',
            name: 'LiveRevenueLogList',
            component: () => import('@/views/log/live/revenue-log-list.vue'),
            meta: {title: '直播收益流水', parentTitle: '直播日志'},
        },
        {
            path: 'live/daily-effective-live-list',
            name: 'LiveDailyEffectiveLiveList',
            component: () => import('@/views/log/live/daily-effective-live-list.vue'),
            meta: {title: '每日流水', parentTitle: '直播日志'},
        },
        {
            path: 'live/weekly-unsettled-live-list',
            name: 'LiveWeeklyUnsettledLiveList',
            component: () => import('@/views/log/live/weekly-unsettled-live-list.vue'),
            meta: {title: '本周流水', parentTitle: '直播日志'},
        },
        {
            path: 'live/video-call-log-list',
            name: 'VideoCallLogList',
            component: () => import('@/views/log/live/video-call-log-list.vue'),
            meta: {title: '视频通话日志', parentTitle: '直播日志'},
        },
        {
            path: 'call/video-call-log-list',
            redirect: '/log/live/video-call-log-list',
        },
    ],
}
