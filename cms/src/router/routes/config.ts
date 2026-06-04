import type {RouteRecordRaw} from 'vue-router'

/** views/config */
export const configRoutes: RouteRecordRaw = {
    path: '/config',
    meta: {title: '系统配置', icon: 'Setting'},
    redirect: '/config/global',
    children: [
        {
            path: 'global',
            name: 'GlobalConfig',
            component: () => import('@/views/config/global.vue'),
            meta: {title: '全局配置'},
        },
        {
            path: 'app-token',
            name: 'AppTokenConfig',
            component: () => import('@/views/config/app-token.vue'),
            meta: {title: 'App Token'},
        },
        {
            path: 'text-moderation',
            name: 'TextModerationCfgManagement',
            component: () => import('@/views/config/text-moderation.vue'),
            meta: {title: '敏感词过滤'},
        },
    ],
}
