import {definePageMessagesFromEn} from './_define'

const zh = {
  afLink: 'AF 链接',
  afLinkPlaceholder: 'https://... (OneLink，勿带 inviteCode/roomId)',
  afLinkTip: 'App 第二行：在链接后用 & 拼接 inviteCode=本人邀请码&roomId=请求房间ID。',
  shareTextTab: '分享文案(英/西/葡/印/印尼，默认英文)',
  shareTextPlaceholder: '留空则使用 Go 内置默认文案',
  shareTextBuiltinTip: '各语言留空时，App 与预览使用 Go 内置文案(room_share_default_text.go)。',
  shareTextMaxLength: '单语言文案最多 4000 字符',
  previewTitle: '预览(英文+示例参数)',
  previewEmpty: '（保存后可预览）',
  saveSuccess: '直播间分享配置已保存',
  needAfLink: '请填写 AF 链接',
}

export const liveRoomShareMessages = definePageMessagesFromEn(zh, {
  afLink: 'AF link',
  afLinkPlaceholder: 'https://... (OneLink; no inviteCode/roomId)',
  afLinkTip: 'Line 2: append &inviteCode=&roomId= from the request.',
  shareTextTab: 'Share copy',
  shareTextPlaceholder: 'Leave blank to use Go built-in default copy',
  shareTextBuiltinTip: 'When a locale is empty, App and preview use built-in copy from Go (room_share_default_text.go).',
  shareTextMaxLength: 'Max 4000 chars per locale',
  previewTitle: 'Preview (EN + sample params)',
  previewEmpty: '(Save to preview)',
  saveSuccess: 'Live room share settings saved',
  needAfLink: 'AF link required',
})