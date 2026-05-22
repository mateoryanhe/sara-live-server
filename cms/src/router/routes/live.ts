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
    ],
}
