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
            path: 'account-cfg',
            name: 'AccountCfgManagement',
            component: () => import('@/views/config/account-cfg.vue'),
            meta: {title: '账号配置'},
        },
        {
            path: 'simulator-cpu-keyword',
            name: 'SimulatorCpuKeywordManagement',
            component: () => import('@/views/config/simulator-cpu-keyword-list.vue'),
            meta: {title: '模拟器CPU关键词'},
        },
        {
            path: 'preload-cfg',
            name: 'PreloadCfgManagement',
            component: () => import('@/views/config/preload-cfg.vue'),
            meta: {title: '预热配置'},
        },
        {
            path: 'text-moderation',
            name: 'TextModerationCfgManagement',
            component: () => import('@/views/config/text-moderation.vue'),
            meta: {title: '敏感词过滤'},
        },
        {
            path: 'privacy-policy',
            name: 'PrivacyPolicyCfgManagement',
            component: () => import('@/views/config/privacy-policy.vue'),
            meta: {title: '隐私政策'},
        },
        {
            path: 'google-play',
            name: 'GooglePlayCfgManagement',
            component: () => import('@/views/config/google-play.vue'),
            meta: {title: 'Google Play'},
        },
        {
            path: 'upload-resource',
            name: 'UploadResourceCfgManagement',
            component: () => import('@/views/config/upload-resource.vue'),
            meta: {title: '资源域名'},
        },
        {
            path: 'data-sync',
            name: 'DataSyncCfgManagement',
            component: () => import('@/views/config/data-sync.vue'),
            meta: {title: '数据同步'},
        },
        {
            path: 'resource-monitor',
            name: 'ResourceMonitor',
            component: () => import('@/views/config/resource-monitor.vue'),
            meta: {title: '资源监控'},
        },
        {
            path: 'server-log',
            name: 'ServerLogExplorer',
            component: () => import('@/views/config/server-log-explorer.vue'),
            meta: {title: '服务器日志'},
        },
    ],
}
