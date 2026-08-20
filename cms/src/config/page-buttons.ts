/** 页面按钮权限定义：module 存库格式为 PageName 或 PageName:action */

export interface PageButtonDef {
    key: string
    label: string
}

export const PERMISSION_SEP = ':'

export function buttonPermissionKey(page: string, action: string): string {
    return `${page}${PERMISSION_SEP}${action}`
}

export function isButtonPermissionKey(module: string): boolean {
    return module.includes(PERMISSION_SEP)
}

export function getPageFromPermissionKey(module: string): string {
    const idx = module.indexOf(PERMISSION_SEP)
    return idx === -1 ? module : module.slice(0, idx)
}

export const BTN = {
    view: {key: 'view', label: '访问页面'},
    search: {key: 'search', label: '查询'},
    create: {key: 'create', label: '新增'},
    edit: {key: 'edit', label: '编辑'},
    delete: {key: 'delete', label: '删除'},
    save: {key: 'save', label: '保存'},
    export: {key: 'export', label: '导出'},
} as const satisfies Record<string, PageButtonDef>

/** 列表 CRUD 默认按钮 */
export const DEFAULT_CRUD_BUTTONS: PageButtonDef[] = [
    BTN.view,
    BTN.search,
    BTN.create,
    BTN.edit,
    BTN.delete,
]

/** 只读列表（日志/流水） */
export const DEFAULT_READONLY_BUTTONS: PageButtonDef[] = [
    BTN.view,
    BTN.search,
    BTN.export,
]

/** 配置页（无表格 CRUD） */
export const DEFAULT_CONFIG_BUTTONS: PageButtonDef[] = [
    BTN.view,
    BTN.save,
]

/** 仅访问 */
export const DEFAULT_VIEW_BUTTONS: PageButtonDef[] = [BTN.view]

const READONLY_PAGES = new Set([
    'GoldCurrencyLogList',
    'DiamondCurrencyLogList',
    'LiveRevenueLogList',
    'AnchorIncomeSettlementLogList',
    'GuildIncomeSettlementLogList',
    'GuildProfileAnchorSettlementLogList',
    'LiveRecordList',
    'VideoCallLogList',
    'GameWinLogListManagement',
    'GameBetLogListManagement',
    'ResourceMonitor',
    'ServerLogExplorer',
    'ShortVideoWatchManagement',
])

/** 只读列表但不含导出 */
const READONLY_NO_EXPORT_PAGES = new Set([
    'ServerLogExplorer',
])

const CONFIG_PAGES = new Set([
    'AppTokenConfig',
    'AccountCfgManagement',
    'SimulatorCpuKeywordManagement',
    'ServerRuntimeCfgManagement',
    'PreloadCfgManagement',
    'TextModerationCfgManagement',
    'PrivacyPolicyCfgManagement',
    'GooglePlayCfgManagement',
    'UploadResourceCfgManagement',
    'DataSyncCfgManagement',
    'AgoraCfgManagement',
    'LiveCfgManagement',
    'RandomNicknameManagement',
    'CustomerServiceCfgManagement',
    'WalletExchangeCfgManagement',
    'LiveRevenueShareCfgManagement',
    'ShortVideoCfgManagement',
    'GamePlatformCfgManagement',
])

