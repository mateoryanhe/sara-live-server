import type {RouteRecordRaw} from 'vue-router'

/** views/tracking — 埋点事件统计 */
export const trackingRoutes: RouteRecordRaw = {
    path: '/tracking',
    meta: {title: '埋点', icon: 'DataLine'},
    children: [
        {
            path: 'hot-live-room-join',
            name: 'TrackingHotLiveRoomJoinManagement',
            component: () => import('@/views/tracking/hot-live-room-join.vue'),
            meta: {title: '秀场直播间'},
        },
        {
            path: 'call-1v1-initiate',
            name: 'TrackingCall1v1InitiateManagement',
            component: () => import('@/views/tracking/call-1v1-initiate.vue'),
            meta: {title: '1v1视频通话'},
        },
        {
            path: 'mini-game-round-start',
            name: 'TrackingMiniGameRoundStartManagement',
            component: () => import('@/views/tracking/mini-game-round-start.vue'),
            meta: {title: '游戏直播间'},
        },
        {
            path: 'short-video',
            name: 'TrackingShortVideoManagement',
            component: () => import('@/views/tracking/short-video.vue'),
            meta: {title: '短视频'},
        },
        {
            path: 'voice-chat-room',
            name: 'TrackingVoiceChatRoomManagement',
            component: () => import('@/views/tracking/voice-chat-room.vue'),
            meta: {title: '语聊房'},
        },
        {
            path: 'recharge',
            name: 'TrackingRechargeManagement',
            component: () => import('@/views/tracking/recharge.vue'),
            meta: {title: '充值'},
        },
    ],
}
