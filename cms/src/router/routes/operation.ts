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
            path: 'gift/gift-list',
            name: 'GiftManagement',
            component: () => import('@/views/operation/gift/gift-list.vue'),
            meta: {title: '礼物列表'},
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
    ],
}
