import type {RouteRecordRaw} from 'vue-router'

/** views/activity */
export const activityRoutes: RouteRecordRaw = {
    path: '/activity',
    meta: {title: '活动管理', icon: 'Present'},
    redirect: '/activity/first-recharge-activity-cfg',
    children: [
        {
            path: 'first-recharge-activity-cfg',
            name: 'FirstRechargeActivityManagement',
            component: () => import('@/views/activity/first-recharge-activity-cfg.vue'),
            meta: {title: '首充活动配置'},
        },
    ],
}
