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
            path: 'gift/gift-edit/:id',
            name: 'GiftManagementEdit',
            component: () => import('@/views/live/gift/gift-edit.vue'),
            meta: {title: '礼物编辑', hidden: true, parentPermission: 'GiftManagement'},
        },
        {
            path: 'gift/lucky-gift-pool-cfg',
            name: 'LuckyGiftPoolCfgManagement',
            component: () => import('@/views/live/gift/lucky-gift-pool-cfg.vue'),
            meta: {title: '幸运礼物奖池配置'},
        },
        {
            path: 'gift/lucky-gift-pool-pending-balance',
            name: 'LuckyGiftPoolPendingBalanceList',
            component: () => import('@/views/live/gift/lucky-gift-pool-pending-balance-list.vue'),
            meta: {title: '奖池临时入账'},
        },
        {
            path: 'agora-cfg',
            name: 'AgoraCfgManagement',
            component: () => import('@/views/live/agora-cfg.vue'),
            meta: {title: '声网配置'},
        },
        {
            path: 'private-room-billing/billing-list',
            name: 'PrivateRoomBillingManagement',
            component: () => import('@/views/live/private-room-billing/billing-list.vue'),
            meta: {title: '1v1通话计费'},
        },
        {
            path: 'live-config/live-config',
            name: 'LiveCfgManagement',
            component: () => import('@/views/live/live-config/live-config.vue'),
            meta: {title: '直播配置'},
        },
        {
            path: 'live-room-tag/live-room-tag-list',
            name: 'LiveRoomTagManagement',
            component: () => import('@/views/live/live-room-tag/live-room-tag-list.vue'),
            meta: {title: '直播间标签'},
        },
        {
            path: 'voice-chat-mic-seat-tier/tier-list',
            name: 'VoiceChatMicSeatTierManagement',
            component: () => import('@/views/live/voice-chat-mic-seat-tier/tier-list.vue'),
            meta: {title: '麦位挡位'},
        },
        {
            path: 'anchor-video-price-tier/tier-list',
            name: 'AnchorVideoPriceTierManagement',
            component: () => import('@/views/live/anchor-video-price-tier/tier-list.vue'),
            meta: {title: '主播相册短视频价格'},
        },
    ],
}
