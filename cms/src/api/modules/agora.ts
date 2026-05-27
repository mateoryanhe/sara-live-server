import {request} from '../request'
import type {GetAgoraCfgRes, SaveAgoraCfgReq, SaveAgoraCfgRes} from '@/types/api'

export const agoraApi = {
    getAgoraCfg: () => {
        return request.post<GetAgoraCfgRes>('/agora/getAgoraCfg', {})
    },

    saveAgoraCfg: (data: SaveAgoraCfgReq) => {
        return request.post<SaveAgoraCfgRes>('/agora/saveAgoraCfg', data)
    },
}

export default agoraApi
