import type {RouteRecordRaw} from 'vue-router'

/** views/operation */
export const operationRoutes: RouteRecordRaw = {
    path: '/operation',
    meta: {title: '运营管理', icon: 'Stamp'},
    redirect: '/operation/banner/banner-list',
    children: [
        {
            path: 'banner/banner-list',
            name: 'BannerManagement',
            component: () => import('@/views/operation/banner/banner-list.vue'),
            meta: {title: '首页Banner'},
        },
        {
            path: 'guild/guild-list',
            name: 'GuildManagement',
            component: () => import('@/views/operation/guild/guild-list.vue'),
            meta: {title: '工会管理'},
        },
        {
            path: 'recharge/recharge-cfg-list',
            name: 'RechargeCfgManagement',
            component: () => import('@/views/operation/recharge/recharge-cfg-list.vue'),
            meta: {title: '充值配置'},
        },
        {
            path: 'vip/vip-cfg-list',
            name: 'VipCfgManagement',
            component: () => import('@/views/operation/vip/vip-cfg-list.vue'),
            meta: {title: 'VIP配置'},
        },
        {
            path: 'app-pkg/app-pkg-list',
            name: 'AppPkgManagement',
            component: () => import('@/views/operation/app-pkg/app-pkg-list.vue'),
            meta: {title: 'App包管理'},
        },
        {
            path: 'random-nickname/random-nickname-cfg',
            name: 'RandomNicknameManagement',
            component: () => import('@/views/operation/random-nickname/random-nickname-cfg.vue'),
            meta: {title: '随机昵称库'},
        },
        {
            path: 'customer-service/customer-service-cfg',
            name: 'CustomerServiceCfgManagement',
            component: () => import('@/views/operation/customer-service/customer-service-cfg.vue'),
            meta: {title: '客服联系配置'},
        },
    ],
}
