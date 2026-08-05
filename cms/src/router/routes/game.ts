import type {RouteRecordRaw} from 'vue-router'

/** views/game */
export const gameRoutes: RouteRecordRaw = {
    path: '/game',
    meta: {title: '游戏管理', icon: 'Cpu'},
    redirect: '/game/game-platform-cfg',
    children: [
        {
            path: 'game-platform-cfg',
            name: 'GamePlatformCfgManagement',
            component: () => import('@/views/game/game-platform-cfg.vue'),
            meta: {title: '平台接入配置'},
        },
        {
            path: 'game-list',
            name: 'GameVendorGameListManagement',
            component: () => import('@/views/game/game-list.vue'),
            meta: {title: '游戏列表'},
        },
    ],
}
