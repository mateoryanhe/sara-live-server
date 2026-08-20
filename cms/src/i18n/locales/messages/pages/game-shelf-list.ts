import {definePageMessagesFromEn} from './_define'

export const gameShelfListMessages = definePageMessagesFromEn(
  {
    noteTitle: '说明',
    tipLine1: '展示已上架游戏（读 game_cfgs 永久缓存），App 端仅展示此列表中的游戏。可在「游戏库」中搜索厂商游戏并上架。',
    fetchFailed: '获取上架游戏列表失败',
  },
  {
    noteTitle: 'Note',
    tipLine1: 'Shows published games (from game_cfgs permanent cache). The App displays only these games. Search vendor games and publish them in Game Library.',
    fetchFailed: 'Failed to load published game list',
  },
)
