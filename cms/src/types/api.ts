// 通用响应类型
export interface ApiResponse<T = any> {
    code: number
    data: T
    message?: string
}

// 通用响应类型（包含null情况）
export interface ApiResponseWithNull<T = any> {
    message: string
    data: T
    code: number
}

// 登录响应类型
export interface LoginResponse {
    authId: string
    token: string
}

// 分页查询响应类型
export interface PageQuery {
    pageIndex: number
    pageSize: number
}

// 分页响应类型
export interface PageResponse<T = any> {
    total: number
    data: T[]
}

export interface SysStat {
    totalGold: number
    totalGoldConsume: number
    totalDiamondConsume: number
    totalRecharge: number
    totalWithdraw: number
    totalRegisterUser: string | number
    todayRecharge: number
    todayGoldConsume: number
    todayDiamondConsume: number
    todayRegisterUser: string | number
    onlineCount: string | number
}

export interface UserStatTrendPoint {
    time: string
    activeUserCount: number
    registerUserCount: number
    /** 柱形图指标, key 见 user-stat-bar-series.ts */
    barMetrics?: Record<string, number>
}

export interface UserStatTrend {
    daily: UserStatTrendPoint[]
    weekly: UserStatTrendPoint[]
    monthly: UserStatTrendPoint[]
}

export interface ResourceMetricPoint {
    time: string
    procMemMb: number
    procHeapAllocMb: number
    procHeapInuseMb: number
    procHeapSysMb: number
    procHeapUsedPercent: number
    procHeapIdlePercent: number
    procCpuPercent: number
    sysMemUsedMb: number
    sysMemTotalMb: number
    sysMemUsedPercent: number
    sysCpuPercent: number
    onlineCount: number | string
}

export interface ResourceMetricTrend {
    points: ResourceMetricPoint[]
}

export interface ResourceMetricTrendQuery {
    startTime?: string
    endTime?: string
    limit?: number
}

// 登录请求类型
export interface LoginReq {
    userName: string
    pwd: string
}

// 登录响应数据类型
export interface LoginRes {
    authId: number
    token: string
    admin: boolean
    modules: Permission[]
}

// 账号相关类型
export interface BanReq {
    accountId: string
    openId: string
    channel: number
    banApplyTime?: string
}

export interface BanAnchorReq {
    accountId: string
    openId?: string
    channel?: number
    banApplyTime: string
    banReason: string
}

export interface UnBanAnchorReq {
    accountId: string
    openId?: string
    channel?: number
}

export interface UnBanReq {
    accountId: string
    openId: string
    channel: number
}

export interface CancelReq {
    accountId: string
    openId: string
    channel: number
}

export interface UnCancelReq {
    accountId: string
    openId: string
    channel: number
}

export interface QueryUserInfoReq extends PageQuery {
    key?: string
    startTime?: string
    endTime?: string
    rechargeWhitelist?: number
    isAnchor?: number
}

export interface UserInfo {
    id: string
    createdAt?: string | null
    openId: string
    ip: string
    registerIp?: string
    registerCountry?: string
    loginCountry?: string
    channel: number
    ban: boolean
    banTime?: string | null
    banApplyTime?: string | null
    cancel: boolean
    phoneAreaCode?: string
    // 来自 user_infos 表(LEFT JOIN，可能为空)
    nickname?: string
    phone?: string
    avatar?: string
    remark?: string
    gold?: number
    diamond?: number
    shareCode?: string
    guildId?: string | number
    userType?: number
    isAnchor?: boolean
    vipLevel?: number
    lastLoginTime?: string | null
    deviceType?: string
    packageName?: string
    appVersion?: string
    canRank?: boolean
    rechargeWhitelist?: boolean
}

export interface SetAnchorReq {
    accountId: string
}

export interface SetSeniorAnchorReq {
    accountId: string
}

export interface BatchSetAnchorReq {
    ids: string[]
}

export interface BatchSetSeniorAnchorReq {
    ids: string[]
}

export interface BatchSetAnchorRes {
    successCount: number
    failCount: number
    failIds?: string[]
}

export interface ExitGuildReq {
    anchorId: string
}

export interface ExitGuildRes {
    success: boolean
}

export interface SetLiveRoomStatusReq {
    anchorId: string
    /** 0=下架, 1=上架 */
    status: number
}

export interface SetLiveRoomStatusRes {
    success: boolean
    status: number
}

export interface SetUserTypeReq {
    accountId: string
    userType: number
}

export interface SetCanRankReq {
    accountId: string
    canRank: boolean
}

export interface SetRechargeWhitelistReq {
    accountId: string
    rechargeWhitelist: boolean
}

export interface QueryAnchorListReq extends PageQuery {
    key?: string
    guildId?: string | number
    platformOnly?: boolean
    guildOnly?: boolean
    liveStatus?: number
}

export interface QueryOffShelfLiveRoomListReq extends PageQuery {
    key?: string
}

export interface OffShelfLiveRoomItem {
    id: string
    nickname?: string
    phone?: string
    avatar?: string
    userType?: number
    guildId?: string | number
    roomTitle?: string
    roomId?: string
    category?: number
    ban?: boolean
    banApplyTime?: string | null
    banReason?: string
    status?: number
    updatedAt?: string | null
    createdAt?: string | null
}

export interface AnchorListItem {
    id: string
    nickname?: string
    phone?: string
    avatar?: string
    /** 1=普通主播, 7=高级主播 */
    userType?: number
    guildId?: string | number
    guildName?: string
    ip?: string
    roomTitle?: string
    roomId?: string
    category?: number
    privateInviteType?: number
    ticket?: number
    billing?: number
    liveStatus?: number
    totalIncome?: number
    totalGiftIncome?: number
    totalPaidDanmakuIncome?: number
    totalPrivateRoomTicketIncome?: number
    totalPrivateRoomWatchIncome?: number
    totalVideoCallIncome?: number
    totalVideoCallTicketIncome?: number
    totalVideoCallBillingIncome?: number
    ban?: boolean
    banApplyTime?: string | null
    banReason?: string
    /** 0=下架, 1=上架 */
    status?: number
    createdAt?: string | null
    registeredAt?: string | null
}

