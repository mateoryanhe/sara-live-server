import {request} from '../request'
import type {
    CoinMerchantTransferLogItem,
    CoinMerchantTransferLogQuery,
    PageResponse,
} from '@/types/api'

export const coinMerchantTransferLogApi = {
    getList: (params: CoinMerchantTransferLogQuery) => {
        return request.post<PageResponse<CoinMerchantTransferLogItem>>(
            '/gold/cmsCoinMerchantTransferLogList',
            params,
        )
    },
}

export default coinMerchantTransferLogApi
