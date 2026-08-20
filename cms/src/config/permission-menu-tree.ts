import type {RouteRecordRaw} from 'vue-router'
import {layoutRouteGroups} from '@/router/routes'
import {
    buttonPermissionKey,
    getPageButtons,
    type PageButtonDef,
} from './page-buttons'

/** 权限树节点 */
export type PermissionMenuNode =
    | PermissionMenuGroup
    | PermissionPageNode

export interface PermissionMenuGroup {
    kind: 'group'
    id: string
    /** i18n key，如 menu.UserManagement */
    titleKey: string
    children: PermissionMenuNode[]
}

/** 页面权限（按钮可分组；子页挂在 children） */
export interface PermissionPageNode {
    kind: 'page'
    pageName: string
    titleKey?: string
    buttonGroups?: readonly PermissionButtonGroupDef[]
    subPages?: readonly PermissionSubPageDef[]
}

export interface PermissionButtonGroupDef {
    id?: string
    titleKey: string
    buttonKeys: readonly string[]
}

export interface PermissionSubPageDef {
    pageName: string
    titleKey?: string
}

export interface PermissionModuleNode {
    id: string
    name: string
    children?: PermissionModuleNode[]
}

const USER_LIST_BUTTON_GROUPS: readonly PermissionButtonGroupDef[] = [
    {id: 'access', titleKey: 'pages.moduleList.groupAccess', buttonKeys: ['view', 'search']},
    {id: 'navigate', titleKey: 'pages.moduleList.groupNavigate', buttonKeys: ['viewDetail', 'viewAnchorDetail']},
    {id: 'anchor', titleKey: 'pages.moduleList.groupAnchor', buttonKeys: ['setAnchor', 'setSeniorAnchor', 'setAnchorType']},
    {id: 'currency', titleKey: 'pages.moduleList.groupCurrency', buttonKeys: ['goldAdd', 'goldSub', 'diamondAdd', 'diamondSub']},
    {
        id: 'account',
        titleKey: 'pages.moduleList.groupAccount',
        buttonKeys: ['ban', 'rankOff', 'rankOn', 'rechargeWhitelistOn', 'rechargeWhitelistOff', 'cancel', 'setUserType'],
    },
]

const GUILD_BUTTON_GROUPS: readonly PermissionButtonGroupDef[] = [
    {
        id: 'list',
        titleKey: 'pages.moduleList.groupGuildList',
        buttonKeys: [
            'view',
            'search',
            'create',
            'edit',
            'offShelf',
            'viewMembers',
            'viewDetail',
            'viewUserDetail',
            'viewAnchorSettlementLogs',
            'joinGuildAnchor',
            'batchSetAnchor',
            'batchSetSeniorAnchor',
        ],
    },
    {
        id: 'members',
        titleKey: 'pages.moduleList.groupGuildMembers',
        buttonKeys: ['viewDetail', 'viewUserDetail', 'ban', 'unban', 'exitGuild', 'setAnchorType'],
    },
]

function page(
    pageName: string,
    options?: {
        titleKey?: string
        buttonGroups?: readonly PermissionButtonGroupDef[]
        subPages?: readonly PermissionSubPageDef[]
    },
): PermissionPageNode {
    return {
        kind: 'page',
        pageName,
        titleKey: options?.titleKey,
        buttonGroups: options?.buttonGroups,
        subPages: options?.subPages,
    }
}

/**
 * CMS 权限树（扁平分组 + 子页挂靠父页 + 按钮分组）
 * 结构比侧边栏少一层，便于角色勾选；module 存库格式不变。
 */