/** 各页面自定义按钮（未列出的页面按类型使用默认集） */
export const PAGE_BUTTON_OVERRIDES: Record<string, PageButtonDef[]> = {
    Dashboard: DEFAULT_VIEW_BUTTONS,
    UserList: [
        BTN.view,
        BTN.search,
        {key: 'viewDetail', label: '查看详情'},
        {key: 'viewAnchorDetail', label: '查看主播详情'},
        {key: 'setAnchor', label: '设为主播'},
        {key: 'setSeniorAnchor', label: '设为高级主播'},
        {key: 'goldAdd', label: '加金币'},
        {key: 'goldSub', label: '减金币'},
        {key: 'diamondAdd', label: '加钻石'},
        {key: 'diamondSub', label: '减钻石'},
        {key: 'ban', label: '封号/解封'},
        {key: 'rankOff', label: '下榜'},
        {key: 'rankOn', label: '上榜'},
        {key: 'rechargeWhitelistOn', label: '加入充值白名单'},
        {key: 'rechargeWhitelistOff', label: '移出充值白名单'},
        {key: 'cancel', label: '注销/取消注销'},
        {key: 'setUserType', label: '修改用户类型'},
        {key: 'setAnchorType', label: '设置主播类型'},
    ],
    AnchorListManagement: [BTN.view, BTN.search, BTN.edit, {key: 'viewDetail', label: '查看详情'}, {key: 'viewUserDetail', label: '查看用户详情'}, {key: 'viewGuildDetail', label: '查看工会详情'}, {key: 'offShelf', label: '下架'}],
    PlatformAnchorList: [
        BTN.view,
        BTN.search,
        {key: 'viewDetail', label: '查看详情'},
        {key: 'viewUserDetail', label: '查看用户详情'},
        {key: 'offShelf', label: '下架'},
        {key: 'ban', label: '封禁主播'},
        {key: 'unban', label: '解封主播'},
        {key: 'setAnchorType', label: '设置主播类型'},
    ],
    LiveRoomRecycleBinManagement: [BTN.view, BTN.search, {key: 'onShelf', label: '上架'}],
    RoleManagement: [
        BTN.view,
        BTN.search,
        BTN.create,
        BTN.edit,
        BTN.delete,
        {key: 'permission', label: '权限配置'},
    ],
    CMSUserManagement: [...DEFAULT_CRUD_BUTTONS],
    ModuleList: [BTN.view, BTN.save],
    BanUser: [BTN.view, BTN.save],
    GiftManagement: [
        BTN.view,
        BTN.create,
        BTN.edit,
        BTN.delete,
        {key: 'sort', label: '排序'},
        {key: 'sync', label: '同步数据'},
    ],
    GuildManagement: [
        BTN.view,
        BTN.search,
        BTN.create,
        BTN.edit,
        {key: 'offShelf', label: '下架'},
        {key: 'viewMembers', label: '查看成员'},
        {key: 'viewDetail', label: '查看详情'},
        {key: 'viewUserDetail', label: '查看用户详情'},
        {key: 'viewAnchorSettlementLogs', label: '查看名下主播结算流水'},
        {key: 'ban', label: '封禁主播'},
        {key: 'unban', label: '解封主播'},
        {key: 'exitGuild', label: '退出工会'},
        {key: 'setAnchorType', label: '设置主播类型'},
        {key: 'joinGuildAnchor', label: '加入工会'},
        {key: 'batchSetAnchor', label: '导入普通主播'},
        {key: 'batchSetSeniorAnchor', label: '导入高级主播'},
    ],
    AnchorDetail: [
        BTN.view,
        {key: 'dailyEffectiveLive', label: '查看每日流水'},
        {key: 'exportDailyEffectiveLive', label: '导出每日流水'},
        {key: 'liveRecord', label: '查看直播记录'},
        {key: 'exportLiveRecord', label: '导出直播记录'},
        {key: 'settlementLog', label: '查看结算流水'},
        {key: 'exportSettlementLog', label: '导出结算流水'},
    ],
    GuildDetail: [
        BTN.view,
        {key: 'viewUserDetail', label: '查看用户详情'},
        {key: 'dailyEffectiveLive', label: '查看工会每日流水'},
        {key: 'exportDailyEffectiveLive', label: '导出工会每日流水'},
        {key: 'incomeArchive', label: '查看下架归档'},
        {key: 'settlementLog', label: '查看结算流水'},
        {key: 'exportSettlementLog', label: '导出结算流水'},
        {key: 'anchorSettlementLog', label: '查看名下主播结算流水'},
        {key: 'exportAnchorSettlementLog', label: '导出名下主播结算流水'},
        {key: 'anchorDailyEffectiveLive', label: '查看名下主播每日流水'},
        {key: 'exportAnchorDailyEffectiveLive', label: '导出名下主播每日流水'},
    ],
    GuildRecycleBinManagement: [BTN.view, BTN.search, {key: 'onShelf', label: '上架'}],
    GuildProfileManagement: [
        ...DEFAULT_VIEW_BUTTONS,
        {key: 'viewUserDetail', label: '查看用户详情'},
        {key: 'viewAnchors', label: '查询名下主播'},
        {key: 'viewDetail', label: '查看详情'},
        {key: 'viewAnchorSettlementLogs', label: '查看名下主播结算流水'},
    ],
    BannerManagement: [
        ...DEFAULT_CRUD_BUTTONS,
        {key: 'sort', label: '排序'},
        {key: 'toggle', label: '上下架'},
        {key: 'sync', label: '同步数据'},
    ],
    ActivityMessageManagement: [
        ...DEFAULT_CRUD_BUTTONS,
        {key: 'publish', label: '发布'},
        {key: 'sync', label: '同步数据'},
    ],
    GameVendorGameListManagement: [
        BTN.view,
        BTN.search,
        BTN.edit,
        {key: 'shelf', label: '上下架'},
        {key: 'sync', label: '同步厂商'},
    ],
    GameShelfListManagement: [
        BTN.view,
        BTN.search,
        {key: 'shelf', label: '下架'},
    ],
    VipCfgManagement: [
        ...DEFAULT_CRUD_BUTTONS,
        {key: 'sync', label: '同步数据'},
    ],
    RechargeCfgManagement: [
        ...DEFAULT_CRUD_BUTTONS,
        {key: 'sync', label: '同步数据'},
    ],
    RechargeOrderList: [
        BTN.view,
        BTN.search,
        {key: 'manualCreateOrder', label: '人工创建订单'},
        {key: 'manualRecharge', label: '人工补单'},
    ],
    LiveRecordList: [
        BTN.view,
        BTN.search,
        BTN.export,
        {key: 'viewUserDetail', label: '查看用户详情'},
    ],
    LiveRevenueLogList: [
        BTN.view,
        BTN.search,
        BTN.export,
        {key: 'viewUserDetail', label: '查看用户详情'},
    ],
}

export function getPageButtons(pageName: string, metaButtons?: PageButtonDef[]): PageButtonDef[] {
    if (metaButtons?.length) {
        return metaButtons
    }
    if (PAGE_BUTTON_OVERRIDES[pageName]) {
        return PAGE_BUTTON_OVERRIDES[pageName]
    }
    if (READONLY_PAGES.has(pageName)) {
        if (READONLY_NO_EXPORT_PAGES.has(pageName)) {
            return [BTN.view, BTN.search]
        }
        return DEFAULT_READONLY_BUTTONS
    }
    if (CONFIG_PAGES.has(pageName)) {
        return DEFAULT_CONFIG_BUTTONS
    }
    return DEFAULT_CRUD_BUTTONS
}
