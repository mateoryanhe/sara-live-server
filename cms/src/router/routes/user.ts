import type {RouteRecordRaw} from 'vue-router'

/** views/user */
export const userRoutes: RouteRecordRaw = {
    path: '/user',
    meta: {title: '用户管理', icon: 'User'},
    redirect: '/user/account/user-list',
    children: [
        {
            path: 'account/user-list',
            name: 'UserList',
            component: () => import('@/views/user/account/user-list.vue'),
            meta: {title: '用户列表'},
        },
        {
            path: 'anchor/anchor-list',
            name: 'AnchorListManagement',
            component: () => import('@/views/user/anchor/anchor-list.vue'),
            meta: {title: '主播列表'},
        },
        {
            path: 'account/ban-user',
            name: 'BanUser',
            component: () => import('@/views/user/account/ban-user.vue'),
            meta: {title: '封禁用户', hidden: true},
        },
        {
            path: 'recharge-order/recharge-order-list',
            name: 'RechargeOrderList',
            component: () => import('@/views/user/recharge-order/recharge-order-list.vue'),
            meta: {title: '充值订单'},
        },
    ],
}
