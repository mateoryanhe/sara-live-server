import fs from 'node:fs'
import path from 'node:path'

const pages = [
  'gift-list', 'vip-cfg-list', 'anchor-list', 'bot-anchor-list', 'recharge-order-list', 'currency-log-list', 'ban-user',
  'banner-list', 'activity-message-list', 'guild-list', 'guild-profile', 'recharge-cfg-list', 'app-pkg-list',
  'random-nickname-cfg', 'customer-service-cfg', 'wallet-exchange-cfg', 'agora-cfg', 'ticket-list', 'billing-list',
  'live-config', 'live-room-tag-list', 'revenue-log-list', 'live-record-list', 'video-call-log-list', 'role-list',
  'cmsuser-list', 'module-list', 'game-list', 'game-platform-cfg', 'game-bet-log-list', 'game-win-log-list',
  'short-video-list', 'short-video-cfg', 'short-video-category-list', 'short-video-watch-list', 'account-cfg',
  'app-token', 'preload-cfg', 'text-moderation', 'privacy-policy', 'google-play', 'upload-resource', 'data-sync',
  'resource-monitor', 'server-log-explorer', 'dashboard',
]

const dir = 'src/i18n/locales/messages/pages'

for (const file of pages) {
  const camel = file.replace(/-([a-z])/g, (_, c) => c.toUpperCase())
  const exportName = `${camel}Messages`
  const target = path.join(dir, `${file}.ts`)
  if (fs.existsSync(target)) {
    continue
  }
  fs.writeFileSync(
    target,
    `import {definePageMessagesFromEn} from './_define'

const zh = {} as Record<string, string>
const en = {} as Record<string, string>

export const ${exportName} = definePageMessagesFromEn(zh, en)
`,
  )
}

console.log('created stubs')
