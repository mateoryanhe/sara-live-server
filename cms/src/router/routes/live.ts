import type {RouteRecordRaw} from 'vue-router'

/** views/live */
export const liveRoutes: RouteRecordRaw = {
    path: '/live',
    meta: {title: '直播管理', icon: 'VideoPlay'},
    redirect: '/live/gift/gift-list',
    children: [
        {
            path: 'gift/gift-list',
            name: 'GiftManagement',
            component: () => import('@/views/live/gift/gift-list.vue'),
            meta: {title: '礼物管理'},
        },
        {
            path: 'agora-cfg',
            name: 'AgoraCfgManagement',
            component: () => import('@/views/live/agora-cfg.vue'),
            meta: {title: '声网配置'},
        },
        {
            path: 'ticket/ticket-list',
            name: 'TicketManagement',
            component: () => import('@/views/live/ticket/ticket-list.vue'),
            meta: {title: '门票管理'},
        },
        {
            path: 'private-room-billing/billing-list',
            name: 'PrivateRoomBillingManagement',
            component: () => import('@/views/live/private-room-billing/billing-list.vue'),
            meta: {title: '私密直播间计费'},
        },
        {
            path: 'live-config/live-config',
            name: 'LiveCfgManagement',
            component: () => import('@/views/live/live-config/live-config.vue'),
            meta: {title: '直播配置'},
        },
    ],
}
