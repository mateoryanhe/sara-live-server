/** 礼物分类 live_gifts.category_type */
export const GIFT_CATEGORY_NORMAL = 0
export const GIFT_CATEGORY_LUCKY = 1

export const GIFT_CATEGORY_VALUES = [GIFT_CATEGORY_NORMAL, GIFT_CATEGORY_LUCKY] as const

export function isKnownGiftCategoryType(categoryType: number | undefined | null): boolean {
  return categoryType === GIFT_CATEGORY_NORMAL || categoryType === GIFT_CATEGORY_LUCKY
}
