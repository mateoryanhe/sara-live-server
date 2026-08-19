import type {RouteRecordRaw} from 'vue-router'
import {layoutRouteGroups} from '@/router/routes'
import {
    buttonPermissionKey,
    getPageButtons,
    type PageButtonDef,
} from './page-buttons'

/** 权限树节点：与侧边栏分组结构一致 */
export type PermissionMenuNode =
    | PermissionMenuGroup
    | PermissionPageNode
    | PermissionPageSliceNode
    | PermissionSubPageNode

export interface PermissionMenuGroup {
    kind: 'group'
    id: string
    /** i18n key，如 menu.UserManagement */
    titleKey: string
    children: PermissionMenuNode[]
}

/** 整页权限（含全部或指定按钮） */
export interface PermissionPageNode {
    kind: 'page'
    pageName: string
    titleKey?: string
}

/** 同一页面按钮拆开展示（module 仍为 PageName:action） */
export interface PermissionPageSliceNode {
    kind: 'pageSlice'
    pageName: string
    titleKey: string
    buttonKeys: string[]
}

/** 隐藏子页（独立 module 命名空间，如 UserDetail） */
export interface PermissionSubPageNode {
    kind: 'subPage'
    pageName: string
    titleKey?: string
}

export interface PermissionModuleNode {
    id: string
    name: string
    children?: PermissionModuleNode[]
}

const GUILD_LIST_BUTTON_KEYS = [
    'view',
    'search',
    'create',
    'edit',
    'offShelf',
    'viewMembers',
    'viewDetail',
    'viewAnchorSettlementLogs',
    'joinGuildAnchor',
    'batchSetAnchor',
    'batchSetSeniorAnchor',
] as const

const GUILD_MEMBER_BUTTON_KEYS = [
    'viewDetail',
    'ban',
    'unban',
    'exitGuild',
    'setAnchorType',
] as const

/**
 * CMS 权限树结构（与 layout 侧边栏分组、顺序一致）
 * 修改导航结构时请同步更新此配置。
 */
