import type {RouteRecordRaw} from 'vue-router'

/** views/withdraw */
export const withdrawRoutes: RouteRecordRaw = {
    path: '/withdraw',
    meta: {title: '提现', icon: 'Wallet'},
    redirect: '/withdraw/user-withdraw-log-list',
    children: [
        {
            path: 'user-withdraw-log-list',
            name: 'UserWithdrawLogList',
            component: () => import('@/views/withdraw/user-withdraw-log-list.vue'),
            meta: {title: '用户提现记录'},
        },
        {
            path: 'user-withdraw-cfg',
            name: 'UserWithdrawCfgManagement',
            component: () => import('@/views/withdraw/user-withdraw-cfg.vue'),
            meta: {title: '用户提现配置'},
        },
        {
            path: 'user-withdraw-tier-list',
            name: 'UserWithdrawTierListManagement',
            component: () => import('@/views/withdraw/user-withdraw-tier-list.vue'),
            meta: {title: '用户提现档位'},
        },
        {
            path: 'user-withdraw-payout-country-cfg',
            name: 'UserWithdrawPayoutCountryCfgManagement',
            component: () => import('@/views/operation/recharge/payment-country-cfg.vue'),
            meta: {title: '提现代付国家配置', paymentCountryCfgScope: 'userWithdrawPayout'},
        },
        {
            path: 'user-withdraw-payout-country-cfg/:continent',
            name: 'UserWithdrawPayoutCountryCfgContinent',
            component: () => import('@/views/operation/recharge/payment-country-continent.vue'),
            meta: {
                title: '提现代付国家/地区配置',
                hidden: true,
                paymentCountryCfgScope: 'userWithdrawPayout',
                parentPermission: ['UserWithdrawPayoutCountryCfgManagement'],
            },
        },
        {
            path: 'user-withdraw-payout-country-cfg/:continent/:countryCode',
            name: 'UserWithdrawPayoutCountryCfgEdit',
            component: () => import('@/views/operation/recharge/payment-country-edit.vue'),
            meta: {
                title: '提现代付配置编辑',
                hidden: true,
                paymentCountryCfgScope: 'userWithdrawPayout',
                parentPermission: ['UserWithdrawPayoutCountryCfgManagement'],
            },
        },
    ],
}
