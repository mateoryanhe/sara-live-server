import {request} from '../request'
import type {GetFirebaseCfgRes, SaveFirebaseCfgReq, SaveFirebaseCfgRes} from '@/types/api'

export const firebaseApi = {
    getFirebaseCfg: () => {
        return request.post<GetFirebaseCfgRes>('/firebase/getFirebaseCfg', {})
    },

    saveFirebaseCfg: (data: SaveFirebaseCfgReq) => {
        return request.post<SaveFirebaseCfgRes>('/firebase/saveFirebaseCfg', data)
    },
}

export default firebaseApi