export interface LiveRoomIncomeAmounts {
    totalIncome?: number
    totalGiftIncome?: number
    totalPaidDanmakuIncome?: number
    totalPrivateRoomTicketIncome?: number
    totalPrivateRoomWatchIncome?: number
    totalVideoCallIncome?: number
    totalVideoCallTicketIncome?: number
    totalVideoCallBillingIncome?: number
    totalLiveDuration?: number
}

export interface AnchorLiveRoomDetail {
    id: string
    guildId?: string
    title?: string
    cover?: string
    notice?: string
    liveRecordId?: string
    heartTime?: string | null
    ban?: boolean
    banApplyTime?: string | null
    banReason?: string
    status?: number
    liveStatus?: number
    category?: number
    privateInviteType?: number
    ticket?: number
    billing?: number
    createdAt?: string | null
    updatedAt?: string | null
}

export interface LiveRoomIncomeUnsettledDetail extends LiveRoomIncomeAmounts {
    updatedAt?: string | null
}

export interface LiveRoomIncomeSettledDetail extends LiveRoomIncomeAmounts {
    settlementSalary?: number
    settlementShareAmount?: number
    updatedAt?: string | null
}

export interface LiveRoomIncomeTotalDetail extends LiveRoomIncomeAmounts {
    settlementSalary?: number
    settlementShareAmount?: number
    updatedAt?: string | null
}

export interface LiveRoomIncomeArchiveItem extends LiveRoomIncomeAmounts {
    id: string
    roomId?: string
    guildId?: string
    settlementSalary?: number
    createdAt?: string | null
}

export interface AnchorDetail {
    anchor?: AnchorListItem
    liveRoom?: AnchorLiveRoomDetail
    incomeUnsettled?: LiveRoomIncomeUnsettledDetail | null
    incomeSettled?: LiveRoomIncomeSettledDetail | null
    incomeTotal?: LiveRoomIncomeTotalDetail | null
    incomeArchives?: LiveRoomIncomeArchiveItem[]
}

export interface GuildIncomeArchiveItem extends LiveRoomIncomeAmounts {
    id: string
    guildId?: string
    createdAt?: string | null
}

export interface GuildDetailIncome {
    incomeUnsettled?: LiveRoomIncomeUnsettledDetail | null
    incomeSettled?: LiveRoomIncomeSettledDetail | null
    incomeTotal?: LiveRoomIncomeTotalDetail | null
}

export interface GuildIncomeArchivesRes {
    list: GuildIncomeArchiveItem[]
}

export interface AnchorDailyEffectiveLiveQuery extends PageQuery {
    anchorId: string | number
    liveDateStart?: string
    liveDateEnd?: string
    settled?: number
}

export interface AnchorDailyEffectiveLiveItem extends LiveRoomIncomeAmounts {
    id: string
    roomId?: string
    liveDate?: string
    liveDuration?: number
    settled?: boolean
    createdAt?: string | null
    updatedAt?: string | null
}

export interface GuildDailyEffectiveLiveQuery extends PageQuery {
    guildId: string | number
    liveDateStart?: string
    liveDateEnd?: string
    settled?: number
}

export interface GuildDailyEffectiveLiveItem extends LiveRoomIncomeAmounts {
    id: string
    guildId?: string
    liveDate?: string
    liveDuration?: number
    settled?: boolean
    createdAt?: string | null
    updatedAt?: string | null
}

export interface GuildAnchorDailyEffectiveLiveQuery extends PageQuery {
    guildId: string | number
    roomId?: string
    liveDateStart?: string
    liveDateEnd?: string
    settled?: number
}

export interface GuildAnchorDailyEffectiveLiveItem extends LiveRoomIncomeAmounts {
    id: string
    roomId?: string
    roomNickname?: string
    liveDate?: string
    liveDuration?: number
    settled?: boolean
    createdAt?: string | null
    updatedAt?: string | null
}

export interface UserAccountDetailItem {
    id: string
    openId?: string
    ip?: string
    registerIp?: string
    registerCountry?: string
    loginCountry?: string
    channel?: number
    phoneAreaCode?: string
    ban?: boolean
    banTime?: string | null
    banApplyTime?: string | null
    cancel?: boolean
    createdAt?: string | null
}

export interface UserProfileDetailItem {
    nickname?: string
    phone?: string
    avatar?: string
    remark?: string
    shareCode?: string
    userType?: number
    isAnchor?: boolean
    inviterId?: string
    vipLevel?: number
    lastLoginTime?: string | null
    liveRoomId?: string
    liveRoomVer?: string
    gender?: number
    birthday?: string | null
    botAnchorStatus?: number
    guildId?: string
    updatedAt?: string | null
}

export interface UserWalletDetailItem {
    gold?: number
    diamond?: number
}

export interface UserExtDetailItem {
    canRank?: boolean
    prettyId?: string
    packageName?: string
    appVersion?: string
    followCount?: number
    followerCount?: number
    cancelCode?: string
    cancelCodeExpireAt?: string | null
    rechargeWhitelist?: boolean
    updatedAt?: string | null
}

export interface UserCumulativeStatDetailItem {
    totalRecharge?: number
    totalWithdraw?: number
    totalPayCount?: number
    totalDiamondConsume?: number
    totalGoldConsume?: number
    updatedAt?: string | null
}

