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
    ],
}
