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

/** 挡位权重占比(%)；总权重为 0 时返回 0 */
export function luckyGiftTierWeightSharePercent(
  tiers: Array<Pick<LuckyGiftMultiplierTier, 'weight'> | null | undefined>,
  weight: number,
): number {
  const total = (tiers || []).reduce((sum, tier) => {
    const w = Number(tier?.weight)
    return sum + (Number.isFinite(w) && w > 0 ? w : 0)
  }, 0)
  if (total <= 0) {
    return 0
  }
  const w = Number(weight)
  if (!Number.isFinite(w) || w <= 0) {
    return 0
  }
  return (w / total) * 100
}