export interface UserLoginDeviceDetailItem {
    deviceType?: string
    deviceModel?: string
    cpuModel?: string
    osVersion?: string
    appVersion?: string
    deviceId?: string
    updatedAt?: string | null
}

export interface UserDetail {
    account?: UserAccountDetailItem | null
    profile?: UserProfileDetailItem | null
    wallet?: UserWalletDetailItem | null
    userExt?: UserExtDetailItem | null
    cumulativeStat?: UserCumulativeStatDetailItem | null
    loginDevice?: UserLoginDeviceDetailItem | null
}

export interface QueryBotAnchorListReq extends PageQuery {
    key?: string
}

export interface BotAnchorListItem {
    id: string
    nickname?: string
    avatar?: string
    guildId?: string | number
    roomId?: string
    roomTitle?: string
    category?: number
    tagId?: string | number
    tagName?: string
    cloudPlayerVideo?: string
    cloudPlayerVideoFile?: string
    pushStream?: boolean
    isTest?: boolean
    botAnchorStatus?: number
    liveStatus?: number
    createdAt?: string | null
    updatedAt?: string | null
}

export interface CreateBotAnchorReq {
    nickname: string
    avatar?: string
    guildId?: string | number
    roomTitle?: string
    category: number
    tagId?: string | number
    cloudPlayerVideo?: string
    pushStream?: boolean
    isTest?: boolean
}

export interface UpdateBotAnchorReq {
    id: string
    nickname: string
    avatar?: string
    roomTitle?: string
    category: number
    tagId?: string | number
    cloudPlayerVideo?: string
    pushStream?: boolean
    isTest?: boolean
}

export interface SetBotAnchorStatusReq {
    id: string
    status: number
}

export interface StartBotAnchorLiveReq {
    id: string
}

export interface StopBotAnchorLiveReq {
    id: string
}

export interface BatchStartBotAnchorLiveReq {
    ids: string[]
}

export interface BatchStopBotAnchorLiveReq {
    ids: string[]
}

export interface BatchBotAnchorLiveRes {
    successCount: number
    failCount: number
    failIds?: string[]
}

export interface CurrencyLogQuery extends PageQuery {
    userId?: string
    currencyType: number
}

export interface CurrencyLogItem {
    id: string
    userId: string
    nickname?: string
    action: number
    amount: number
    before: number
    after: number
    reason: number
    reasonText?: string
    createdAt?: string | null
}

export interface GameWinLogQuery extends PageQuery {
    userId?: string
    gameCode?: string
    orderId?: string
    platformType?: string
}

export interface GameWinLogItem {
    id: string
    userId: string
    nickname?: string
    gameCode: string
    nameEn?: string
    cover?: string
    amount: number
    platformType: string
    orderId: string
    createdAt?: string | null
}

export interface GameBetLogQuery extends PageQuery {
    userId?: string
    gameCode?: string
    orderId?: string
    platformType?: string
}

export interface GameBetLogItem {
    id: string
    userId: string
    nickname?: string
    gameCode: string
    nameEn?: string
    cover?: string
    amount: number
    platformType: string
    orderId: string
    liveRoomId?: string
    liveRecordId?: string
    liveRoomTitle?: string
    anchorNickname?: string
    createdAt?: string | null
}


// App Token相关类型
export interface AppToken {
    id: string
    token: string
    expireAt?: string | null
    expired?: boolean
}

export interface GetAppTokenReq {
    userId?: string
    pageIndex?: number
    pageSize?: number
}

export interface SaveAppTokenReq {
    id: string
    token?: string
    expireAt?: string | null
}

export interface AccountCfg {
    id?: string
    cancelAccountByCodeEnabled: boolean
    blockSimulatorLogin: boolean
    createdAt?: string
    updatedAt?: string
}

export interface SimulatorCpuKeyword {
    id: string
    keyword: string
    remark: string
    createdAt: string
    updatedAt: string
}

export interface SimulatorCpuKeywordQuery extends PageQuery {
    key?: string
}

export interface GetAccountCfgRes {
    cfg?: AccountCfg | null
}

export interface SaveAccountCfgReq {
    id?: string | number
    cancelAccountByCodeEnabled: boolean
    blockSimulatorLogin: boolean
}

export interface SaveAccountCfgRes {
    success: boolean
    id?: string
}

export interface FirstRechargePrivilegeItem {
    icon?: string
    iconName?: string
    descEn?: string
    descEs?: string
    descPt?: string
    descHi?: string
    descId?: string
}

export interface FirstRechargeActivityCfg {
    id?: string
    enabled: boolean
    icon?: string
    iconName?: string
    titleEn?: string
    titleEs?: string
    titlePt?: string
    titleHi?: string
    titleId?: string
    rechargeBtnTextEn?: string
    rechargeBtnTextEs?: string
    rechargeBtnTextPt?: string
    rechargeBtnTextHi?: string
    rechargeBtnTextId?: string
    privileges?: FirstRechargePrivilegeItem[]
    createdAt?: string
    updatedAt?: string
}

export interface GetFirstRechargeActivityCfgRes {
    cfg?: FirstRechargeActivityCfg | null
}

export interface SaveFirstRechargeActivityCfgReq {
    id?: string | number
    enabled: boolean
    icon?: string
    titleEn?: string
    titleEs?: string
    titlePt?: string
    titleHi?: string
    titleId?: string
    rechargeBtnTextEn?: string
    rechargeBtnTextEs?: string
    rechargeBtnTextPt?: string
    rechargeBtnTextHi?: string
    rechargeBtnTextId?: string
    privileges?: FirstRechargePrivilegeItem[]
}

export interface SaveFirstRechargeActivityCfgRes {
    success: boolean
    id?: string
}

export interface PreloadCfg {
    id?: string
    recentLoginLimit: number
    hotRestartAuth?: string
    memoryLimitM?: number
    ipGeoDbPath?: string
    createdAt?: string
    updatedAt?: string
}

