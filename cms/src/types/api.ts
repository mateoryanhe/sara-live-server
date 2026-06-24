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
    banApplyTime?: string
}

export interface BanAnchorReq {
    accountId: string
    banApplyTime: string
    banReason: string
}

export interface UnBanAnchorReq {
    accountId: string
}

export interface UnBanReq {
    accountId: string
}

export interface CancelReq {
    accountId: string
}

export interface UnCancelReq {
    accountId: string
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
    // 来自 user_infos 表(LEFT JOIN，可能为空)
    nickname?: string
    phone?: string
    avatar?: string
    remark?: string
    gold?: number
    diamond?: number
    shareCode?: string
    guildId?: string | number
    isAnchor?: boolean
    vipLevel?: number
    deviceType?: string
}

export interface SetAnchorReq {
    accountId: string
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
    liveStatus?: number
    totalIncome?: number
    totalGiftIncome?: number
    totalPaidDanmakuIncome?: number
    totalPrivateRoomTicketIncome?: number
    totalPrivateRoomWatchIncome?: number
    ban?: boolean
    banApplyTime?: string | null
    banReason?: string
    createdAt?: string | null
    registeredAt?: string | null
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


// 全局配置相关类型
export interface GlobalCfg {
    id: string  // 根据API返回的实际数据，ID 是字符串类型
    module: string
    moduleName: string
    key: string
    value: string
    desc: string
}

export interface GetGlobalCfgReq {
    module?: string
    moduleName?: string
}

export interface SaveGlobalCfgReq extends GlobalCfg {
}

export interface DelGlobalCfgReq extends GlobalCfg {
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
    createdAt: string
    updatedAt: string
}

export interface VipCfgQuery extends PageQuery {
    levelName?: string
    statusFilter?: number
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
    direction: number
    sort: number
    status: number
    createdAt: string
    updatedAt: string
}

export interface BannerQuery extends PageQuery {
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
    diamondPerMinute: number
    categoryId: number
    source: number
    authorId: string
    authorNickname: string
    duration: number
    freeWatchSeconds: number
    likeCount: number
    viewCount: number
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
    billedSeconds: number
    watchSeconds: number
    createdAt: string
    updatedAt: string
}

export interface ShortVideoWatchQuery extends PageQuery {
    userId?: string
    startTime?: number
    endTime?: number
}

export interface ShortVideoCfg {
    id: string
    maxFileSize: number
    maxCoverFileSize: number
    maxDuration: number
    freeWatchSeconds: number
    entryEnabled: number
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

export interface AgoraCfg {
    id: string
    appId: string
    appCertificate: string
    restCustomerId: string
    restCustomerSecret: string
    tokenExpireSeconds: number
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
    tokenExpireSeconds: number
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

export interface UploadResourceCfg {
    id: string
    resourceDomain: string
    defaultAvatarUrl: string
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
    defaultAvatarUrl: string
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
    totalGameBet: number
    totalGiftSender: number
    totalNewFollower: number
    createdAt?: string | null
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