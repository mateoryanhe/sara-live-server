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
            path: 'account/user-detail',
            name: 'UserDetail',
            component: () => import('@/views/user/account/user-detail.vue'),
            meta: {
                title: '用户详情',
                hidden: true,
                parentPermission: ['UserList'],
            },
        },
        {
            path: 'anchor/anchor-list',
            name: 'AnchorListManagement',
            component: () => import('@/views/user/anchor/anchor-list.vue'),
            meta: {title: '主播列表'},
        },
        {
            path: 'anchor/live-room-recycle-bin',
            name: 'LiveRoomRecycleBinManagement',
            component: () => import('@/views/user/anchor/live-room-recycle-bin.vue'),
            meta: {title: '主播间回收站'},
        },
        {
            path: 'anchor/anchor-detail',
            name: 'AnchorDetail',
            component: () => import('@/views/user/anchor/anchor-detail.vue'),
            meta: {
                title: '主播详情',
                hidden: true,
                parentPermission: ['AnchorListManagement', 'GuildManagement', 'GuildProfileManagement'],
            },
        },
        {
            path: 'bot-anchor/bot-anchor-list',
            name: 'BotAnchorManagement',
            component: () => import('@/views/user/bot-anchor/bot-anchor-list.vue'),
            meta: {title: '机器人主播'},
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
        {
            path: 'currency-log/gold-log-list',
            name: 'GoldCurrencyLogList',
            component: () => import('@/views/user/currency-log/currency-log-list.vue'),
            meta: {title: '金币流水', currencyType: 1},
        },
        {
            path: 'currency-log/diamond-log-list',
            name: 'DiamondCurrencyLogList',
            component: () => import('@/views/user/currency-log/currency-log-list.vue'),
            meta: {title: '钻石流水', currencyType: 2},
        },
    ],
}
