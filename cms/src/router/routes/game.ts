import type {RouteRecordRaw} from 'vue-router'

/** views/game */
export const gameRoutes: RouteRecordRaw = {
    path: '/game',
    meta: {title: '游戏管理', icon: 'Cpu'},
    redirect: '/game/game-shelf-list',
    children: [
        {
            path: 'game-shelf-list',
            name: 'GameShelfListManagement',
            component: () => import('@/views/game/game-shelf-list.vue'),
            meta: {title: '上架游戏列表'},
        },
        {
            path: 'game-vendor-config',
            name: 'GameVendorConfig',
            component: () => import('@/views/game/game-vendor-config.vue'),
            meta: {
                title: '第三方配置',
                hidden: true,
                parentPermission: ['GameShelfListManagement'],
            },
        },
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
            path: 'game-bet-log-list',
            redirect: '/log/game/game-bet-log-list',
        },
        {
            path: 'game-win-log-list',
            redirect: '/log/game/game-win-log-list',
        },
    ],
}
