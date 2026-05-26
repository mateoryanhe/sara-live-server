import type {RouteRecordRaw} from 'vue-router'

/** views/game */
export const gameRoutes: RouteRecordRaw = {
    path: '/game',
    meta: {title: '游戏管理', icon: 'Cpu'},
    redirect: '/game/game-cfg-list',
    children: [
        {
            path: 'game-cfg-list',
            name: 'GameCfgManagement',
            component: () => import('@/views/game/game-cfg-list.vue'),
            meta: {title: '游戏配置'},
        },
    ],
}
