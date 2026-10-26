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
            meta: {title: '推广渠道URL列表'},
        },
        {
            path: 'channel-daily-stat',
            name: 'PromotionChannelDailyStatManagement',
            component: () => import('@/views/promotion/promotion-channel-daily-stat-list.vue'),
            meta: {title: '推广渠道日统计'},
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
