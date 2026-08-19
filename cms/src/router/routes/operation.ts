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
            path: 'activity-message/activity-message-list',
            name: 'ActivityMessageManagement',
            component: () => import('@/views/operation/activity-message/activity-message-list.vue'),
            meta: {title: '活动消息'},
        },
        {
            path: 'guild/guild-list',
            name: 'GuildManagement',
            component: () => import('@/views/operation/guild/guild-list.vue'),
            meta: {title: '工会管理'},
        },
        {
            path: 'guild/guild-detail',
            name: 'GuildDetail',
            component: () => import('@/views/operation/guild/guild-detail.vue'),
            meta: {
                title: '工会详情',
                hidden: true,
                parentPermission: ['GuildManagement', 'GuildProfileManagement', 'AnchorListManagement'],
            },
        },
        {
            path: 'guild/platform-anchor-list',
            name: 'PlatformAnchorList',
            component: () => import('@/views/operation/guild/platform-anchor-list.vue'),
            meta: {title: '平台主播'},
        },
        {
            path: 'guild/guild-recycle-bin',
            name: 'GuildRecycleBinManagement',
            component: () => import('@/views/operation/guild/guild-recycle-bin.vue'),
            meta: {title: '工会垃圾库'},
        },
        {
            path: 'guild/guild-anchor-import-result',
            name: 'GuildAnchorImportResult',
            component: () => import('@/views/operation/guild/guild-anchor-import-result.vue'),
            meta: {title: '工会主播导入结果', hidden: true, parentPermission: 'GuildManagement'},
        },
        {
            path: 'guild/guild-members',
            name: 'GuildMembers',
            component: () => import('@/views/operation/guild/guild-members.vue'),
            meta: {title: '工会成员', hidden: true, parentPermission: 'GuildManagement'},
        },
        {
            path: 'guild/guild-profile-members',
            name: 'GuildProfileMembers',
            component: () => import('@/views/operation/guild/guild-members.vue'),
            meta: {title: '名下主播', hidden: true, parentPermission: 'GuildProfileManagement'},
        },
        {
            path: 'guild/guild-profile',
            name: 'GuildProfileManagement',
            component: () => import('@/views/operation/guild/guild-profile.vue'),
            meta: {title: '工会数据查询'},
        },
        {
            path: 'guild/guild-anchor-income-settlement-log-list',
            name: 'GuildProfileAnchorSettlementLogList',
            component: () => import('@/views/operation/guild/guild-anchor-income-settlement-log-list.vue'),
            meta: {title: '名下主播结算流水'},
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
        {
            path: 'wallet/wallet-exchange-cfg',
            name: 'WalletExchangeCfgManagement',
            component: () => import('@/views/operation/wallet/wallet-exchange-cfg.vue'),
            meta: {title: '金币兑换配置'},
        },
        {
            path: 'salary/anchor-salary-cfg-list',
            name: 'AnchorSalaryCfgManagement',
            component: () => import('@/views/operation/salary/anchor-salary-cfg-list.vue'),
            meta: {title: '主播结算薪资'},
        },
        {
            path: 'salary/live-revenue-share-cfg',
            name: 'LiveRevenueShareCfgManagement',
            component: () => import('@/views/operation/salary/live-revenue-share-cfg.vue'),
            meta: {title: '流水分佣配置'},
        },
        {
            path: 'salary/anchor-income-settlement-log-list',
            name: 'AnchorIncomeSettlementLogList',
            component: () => import('@/views/operation/salary/anchor-income-settlement-log-list.vue'),
            meta: {title: '主播结算流水'},
        },
        {
            path: 'salary/guild-income-settlement-log-list',
            name: 'GuildIncomeSettlementLogList',
            component: () => import('@/views/operation/salary/guild-income-settlement-log-list.vue'),
            meta: {title: '工会结算流水'},
        },
    ],
}