export const PERMISSION_MENU_TREE: PermissionMenuNode[] = [
    {
        kind: 'group',
        id: 'dashboard',
        titleKey: 'menu.Dashboard',
        children: [page('Dashboard')],
    },
    {
        kind: 'group',
        id: 'user',
        titleKey: 'menu.UserManagement',
        children: [
            page('UserList', {buttonGroups: USER_LIST_BUTTON_GROUPS, subPages: [{pageName: 'UserDetail'}]}),
            page('AnchorListManagement', {subPages: [{pageName: 'AnchorDetail'}]}),
            page('LiveRoomRecycleBinManagement'),
            page('BotAnchorManagement'),
            page('RechargeOrderList'),
        ],
    },
    {
        kind: 'group',
        id: 'operation',
        titleKey: 'menu.OperationManagement',
        children: [
            page('BannerManagement'),
            page('ActivityMessageManagement'),
            page('RechargeCfgManagement'),
            page('VipCfgManagement'),
            page('WalletExchangeCfgManagement'),
            page('GuildManagement', {buttonGroups: GUILD_BUTTON_GROUPS, subPages: [{pageName: 'GuildDetail'}]}),
            page('PlatformAnchorList'),
            page('GuildRecycleBinManagement'),
            page('GuildProfileManagement'),
            page('GuildProfileAnchorSettlementLogList'),
            page('AnchorSalaryCfgManagement'),
            page('LiveRevenueShareCfgManagement'),
            page('AppPkgManagement'),
            page('RandomNicknameManagement'),
            page('CustomerServiceCfgManagement'),
        ],
    },
    {
        kind: 'group',
        id: 'live',
        titleKey: 'menu.LiveManagement',
        children: [
            page('GiftManagement'),
            page('AgoraCfgManagement'),
            page('TicketManagement'),
            page('PrivateRoomBillingManagement'),
            page('LiveCfgManagement'),
            page('LiveRoomTagManagement'),
        ],
    },
    {
        kind: 'group',
        id: 'log',
        titleKey: 'menu.LogManagement',
        children: [
            page('LiveRevenueLogList'),
            page('LiveRecordList'),
            page('VideoCallLogList'),
            page('GoldCurrencyLogList'),
            page('DiamondCurrencyLogList'),
            page('AnchorIncomeSettlementLogList'),
            page('GuildIncomeSettlementLogList'),
        ],
    },
    {
        kind: 'group',
        id: 'shortvideo',
        titleKey: 'menu.ShortVideoGroup',
        children: [
            page('ShortVideoManagement'),
            page('ShortVideoCategoryManagement'),
            page('ShortVideoCfgManagement'),
            page('ShortVideoWatchManagement'),
        ],
    },
    {
        kind: 'group',
        id: 'game',
        titleKey: 'menu.GameManagement',
        children: [
            page('GameShelfListManagement'),
            page('GamePlatformCfgManagement'),
            page('GameVendorGameListManagement'),
            page('GameBetLogListManagement'),
            page('GameWinLogListManagement'),
        ],
    },
    {
        kind: 'group',
        id: 'config',
        titleKey: 'menu.ConfigManagement',
        children: [
            page('AppTokenConfig'),
            page('AccountCfgManagement'),
            page('ServerRuntimeCfgManagement'),
            page('SimulatorCpuKeywordManagement'),
            page('TextModerationCfgManagement'),
            page('PrivacyPolicyCfgManagement'),
            page('GooglePlayCfgManagement'),
            page('UploadResourceCfgManagement'),
            page('DataSyncCfgManagement'),
            page('ResourceMonitor'),
            page('ServerLogExplorer'),
        ],
    },
    {
        kind: 'group',
        id: 'role',
        titleKey: 'menu.RoleManagementGroup',
        children: [
            page('RoleManagement', {subPages: [{pageName: 'ModuleList', titleKey: 'menu.ModuleList'}]}),
            page('CMSUserManagement'),
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

function buildButtonLeafNodes(
    pageName: string,
    buttons: PageButtonDef[],
): PermissionModuleNode[] {
    return buttons.map(btn => ({
        id: buttonPermissionKey(pageName, btn.key),
        name: btn.label,
    }))
}

function buildButtonGroupNode(
    pageName: string,
    group: PermissionButtonGroupDef,
    index: number,
    allButtons: PageButtonDef[],
    t: (key: string) => string,
): PermissionModuleNode | null {
    const buttons = allButtons.filter(btn => group.buttonKeys.includes(btn.key))
    if (buttons.length === 0) {
        return null
    }
    const groupId = group.id || String(index)
    return {
        id: `module_${pageName}_group_${groupId}`,
        name: t(group.titleKey),
        children: buildButtonLeafNodes(pageName, buttons),
    }
}

function buildPageModuleNode(
    pageName: string,
    t: (key: string) => string,
    options?: {
        titleKey?: string
        buttonGroups?: readonly PermissionButtonGroupDef[]
        subPages?: readonly PermissionSubPageDef[]
    },
): PermissionModuleNode {
    const route = routeMetaByPageName.get(pageName)
    const metaButtons = route?.meta?.buttons as PageButtonDef[] | undefined
    const allButtons = getPageButtons(pageName, metaButtons)

    const groupedKeys = new Set(
        (options?.buttonGroups || []).flatMap(group => group.buttonKeys),
    )
    const ungroupedButtons = options?.buttonGroups?.length
        ? allButtons.filter(btn => !groupedKeys.has(btn.key))
        : allButtons

    const children: PermissionModuleNode[] = []

    if (options?.buttonGroups?.length) {
        options.buttonGroups.forEach((group, index) => {
            const groupNode = buildButtonGroupNode(pageName, group, index, allButtons, t)
            if (groupNode) {
                children.push(groupNode)
            }
        })
        if (ungroupedButtons.length > 0) {
            children.push({
                id: `module_${pageName}_group_other`,
                name: t('pages.moduleList.groupOther'),
                children: buildButtonLeafNodes(pageName, ungroupedButtons),
            })
        }
    } else if (ungroupedButtons.length > 0) {
        children.push(...buildButtonLeafNodes(pageName, ungroupedButtons))
    }

    options?.subPages?.forEach(subPage => {
        children.push(buildPageModuleNode(subPage.pageName, t, {titleKey: subPage.titleKey}))
    })

    return {
        id: pageName,
        name: resolvePageTitle(pageName, options?.titleKey, t),
        children: children.length > 0 ? children : undefined,
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
    return buildPageModuleNode(node.pageName, t, {
        titleKey: node.titleKey,
        buttonGroups: node.buttonGroups,
        subPages: node.subPages,
    })
}

/** 生成权限配置树 */
export function buildPermissionModuleTree(t: (key: string) => string): PermissionModuleNode[] {
    return PERMISSION_MENU_TREE
        .map(node => buildPermissionMenuNode(node, t))
        .filter((node): node is PermissionModuleNode => node != null)
}
