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
}

export interface UserInfo {
    id: string
    createdAt?: string | null
    openId: string
    ip: string
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
    deviceType?: string
    packageName?: string
    appVersion?: string
    canRank?: boolean
    cancelCode?: string
}

export interface SetAnchorReq {
    accountId: string
}

export interface SetUserTypeReq {
    accountId: string
    userType: number
}

export interface SetCanRankReq {
    accountId: string
    canRank: boolean
}

export interface QueryAnchorListReq extends PageQuery {
    key?: string
}

export interface AnchorListItem {
    id: string
    nickname?: string
    phone?: string
    avatar?: string
    guildId?: string | number
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
    createdAt?: string | null
    registeredAt?: string | null
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
    createdAt?: string
    updatedAt?: string
}

export interface GetAccountCfgRes {
    cfg?: AccountCfg | null
}

export interface SaveAccountCfgReq {
    id?: string | number
    cancelAccountByCodeEnabled: boolean
}

export interface SaveAccountCfgRes {
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
    createdAt?: string
    updatedAt?: string
}

// 直播工会相关类型
export interface Guild {
    id: string
    name: string
    leaderId: string
    contact: string
    description: string
    status: number
    createdAt: string
    updatedAt: string
}

export interface GuildQuery extends PageQuery {
    name?: string
}

// 礼物相关类型
export interface Gift {
    id: string
    name: string
    nameEn: string
    nameEs: string
    namePt: string
    nameHi: string
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
    typeFilter?: number
    statusFilter?: number
}

// VIP配置
export interface VipCfg {
    id: string
    level: number
    levelName: string
    status: number
    upgradeRechargeLimit: number
    minWithdrawAmount: number
    maxWithdrawAmount: number
    fee: number
    animation?: string
    animationName?: string
    createdAt: string
    updatedAt: string
}

export interface VipCfgQuery extends PageQuery {
    levelName?: string
    statusFilter?: number
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

// 游戏配置
export interface GameCfg {
    id: string
    name: string
    code: string
    liveCover: string
    liveCoverUrl?: string
    link: string
    sort: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface GameCfgQuery extends PageQuery {
    name?: string
    code?: string
    statusFilter?: number
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
}

export interface ShortVideoWatchRecord {
    id: string
    userId: string
    nickname: string
    videoId: string
    videoTitle: string
    paidTime: string
    freeTime: number
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

export interface ShortVideoCategory {
    id: string
    name: string
    sort: number
    createdAt: string
    updatedAt: string
}

export interface ShortVideoCategoryQuery extends PageQuery {
}

export interface LiveRoomTag {
    id: string
    name: string
    sort: number
    createdAt: string
    updatedAt: string
}

export interface LiveRoomTagQuery extends PageQuery {
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
    privacyPolicyUrl: string
    termsOfServiceUrl: string
    createdAt: string
    updatedAt: string
}

export interface GetPrivacyPolicyCfgRes {
    cfg: PrivacyPolicyCfg | null
}

export interface SavePrivacyPolicyCfgReq {
    id?: number
    privacyPolicyUrl: string
    termsOfServiceUrl: string
}

export interface SavePrivacyPolicyCfgRes {
    success: boolean
    id: string
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
    startTime?: number
    endTime?: number
}

export interface LiveRecordItem {
    id: string
    anchorId: string
    nickname?: string
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
    exportSubDir: string
    exportStaticPrefix?: string
    exportAbsDir?: string
    exportUrlPrefix: string
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
}

export interface AccessLogQuery {
    pageIndex?: number
    pageSize?: number
    startDate: string
    endDate: string
    traceId?: string
    url?: string
    ip?: string
    statusCode?: number
    minHandlerMs?: number
    maxHandlerMs?: number
}

export interface AccessLogItem {
    time: string
    traceId: string
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

// CMS用户相关类型
export interface CMSUser {
    id: string
    name: string
    pwd: string
    status: number
    admin: boolean
    roleId: string
    roleName: string
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