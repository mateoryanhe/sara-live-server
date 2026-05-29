import type {RouteRecordRaw} from 'vue-router'

/** views/log */
export const logRoutes: RouteRecordRaw = {
    path: '/log',
    meta: {title: '日志', icon: 'Document'},
    redirect: '/log/live/revenue-log-list',
    children: [
        {
            path: 'live/revenue-log-list',
            name: 'LiveRevenueLogList',
            component: () => import('@/views/log/live/revenue-log-list.vue'),
            meta: {title: '直播收益流水', parentTitle: '直播日志'},
        },
        {
            path: 'live/live-record-list',
            name: 'LiveRecordList',
            component: () => import('@/views/log/live/live-record-list.vue'),
            meta: {title: '直播记录', parentTitle: '直播日志'},
        },
    ],
}
