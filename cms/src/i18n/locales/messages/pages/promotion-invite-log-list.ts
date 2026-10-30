import {definePageMessagesFromEn} from './_define'

const zh = {
  keywordPlaceholder: '用户ID / 昵称（被邀请人或邀请人）',
  dateRange: '邀请时间',
  dateRangeSeparator: '至',
  invitee: '被邀请人',
  inviter: '邀请人',
  inviterId: '邀请人ID',
  inviteTime: '邀请时间',
  inviterRewardGold: '邀请人返现金币',
  fetchFailed: '获取邀请列表失败',
} as const

export const promotionInviteLogListMessages = definePageMessagesFromEn(zh, {
  keywordPlaceholder: 'User ID / nickname (invitee or inviter)',
  dateRange: 'Invite time',
  dateRangeSeparator: 'to',
  invitee: 'Invitee',
  inviter: 'Inviter',
  inviterId: 'Inviter ID',
  inviteTime: 'Invite time',
  inviterRewardGold: 'Inviter cashback gold',
  fetchFailed: 'Failed to load invite list',
})
