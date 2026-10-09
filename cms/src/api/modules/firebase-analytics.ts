import request from '../request'
import type {
    GetFirebaseAnalyticsCfgRes,
    SaveFirebaseAnalyticsCfgReq,
    SaveFirebaseAnalyticsCfgRes,
} from '@/types/api'

export const firebaseAnalyticsApi = {
    getFirebaseAnalyticsCfg: () => {
        return request.post<GetFirebaseAnalyticsCfgRes>('/firebaseAnalytics/getFirebaseAnalyticsCfg', {})
    },

    saveFirebaseAnalyticsCfg: (data: SaveFirebaseAnalyticsCfgReq) => {
        return request.post<SaveFirebaseAnalyticsCfgRes>('/firebaseAnalytics/saveFirebaseAnalyticsCfg', data)
    },
}
