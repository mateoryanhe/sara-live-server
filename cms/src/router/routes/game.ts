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
            meta: {title: '游戏库'},
        },
        {
            path: 'game-shelf-list',
            name: 'GameShelfListManagement',
            component: () => import('@/views/game/game-shelf-list.vue'),
            meta: {title: '上架游戏列表'},
        },
        {
            path: 'game-win-log-list',
            name: 'GameWinLogListManagement',
            component: () => import('@/views/game/game-win-log-list.vue'),
            meta: {title: '派彩记录'},
        },
        {
            path: 'game-bet-log-list',
            name: 'GameBetLogListManagement',
            component: () => import('@/views/game/game-bet-log-list.vue'),
            meta: {title: '下注记录'},
        },
    ],
}