export interface GetPreloadCfgRes {
    cfg?: PreloadCfg | null
}

export interface SavePreloadCfgReq {
    id?: string | number
    recentLoginLimit: number
    hotRestartAuth: string
    memoryLimitM: number
    ipGeoDbPath: string
}

export interface SavePreloadCfgRes {
    success: boolean
    id?: string
}

// 角色相关类型
export interface Role {
    id: string
    name: string
    description: string
    status: number
    createdAt: string
    updatedAt: string
}

// 权限相关类型
export interface Permission {
    id: string
    module: string
    roleId: number
    apiPath?: string
    createdAt?: string
    updatedAt?: string
}

// 直播工会相关类型
export interface Guild {
    id: string
    name: string
    leaderId: string
    leaderName?: string
    description: string
    status: number
    createdAt: string
    updatedAt: string
}

export interface GuildQuery extends PageQuery {
    name?: string
}

export interface MyGuildProfile {
    id: string
    name: string
    description: string
    updatedAt?: string
}

export interface UpdateMyGuildProfileReq {
    id: string | number
    name: string
    description?: string
}

export interface ImportGuildAnchorRow {
    userId: string
}

export interface SetGuildAnchorTypeReq {
    guildId: string | number
    userId: string | number
    /** 1 普通主播, 7 高级主播 */
    anchorType: 1 | 7
}

export interface SetPlatformAnchorTypeReq {
    userId: string | number
    /** 1 普通主播, 7 高级主播 */
    anchorType: 1 | 7
}

export interface SetPlatformAnchorTypeRes {
    success: boolean
}

export interface SetGuildAnchorTypeRes {
    success: boolean
}

export interface JoinGuildAnchorReq {
    guildId: string | number
    userId: string | number
    /** 1 普通主播, 7 高级主播 */
    anchorType: 1 | 7
}

export interface JoinGuildAnchorRes {
    success: boolean
    reason?: number
    nickname?: string
}

export interface ImportGuildAnchorsReq {
    guildId: string | number
    /** 1 普通主播, 7 高级主播 */
    anchorType: 1 | 7
    rows: ImportGuildAnchorRow[]
}

export interface ImportGuildAnchorFailItem {
    userId: string
    nickname: string
    /** 1用户不存在 2注销码错误(废弃) 3注销码过期(废弃) 4已加入工会 5无法设为主播 6主播间缓存已存在 */
    reason: number
}

export interface ImportGuildAnchorsRes {
    successCount: number
    failCount: number
    fails: ImportGuildAnchorFailItem[]
}

export interface GuildAnchorImportResultState {
    guildId: string
    guildName: string
    anchorType: 1 | 7
    successCount: number
    failCount: number
    fails: ImportGuildAnchorFailItem[]
}

export interface MyGuildProfileListRes {
    list: MyGuildProfile[]
}

// 礼物相关类型
export interface Gift {
    id: string
    name: string
    nameEn: string
    nameEs: string
    namePt: string
    nameHi: string
    nameId: string
    icon: string
    iconName: string
    animation: string
    animationName: string
    price: number
    category: string
    sort: number
    status: number
    publishedAt?: string | null
    description: string
    createdAt: string
    updatedAt: string
}

export interface GiftQuery extends PageQuery {
    name?: string
    category?: string
    statusFilter?: number
}

// 充值配置
export interface RechargeCfg {
    id: string
    name: string
    packageName: string
    cfgType: number
    icon: string
    iconName: string
    /** 基础到账金币数(接口字段仍为 diamond) */
    diamond: number
    /** 额外赠送金币数(接口字段仍为 extraDiamond) */
    extraDiamond: number
    price: number
    /** 固定 USD，仅列表展示用 */
    currency?: string
    productId: string
    sort: number
    status: number
    description: string
    createdAt: string
    updatedAt: string
}

export interface RechargeCfgQuery extends PageQuery {
    name?: string
    packageName?: string
    typeFilter?: number
    statusFilter?: number
}

// VIP配置
export interface VipCfg {
    id: string
    level: number
    levelName: string
    levelIcon?: string
    levelIconName?: string
    animationSwitch?: number
    commentEffectSwitch?: number
    customerServiceSwitch?: number
    upgradeRechargeLimit: number
    animation?: string
    animationName?: string
    animationIcon?: string
    animationIconName?: string
    animationTitleEn?: string
    animationTitleEs?: string
    animationTitlePt?: string
    animationTitleHi?: string
    animationTitleId?: string
    animationDescEn?: string
    animationDescEs?: string
    animationDescPt?: string
    animationDescHi?: string
    animationDescId?: string
    commentEffect?: string
    commentEffectName?: string
    commentEffectIcon?: string
    commentEffectIconName?: string
    commentEffectTitleEn?: string
    commentEffectTitleEs?: string
    commentEffectTitlePt?: string
    commentEffectTitleHi?: string
    commentEffectTitleId?: string
    commentEffectDescEn?: string
    commentEffectDescEs?: string
    commentEffectDescPt?: string
    commentEffectDescHi?: string
    commentEffectDescId?: string
    customerServiceIcon?: string
    customerServiceIconName?: string
    customerServiceTitleEn?: string
    customerServiceTitleEs?: string
    customerServiceTitlePt?: string
    customerServiceTitleHi?: string
    customerServiceTitleId?: string
    customerServiceDescEn?: string
    customerServiceDescEs?: string
    customerServiceDescPt?: string
    customerServiceDescHi?: string
    customerServiceDescId?: string
    createdAt: string
    updatedAt: string
}

export interface VipCfgQuery extends PageQuery {
    levelName?: string
}

