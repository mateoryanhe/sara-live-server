import request from '../request'
import type {GetWalletExchangeCfgRes, SaveWalletExchangeCfgReq, SaveWalletExchangeCfgRes} from '@/types/api'

export const walletApi = {
    getWalletExchangeCfg: () => {
        return request.post<GetWalletExchangeCfgRes>('/wallet/getWalletExchangeCfg', {})
    },

    saveWalletExchangeCfg: (data: SaveWalletExchangeCfgReq) => {
        return request.post<SaveWalletExchangeCfgRes>('/wallet/saveWalletExchangeCfg', data)
    },
}

export default walletApi
