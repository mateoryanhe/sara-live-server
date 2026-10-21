import type {CmsLocale} from '../../types'
import {convertLocaleDeep} from '../../zh-tw-convert'
import {userListMessages} from '../user-list'

import {giftListMessages} from './gift-list'
import {luckyGiftMultiplierCfgMessages} from './lucky-gift-multiplier-cfg'
import {vipCfgListMessages} from './vip-cfg-list'
import {anchorListMessages} from './anchor-list'
import {liveRoomRecycleBinMessages} from './live-room-recycle-bin'
import {oneToOneRoomListMessages} from './one-to-one-room-list'
import {botAnchorListMessages} from './bot-anchor-list'
import {rechargeOrderListMessages} from './recharge-order-list'
import {currencyLogListMessages} from './currency-log-list'
import {coinMerchantTransferLogListMessages} from './coin-merchant-transfer-log-list'
import {banUserMessages} from './ban-user'
import {bannerListMessages} from './banner-list'
import {activityMessageListMessages} from './activity-message-list'
import {guildListMessages} from './guild-list'
import {guildImmediateSettlementMessages} from './guild-immediate-settlement'
import {platformAnchorImmediateSettlementMessages} from './platform-anchor-immediate-settlement'
import {platformAnchorResetUnsettledIncomeMessages} from './platform-anchor-reset-unsettled-income'
import {guildResetUnsettledIncomeMessages} from './guild-reset-unsettled-income'
import {guildRecycleBinMessages} from './guild-recycle-bin'
import {guildAnchorDailyLiveListMessages} from './guild-anchor-daily-live-list'
import {guildAnchorImportResultMessages} from './guild-anchor-import-result'
import {guildMembersMessages} from './guild-members'
import {guildProfileMessages} from './guild-profile'
import {guildProfileMembersMessages} from './guild-profile-members'
import {guildCmsUserListMessages} from './guild-cms-user-list'
import {guildVisibilityListMessages} from './guild-visibility-list'
import {platformAnchorVisibilityListMessages} from './platform-anchor-visibility-list'
import {coinMerchantListMessages} from './coin-merchant-list'
import {coinMerchantRechargeCfgListMessages} from './coin-merchant-recharge-cfg-list'
import {guildAnchorIncomeSettlementLogListMessages} from './guild-anchor-income-settlement-log-list'
import {rechargeCfgListMessages} from './recharge-cfg-list'
import {appPkgListMessages} from './app-pkg-list'
import {randomNicknameCfgMessages} from './random-nickname-cfg'
import {customerServiceCfgMessages} from './customer-service-cfg'
import {walletExchangeCfgMessages} from './wallet-exchange-cfg'
import {effectiveLiveCfgMessages} from './effective-live-cfg'
import {platformAnchorSettlementCfgMessages} from './platform-anchor-settlement-cfg'
import {luckyGiftSettlementCfgMessages} from './lucky-gift-settlement-cfg'
import {paymentCountryCfgMessages} from './payment-country-cfg'
import {agoraCfgMessages} from './agora-cfg'
import {billingListMessages} from './billing-list'
import {liveConfigMessages} from './live-config'
import {liveRoomTagListMessages} from './live-room-tag-list'
import {revenueLogListMessages} from './revenue-log-list'
import {liveRecordListMessages} from './live-record-list'
import {liveDailyEffectiveLiveListMessages} from './live-daily-effective-live-list'
import {liveWeeklyUnsettledLiveListMessages} from './live-weekly-unsettled-live-list'
import {videoCallLogListMessages} from './video-call-log-list'
import {roleListMessages} from './role-list'
import {cmsUserListMessages} from './cmsuser-list'
import {moduleListMessages} from './module-list'
import {gameListMessages} from './game-list'
import {gameShelfListMessages} from './game-shelf-list'
import {gamePlatformCfgMessages} from './game-platform-cfg'
import {gameBetLogListMessages} from './game-bet-log-list'
import {gameWinLogListMessages} from './game-win-log-list'
import {shortVideoListMessages} from './short-video-list'
import {shortVideoCfgMessages} from './short-video-cfg'
import {shortVideoCategoryListMessages} from './short-video-category-list'
import {shortVideoPriceTierListMessages} from './short-video-price-tier-list'
import {shortVideoWatchListMessages} from './short-video-watch-list'
import {shortVideoAuthorSettlementLogListMessages} from './short-video-author-settlement-log-list'
import {accountCfgMessages} from './account-cfg'
import {deviceRegisterRiskCfgMessages} from './device-register-risk-cfg'
import {userWithdrawCfgMessages} from './user-withdraw-cfg'
import {userWithdrawLogListMessages} from './user-withdraw-log-list'
import {userWithdrawTierListMessages} from './user-withdraw-tier-list'
import {appVersionCfgMessages} from './app-version-cfg'
import {firstRechargeActivityCfgMessages} from './first-recharge-activity-cfg'
import {inviteRechargeRewardCfgMessages} from './invite-recharge-reward-cfg'
import {promotionChannelUrlListMessages} from './promotion-channel-url-list'
import {promotionChannelDailyStatListMessages} from './promotion-channel-daily-stat-list'
import {simulatorCpuKeywordListMessages} from './simulator-cpu-keyword-list'
import {simulatorDeviceWhitelistListMessages} from './simulator-device-whitelist-list'
import {anchorSalarySocialShareCfgListMessages} from './anchor-salary-social-share-cfg-list'
import {anchorNoSalaryShareCfgMessages} from './anchor-no-salary-share-cfg'
import {anchorGameShareCfgListMessages} from './anchor-game-share-cfg-list'
import {anchorIncomeSettlementLogListMessages} from './anchor-income-settlement-log-list'
import {guildIncomeSettlementLogListMessages} from './guild-income-settlement-log-list'
import {guildTransferListMessages} from './guild-transfer-list'
import {coinMerchantGuildListMessages} from './coin-merchant-guild-list'
import {coinMerchantGuildTransferListMessages} from './coin-merchant-guild-transfer-list'
import {appTokenMessages} from './app-token'
import {serverRuntimeCfgMessages} from './server-runtime-cfg'
import {staticCacheCfgMessages} from './static-cache-cfg'
import {preloadCfgMessages} from './preload-cfg'
import {textModerationMessages} from './text-moderation'
import {privacyPolicyMessages} from './privacy-policy'
import {googlePlayMessages} from './google-play'
import {firebaseMessages} from './firebase'
import {metaPixelMessages} from './meta-pixel'
import {haipayMessages} from './haipay'
import {cfEmailMessages} from './cf-email'
import {dbBackupMessages} from './db-backup'
import {uploadResourceMessages} from './upload-resource'
import {cmsDomainMappingMessages} from './cms-domain-mapping'
import {cmsDomainAliasMessages} from './cms-domain-alias'
import {h5LiveDeployMessages} from './h5-live-deploy'
import {coinMerchantDeployMessages} from './coin-merchant-deploy'
import {countryFlagDeployMessages} from './country-flag-deploy'
import {thirdPayDeployMessages} from './third-pay-deploy'
import {userWithdrawDeployMessages} from './user-withdraw-deploy'
import {officialSiteDeployMessages} from './official-site-deploy'
import {dataSyncMessages} from './data-sync'
import {resourceMonitorMessages} from './resource-monitor'
import {serverLogExplorerMessages} from './server-log-explorer'
import {dashboardMessages} from './dashboard'
import {trackingHotLiveRoomJoinMessages} from './tracking-hot-live-room-join'
import {trackingCall1v1InitiateMessages} from './tracking-call-1v1-initiate'
import {trackingMiniGameRoundStartMessages} from './tracking-mini-game-round-start'
import {trackingShortVideoMessages} from './tracking-short-video'
import {trackingVoiceChatRoomMessages} from './tracking-voice-chat-room'
import {trackingRechargeMessages} from './tracking-recharge'

