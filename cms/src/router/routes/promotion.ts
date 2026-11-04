import type {RouteRecordRaw} from 'vue-router'

/** views/promotion */
export const promotionRoutes: RouteRecordRaw = {
    path: '/promotion',
    meta: {title: '推广', icon: 'Promotion'},
    redirect: '/promotion/channel-url-list',
    children: [
        {
            path: 'channel-url-list',
            name: 'PromotionChannelUrlListManagement',
            component: () => import('@/views/promotion/promotion-channel-url-list.vue'),
            meta: {title: 'H5推广渠道'},
        },
        {
            path: 'channel-daily-stat',
            name: 'PromotionChannelDailyStatManagement',
            component: () => import('@/views/promotion/promotion-channel-daily-stat-list.vue'),
            meta: {title: 'H5推广日统计'},
        },
        {
            path: 'apk-channel-list',
            name: 'ApkPromotionChannelListManagement',
            component: () => import('@/views/promotion/apk-promotion-channel-list.vue'),
            meta: {title: '落地页推广'},
        },
        {
            path: 'apk-channel-daily-stat',
            name: 'ApkPromotionChannelDailyStatManagement',
            component: () => import('@/views/promotion/apk-promotion-channel-daily-stat-list.vue'),
            meta: {title: '落地页下载日统计'},
        },
        {
            path: 'h5-app-promotion-daily-stat',
            name: 'H5AppPromotionDailyStatManagement',
            component: () => import('@/views/promotion/h5-app-promotion-daily-stat-list.vue'),
            meta: {title: 'H5下载App推广日志报表'},
        },
        {
            path: 'invite-log-list',
            name: 'PromotionInviteLogListManagement',
            component: () => import('@/views/promotion/promotion-invite-log-list.vue'),
            meta: {title: '邀请列表'},
        },
        {
            path: 'share-app-download-cfg',
            name: 'PromotionShareAppDownloadCfgManagement',
            component: () => import('@/views/promotion/share-app-download-cfg.vue'),
            meta: {title: '分享下载文案'},
        },
    ],
}
