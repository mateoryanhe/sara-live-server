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
            meta: {title: '进入秀场直播间'},
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
            meta: {title: '点击开始游戏'},
        },
    ],
}