export const PERMISSION_MENU_TREE: PermissionMenuNode[] = [
    {
        kind: 'group',
        id: 'dashboard',
        titleKey: 'menu.Dashboard',
        children: [{kind: 'page', pageName: 'Dashboard'}],
    },
    {
        kind: 'group',
        id: 'user',
        titleKey: 'menu.UserManagement',
        children: [
            {kind: 'page', pageName: 'UserList'},
            {kind: 'subPage', pageName: 'UserDetail', titleKey: 'menu.UserDetail'},
            {kind: 'page', pageName: 'AnchorListManagement'},
            {kind: 'subPage', pageName: 'AnchorDetail', titleKey: 'menu.AnchorDetail'},
            {kind: 'page', pageName: 'LiveRoomRecycleBinManagement'},
            {kind: 'page', pageName: 'BotAnchorManagement'},
            {kind: 'page', pageName: 'RechargeOrderList'},
        ],
    },
    {
        kind: 'group',
        id: 'operation',
        titleKey: 'menu.OperationManagement',
        children: [
            {
                kind: 'group',
                id: 'operation-content',
                titleKey: 'menu.OperationContentGroup',
                children: [
                    {kind: 'page', pageName: 'BannerManagement'},
                    {kind: 'page', pageName: 'ActivityMessageManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'operation-recharge',
                titleKey: 'menu.OperationRechargeGroup',
                children: [
                    {kind: 'page', pageName: 'RechargeCfgManagement'},
                    {kind: 'page', pageName: 'VipCfgManagement'},
                    {kind: 'page', pageName: 'WalletExchangeCfgManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'operation-guild',
                titleKey: 'menu.OperationGuildGroup',
                children: [
                    {
                        kind: 'pageSlice',
                        pageName: 'GuildManagement',
                        titleKey: 'menu.GuildManagement',
                        buttonKeys: [...GUILD_LIST_BUTTON_KEYS],
                    },
                    {
                        kind: 'pageSlice',
                        pageName: 'GuildManagement',
                        titleKey: 'menu.GuildMembers',
                        buttonKeys: [...GUILD_MEMBER_BUTTON_KEYS],
                    },
                    {kind: 'subPage', pageName: 'GuildDetail', titleKey: 'menu.GuildDetail'},
                    {kind: 'page', pageName: 'PlatformAnchorList'},
                    {kind: 'page', pageName: 'GuildRecycleBinManagement'},
                    {kind: 'page', pageName: 'GuildProfileManagement'},
                    {kind: 'page', pageName: 'GuildProfileAnchorSettlementLogList'},
                ],
            },
            {
                kind: 'group',
                id: 'operation-settlement',
                titleKey: 'menu.OperationSettlementGroup',
                children: [
                    {kind: 'page', pageName: 'AnchorSalaryCfgManagement'},
                    {kind: 'page', pageName: 'LiveRevenueShareCfgManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'operation-app',
                titleKey: 'menu.OperationAppGroup',
                children: [
                    {kind: 'page', pageName: 'AppPkgManagement'},
                    {kind: 'page', pageName: 'RandomNicknameManagement'},
                    {kind: 'page', pageName: 'CustomerServiceCfgManagement'},
                ],
            },
        ],
    },
    {
        kind: 'group',
        id: 'live',
        titleKey: 'menu.LiveManagement',
        children: [
            {kind: 'page', pageName: 'GiftManagement'},
            {kind: 'page', pageName: 'AgoraCfgManagement'},
            {kind: 'page', pageName: 'TicketManagement'},
            {kind: 'page', pageName: 'PrivateRoomBillingManagement'},
            {kind: 'page', pageName: 'LiveCfgManagement'},
            {kind: 'page', pageName: 'LiveRoomTagManagement'},
        ],
    },
    {
        kind: 'group',
        id: 'log',
        titleKey: 'menu.LogManagement',
        children: [
            {
                kind: 'group',
                id: 'log-live',
                titleKey: 'menu.LiveLogGroup',
                children: [
                    {kind: 'page', pageName: 'LiveRevenueLogList'},
                    {kind: 'page', pageName: 'LiveRecordList'},
                ],
            },
            {
                kind: 'group',
                id: 'log-call',
                titleKey: 'menu.CallLogGroup',
                children: [{kind: 'page', pageName: 'VideoCallLogList'}],
            },
            {
                kind: 'group',
                id: 'log-user',
                titleKey: 'menu.UserLogGroup',
                children: [
                    {kind: 'page', pageName: 'GoldCurrencyLogList'},
                    {kind: 'page', pageName: 'DiamondCurrencyLogList'},
                ],
            },
            {
                kind: 'group',
                id: 'log-settlement',
                titleKey: 'menu.SettlementLogGroup',
                children: [
                    {kind: 'page', pageName: 'AnchorIncomeSettlementLogList'},
                    {kind: 'page', pageName: 'GuildIncomeSettlementLogList'},
                ],
            },
        ],
    },
    {
        kind: 'group',
        id: 'shortvideo',
        titleKey: 'menu.ShortVideoGroup',
        children: [
            {kind: 'page', pageName: 'ShortVideoManagement'},
            {kind: 'page', pageName: 'ShortVideoCategoryManagement'},
            {kind: 'page', pageName: 'ShortVideoCfgManagement'},
            {kind: 'page', pageName: 'ShortVideoWatchManagement'},
        ],
    },
    {
        kind: 'group',
        id: 'game',
        titleKey: 'menu.GameManagement',
        children: [
            {kind: 'page', pageName: 'GamePlatformCfgManagement'},
            {kind: 'page', pageName: 'GameVendorGameListManagement'},
            {kind: 'page', pageName: 'GameBetLogListManagement'},
            {kind: 'page', pageName: 'GameWinLogListManagement'},
        ],
    },
    {
        kind: 'group',
        id: 'config',
        titleKey: 'menu.ConfigManagement',
        children: [
            {
                kind: 'group',
                id: 'config-basic',
                titleKey: 'menu.ConfigBasicGroup',
                children: [
                    {kind: 'page', pageName: 'AppTokenConfig'},
                    {kind: 'page', pageName: 'AccountCfgManagement'},
                    {kind: 'page', pageName: 'PreloadCfgManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'config-security',
                titleKey: 'menu.ConfigSecurityGroup',
                children: [
                    {kind: 'page', pageName: 'SimulatorCpuKeywordManagement'},
                    {kind: 'page', pageName: 'TextModerationCfgManagement'},
                    {kind: 'page', pageName: 'PrivacyPolicyCfgManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'config-platform',
                titleKey: 'menu.ConfigPlatformGroup',
                children: [
                    {kind: 'page', pageName: 'GooglePlayCfgManagement'},
                    {kind: 'page', pageName: 'UploadResourceCfgManagement'},
                    {kind: 'page', pageName: 'DataSyncCfgManagement'},
                ],
            },
            {
                kind: 'group',
                id: 'config-ops',
                titleKey: 'menu.ConfigOpsGroup',
                children: [
                    {kind: 'page', pageName: 'ResourceMonitor'},
                    {kind: 'page', pageName: 'ServerLogExplorer'},
                ],
            },
        ],
    },
    {
        kind: 'group',
        id: 'role',
        titleKey: 'menu.RoleManagementGroup',
        children: [
            {kind: 'page', pageName: 'RoleManagement'},
            {kind: 'subPage', pageName: 'ModuleList', titleKey: 'menu.ModuleList'},
            {kind: 'page', pageName: 'CMSUserManagement'},
        ],
    },
]

const routeMetaByPageName = new Map<string, RouteRecordRaw>()

function collectRoutes(group: RouteRecordRaw) {
    group.children?.forEach(child => {
        if (child.name) {
            routeMetaByPageName.set(String(child.name), child)
        }
    })
}

layoutRouteGroups.forEach(collectRoutes)

function defaultTitleKey(pageName: string): string {
    return `menu.${pageName}`
}

function resolvePageTitle(pageName: string, titleKey: string | undefined, t: (key: string) => string): string {
    const key = titleKey || defaultTitleKey(pageName)
    const translated = t(key)
    if (translated !== key) {
        return translated
    }
    const route = routeMetaByPageName.get(pageName)
    return String(route?.meta?.title || pageName)
}

function buildButtonNodes(
    pageName: string,
    t: (key: string) => string,
    metaButtons?: PageButtonDef[],
    buttonKeys?: readonly string[],
): PermissionModuleNode[] {
    const allButtons = getPageButtons(pageName, metaButtons)
    const filtered = buttonKeys?.length
        ? allButtons.filter(btn => buttonKeys.includes(btn.key))
        : allButtons
    return filtered.map(btn => ({
        id: buttonPermissionKey(pageName, btn.key),
        name: btn.label,
    }))
}

function buildPageModuleNode(
    pageName: string,
    t: (key: string) => string,
    titleKey?: string,
    buttonKeys?: readonly string[],
): PermissionModuleNode {
    const route = routeMetaByPageName.get(pageName)
    const metaButtons = route?.meta?.buttons as PageButtonDef[] | undefined
    const buttonNodes = buildButtonNodes(pageName, t, metaButtons, buttonKeys)
    return {
        id: pageName,
        name: resolvePageTitle(pageName, titleKey, t),
        children: buttonNodes.length > 0 ? buttonNodes : undefined,
    }
}

function buildPermissionMenuNode(node: PermissionMenuNode, t: (key: string) => string): PermissionModuleNode | null {
    if (node.kind === 'group') {
        const children = node.children
            .map(child => buildPermissionMenuNode(child, t))
            .filter((child): child is PermissionModuleNode => child != null)
        if (children.length === 0) {
            return null
        }
        return {
            id: `module_${node.id}`,
            name: t(node.titleKey),
            children,
        }
    }
    if (node.kind === 'pageSlice') {
        const route = routeMetaByPageName.get(node.pageName)
        const metaButtons = route?.meta?.buttons as PageButtonDef[] | undefined
        const buttonNodes = buildButtonNodes(node.pageName, t, metaButtons, node.buttonKeys)
        return {
            id: `module_${node.pageName}_${node.titleKey.replace(/\./g, '_')}`,
            name: resolvePageTitle(node.pageName, node.titleKey, t),
            children: buttonNodes.length > 0 ? buttonNodes : undefined,
        }
    }
    if (node.kind === 'subPage') {
        return buildPageModuleNode(node.pageName, t, node.titleKey)
    }
    return buildPageModuleNode(node.pageName, t, node.titleKey)
}

/** 生成与侧边栏对齐的权限配置树 */
export function buildPermissionModuleTree(t: (key: string) => string): PermissionModuleNode[] {
    return PERMISSION_MENU_TREE
        .map(node => buildPermissionMenuNode(node, t))
        .filter((node): node is PermissionModuleNode => node != null)
}
