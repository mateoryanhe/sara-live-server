import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {diamondApi, goldApi} from '@/api'
import {formatWalletBalance} from '@/utils/number-format'

/** 与 go-src/constants/currency ReasonGmAdjustReset 一致 */
export const GM_ADJUST_REASON_RESET = 35

export type WalletCurrencyType = 'gold' | 'diamond'

export function walletBalanceNumber(val: unknown): number {
    const n = Number(val)
    return Number.isFinite(n) ? n : 0
}

export function useWalletCurrencyReset() {
    const {t} = useI18n()

    const confirmDangerousActionTwice = async (
        firstMessage: string,
        firstTitle: string,
        secondMessage: string,
        secondTitle: string,
    ) => {
        const boxOpts = {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning' as const,
        }
        await ElMessageBox.confirm(firstMessage, firstTitle, boxOpts)
        await ElMessageBox.confirm(secondMessage, secondTitle, boxOpts)
    }

    /** 将用户指定货币余额扣减至 0；余额已为 0 时返回 false */
    const resetWalletCurrency = async (
        userId: string | number,
        type: WalletCurrencyType,
        balance: number,
    ): Promise<boolean> => {
        if (balance <= 0) {
            ElMessage.warning(
                type === 'gold'
                    ? t('pages.userList.resetGoldAlreadyZero')
                    : t('pages.userList.resetDiamondAlreadyZero'),
            )
            return false
        }
        const balanceText = formatWalletBalance(balance)
        const isGold = type === 'gold'
        await confirmDangerousActionTwice(
            isGold
                ? t('pages.userList.resetGoldConfirm1', {id: userId, balance: balanceText})
                : t('pages.userList.resetDiamondConfirm1', {id: userId, balance: balanceText}),
            isGold ? t('pages.userList.resetGoldTitle') : t('pages.userList.resetDiamondTitle'),
            isGold
                ? t('pages.userList.resetGoldConfirm2', {id: userId})
                : t('pages.userList.resetDiamondConfirm2', {id: userId}),
            isGold ? t('pages.userList.resetGoldTitleFinal') : t('pages.userList.resetDiamondTitleFinal'),
        )
        const api = isGold ? goldApi : diamondApi
        await api.sub({
            userId,
            amount: balance,
            reason: GM_ADJUST_REASON_RESET,
        })
        ElMessage.success(isGold ? t('pages.userList.resetGoldSuccess') : t('pages.userList.resetDiamondSuccess'))
        return true
    }

    return {resetWalletCurrency, walletBalanceNumber}
}