// App包管理
export interface AppPkg {
    id: string
    packageName: string
    secretKey: string
    privacyPolicyUrl?: string
    termsOfServiceUrl?: string
    remark: string
    createdAt: string
    updatedAt: string
}

export interface AppPkgQuery extends PageQuery {
    packageName?: string
}

// 首页 Banner
export interface Banner {
    id: string
    title: string
    image: string
    imageName: string
    link: string
    scene: number
    direction: number
    sort: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface BannerQuery extends PageQuery {
    title?: string
    sceneFilter?: number
    statusFilter?: number
}

export interface ActivityMessage {
    id: string
    iconEn: string
    iconEnName: string
    iconEs: string
    iconEsName: string
    iconPt: string
    iconPtName: string
    iconHi: string
    iconHiName: string
    iconId: string
    iconIdName: string
    bgEn: string
    bgEnName: string
    bgEs: string
    bgEsName: string
    bgPt: string
    bgPtName: string
    bgHi: string
    bgHiName: string
    bgId: string
    bgIdName: string
    titleEn: string
    titleEs: string
    titlePt: string
    titleHi: string
    titleId: string
    contentEn: string
    contentEs: string
    contentPt: string
    contentHi: string
    contentId: string
    urlEn: string
    urlEs: string
    urlPt: string
    urlHi: string
    urlId: string
    status: number
    publishedAt: string
    createdAt: string
    updatedAt: string
}

export interface ActivityMessageQuery extends PageQuery {
    title?: string
    statusFilter?: number
}

export interface Ticket {
    id: string
    price: number
    sort: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface TicketQuery extends PageQuery {
    statusFilter?: number
}

export interface PrivateRoomBilling {
    id: string
    pricePerMinute: number
    sort: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface PrivateRoomBillingQuery extends PageQuery {
    statusFilter?: number
}

export interface ShortVideo {
    id: string
    title: string
    video: string
    videoName: string
    cover: string
    coverName: string
    sort: number
    status: number
    isPaid: number
    payDiamond: number
    categoryId: number
    source: number
    authorId: string
    authorNickname: string
    authorType?: number
    duration: number
    freeWatchSeconds: number
    likeCount: number
    viewCount: number
    watchCount: number
    totalDiamondIncome: number
    createdAt: string
    updatedAt: string
}

export interface ShortVideoQuery extends PageQuery {
    title?: string
    authorNickname?: string
    statusFilter?: number
    sortField?: '' | 'viewCount' | 'totalDiamondIncome'
}

export interface ShortVideoWatchRecord {
    id: string
    userId: string
    nickname: string
    videoId: string
    videoTitle: string
    paidTime: string
    createdAt: string
    updatedAt: string
}

export interface ShortVideoWatchQuery extends PageQuery {
    userId?: string
    startTime?: number
    endTime?: number
}

export interface CreateShortVideoReq {
    video: string
    coverName?: string
    title: string
    sort?: number
    isPaid: number
    payDiamond?: number
    categoryId?: number
    source: number
    duration: number
    freeWatchSeconds?: number
    authorNickname?: string
}

export interface CreateShortVideoRes {
    id: string
    video: string
    cover?: string
    authorId?: string
}

export interface ShortVideoCfg {
    id: string
    maxFileSize: number
    maxCoverFileSize: number
    maxDuration: number
    freeWatchSeconds: number
    entryEnabled: number
    anchorDailyUploadLimit: number
    normalUserDailyUploadLimit: number
    createdAt: string
    updatedAt: string
}

export interface GetShortVideoCfgRes {
    cfg: ShortVideoCfg | null
}

export interface SaveShortVideoCfgReq {
    id?: string
    maxFileSize: number
    maxCoverFileSize: number
    maxDuration: number
    freeWatchSeconds: number
    entryEnabled: number
    anchorDailyUploadLimit: number
    normalUserDailyUploadLimit: number
}

export interface SaveShortVideoCfgRes {
    success: boolean
    id: string
}

export interface ShortVideoStorageStat {
    totalCount: number
    imageDirPath: string
    imageDirUsedBytes: number
    diskTotalBytes: number
    diskFreeBytes: number
    diskFreeRatio: number
}

export interface ShortVideoCategory {
    id: string
    name: string
    sort: number
    createdAt: string
    updatedAt: string
}

export interface ShortVideoCategoryQuery extends PageQuery {
}

export interface ShortVideoPriceTier {
    id: string
    price: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface ShortVideoPriceTierQuery extends PageQuery {
    statusFilter?: number
}

export interface LiveRoomTag {
    id: string
    name: string
    sort: number
    isSpecial: boolean
    createdAt: string
    updatedAt: string
}

export interface LiveRoomTagQuery extends PageQuery {
}

export interface AnchorSalaryCfg {
    id: string
    weeklyWorkDays: number
    dailyLiveDurationMinutes: number
    salaryAmount: number
    sort: number
    createdAt: string
    updatedAt: string
}

export interface AnchorSalaryCfgQuery extends PageQuery {
}

export interface IncomeSettlementLogAmounts {
    totalIncome: number
    totalGiftIncome: number
    totalPaidDanmakuIncome: number
    totalPrivateRoomTicketIncome: number
    totalPrivateRoomWatchIncome: number
    totalVideoCallIncome: number
    totalVideoCallTicketIncome: number
    totalVideoCallBillingIncome: number
    totalLiveDuration: number
    settlementSalary: number
    settlementShareAmount?: number
    anchorSharePercent?: number
    guildSharePercent?: number
}

export interface AnchorIncomeSettlementLogQuery extends PageQuery {
    roomId?: string
    startTime?: number
    endTime?: number
}

export interface AnchorIncomeSettlementLogItem extends IncomeSettlementLogAmounts {
    id: string
    roomId: string
    roomNickname?: string
    guildId?: string
    guildName?: string
    createdAt?: string | null
}

export interface MyGuildAnchorIncomeSettlementLogQuery extends PageQuery {
    guildId?: string
    roomId?: string
    startTime?: number
    endTime?: number
}

export interface GuildIncomeSettlementLogQuery extends PageQuery {
    guildId?: string
    startTime?: number
    endTime?: number
}

export interface GuildIncomeSettlementLogItem extends IncomeSettlementLogAmounts {
    id: string
    guildId: string
    guildName?: string
    createdAt?: string | null
}

export interface AgoraCfg {
    id: string
    appId: string
    appCertificate: string
    restCustomerId: string
    restCustomerSecret: string
    cloudPlayerRegion?: string
    tokenExpireSeconds: number
    tokenRefreshSeconds: number
    createdAt: string
    updatedAt: string
}

export interface GetAgoraCfgRes {
    cfg: AgoraCfg | null
}

export interface SaveAgoraCfgReq {
    id?: number
    appId: string
    appCertificate: string
    restCustomerId: string
    restCustomerSecret: string
    cloudPlayerRegion?: string
    tokenExpireSeconds: number
    tokenRefreshSeconds: number
}

export interface SaveAgoraCfgRes {
    success: boolean
    id: string
}

export interface LiveCfg {
    id: string
    paidDanmakuPrice: number
    privateRoomFreeWatchSeconds: number
    createdAt: string
    updatedAt: string
}

export interface GetLiveCfgRes {
    cfg: LiveCfg | null
}

export interface SaveLiveCfgReq {
    id?: number
    paidDanmakuPrice: number
    privateRoomFreeWatchSeconds: number
}

export interface SaveLiveCfgRes {
    success: boolean
    id: string
}

export interface TextModerationCfg {
    id: string
    enabled: boolean
    accessKeyId: string
    accessKeySecret: string
    regionId: string
    endpoint: string
    chatService: string
    nicknameService: string
    commentService: string
    createdAt: string
    updatedAt: string
}

export interface GetTextModerationCfgRes {
    cfg: TextModerationCfg | null
}

export interface SaveTextModerationCfgReq {
    id?: number
    enabled: boolean
    accessKeyId: string
    accessKeySecret: string
    regionId: string
    endpoint: string
    chatService: string
    nicknameService: string
    commentService: string
}

export interface SaveTextModerationCfgRes {
    success: boolean
    id: string
}

export interface PrivacyPolicyCfg {
    id: string
    apiBase: string
    privacyPolicyUrl: string
    termsOfServiceUrl: string
    creatorTermsUrl: string
    roomOwnerTermsUrl: string
    vipDescUrl: string
    aboutSiteUrl: string
    safetyCenterUrl: string
    createdAt: string
    updatedAt: string
}

export interface GetPrivacyPolicyCfgRes {
    cfg: PrivacyPolicyCfg | null
}

export interface SavePrivacyPolicyCfgReq {
    id?: number
    apiBase: string
    privacyPolicyUrl: string
    termsOfServiceUrl: string
    creatorTermsUrl: string
    roomOwnerTermsUrl: string
    vipDescUrl: string
    aboutSiteUrl: string
    safetyCenterUrl: string
}

export interface SavePrivacyPolicyCfgRes {
    success: boolean
    id: string
}

export interface CustomerServiceCfg {
    id: string
    telegramUrl: string
    facebookUrl: string
    whatsappUrl: string
    createdAt: string
    updatedAt: string
}

export interface GetCustomerServiceCfgRes {
    cfg: CustomerServiceCfg | null
}

export interface SaveCustomerServiceCfgReq {
    id?: number
    telegramUrl: string
    facebookUrl: string
    whatsappUrl: string
}

export interface SaveCustomerServiceCfgRes {
    success: boolean
    id: string
}

export interface WalletExchangeCfg {
    id: string
    goldToDiamondRate: number
    exchangeFeePercent: number
    createdAt: string
    updatedAt: string
}

export interface GetWalletExchangeCfgRes {
    cfg: WalletExchangeCfg | null
}

export interface SaveWalletExchangeCfgReq {
    id?: number
    goldToDiamondRate: number
    exchangeFeePercent: number
}

export interface SaveWalletExchangeCfgRes {
    success: boolean
    id: string
}

export interface LiveRevenueShareCfg {
    id: string
    anchorSharePercent: number
    guildSharePercent: number
    createdAt: string
    updatedAt: string
}

export interface GetLiveRevenueShareCfgRes {
    cfg: LiveRevenueShareCfg | null
}

export interface SaveLiveRevenueShareCfgReq {
    id?: number
    anchorSharePercent: number
    guildSharePercent: number
}

export interface SaveLiveRevenueShareCfgRes {
    success: boolean
    id: string
}

export interface GooglePlayCfg {
    id: string
    enabled: boolean
    packageName: string
    serviceAccountJson: string
    createdAt: string
    updatedAt: string
}

export interface GetGooglePlayCfgRes {
    cfg: GooglePlayCfg | null
}

export interface SaveGooglePlayCfgReq {
    id?: number
    enabled: boolean
    serviceAccountJson: string
}

export interface SaveGooglePlayCfgRes {
    success: boolean
    id: string
}

export interface DataSyncCfg {
    id: string
    targetApiBase: string
    token: string
    createdAt: string
    updatedAt: string
}

export interface GetDataSyncCfgRes {
    cfg: DataSyncCfg | null
}

export interface SaveDataSyncCfgReq {
    id?: number
    targetApiBase: string
    token: string
}

export interface SaveDataSyncCfgRes {
    success: boolean
    id: string
}

export interface SyncVipCfgReq {
    ids: number[]
}

export interface SyncVipCfgRes {
    success: boolean
    rowCount: number
    fileCount: number
    message: string
}

export interface SyncActivityMessageReq {
    ids: number[]
}

export interface SyncActivityMessageRes {
    success: boolean
    rowCount: number
    fileCount: number
    message: string
}

export interface SyncIdsReq {
    ids: number[]
}

export interface SyncBatchRes {
    success: boolean
    rowCount: number
    fileCount: number
    message: string
}

export interface SyncFirstRechargeActivityCfgRes {
    success: boolean
    rowCount: number
    fileCount: number
    message: string
}

export interface GamePlatformCfg {
    id: string
    vendorUrl: string
    token: string
    secretKey: string
    iconUrl: string
    createdAt: string
    updatedAt: string
}

export interface GetGamePlatformCfgRes {
    cfg: GamePlatformCfg | null
}

export interface SaveGamePlatformCfgReq {
    id?: number
    vendorUrl: string
    token: string
    secretKey: string
    iconUrl?: string
}

export interface SaveGamePlatformCfgRes {
    success: boolean
    id: string
}

export interface VendorGame {
    gameCode: string
    name: string
    nameEn: string
    category: string
    cover: string
    platform: string
    onShelf?: boolean
}

export interface VendorGameQuery extends PageQuery {
    gameCode?: string
    name?: string
    platform?: string
    category?: string
}

export interface GameShelfItem {
    id: string
    gameCode: string
    name: string
    nameEn: string
    cover: string
    liveGameName: string
    liveGameCover: string
    liveGameCoverUrl: string
    platform: string
}

export interface UpdateGameShelfReq {
    gameCode: string
    liveGameName?: string
    liveGameCover?: string
}

export interface UpdateGameShelfRes {
    success: boolean
}

export interface GetMultiplayerConfigUrlReq {
    gameCode: string
    platform?: string
}

export interface GetMultiplayerConfigUrlRes {
    configUrl: string
    expireInMs: number
}

export interface CMSGameStartLinkReq {
    userId: string
    gameCode: string
    platform?: string
}

export interface CMSGameStartLinkRes {
    link: string
}

export interface GameShelfQuery extends PageQuery {
    gameCode?: string
    name?: string
    platform?: string
}

export interface BatchAddGameShelfItem {
    gameCode: string
    platform: string
}

export interface BatchAddGameShelfReq {
    items: BatchAddGameShelfItem[]
}

export interface BatchAddGameShelfRes {
    success: boolean
    successCount: number
    skipCount: number
}

export interface BatchDeleteGameShelfReq {
    gameCodes: string[]
}

export interface BatchDeleteGameShelfRes {
    success: boolean
    successCount: number
}

export interface DeleteGameShelfReq {
    id?: number
    gameCode?: string
}

export interface DeleteGameShelfRes {
    success: boolean
}

export interface AddGameShelfReq {
    gameCode: string
    platform: string
}

export interface AddGameShelfRes {
    success: boolean
    id: string
}

export interface ReloadVendorGameCacheRes {
    success: boolean
    count: number
}

export interface UploadResourceCfg {
    id: string
    resourceDomain: string
    appImageMaxSizeMB: number
    imageModerationEnabled: boolean
    imageModerationAccessKeyId: string
    imageModerationAccessKeySecret: string
    imageModerationRegionId: string
    imageModerationEndpoint: string
    imageModerationService: string
    createdAt: string
    updatedAt: string
}

export interface GetUploadResourceCfgRes {
    cfg: UploadResourceCfg | null
}

export interface SaveUploadResourceCfgReq {
    id?: number
    resourceDomain: string
    appImageMaxSizeMB: number
    imageModerationEnabled: boolean
    imageModerationAccessKeyId: string
    imageModerationAccessKeySecret: string
    imageModerationRegionId: string
    imageModerationEndpoint: string
    imageModerationService: string
}

export interface SaveUploadResourceCfgRes {
    success: boolean
    id: string
}

// 充值订单
export interface RechargeOrder {
    id: string
    userId: string
    nickname: string
    cfgId: string
    price: number
    currency: string
    gold: number
    status: number
    source: number
    payChannel: number
    thirdOrderId: string
    remark: string
    operatorId: string
    createdAt: number
    paidAt: number
}

export interface RechargeOrderQuery extends PageQuery {
    orderId?: string
    userId?: string
    statusFilter?: number
    source?: number
    startTime?: number
    endTime?: number
}

export interface LiveRevenueLogQuery extends PageQuery {
    receiverId?: string
    receiverIds?: string[]
    revenueType?: number
    startTime?: number
    endTime?: number
}

export interface LiveRevenueLogItem {
    id: string
    revenueType: number
    revenueTypeText?: string
    roomId: string
    liveRecordId: string
    senderId: string
    senderNickname?: string
    receiverId: string
    receiverNickname?: string
    receiverAvatar?: string
    bizId: string
    bizName?: string
    count: number
    unitPrice: number
    totalAmount: number
    status?: number
    statusText?: string
    createdAt?: string | null
}

export interface LiveRecordQuery extends PageQuery {
    anchorId?: string
    anchorIds?: string[]
    startTime?: number
    endTime?: number
}

export interface LiveRecordItem {
    id: string
    anchorId: string
    nickname?: string
    avatar?: string
    startTime?: string | null
    endTime?: string | null
    totalAudience: number
    totalLiveDuration: number
    totalIncome: number
    totalGiftIncome: number
    totalPaidDanmakuIncome: number
    totalPrivateRoomIncome: number
    totalPrivateRoomTicketIncome: number
    totalPrivateRoomWatchIncome: number
    totalVideoCallIncome: number
    totalVideoCallTicketIncome: number
    totalVideoCallBillingIncome: number
    totalGameBet: number
    totalGiftSender: number
    totalNewFollower: number
    createdAt?: string | null
}

export interface VideoCallLogQuery extends PageQuery {
    callerId?: string
    receiverId?: string
    source?: number
    startTime?: number
    endTime?: number
}

export interface VideoCallLogItem {
    id: string
    callerId: string
    callerNickname?: string
    receiverId: string
    receiverNickname?: string
    callType: number
    callTypeText?: string
    source: number
    sourceText?: string
    statusText?: string
    callStartTime?: string | null
    answerTime?: string | null
    callerHeartTime?: string | null
    receiverHeartTime?: string | null
    orderEndTime?: string | null
    callDuration: number
    ticketPrice: number
    pricePerMinute: number
    totalCost: number
    billingDuration: number
    chargeTime?: string | null
    createdAt?: string | null
}

export interface LogPathsConfig {
    serverTime?: string
    logDir: string
    accessPrefix: string
    detailPrefix: string
    errorPrefix: string
    fileExportStaticPrefix?: string
    fileExportAbsDir?: string
    fileExportUrlPrefix: string
    fileExportTtlMinutes?: number
    linuxOnly?: boolean
}

export interface LogQueryExportResult {
    exportId: string
    fileName: string
    fileUrl: string
    total: number
    pageIndex: number
    pageSize: number
}

export interface DetailLogQuery {
    pageIndex?: number
    pageSize?: number
    startDate: string
    endDate: string
    traceId?: string
    reqId?: string
    authId?: string
    url?: string
    keyword?: string
}

/** syndb 刷盘日志明细项 */
export interface SyndbFlushDetailItem {
    table: string
    col: string
    rows: number
    reason: string
    waitMs: number
}

/** syndb 刷盘日志结构化数据 */
export interface SyndbFlushLog {
    reason: string
    sysCpu?: number
    cpuIdle?: number
    idleThreshold?: number
    batchLimit: number
    queues: number
    rows: number
    idleQueues: number
    forceQueues: number
    costMs: number
    details: SyndbFlushDetailItem[]
}

export interface DetailLogItem {
    time: string
    level: string
    traceId: string
    reqId: string
    authId: string
    url: string
    elapsedMs?: number
    message: string
    raw: string
    syndbFlush?: SyndbFlushLog
}

export interface AccessLogQuery {
    pageIndex?: number
    pageSize?: number
    startDate: string
    endDate: string
    traceId?: string
    authId?: string
    url?: string
    ip?: string
    statusCode?: number
    minHandlerMs?: number
    maxHandlerMs?: number
}

export interface AccessLogItem {
    time: string
    traceId: string
    authId?: string
    statusCode: number
    method: string
    url: string
    handlerMs: number
    ip: string
    userAgent: string
    raw: string
}

export interface TraceLogDetail {
    traceId: string
    startDate: string
    endDate: string
    detailLogs: DetailLogItem[]
    accessLogs: AccessLogItem[]
    errorLogs: ErrorLogItem[]
}

export interface ErrorLogQuery {
    pageIndex?: number
    pageSize?: number
    startDate: string
    endDate: string
    traceId?: string
    url?: string
    ip?: string
    statusCode?: number
    keyword?: string
}

export interface ErrorLogItem {
    time: string
    level: string
    traceId: string
    statusCode: number
    method: string
    url: string
    handlerMs: number
    ip: string
    authId?: string
    errorCode: number
    errorMessage: string
    detail: string
    stack: string
    raw: string
}

export interface TopStatItem {
    key: string
    count: number
}

export interface AccessLogStats {
    urlTop: TopStatItem[]
    ipTop: TopStatItem[]
}

export interface AccessTrendPoint {
    time: string
    count: number
}

export interface AccessTrendData {
    intervalMinutes: number
    points: AccessTrendPoint[]
    totalCount: number
    peakTime: string
    peakCount: number
}

export interface AccessTrendQuery {
    startDate: string
    endDate: string
    traceId?: string
    authId?: string
    url?: string
    ip?: string
    statusCode?: number
    minHandlerMs?: number
    maxHandlerMs?: number
    intervalMinutes?: number
}

export type LogQueryJobStatus = 'pending' | 'running' | 'done' | 'failed'

export interface LogQueryJobSubmitResult {
    jobId: string
    queuePosition: number
}

export interface LogQueryJobResult<T = unknown> {
    jobId: string
    queryType: string
    status: LogQueryJobStatus
    queuePosition: number
    errorMessage?: string
    result?: T
}

export type CMSExportJobStatus = 'pending' | 'running' | 'done' | 'failed'

export interface CMSExportJobSubmitResult {
    jobId: string
    queuePosition: number
}

export interface CMSExportJobProgress {
    exportedRows: number
    totalRows: number
}

export interface CMSExportResult {
    exportId: string
    fileName: string
    fileUrl: string
    total: number
    pageIndex?: number
    pageSize?: number
}

export interface CMSExportJobResult {
    jobId: string
    exportType: string
    status: CMSExportJobStatus
    queuePosition: number
    errorMessage?: string
    progress?: CMSExportJobProgress
    result?: CMSExportResult
}

// CMS用户相关类型
export interface CMSUser {
    id: string
    name: string
    pwd: string
    status: number
    admin: boolean
    roleId: string
    roleName: string
    remark: string
    createdAt: string
    updatedAt: string
}

// 随机昵称库
export interface RandomNicknameLangItem {
    lang: number
    langCode: string
    langLabel: string
    count: number
    samples: string[]
}

export interface GetRandomNicknameCfgRes {
    useDB: boolean
    langs: RandomNicknameLangItem[]
}

export interface ImportRandomNicknamesReq {
    lang: number
    content: string
    replace: boolean
}

export interface ImportRandomNicknamesRes {
    imported: number
    total: number
}

export interface ClearRandomNicknamesReq {
    lang: number
}

export interface ClearRandomNicknamesRes {
    success: boolean
    total: number
}