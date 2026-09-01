import {getPageFromPermissionKey, PERMISSION_SEP} from './page-buttons'

/** 权限 module -> CMS 后端请求路径(与 request.post 第一个参数一致) */
const PERMISSION_API_PATHS: Record<string, string> = {
    Dashboard: '/sysStat/getSysStat',
    'Dashboard:view': '/sysStat/getSysStat',

    UserList: '/account/getUserInfo',
    'UserList:view': '/account/getUserInfo',
    'UserList:search': '/account/getUserInfo',
    'UserList:viewDetail': '/account/getUserDetail',
    'UserList:viewAnchorDetail': '/account/getAnchorDetail',
    'UserList:setAnchor': '/account/setAnchor',
    'UserList:setSeniorAnchor': '/account/setSeniorAnchor',
    'UserList:goldAdd': '/gold/add',
    'UserList:goldSub': '/gold/sub',
    'UserList:diamondAdd': '/diamond/add',
    'UserList:diamondSub': '/diamond/sub',
    'UserList:ban': '/account/ban',
    'UserList:rankOff': '/account/setCanRank',
    'UserList:rankOn': '/account/setCanRank',
    'UserList:rechargeWhitelistOn': '/account/setRechargeWhitelist',
    'UserList:rechargeWhitelistOff': '/account/setRechargeWhitelist',
    'UserList:cancel': '/account/cancel',
    'UserList:setUserType': '/account/setUserType',
    'UserList:uploadAvatar': '/account/setUserAvatar',
    'UserList:setAnchorType': '/account/setPlatformAnchorType',
    'UserList:openGame': '/gamePlatform/gameShelfList',
    UserDetail: '/account/getUserDetail',
    'UserDetail:goldLog': '/currencyLog/cmsCurrencyLogList',
    'UserDetail:exportGoldLog': '/cmsExport/submitJob',
    'UserDetail:diamondLog': '/currencyLog/cmsCurrencyLogList',
    'UserDetail:exportDiamondLog': '/cmsExport/submitJob',

    AnchorListManagement: '/account/getAnchorList',
    'AnchorListManagement:view': '/account/getAnchorList',
    'AnchorListManagement:search': '/account/getAnchorList',
    'AnchorListManagement:edit': '/account/setAnchor',
    'AnchorListManagement:viewDetail': '/account/getAnchorDetail',
    'AnchorListManagement:viewUserDetail': '/account/getUserDetail',
    'AnchorListManagement:viewGuildDetail': '/guild/getGuildDetail',
    'AnchorListManagement:offShelf': '/account/setLiveRoomStatus',
    'AnchorListManagement:uploadRoomCover': '/account/setLiveRoomCover',

    PlatformAnchorList: '/account/getAnchorList',
    'PlatformAnchorList:view': '/account/getAnchorList',
    'PlatformAnchorList:search': '/account/getAnchorList',
    'PlatformAnchorList:viewDetail': '/account/getAnchorDetail',
    'PlatformAnchorList:viewUserDetail': '/account/getUserDetail',
    'PlatformAnchorList:offShelf': '/account/setLiveRoomStatus',
    'PlatformAnchorList:ban': '/account/banAnchor',
    'PlatformAnchorList:unban': '/account/unBanAnchor',
    'PlatformAnchorList:setAnchorType': '/account/setPlatformAnchorType',

    AnchorDetail: '/account/getAnchorDetail',
    'AnchorDetail:dailyEffectiveLive': '/account/getAnchorDailyEffectiveLiveList',
    'AnchorDetail:exportDailyEffectiveLive': '/cmsExport/submitJob',
    'AnchorDetail:liveRecord': '/liveRecord/cmsLiveRecordList',
    'AnchorDetail:exportLiveRecord': '/cmsExport/submitJob',
    'AnchorDetail:liveRecordRevenue': '/liveRevenueLog/cmsLiveRevenueLogList',
    'AnchorDetail:settlementLog': '/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList',
    'AnchorDetail:exportSettlementLog': '/cmsExport/submitJob',

    LiveRoomRecycleBinManagement: '/account/getOffShelfLiveRoomList',
    'LiveRoomRecycleBinManagement:view': '/account/getOffShelfLiveRoomList',
    'LiveRoomRecycleBinManagement:search': '/account/getOffShelfLiveRoomList',
    'LiveRoomRecycleBinManagement:onShelf': '/account/setLiveRoomStatus',

    BotAnchorManagement: '/botAnchor/getBotAnchorList',
    'BotAnchorManagement:view': '/botAnchor/getBotAnchorList',
    'BotAnchorManagement:search': '/botAnchor/getBotAnchorList',
    'BotAnchorManagement:create': '/botAnchor/createBotAnchor',
    'BotAnchorManagement:edit': '/botAnchor/updateBotAnchor',
    'BotAnchorManagement:delete': '/botAnchor/setBotAnchorStatus',

    BanUser: '/account/ban',
    'BanUser:view': '/account/ban',
    'BanUser:save': '/account/ban',

    GuildMembers: '/account/getAnchorList',
    'GuildMembers:view': '/account/getAnchorList',

    GuildAnchorImportResult: '/guild/importGuildAnchors',
    'GuildAnchorImportResult:view': '/guild/importGuildAnchors',

    GuildProfileMembers: '/guild/getMyGuildAnchorList',
    'GuildProfileMembers:view': '/guild/getMyGuildAnchorList',

    GuildProfileAnchorDailyLive: '/guild/getMyGuildAnchorDailyEffectiveLiveList',
    'GuildProfileAnchorDailyLive:view': '/guild/getMyGuildAnchorDailyEffectiveLiveList',
    'GuildProfileAnchorDailyLive:search': '/guild/getMyGuildAnchorDailyEffectiveLiveList',
    'GuildProfileAnchorDailyLive:export': '/cmsExport/submitJob',

    AnchorLiveRecordDetail: '/liveRecord/getLiveRecordList',
    'AnchorLiveRecordDetail:view': '/liveRecord/getLiveRecordList',
    'AnchorLiveRecordDetail:search': '/liveRecord/getLiveRecordList',

    RechargeOrderList: '/rechargeOrder/rechargeOrderList',
    'RechargeOrderList:view': '/rechargeOrder/rechargeOrderList',
    'RechargeOrderList:search': '/rechargeOrder/rechargeOrderList',
    'RechargeOrderList:viewUserDetail': '/account/getUserDetail',
    'RechargeOrderList:manualCreateOrder': '/rechargeOrder/manualCreateOrder',
    'RechargeOrderList:manualRecharge': '/rechargeOrder/manualRecharge',

    GoldCurrencyLogList: '/currencyLog/cmsCurrencyLogList',
    'GoldCurrencyLogList:view': '/currencyLog/cmsCurrencyLogList',
    'GoldCurrencyLogList:search': '/currencyLog/cmsCurrencyLogList',
    'GoldCurrencyLogList:export': '/currencyLog/cmsCurrencyLogList',
    'GoldCurrencyLogList:viewUserDetail': '/account/getUserDetail',

    DiamondCurrencyLogList: '/currencyLog/cmsCurrencyLogList',
    'DiamondCurrencyLogList:view': '/currencyLog/cmsCurrencyLogList',
    'DiamondCurrencyLogList:search': '/currencyLog/cmsCurrencyLogList',
    'DiamondCurrencyLogList:export': '/currencyLog/cmsCurrencyLogList',
    'DiamondCurrencyLogList:viewUserDetail': '/account/getUserDetail',

    BannerManagement: '/banner/bannerList',
    'BannerManagement:view': '/banner/bannerList',
    'BannerManagement:search': '/banner/bannerList',
    'BannerManagement:create': '/banner/createBanner',
    'BannerManagement:edit': '/banner/updateBanner',
    'BannerManagement:delete': '/banner/deleteBanner',
    'BannerManagement:sort': '/banner/updateBanner',
    'BannerManagement:toggle': '/banner/onShelfBanner',
    'BannerManagement:sync': '/dataSync/syncBanner',

    ActivityMessageManagement: '/activityMessage/activityMessageList',
    'ActivityMessageManagement:view': '/activityMessage/activityMessageList',
    'ActivityMessageManagement:search': '/activityMessage/activityMessageList',
    'ActivityMessageManagement:create': '/activityMessage/createActivityMessage',
    'ActivityMessageManagement:edit': '/activityMessage/updateActivityMessage',
    'ActivityMessageManagement:delete': '/activityMessage/deleteActivityMessage',
    'ActivityMessageManagement:publish': '/activityMessage/publishActivityMessage',
    'ActivityMessageManagement:sync': '/dataSync/syncActivityMessage',

    GuildManagement: '/guild/guildList',
    'GuildManagement:view': '/guild/guildList',
    'GuildManagement:search': '/guild/guildList',
    'GuildManagement:create': '/guild/createGuild',
    'GuildManagement:edit': '/guild/updateGuild',
    'GuildManagement:offShelf': '/guild/deleteGuild',
    'GuildManagement:viewMembers': '/account/getAnchorList',
    'GuildManagement:viewDetail': '/guild/getGuildDetail',
    'GuildManagement:viewUserDetail': '/account/getUserDetail',
    'GuildManagement:viewAnchorSettlementLogs': '/guild/cmsGuildAnchorIncomeSettlementLogList',
    'GuildManagement:ban': '/account/banAnchor',
    'GuildManagement:unban': '/account/unBanAnchor',
    'GuildManagement:exitGuild': '/account/exitGuild',
    'GuildManagement:setAnchorType': '/guild/setGuildAnchorType',
    'GuildManagement:joinGuildAnchor': '/guild/joinGuildAnchor',
    'GuildManagement:batchSetAnchor': '/guild/importGuildAnchors',
    'GuildManagement:batchSetSeniorAnchor': '/guild/importGuildAnchors',

    GuildCMSUserManagement: '/cmsuser/cmsUserList',
    'GuildCMSUserManagement:view': '/cmsuser/cmsUserList',
    'GuildCMSUserManagement:search': '/cmsuser/cmsUserList',
    'GuildCMSUserManagement:create': '/cmsuser/createGuildCMSUser',
    'GuildCMSUserManagement:resetPassword': '/cmsuser/updateCMSUser',

    GuildAnchorDailyLiveManagement: '/guild/cmsMyGuildAnchorDailyEffectiveLiveList',
    'GuildAnchorDailyLiveManagement:view': '/guild/cmsMyGuildAnchorDailyEffectiveLiveList',
    'GuildAnchorDailyLiveManagement:search': '/guild/cmsMyGuildAnchorDailyEffectiveLiveList',
    'GuildAnchorDailyLiveManagement:export': '/cmsExport/submitJob',

    GuildDetail: '/guild/getGuildDetail',
    'GuildDetail:viewUserDetail': '/account/getUserDetail',
    'GuildDetail:incomeArchive': '/guild/getGuildIncomeArchives',
    'GuildDetail:settlementLog': '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList',
    'GuildDetail:exportSettlementLog': '/cmsExport/submitJob',
    'GuildDetail:anchorSettlementLog': '/guild/cmsGuildAnchorIncomeSettlementLogList',
    'GuildDetail:exportAnchorSettlementLog': '/cmsExport/submitJob',
    'GuildDetail:anchorDailyEffectiveLive': '/guild/cmsGuildAnchorDailyEffectiveLiveList',
    'GuildDetail:exportAnchorDailyEffectiveLive': '/cmsExport/submitJob',
    'GuildDetail:dailyEffectiveLive': '/guild/getGuildDailyEffectiveLiveList',
    'GuildDetail:exportDailyEffectiveLive': '/cmsExport/submitJob',

    GuildRecycleBinManagement: '/guild/offShelfGuildList',
    'GuildRecycleBinManagement:view': '/guild/offShelfGuildList',
    'GuildRecycleBinManagement:search': '/guild/offShelfGuildList',
    'GuildRecycleBinManagement:onShelf': '/guild/onShelfGuild',

    GuildProfileManagement: '/guild/getMyGuildProfile',
    'GuildProfileManagement:view': '/guild/getMyGuildProfile',
    'GuildProfileManagement:viewUserDetail': '/account/getUserDetail',
    'GuildProfileManagement:viewAnchors': '/guild/getMyGuildAnchorList',
    'GuildProfileManagement:dailyEffectiveLive': '/guild/getMyGuildAnchorDailyEffectiveLiveList',
    'GuildProfileManagement:viewDetail': '/account/getAnchorDetail',

    RechargeCfgManagement: '/rechargeCfg/rechargeCfgList',
    'RechargeCfgManagement:view': '/rechargeCfg/rechargeCfgList',
    'RechargeCfgManagement:search': '/rechargeCfg/rechargeCfgList',
    'RechargeCfgManagement:create': '/rechargeCfg/createRechargeCfg',
    'RechargeCfgManagement:edit': '/rechargeCfg/updateRechargeCfg',
    'RechargeCfgManagement:delete': '/rechargeCfg/deleteRechargeCfg',
    'RechargeCfgManagement:sync': '/dataSync/syncRechargeCfg',

    VipCfgManagement: '/vipCfg/vipCfgList',
    'VipCfgManagement:view': '/vipCfg/vipCfgList',
    'VipCfgManagement:search': '/vipCfg/vipCfgList',
    'VipCfgManagement:create': '/vipCfg/createVipCfg',
    'VipCfgManagement:edit': '/vipCfg/updateVipCfg',
    'VipCfgManagement:delete': '/vipCfg/deleteVipCfg',
    'VipCfgManagement:sync': '/dataSync/syncVipCfg',
    'VipCfgManagement:syncAssets': '/dataSync/syncVipCfgAssets',

    AppPkgManagement: '/appPkg/appPkgList',
    'AppPkgManagement:view': '/appPkg/appPkgList',
    'AppPkgManagement:search': '/appPkg/appPkgList',
    'AppPkgManagement:create': '/appPkg/createAppPkg',
    'AppPkgManagement:edit': '/appPkg/updateAppPkg',
    'AppPkgManagement:delete': '/appPkg/deleteAppPkg',

    RandomNicknameManagement: '/randomNickname/getRandomNicknameCfg',
    'RandomNicknameManagement:view': '/randomNickname/getRandomNicknameCfg',
    'RandomNicknameManagement:save': '/randomNickname/importRandomNicknames',

    CustomerServiceCfgManagement: '/customerService/getCustomerServiceCfg',
    'CustomerServiceCfgManagement:view': '/customerService/getCustomerServiceCfg',
    'CustomerServiceCfgManagement:save': '/customerService/saveCustomerServiceCfg',

    WalletExchangeCfgManagement: '/wallet/getWalletExchangeCfg',
    'WalletExchangeCfgManagement:view': '/wallet/getWalletExchangeCfg',
    'WalletExchangeCfgManagement:save': '/wallet/saveWalletExchangeCfg',

    LiveRevenueShareCfgManagement: '/liveRevenueShareCfg/getLiveRevenueShareCfg',
    'LiveRevenueShareCfgManagement:view': '/liveRevenueShareCfg/getLiveRevenueShareCfg',
    'LiveRevenueShareCfgManagement:save': '/liveRevenueShareCfg/saveLiveRevenueShareCfg',

    GiftManagement: '/gift/giftList',
    'GiftManagement:view': '/gift/giftList',
    'GiftManagement:create': '/gift/createGift',
    'GiftManagement:edit': '/gift/updateGift',
    'GiftManagement:delete': '/gift/deleteGift',
    'GiftManagement:sort': '/gift/updateGift',
    'GiftManagement:sync': '/dataSync/syncGift',
    'GiftManagement:syncAssets': '/dataSync/syncGiftAssets',

    AgoraCfgManagement: '/agora/getAgoraCfg',
    'AgoraCfgManagement:view': '/agora/getAgoraCfg',
    'AgoraCfgManagement:save': '/agora/saveAgoraCfg',

    TicketManagement: '/ticket/ticketList',
    'TicketManagement:view': '/ticket/ticketList',
    'TicketManagement:search': '/ticket/ticketList',
    'TicketManagement:create': '/ticket/createTicket',
    'TicketManagement:edit': '/ticket/updateTicket',
    'TicketManagement:delete': '/ticket/deleteTicket',

    PrivateRoomBillingManagement: '/privateRoomBilling/billingList',
    'PrivateRoomBillingManagement:view': '/privateRoomBilling/billingList',
    'PrivateRoomBillingManagement:search': '/privateRoomBilling/billingList',
    'PrivateRoomBillingManagement:create': '/privateRoomBilling/createBilling',
    'PrivateRoomBillingManagement:edit': '/privateRoomBilling/updateBilling',
    'PrivateRoomBillingManagement:delete': '/privateRoomBilling/deleteBilling',

    LiveCfgManagement: '/liveCfg/getLiveCfg',
    'LiveCfgManagement:view': '/liveCfg/getLiveCfg',
    'LiveCfgManagement:save': '/liveCfg/saveLiveCfg',

    LiveRoomTagManagement: '/liveRoomTag/liveRoomTagList',
    'LiveRoomTagManagement:view': '/liveRoomTag/liveRoomTagList',
    'LiveRoomTagManagement:search': '/liveRoomTag/liveRoomTagList',
    'LiveRoomTagManagement:create': '/liveRoomTag/createLiveRoomTag',
    'LiveRoomTagManagement:edit': '/liveRoomTag/updateLiveRoomTag',
    'LiveRoomTagManagement:delete': '/liveRoomTag/deleteLiveRoomTag',

    LiveRevenueLogList: '/liveRevenueLog/cmsLiveRevenueLogList',
    'LiveRevenueLogList:view': '/liveRevenueLog/cmsLiveRevenueLogList',
    'LiveRevenueLogList:search': '/liveRevenueLog/cmsLiveRevenueLogList',
    'LiveRevenueLogList:export': '/cmsExport/submitJob',
    'LiveRevenueLogList:viewUserDetail': '/account/getUserDetail',

    LiveRecordList: '/liveRecord/cmsLiveRecordList',
    'LiveRecordList:view': '/liveRecord/cmsLiveRecordList',
    'LiveRecordList:search': '/liveRecord/cmsLiveRecordList',
    'LiveRecordList:export': '/cmsExport/submitJob',
    'LiveRecordList:viewUserDetail': '/account/getUserDetail',

    LiveDailyEffectiveLiveList: '/liveRecord/cmsDailyEffectiveLiveList',
    'LiveDailyEffectiveLiveList:view': '/liveRecord/cmsDailyEffectiveLiveList',
    'LiveDailyEffectiveLiveList:search': '/liveRecord/cmsDailyEffectiveLiveList',
    'LiveDailyEffectiveLiveList:export': '/cmsExport/submitJob',

    VideoCallLogList: '/call/cmsVideoCallLogList',
    'VideoCallLogList:view': '/call/cmsVideoCallLogList',
    'VideoCallLogList:search': '/call/cmsVideoCallLogList',
    'VideoCallLogList:export': '/cmsExport/submitJob',
    'VideoCallLogList:viewUserDetail': '/account/getUserDetail',
    'VideoCallLogList:viewAnchorDetail': '/account/getAnchorDetail',

    ShortVideoManagement: '/shortVideo/shortVideoList',
    'ShortVideoManagement:view': '/shortVideo/shortVideoList',
    'ShortVideoManagement:search': '/shortVideo/shortVideoList',
    'ShortVideoManagement:create': '/shortVideo/createShortVideo',
    'ShortVideoManagement:edit': '/shortVideo/updateShortVideo',
    'ShortVideoManagement:delete': '/shortVideo/deleteShortVideo',
    'ShortVideoManagement:onShelf': '/shortVideo/onShelfShortVideo',
    'ShortVideoManagement:offShelf': '/shortVideo/offShelfShortVideo',
    'ShortVideoManagement:viewStorageStat': '/shortVideo/shortVideoStorageStat',

    ShortVideoCfgManagement: '/shortVideo/getShortVideoCfg',
    'ShortVideoCfgManagement:view': '/shortVideo/getShortVideoCfg',
    'ShortVideoCfgManagement:save': '/shortVideo/saveShortVideoCfg',

    ShortVideoCategoryManagement: '/shortVideo/shortVideoCategoryList',
    'ShortVideoCategoryManagement:view': '/shortVideo/shortVideoCategoryList',
    'ShortVideoCategoryManagement:search': '/shortVideo/shortVideoCategoryList',
    'ShortVideoCategoryManagement:create': '/shortVideo/createShortVideoCategory',
    'ShortVideoCategoryManagement:edit': '/shortVideo/updateShortVideoCategory',
    'ShortVideoCategoryManagement:delete': '/shortVideo/deleteShortVideoCategory',

    ShortVideoPriceTierManagement: '/shortVideo/shortVideoPriceTierList',
    'ShortVideoPriceTierManagement:view': '/shortVideo/shortVideoPriceTierList',
    'ShortVideoPriceTierManagement:search': '/shortVideo/shortVideoPriceTierList',
    'ShortVideoPriceTierManagement:create': '/shortVideo/createShortVideoPriceTier',
    'ShortVideoPriceTierManagement:edit': '/shortVideo/updateShortVideoPriceTier',
    'ShortVideoPriceTierManagement:delete': '/shortVideo/deleteShortVideoPriceTier',

    ShortVideoWatchManagement: '/shortVideo/shortVideoWatchList',
    'ShortVideoWatchManagement:view': '/shortVideo/shortVideoWatchList',
    'ShortVideoWatchManagement:search': '/shortVideo/shortVideoWatchList',
    'ShortVideoWatchManagement:export': '/shortVideo/shortVideoWatchList',

    GamePlatformCfgManagement: '/gamePlatform/getGamePlatformCfg',
    'GamePlatformCfgManagement:view': '/gamePlatform/getGamePlatformCfg',
    'GamePlatformCfgManagement:save': '/gamePlatform/saveGamePlatformCfg',

    GameVendorGameListManagement: '/gamePlatform/vendorGameList',
    'GameVendorGameListManagement:view': '/gamePlatform/vendorGameList',
    'GameVendorGameListManagement:search': '/gamePlatform/vendorGameList',
    'GameVendorGameListManagement:edit': '/gamePlatform/addGameShelf',
    'GameVendorGameListManagement:shelf': '/gamePlatform/addGameShelf',
    'GameVendorGameListManagement:sync': '/gamePlatform/reloadVendorGameCache',

    GameShelfListManagement: '/gamePlatform/gameShelfList',
    'GameShelfListManagement:view': '/gamePlatform/gameShelfList',
    'GameShelfListManagement:search': '/gamePlatform/gameShelfList',
    'GameShelfListManagement:edit': '/gamePlatform/updateGameShelf',
    'GameShelfListManagement:shelf': '/gamePlatform/deleteGameShelf',
    'GameShelfListManagement:vendorConfig': '/gamePlatform/getMultiplayerConfigUrl',
    'GameShelfListManagement:startGame': '/gamePlatform/cmsGameStartLink',

    GameVendorConfig: '/gamePlatform/getMultiplayerConfigUrl',

    FirstRechargeActivityManagement: '/firstRechargeActivity/getFirstRechargeActivityCfg',
    'FirstRechargeActivityManagement:view': '/firstRechargeActivity/getFirstRechargeActivityCfg',
    'FirstRechargeActivityManagement:save': '/firstRechargeActivity/saveFirstRechargeActivityCfg',
    'FirstRechargeActivityManagement:sync': '/dataSync/syncFirstRechargeActivityCfg',

    GameWinLogListManagement: '/gameWinLog/cmsGameWinLogList',
    'GameWinLogListManagement:view': '/gameWinLog/cmsGameWinLogList',
    'GameWinLogListManagement:search': '/gameWinLog/cmsGameWinLogList',
    'GameWinLogListManagement:export': '/gameWinLog/cmsGameWinLogList',

    GameBetLogListManagement: '/gameBetLog/cmsGameBetLogList',
    'GameBetLogListManagement:view': '/gameBetLog/cmsGameBetLogList',
    'GameBetLogListManagement:search': '/gameBetLog/cmsGameBetLogList',
    'GameBetLogListManagement:export': '/gameBetLog/cmsGameBetLogList',

    AppTokenConfig: '/appToken/getAppToken',
    'AppTokenConfig:view': '/appToken/getAppToken',
    'AppTokenConfig:save': '/appToken/saveAppToken',

    AccountCfgManagement: '/accountCfg/getAccountCfg',
    'AccountCfgManagement:view': '/accountCfg/getAccountCfg',
    'AccountCfgManagement:save': '/accountCfg/saveAccountCfg',

    AppVersionCfgManagement: '/appVersionCfg/getAppVersionCfg',
    'AppVersionCfgManagement:view': '/appVersionCfg/getAppVersionCfg',
    'AppVersionCfgManagement:save': '/appVersionCfg/saveAppVersionCfg',

    SimulatorCpuKeywordManagement: '/simulatorCpuKeyword/simulatorCpuKeywordList',
    'SimulatorCpuKeywordManagement:view': '/simulatorCpuKeyword/simulatorCpuKeywordList',
    'SimulatorCpuKeywordManagement:search': '/simulatorCpuKeyword/simulatorCpuKeywordList',
    'SimulatorCpuKeywordManagement:create': '/simulatorCpuKeyword/createSimulatorCpuKeyword',
    'SimulatorCpuKeywordManagement:edit': '/simulatorCpuKeyword/updateSimulatorCpuKeyword',
    'SimulatorCpuKeywordManagement:delete': '/simulatorCpuKeyword/deleteSimulatorCpuKeyword',

    AnchorSalaryCfgManagement: '/anchorSalaryCfg/anchorSalaryCfgList',
    'AnchorSalaryCfgManagement:view': '/anchorSalaryCfg/anchorSalaryCfgList',
    'AnchorSalaryCfgManagement:search': '/anchorSalaryCfg/anchorSalaryCfgList',
    'AnchorSalaryCfgManagement:create': '/anchorSalaryCfg/createAnchorSalaryCfg',
    'AnchorSalaryCfgManagement:edit': '/anchorSalaryCfg/updateAnchorSalaryCfg',
    'AnchorSalaryCfgManagement:delete': '/anchorSalaryCfg/deleteAnchorSalaryCfg',

    AnchorIncomeSettlementLogList: '/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList',
    'AnchorIncomeSettlementLogList:view': '/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList',
    'AnchorIncomeSettlementLogList:search': '/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList',
    'AnchorIncomeSettlementLogList:export': '/cmsExport/submitJob',

    GuildIncomeSettlementLogList: '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList',
    'GuildIncomeSettlementLogList:view': '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList',
    'GuildIncomeSettlementLogList:search': '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList',
    'GuildIncomeSettlementLogList:export': '/cmsExport/submitJob',

    ServerRuntimeCfgManagement: '/preloadCfg/getPreloadCfg',
    'ServerRuntimeCfgManagement:view': '/preloadCfg/getPreloadCfg',
    'ServerRuntimeCfgManagement:save': '/preloadCfg/savePreloadCfg',
    PreloadCfgManagement: '/preloadCfg/getPreloadCfg',
    'PreloadCfgManagement:view': '/preloadCfg/getPreloadCfg',
    'PreloadCfgManagement:save': '/preloadCfg/savePreloadCfg',

    TextModerationCfgManagement: '/textModeration/getTextModerationCfg',
    'TextModerationCfgManagement:view': '/textModeration/getTextModerationCfg',
    'TextModerationCfgManagement:save': '/textModeration/saveTextModerationCfg',

    PrivacyPolicyCfgManagement: '/privacyPolicy/getPrivacyPolicyCfg',
    'PrivacyPolicyCfgManagement:view': '/privacyPolicy/getPrivacyPolicyCfg',
    'PrivacyPolicyCfgManagement:save': '/privacyPolicy/savePrivacyPolicyCfg',

    GooglePlayCfgManagement: '/googlePlay/getGooglePlayCfg',
    'GooglePlayCfgManagement:view': '/googlePlay/getGooglePlayCfg',
    'GooglePlayCfgManagement:save': '/googlePlay/saveGooglePlayCfg',

    UploadResourceCfgManagement: '/upload/getUploadResourceCfg',
    'UploadResourceCfgManagement:view': '/upload/getUploadResourceCfg',
    'UploadResourceCfgManagement:save': '/upload/saveUploadResourceCfg',

    DataSyncCfgManagement: '/dataSync/getDataSyncCfg',
    'DataSyncCfgManagement:view': '/dataSync/getDataSyncCfg',
    'DataSyncCfgManagement:save': '/dataSync/saveDataSyncCfg',

    H5LiveDeployManagement: '/h5LiveDeploy/getH5LiveDeployInfo',
    'H5LiveDeployManagement:view': '/h5LiveDeploy/getH5LiveDeployInfo',
    'H5LiveDeployManagement:save': '/h5LiveDeploy/saveH5LiveDeployCfg',
    'H5LiveDeployManagement:deploy': '/h5LiveDeploy/deployZip',

    ResourceMonitor: '/resourceMetric/getResourceMetricMemoryTrend',
    'ResourceMonitor:view': '/resourceMetric/getResourceMetricMemoryTrend',
    'ResourceMonitor:search': '/resourceMetric/getResourceMetricMemoryTrend',
    'ResourceMonitor:export': '/resourceMetric/getResourceMetricMemoryTrend',

    ServerLogExplorer: '/logQuery/getLogPaths',
    'ServerLogExplorer:view': '/logQuery/getLogPaths',
    'ServerLogExplorer:search': '/logQuery/submitJob',

    RoleManagement: '/role/roleList',
    'RoleManagement:view': '/role/roleList',
    'RoleManagement:search': '/role/roleList',
    'RoleManagement:create': '/role/createRole',
    'RoleManagement:edit': '/role/updateRole',
    'RoleManagement:delete': '/role/deleteRole',
    'RoleManagement:permission': '/role/permissionList',

    ModuleList: '/role/permissionList',
    'ModuleList:view': '/role/permissionList',
    'ModuleList:save': '/role/createPermission',

    CMSUserManagement: '/cmsuser/cmsUserList',
    'CMSUserManagement:view': '/cmsuser/cmsUserList',
    'CMSUserManagement:search': '/cmsuser/cmsUserList',
    'CMSUserManagement:create': '/cmsuser/createCMSUser',
    'CMSUserManagement:edit': '/cmsuser/updateCMSUser',
    'CMSUserManagement:delete': '/cmsuser/deleteCMSUser',
}

/** 根据权限 module 键解析 CMS 请求路径 */
export function getPermissionApiPath(moduleKey: string): string {
    const key = moduleKey.trim()
    if (!key || key.startsWith('module_')) {
        return ''
    }
    if (PERMISSION_API_PATHS[key]) {
        return PERMISSION_API_PATHS[key]
    }
    const page = getPageFromPermissionKey(key)
    if (page && PERMISSION_API_PATHS[page]) {
        return PERMISSION_API_PATHS[page]
    }
    if (key.includes(PERMISSION_SEP)) {
        const action = key.slice(key.indexOf(PERMISSION_SEP) + 1)
        if (action === 'view' || action === 'search' || action === 'export') {
            return PERMISSION_API_PATHS[page] || ''
        }
    }
    return ''
}