const pageMessageBuilders = [
  ['userList', userListMessages],
  ['giftList', giftListMessages],
  ['luckyGiftMultiplierCfg', luckyGiftMultiplierCfgMessages],
  ['vipCfgList', vipCfgListMessages],
  ['anchorList', anchorListMessages],
  ['liveRoomRecycleBin', liveRoomRecycleBinMessages],
  ['oneToOneRoomList', oneToOneRoomListMessages],
  ['botAnchorList', botAnchorListMessages],
  ['rechargeOrderList', rechargeOrderListMessages],
  ['currencyLogList', currencyLogListMessages],
  ['coinMerchantTransferLogList', coinMerchantTransferLogListMessages],
  ['banUser', banUserMessages],
  ['bannerList', bannerListMessages],
  ['activityMessageList', activityMessageListMessages],
  ['guildList', guildListMessages],
  ['guildImmediateSettlement', guildImmediateSettlementMessages],
  ['platformAnchorImmediateSettlement', platformAnchorImmediateSettlementMessages],
  ['platformAnchorResetUnsettledIncome', platformAnchorResetUnsettledIncomeMessages],
  ['guildResetUnsettledIncome', guildResetUnsettledIncomeMessages],
  ['guildRecycleBin', guildRecycleBinMessages],
  ['guildAnchorDailyLiveList', guildAnchorDailyLiveListMessages],
  ['guildAnchorImportResult', guildAnchorImportResultMessages],
  ['guildMembers', guildMembersMessages],
  ['guildProfile', guildProfileMessages],
  ['guildProfileMembers', guildProfileMembersMessages],
  ['guildCmsUserList', guildCmsUserListMessages],
  ['guildVisibilityList', guildVisibilityListMessages],
  ['platformAnchorVisibilityList', platformAnchorVisibilityListMessages],
  ['coinMerchantList', coinMerchantListMessages],
  ['coinMerchantRechargeCfgList', coinMerchantRechargeCfgListMessages],
  ['coinMerchantGuildList', coinMerchantGuildListMessages],
  ['guildAnchorIncomeSettlementLogList', guildAnchorIncomeSettlementLogListMessages],
  ['rechargeCfgList', rechargeCfgListMessages],
  ['appPkgList', appPkgListMessages],
  ['randomNicknameCfg', randomNicknameCfgMessages],
  ['customerServiceCfg', customerServiceCfgMessages],
  ['walletExchangeCfg', walletExchangeCfgMessages],
  ['effectiveLiveCfg', effectiveLiveCfgMessages],
  ['platformAnchorSettlementCfg', platformAnchorSettlementCfgMessages],
  ['luckyGiftSettlementCfg', luckyGiftSettlementCfgMessages],
  ['paymentCountryCfg', paymentCountryCfgMessages],
  ['agoraCfg', agoraCfgMessages],
  ['billingList', billingListMessages],
  ['liveConfig', liveConfigMessages],
  ['liveRoomTagList', liveRoomTagListMessages],
  ['revenueLogList', revenueLogListMessages],
  ['liveRecordList', liveRecordListMessages],
  ['liveDailyEffectiveLiveList', liveDailyEffectiveLiveListMessages],
  ['liveWeeklyUnsettledLiveList', liveWeeklyUnsettledLiveListMessages],
  ['videoCallLogList', videoCallLogListMessages],
  ['roleList', roleListMessages],
  ['cmsUserList', cmsUserListMessages],
  ['moduleList', moduleListMessages],
  ['gameList', gameListMessages],
  ['gameShelfList', gameShelfListMessages],
  ['gamePlatformCfg', gamePlatformCfgMessages],
  ['gameBetLogList', gameBetLogListMessages],
  ['gameWinLogList', gameWinLogListMessages],
  ['shortVideoList', shortVideoListMessages],
  ['shortVideoCfg', shortVideoCfgMessages],
  ['shortVideoCategoryList', shortVideoCategoryListMessages],
  ['shortVideoPriceTierList', shortVideoPriceTierListMessages],
  ['shortVideoWatchList', shortVideoWatchListMessages],
  ['shortVideoAuthorSettlementLogList', shortVideoAuthorSettlementLogListMessages],
  ['accountCfg', accountCfgMessages],
  ['deviceRegisterRiskCfg', deviceRegisterRiskCfgMessages],
  ['userWithdrawCfg', userWithdrawCfgMessages],
  ['userWithdrawLogList', userWithdrawLogListMessages],
  ['userWithdrawTierList', userWithdrawTierListMessages],
  ['appVersionCfg', appVersionCfgMessages],
  ['firstRechargeActivityCfg', firstRechargeActivityCfgMessages],
  ['inviteRechargeRewardCfg', inviteRechargeRewardCfgMessages],
  ['promotionChannelUrlList', promotionChannelUrlListMessages],
  ['promotionChannelDailyStatList', promotionChannelDailyStatListMessages],
  ['simulatorCpuKeywordList', simulatorCpuKeywordListMessages],
  ['simulatorDeviceWhitelistList', simulatorDeviceWhitelistListMessages],
  ['anchorSalarySocialShareCfgList', anchorSalarySocialShareCfgListMessages],
  ['anchorNoSalaryShareCfg', anchorNoSalaryShareCfgMessages],
  ['anchorGameShareCfgList', anchorGameShareCfgListMessages],
  ['anchorIncomeSettlementLogList', anchorIncomeSettlementLogListMessages],
  ['guildIncomeSettlementLogList', guildIncomeSettlementLogListMessages],
  ['guildTransferList', guildTransferListMessages],
  ['coinMerchantGuildTransferList', coinMerchantGuildTransferListMessages],
  ['appToken', appTokenMessages],
  ['preloadCfg', preloadCfgMessages],
  ['serverRuntimeCfg', serverRuntimeCfgMessages],
  ['staticCacheCfg', staticCacheCfgMessages],
  ['textModeration', textModerationMessages],
  ['privacyPolicy', privacyPolicyMessages],
  ['googlePlay', googlePlayMessages],
  ['firebase', firebaseMessages],
  ['haipay', haipayMessages],
  ['metaPixel', metaPixelMessages],
  ['cfEmail', cfEmailMessages],
  ['dbBackup', dbBackupMessages],
  ['uploadResource', uploadResourceMessages],
  ['cmsDomainMapping', cmsDomainMappingMessages],
  ['cmsDomainAlias', cmsDomainAliasMessages],
  ['h5LiveDeploy', h5LiveDeployMessages],
  ['coinMerchantDeploy', coinMerchantDeployMessages],
  ['countryFlagDeploy', countryFlagDeployMessages],
  ['thirdPayDeploy', thirdPayDeployMessages],
  ['userWithdrawDeploy', userWithdrawDeployMessages],
  ['officialSiteDeploy', officialSiteDeployMessages],
  ['dataSync', dataSyncMessages],
  ['resourceMonitor', resourceMonitorMessages],
  ['serverLogExplorer', serverLogExplorerMessages],
  ['dashboard', dashboardMessages],
  ['trackingHotLiveRoomJoin', trackingHotLiveRoomJoinMessages],
  ['trackingCall1v1Initiate', trackingCall1v1InitiateMessages],
  ['trackingMiniGameRoundStart', trackingMiniGameRoundStartMessages],
  ['trackingShortVideo', trackingShortVideoMessages],
  ['trackingVoiceChatRoom', trackingVoiceChatRoomMessages],
  ['trackingRecharge', trackingRechargeMessages],
] as const

export type PageNamespace = typeof pageMessageBuilders[number][0]

type PageMessageBundle = Record<'zh-CN' | 'en' | 'es' | 'pt' | 'hi' | 'id', Record<string, string>>

export function buildPageMessages(locale: CmsLocale): Record<PageNamespace, Record<string, string>> {
  const pages = {} as Record<PageNamespace, Record<string, string>>
  for (const [key, messages] of pageMessageBuilders) {
    const bundle = messages as PageMessageBundle
    if (locale === 'zh-TW') {
      pages[key] = convertLocaleDeep(bundle['zh-CN'])
    } else {
      pages[key] = bundle[locale]
    }
  }
  return pages
}
