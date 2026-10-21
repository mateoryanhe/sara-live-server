/** 幸运礼物返奖倍率挡位类型 */
export type LuckyGiftMultiplierType = 'low' | 'medium' | 'high'

export interface LuckyGiftMultiplierTier {
  /** 前端行 id */
  key: string
  multiplier: number
  /** 抽奖权重(整数,越大越容易命中) */
  weight: number
  type: LuckyGiftMultiplierType
}

export const DEFAULT_LUCKY_GIFT_MULTIPLIER_TIERS: Omit<LuckyGiftMultiplierTier, 'key'>[] = [
  {multiplier: 5, weight: 50, type: 'low'},
  {multiplier: 10, weight: 30, type: 'low'},
  {multiplier: 88, weight: 15, type: 'medium'},
  {multiplier: 200, weight: 4, type: 'high'},
  {multiplier: 500, weight: 1, type: 'high'},
]

export function createDefaultLuckyGiftMultiplierTiers(): LuckyGiftMultiplierTier[] {
  return DEFAULT_LUCKY_GIFT_MULTIPLIER_TIERS.map((row, index) => ({
    key: `default-${index}`,
    ...row,
  }))
}

export function newLuckyGiftMultiplierTierKey() {
  return `tier-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}
