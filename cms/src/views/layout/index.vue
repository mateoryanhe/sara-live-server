<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside class="aside" width="220px">
      <div class="logo">{{ t('common.logo') }}</div>
      <el-menu
          :collapse="isCollapse"
          :default-active="activeMenu"
          :router="true"
          :unique-opened="true"
          class="sidebar-menu"
      >
        <!-- 仪表盘 -->
        <el-menu-item v-if="hasMenuPermission('Dashboard')" index="/dashboard">
          <el-icon>
            <Odometer/>
          </el-icon>
          <span>{{ t('menu.Dashboard') }}</span>
        </el-menu-item>
        <el-sub-menu
            v-if="hasMenuPermission('UserList') || hasMenuPermission('AnchorListManagement') || hasMenuPermission('LiveRoomRecycleBinManagement') || hasMenuPermission('BotAnchorManagement') || hasMenuPermission('RechargeOrderList')"
            index="/user/account">
          <template #title>
            <el-icon>
              <User/>
            </el-icon>
            <span>{{ t('menu.UserManagement') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('UserList')" index="/user/account/user-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>{{ t('menu.UserList') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('AnchorListManagement')" index="/user/anchor/anchor-list">
            <el-icon>
              <VideoPlay/>
            </el-icon>
            <span>{{ t('menu.AnchorListManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('LiveRoomRecycleBinManagement')" index="/user/anchor/live-room-recycle-bin">
            <el-icon>
              <Delete/>
            </el-icon>
            <span>{{ t('menu.LiveRoomRecycleBinManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('BotAnchorManagement')" index="/user/bot-anchor/bot-anchor-list">
            <el-icon>
              <Cpu/>
            </el-icon>
            <span>{{ t('menu.BotAnchorManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('RechargeOrderList')" index="/user/recharge-order/recharge-order-list">
            <el-icon>
              <Wallet/>
            </el-icon>
            <span>{{ t('menu.RechargeOrderList') }}</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('BannerManagement') || hasMenuPermission('ActivityMessageManagement') || hasMenuPermission('GuildManagement') || hasMenuPermission('GuildCMSUserManagement') || hasMenuPermission('GuildAnchorDailyLiveManagement') || hasMenuPermission('PlatformAnchorList') || hasMenuPermission('GuildRecycleBinManagement') || hasMenuPermission('GuildProfileManagement') || hasMenuPermission('RechargeCfgManagement') || hasMenuPermission('VipCfgManagement') || hasMenuPermission('AppPkgManagement') || hasMenuPermission('RandomNicknameManagement') || hasMenuPermission('CustomerServiceCfgManagement') || hasMenuPermission('WalletExchangeCfgManagement') || hasMenuPermission('AnchorSalaryCfgManagement') || hasMenuPermission('LiveRevenueShareCfgManagement')"
            index="/operation">
          <template #title>
            <el-icon>
              <Stamp/>
            </el-icon>
            <span>{{ t('menu.OperationManagement') }}</span>
          </template>
          <el-sub-menu
              v-if="hasMenuPermission('BannerManagement') || hasMenuPermission('ActivityMessageManagement')"
              index="/operation/group/content">
            <template #title>
              <el-icon>
                <Picture/>
              </el-icon>
              <span>{{ t('menu.OperationContentGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('BannerManagement')" index="/operation/banner/banner-list">
              <el-icon>
                <Picture/>
              </el-icon>
              <span>{{ t('menu.BannerManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('ActivityMessageManagement')" index="/operation/activity-message/activity-message-list">
              <el-icon>
                <Bell/>
              </el-icon>
              <span>{{ t('menu.ActivityMessageManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('RechargeCfgManagement') || hasMenuPermission('VipCfgManagement') || hasMenuPermission('WalletExchangeCfgManagement')"
              index="/operation/group/recharge">
            <template #title>
              <el-icon>
                <Wallet/>
              </el-icon>
              <span>{{ t('menu.OperationRechargeGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('RechargeCfgManagement')" index="/operation/recharge/recharge-cfg-list">
              <el-icon>
                <Wallet/>
              </el-icon>
              <span>{{ t('menu.RechargeCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('VipCfgManagement')" index="/operation/vip/vip-cfg-list">
              <el-icon>
                <Medal/>
              </el-icon>
              <span>{{ t('menu.VipCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('WalletExchangeCfgManagement')" index="/operation/wallet/wallet-exchange-cfg">
              <el-icon>
                <Coin/>
              </el-icon>
              <span>{{ t('menu.WalletExchangeCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('GuildManagement') || hasMenuPermission('GuildCMSUserManagement') || hasMenuPermission('GuildAnchorDailyLiveManagement') || hasMenuPermission('GuildRecycleBinManagement') || hasMenuPermission('GuildProfileManagement') || hasMenuPermission('PlatformAnchorList')"
              index="/operation/group/guild">
            <template #title>
              <el-icon>
                <User/>
              </el-icon>
              <span>{{ t('menu.OperationGuildGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('GuildManagement')" index="/operation/guild/guild-list">
              <el-icon>
                <User/>
              </el-icon>
              <span>{{ t('menu.GuildManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('GuildCMSUserManagement')" index="/operation/guild/guild-cms-user-list">
              <el-icon>
                <UserFilled/>
              </el-icon>
              <span>{{ t('menu.GuildCMSUserManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('GuildAnchorDailyLiveManagement')" index="/operation/guild/guild-anchor-daily-live-list">
              <el-icon>
                <Document/>
              </el-icon>
              <span>{{ t('menu.GuildAnchorDailyLiveManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('PlatformAnchorList')" index="/operation/guild/platform-anchor-list">
              <el-icon>
                <VideoPlay/>
              </el-icon>
              <span>{{ t('menu.PlatformAnchorList') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('GuildRecycleBinManagement')" index="/operation/guild/guild-recycle-bin">
              <el-icon>
                <Delete/>
              </el-icon>
              <span>{{ t('menu.GuildRecycleBinManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('GuildProfileManagement')" index="/operation/guild/guild-profile">
              <el-icon>
                <EditPen/>
              </el-icon>
              <span>{{ t('menu.GuildProfileManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('AnchorSalaryCfgManagement') || hasMenuPermission('LiveRevenueShareCfgManagement')"
              index="/operation/group/settlement">
            <template #title>
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.OperationSettlementGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('AnchorSalaryCfgManagement')" index="/operation/salary/anchor-salary-cfg-list">
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.AnchorSalaryCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('LiveRevenueShareCfgManagement')" index="/operation/salary/live-revenue-share-cfg">
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.LiveRevenueShareCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('AppPkgManagement') || hasMenuPermission('RandomNicknameManagement') || hasMenuPermission('CustomerServiceCfgManagement')"
              index="/operation/group/app">
            <template #title>
              <el-icon>
                <Box/>
              </el-icon>
              <span>{{ t('menu.OperationAppGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('AppPkgManagement')" index="/operation/app-pkg/app-pkg-list">
              <el-icon>
                <Box/>
              </el-icon>
              <span>{{ t('menu.AppPkgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('RandomNicknameManagement')" index="/operation/random-nickname/random-nickname-cfg">
              <el-icon>
                <EditPen/>
              </el-icon>
              <span>{{ t('menu.RandomNicknameManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('CustomerServiceCfgManagement')" index="/operation/customer-service/customer-service-cfg">
              <el-icon>
                <Service/>
              </el-icon>
              <span>{{ t('menu.CustomerServiceCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('GiftManagement') || hasMenuPermission('AgoraCfgManagement') || hasMenuPermission('TicketManagement') || hasMenuPermission('PrivateRoomBillingManagement') || hasMenuPermission('LiveCfgManagement') || hasMenuPermission('LiveRoomTagManagement')"
            index="/live">
          <template #title>
            <el-icon>
              <VideoPlay/>
            </el-icon>
            <span>{{ t('menu.LiveManagement') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('GiftManagement')" index="/live/gift/gift-list">
            <el-icon>
              <Present/>
            </el-icon>
            <span>{{ t('menu.GiftManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('AgoraCfgManagement')" index="/live/agora-cfg">
            <el-icon>
              <Setting/>
            </el-icon>
            <span>{{ t('menu.AgoraCfgManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('TicketManagement')" index="/live/ticket/ticket-list">
            <el-icon>
              <Tickets/>
            </el-icon>
            <span>{{ t('menu.TicketManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('PrivateRoomBillingManagement')" index="/live/private-room-billing/billing-list">
            <el-icon>
              <Lock/>
            </el-icon>
            <span>{{ t('menu.PrivateRoomBillingManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('LiveCfgManagement')" index="/live/live-config/live-config">
            <el-icon>
              <VideoCamera/>
            </el-icon>
            <span>{{ t('menu.LiveCfgManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('LiveRoomTagManagement')" index="/live/live-room-tag/live-room-tag-list">
            <el-icon>
              <CollectionTag/>
            </el-icon>
            <span>{{ t('menu.LiveRoomTagManagement') }}</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('LiveRevenueLogList') || hasMenuPermission('LiveRecordList') || hasMenuPermission('VideoCallLogList') || hasMenuPermission('GoldCurrencyLogList') || hasMenuPermission('DiamondCurrencyLogList') || hasMenuPermission('AnchorIncomeSettlementLogList') || hasMenuPermission('GuildIncomeSettlementLogList')"
            index="/log">
          <template #title>
            <el-icon>
              <Document/>
            </el-icon>
            <span>{{ t('menu.LogManagement') }}</span>
          </template>
          <el-sub-menu index="/log/live">
            <template #title>
              <el-icon>
                <VideoPlay/>
              </el-icon>
              <span>{{ t('menu.LiveLogGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('LiveRevenueLogList')" index="/log/live/revenue-log-list">
              <el-icon>
                <Present/>
              </el-icon>
              <span>{{ t('menu.LiveRevenueLogList') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('LiveRecordList')" index="/log/live/live-record-list">
              <el-icon>
                <Monitor/>
              </el-icon>
              <span>{{ t('menu.LiveRecordList') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="/log/call">
            <template #title>
              <el-icon>
                <VideoCamera/>
              </el-icon>
              <span>{{ t('menu.CallLogGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('VideoCallLogList')" index="/log/call/video-call-log-list">
              <el-icon>
                <VideoCamera/>
              </el-icon>
              <span>{{ t('menu.VideoCallLogList') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('GoldCurrencyLogList') || hasMenuPermission('DiamondCurrencyLogList')"
              index="/log/user">
            <template #title>
              <el-icon>
                <Coin/>
              </el-icon>
              <span>{{ t('menu.UserLogGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('GoldCurrencyLogList')" index="/user/currency-log/gold-log-list">
              <el-icon>
                <Coin/>
              </el-icon>
              <span>{{ t('menu.GoldCurrencyLogList') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('DiamondCurrencyLogList')" index="/user/currency-log/diamond-log-list">
              <el-icon>
                <Money/>
              </el-icon>
              <span>{{ t('menu.DiamondCurrencyLogList') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('AnchorIncomeSettlementLogList') || hasMenuPermission('GuildIncomeSettlementLogList')"
              index="/log/settlement">
            <template #title>
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.SettlementLogGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('AnchorIncomeSettlementLogList')" index="/operation/salary/anchor-income-settlement-log-list">
              <el-icon>
                <Document/>
              </el-icon>
              <span>{{ t('menu.AnchorIncomeSettlementLogList') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('GuildIncomeSettlementLogList')" index="/operation/salary/guild-income-settlement-log-list">
              <el-icon>
                <Document/>
              </el-icon>
              <span>{{ t('menu.GuildIncomeSettlementLogList') }}</span>
            </el-menu-item>
          </el-sub-menu>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('ShortVideoManagement') || hasMenuPermission('ShortVideoCategoryManagement') || hasMenuPermission('ShortVideoPriceTierManagement') || hasMenuPermission('ShortVideoCfgManagement') || hasMenuPermission('ShortVideoWatchManagement')"
            index="/shortvideo">
          <template #title>
            <el-icon>
              <VideoCamera/>
            </el-icon>
            <span>{{ t('menu.ShortVideoGroup') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('ShortVideoManagement')" index="/shortvideo/short-video-list">
            <el-icon>
              <VideoCamera/>
            </el-icon>
            <span>{{ t('menu.ShortVideoManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('ShortVideoCategoryManagement')" index="/shortvideo/short-video-category-list">
            <el-icon>
              <Collection/>
            </el-icon>
            <span>{{ t('menu.ShortVideoCategoryManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('ShortVideoPriceTierManagement')" index="/shortvideo/short-video-price-tier-list">
            <el-icon>
              <Money/>
            </el-icon>
            <span>{{ t('menu.ShortVideoPriceTierManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('ShortVideoCfgManagement')" index="/shortvideo/short-video-cfg">
            <el-icon>
              <Setting/>
            </el-icon>
            <span>{{ t('menu.ShortVideoCfgManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('ShortVideoWatchManagement')" index="/shortvideo/short-video-watch-list">
            <el-icon>
              <View/>
            </el-icon>
            <span>{{ t('menu.ShortVideoWatchManagement') }}</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('GamePlatformCfgManagement') || hasMenuPermission('GameVendorGameListManagement') || hasMenuPermission('GameShelfListManagement') || hasMenuPermission('GameBetLogListManagement') || hasMenuPermission('GameWinLogListManagement')"
            index="/game">
          <template #title>
            <el-icon>
              <Cpu/>
            </el-icon>
            <span>{{ t('menu.GameManagement') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('GameShelfListManagement')" index="/game/game-shelf-list">
            <el-icon>
              <Collection/>
            </el-icon>
            <span>{{ t('menu.GameShelfListManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GamePlatformCfgManagement')" index="/game/game-platform-cfg">
            <el-icon>
              <Setting/>
            </el-icon>
            <span>{{ t('menu.GamePlatformCfgManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GameVendorGameListManagement')" index="/game/game-list">
            <el-icon>
              <List/>
            </el-icon>
            <span>{{ t('menu.GameVendorGameListManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GameBetLogListManagement')" index="/game/game-bet-log-list">
            <el-icon>
              <Money/>
            </el-icon>
            <span>{{ t('menu.GameBetLogListManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GameWinLogListManagement')" index="/game/game-win-log-list">
            <el-icon>
              <Coin/>
            </el-icon>
            <span>{{ t('menu.GameWinLogListManagement') }}</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('FirstRechargeActivityManagement')"
            index="/activity">
          <template #title>
            <el-icon>
              <Present/>
            </el-icon>
            <span>{{ t('menu.ActivityManagement') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('FirstRechargeActivityManagement')" index="/activity/first-recharge-activity-cfg">
            <el-icon>
              <Present/>
            </el-icon>
            <span>{{ t('menu.FirstRechargeActivityManagement') }}</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu
            v-if="hasMenuPermission('AppTokenConfig') || hasMenuPermission('AccountCfgManagement') || hasMenuPermission('SimulatorCpuKeywordManagement') || hasMenuPermission('ServerRuntimeCfgManagement') || hasMenuPermission('PreloadCfgManagement') || hasMenuPermission('TextModerationCfgManagement') || hasMenuPermission('PrivacyPolicyCfgManagement') || hasMenuPermission('GooglePlayCfgManagement') || hasMenuPermission('UploadResourceCfgManagement') || hasMenuPermission('DataSyncCfgManagement') || hasMenuPermission('ResourceMonitor') || hasMenuPermission('ServerLogExplorer')"
            index="/config">
          <template #title>
            <el-icon>
              <Setting/>
            </el-icon>
            <span>{{ t('menu.ConfigManagement') }}</span>
          </template>
          <el-sub-menu
              v-if="hasMenuPermission('AppTokenConfig') || hasMenuPermission('AccountCfgManagement') || hasMenuPermission('ServerRuntimeCfgManagement') || hasMenuPermission('PreloadCfgManagement')"
              index="/config/group/basic">
            <template #title>
              <el-icon>
                <Key/>
              </el-icon>
              <span>{{ t('menu.ConfigBasicGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('AppTokenConfig')" index="/config/app-token">
              <el-icon>
                <Key/>
              </el-icon>
              <span>{{ t('menu.AppTokenConfig') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('AccountCfgManagement')" index="/config/account-cfg">
              <el-icon>
                <User/>
              </el-icon>
              <span>{{ t('menu.AccountCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('ServerRuntimeCfgManagement') || hasMenuPermission('PreloadCfgManagement')" index="/config/server-runtime-cfg">
              <el-icon>
                <Cpu/>
              </el-icon>
              <span>{{ t('menu.ServerRuntimeCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('SimulatorCpuKeywordManagement') || hasMenuPermission('TextModerationCfgManagement') || hasMenuPermission('PrivacyPolicyCfgManagement')"
              index="/config/group/security">
            <template #title>
              <el-icon>
                <Lock/>
              </el-icon>
              <span>{{ t('menu.ConfigSecurityGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('SimulatorCpuKeywordManagement')" index="/config/simulator-cpu-keyword">
              <el-icon>
                <Cpu/>
              </el-icon>
              <span>{{ t('menu.SimulatorCpuKeywordManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('TextModerationCfgManagement')" index="/config/text-moderation">
              <el-icon>
                <Document/>
              </el-icon>
              <span>{{ t('menu.TextModerationCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('PrivacyPolicyCfgManagement')" index="/config/privacy-policy">
              <el-icon>
                <Document/>
              </el-icon>
              <span>{{ t('menu.PrivacyPolicyCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('GooglePlayCfgManagement') || hasMenuPermission('UploadResourceCfgManagement') || hasMenuPermission('DataSyncCfgManagement')"
              index="/config/group/platform">
            <template #title>
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.ConfigPlatformGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('GooglePlayCfgManagement')" index="/config/google-play">
              <el-icon>
                <CreditCard/>
              </el-icon>
              <span>{{ t('menu.GooglePlayCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('UploadResourceCfgManagement')" index="/config/upload-resource">
              <el-icon>
                <Picture/>
              </el-icon>
              <span>{{ t('menu.UploadResourceCfgManagement') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('DataSyncCfgManagement')" index="/config/data-sync">
              <el-icon>
                <Refresh/>
              </el-icon>
              <span>{{ t('menu.DataSyncCfgManagement') }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-sub-menu
              v-if="hasMenuPermission('ResourceMonitor') || hasMenuPermission('ServerLogExplorer')"
              index="/config/group/ops">
            <template #title>
              <el-icon>
                <Monitor/>
              </el-icon>
              <span>{{ t('menu.ConfigOpsGroup') }}</span>
            </template>
            <el-menu-item v-if="hasMenuPermission('ResourceMonitor')" index="/config/resource-monitor">
              <el-icon>
                <Monitor/>
              </el-icon>
              <span>{{ t('menu.ResourceMonitor') }}</span>
            </el-menu-item>
            <el-menu-item v-if="hasMenuPermission('ServerLogExplorer')" index="/config/server-log">
              <el-icon>
                <Search/>
              </el-icon>
              <span>{{ t('menu.ServerLogExplorer') }}</span>
            </el-menu-item>
          </el-sub-menu>
        </el-sub-menu>
        <!-- 角色权限管理菜单 -->
        <el-sub-menu
            v-if="hasMenuPermission('RoleManagement') || hasMenuPermission('ModuleList') || hasMenuPermission('CMSUserManagement')"
            index="/role">
          <template #title>
            <el-icon>
              <Lock/>
            </el-icon>
            <span>{{ t('menu.RoleManagementGroup') }}</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('RoleManagement')" index="/role/role-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>{{ t('menu.RoleManagement') }}</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('CMSUserManagement')" index="/role/cmsuser-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>{{ t('menu.CMSUserManagement') }}</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 主内容区域 -->
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button class="collapse-btn" @click="toggleCollapse">
            <el-icon>
              <Fold v-if="!isCollapse"/>
              <Expand v-else/>
            </el-icon>
          </el-button>
        </div>
        <div class="header-right">
          <LanguageSwitcher class="header-lang" compact show-label/>
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ username }}
              <el-icon class="el-icon--right">
                <arrow-down/>
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="logout">{{ t('common.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <LayoutTabs/>
      <el-main class="main-content">
        <router-view v-slot="{ Component, route: currentRoute }">
          <keep-alive :max="15">
            <component :is="Component" v-if="Component" :key="currentRoute.path"/>
          </keep-alive>
        </router-view>
      </el-main>
      <el-footer class="footer">
        <div class="footer-content">
          {{ t('common.footer') }} &copy; {{ new Date().getFullYear() }}
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import {computed, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import LayoutTabs from '@/components/layout/LayoutTabs.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'
import {useLayoutTabs} from '@/composables/useLayoutTabs'
import {ArrowDown, Bell, Box, Coin, Collection, CollectionTag, Cpu, CreditCard, Delete, Document, EditPen, Expand, Fold, Key, List, Lock, Medal, Money, Monitor, Odometer, Picture, Present, Refresh, Search, Service, Setting, Stamp, Tickets, User, UserFilled, VideoCamera, VideoPlay, View, Wallet} from '@element-plus/icons-vue'
import {getIsAdmin, hasPermission} from '@/utils/permission'
import {clearAuthSession} from '@/utils/auth'

const route = useRoute()
const router = useRouter()
const {t} = useI18n()
const isCollapse = ref(false)
const {addTab} = useLayoutTabs()

watch(
    () => route.fullPath,
    () => {
      addTab(route)
    },
    {immediate: true},
)

const activeMenu = computed(() => {
  const {path} = route
  return path
})

const username = computed(() => {
  // 从localStorage获取用户名，如果不存在则显示默认值
  return localStorage.getItem('username') || t('common.admin')
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const logout = () => {
  clearAuthSession()
  router.push('/login')
}

// 检查菜单项是否有权限
const hasMenuPermission = (moduleName: string) => {
  // 管理员拥有所有权限
  if (getIsAdmin()) {
    return true
  }

  // 检查是否有访问该模块的权限
  return hasPermission(moduleName)
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #ffffff;
  color: #333;
  height: 100vh;
  overflow: hidden;
  box-shadow: 1px 0 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #333;
  font-size: 18px;
  font-weight: 600;
  border-bottom: 1px solid #f0f0f0;
  background: #fff;
  transition: all 0.3s ease;
}

.sidebar-menu {
  border: none;
  height: calc(100% - 60px);
  background-color: #ffffff;
  overflow-y: auto;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-lang {
  margin-right: 4px;
}

.collapse-btn {
  margin-right: 20px;
  font-size: 16px;
  color: #409eff;
}

.main-content {
  background-color: #f5f7f8;
  padding: 20px;
  overflow-y: auto;
}

.footer {
  background-color: #fafafa;
  border-top: 1px solid #e6e6e6;
  padding: 0;
  height: 50px;
}

.footer-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  display: flex;
  align-items: center;
  font-size: 14px;
}

/* 简约风格的菜单样式 */
.el-menu {
  border-right: none;
  background-color: #ffffff;
}

.el-sub-menu__title {
  color: #333;
  font-size: 14px;
  height: 48px;
  line-height: 48px;
  display: flex;
  align-items: center;
  padding-left: 20px !important;
  background-color: #fff;
}

.el-sub-menu__title:hover {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

.el-menu-item {
  color: #666;
  font-size: 14px;
  height: 42px;
  line-height: 42px;
  display: flex;
  align-items: center;
  padding-left: 50px !important;
  background-color: #fff;
}

.el-menu-item:hover {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

.el-menu-item.is-active {
  background-color: #ecf5ff !important;
  color: #409eff !important;
  border-left: 3px solid #409eff;
  font-weight: 500;
}

/* 子菜单项样式 */
.el-sub-menu .el-menu-item {
  color: #666;
  padding-left: 65px !important;
  height: 40px;
  line-height: 40px;
}

.el-sub-menu .el-menu-item:hover {
  color: #409eff !important;
}

.el-sub-menu .el-menu-item.is-active {
  color: #409eff !important;
}

/* 菜单收起时的样式 */
.el-aside:not(.el-menu--collapse) {
  width: 220px;
  height: 100%;
}

/* 滚动条样式 */
.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background-color: #e0e0e0;
  border-radius: 2px;
}

.sidebar-menu::-webkit-scrollbar-track {
  background-color: #fff;
}
</style>