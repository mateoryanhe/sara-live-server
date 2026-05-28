import type {RouteRecordRaw} from 'vue-router'

/** views/dashboard */
export const dashboardRoutes: RouteRecordRaw = {
    path: '/dashboard',
    meta: {title: '仪表盘', icon: 'Odometer'},
    children: [
        {
            path: '',
            name: 'Dashboard',
            component: () => import('@/views/dashboard/index.vue'),
            meta: {title: '仪表盘'},
        },
    ],
}
