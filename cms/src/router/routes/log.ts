import type {RouteRecordRaw} from 'vue-router'

/** views/log */
export const logRoutes: RouteRecordRaw = {
    path: '/log',
    meta: {title: '日志', icon: 'Document'},
    redirect: '/log/live/gift-log-list',
    children: [
        {
            path: 'live/gift-log-list',
            name: 'LiveGiftLogList',
            component: () => import('@/views/log/live/gift-log-list.vue'),
            meta: {title: '礼物流水', parentTitle: '直播日志'},
        },
        {
            path: 'live/live-record-list',
            name: 'LiveRecordList',
            component: () => import('@/views/log/live/live-record-list.vue'),
            meta: {title: '直播记录', parentTitle: '直播日志'},
        },
    ],
}
