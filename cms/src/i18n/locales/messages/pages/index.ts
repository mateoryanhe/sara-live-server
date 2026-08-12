import type {CmsLocale} from '../types'
import {userListMessages} from '../user-list'

import {giftListMessages} from './gift-list'
import {vipCfgListMessages} from './vip-cfg-list'
import {anchorListMessages} from './anchor-list'
import {botAnchorListMessages} from './bot-anchor-list'
import {rechargeOrderListMessages} from './recharge-order-list'
import {currencyLogListMessages} from './currency-log-list'
import {banUserMessages} from './ban-user'
import {bannerListMessages} from './banner-list'
import {activityMessageListMessages} from './activity-message-list'
import {guildListMessages} from './guild-list'
import {guildProfileMessages} from './guild-profile'
import {rechargeCfgListMessages} from './recharge-cfg-list'
import {appPkgListMessages} from './app-pkg-list'
import {randomNicknameCfgMessages} from './random-nickname-cfg'
import {customerServiceCfgMessages} from './customer-service-cfg'
import {walletExchangeCfgMessages} from './wallet-exchange-cfg'
import {agoraCfgMessages} from './agora-cfg'
import {ticketListMessages} from './ticket-list'
import {billingListMessages} from './billing-list'
import {liveConfigMessages} from './live-config'
import {liveRoomTagListMessages} from './live-room-tag-list'
import {revenueLogListMessages} from './revenue-log-list'
import {liveRecordListMessages} from './live-record-list'
import {videoCallLogListMessages} from './video-call-log-list'
import {roleListMessages} from './role-list'
import {cmsUserListMessages} from './cmsuser-list'
import {moduleListMessages} from './module-list'
import {gameListMessages} from './game-list'
import {gamePlatformCfgMessages} from './game-platform-cfg'
import {gameBetLogListMessages} from './game-bet-log-list'
import {gameWinLogListMessages} from './game-win-log-list'
import {shortVideoListMessages} from './short-video-list'
import {shortVideoCfgMessages} from './short-video-cfg'
import {shortVideoCategoryListMessages} from './short-video-category-list'
import {shortVideoWatchListMessages} from './short-video-watch-list'
import {accountCfgMessages} from './account-cfg'
import {appTokenMessages} from './app-token'
import {preloadCfgMessages} from './preload-cfg'
import {textModerationMessages} from './text-moderation'
import {privacyPolicyMessages} from './privacy-policy'
import {googlePlayMessages} from './google-play'
import {uploadResourceMessages} from './upload-resource'
import {dataSyncMessages} from './data-sync'
import {resourceMonitorMessages} from './resource-monitor'
import {serverLogExplorerMessages} from './server-log-explorer'
import {dashboardMessages} from './dashboard'

const pageMessageBuilders = [
  ['userList', userListMessages],
  ['giftList', giftListMessages],
  ['vipCfgList', vipCfgListMessages],
  ['anchorList', anchorListMessages],
  ['botAnchorList', botAnchorListMessages],
  ['rechargeOrderList', rechargeOrderListMessages],
  ['currencyLogList', currencyLogListMessages],
  ['banUser', banUserMessages],
  ['bannerList', bannerListMessages],
  ['activityMessageList', activityMessageListMessages],
  ['guildList', guildListMessages],
  ['guildProfile', guildProfileMessages],
  ['rechargeCfgList', rechargeCfgListMessages],
  ['appPkgList', appPkgListMessages],
  ['randomNicknameCfg', randomNicknameCfgMessages],
  ['customerServiceCfg', customerServiceCfgMessages],
  ['walletExchangeCfg', walletExchangeCfgMessages],
  ['agoraCfg', agoraCfgMessages],
  ['ticketList', ticketListMessages],
  ['billingList', billingListMessages],
  ['liveConfig', liveConfigMessages],
  ['liveRoomTagList', liveRoomTagListMessages],
  ['revenueLogList', revenueLogListMessages],
  ['liveRecordList', liveRecordListMessages],
  ['videoCallLogList', videoCallLogListMessages],
  ['roleList', roleListMessages],
  ['cmsUserList', cmsUserListMessages],
  ['moduleList', moduleListMessages],
  ['gameList', gameListMessages],
  ['gamePlatformCfg', gamePlatformCfgMessages],
  ['gameBetLogList', gameBetLogListMessages],
  ['gameWinLogList', gameWinLogListMessages],
  ['shortVideoList', shortVideoListMessages],
  ['shortVideoCfg', shortVideoCfgMessages],
  ['shortVideoCategoryList', shortVideoCategoryListMessages],
  ['shortVideoWatchList', shortVideoWatchListMessages],
  ['accountCfg', accountCfgMessages],
  ['appToken', appTokenMessages],
  ['preloadCfg', preloadCfgMessages],
  ['textModeration', textModerationMessages],
  ['privacyPolicy', privacyPolicyMessages],
  ['googlePlay', googlePlayMessages],
  ['uploadResource', uploadResourceMessages],
  ['dataSync', dataSyncMessages],
  ['resourceMonitor', resourceMonitorMessages],
  ['serverLogExplorer', serverLogExplorerMessages],
  ['dashboard', dashboardMessages],
] as const

export type PageNamespace = typeof pageMessageBuilders[number][0]

export function buildPageMessages(locale: CmsLocale): Record<PageNamespace, Record<string, string>> {
  const pages = {} as Record<PageNamespace, Record<string, string>>
  for (const [key, messages] of pageMessageBuilders) {
    pages[key] = messages[locale]
  }
  return pages
}
