import type {RouteRecordRaw} from 'vue-router'

/** views/config */
export const configRoutes: RouteRecordRaw = {
    path: '/config',
    meta: {title: '系统配置', icon: 'Setting'},
    redirect: '/config/app-token',
    children: [
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
        {
            path: 'upload-resource',
            name: 'UploadResourceCfgManagement',
            component: () => import('@/views/config/upload-resource.vue'),
            meta: {title: '资源域名'},
        },
        {
            path: 'resource-monitor',
            name: 'ResourceMonitor',
            component: () => import('@/views/config/resource-monitor.vue'),
            meta: {title: '资源监控'},
        },
    ],
}
