import {request} from '../request'
import type {GetPrivacyPolicyCfgRes, SavePrivacyPolicyCfgReq, SavePrivacyPolicyCfgRes} from '@/types/api'

export const sofiePrivacyPolicyApi = {
    getSofiePrivacyPolicyCfg: () => {
        return request.post<GetPrivacyPolicyCfgRes>('/sofiePrivacyPolicy/getSofiePrivacyPolicyCfg', {})
    },

    saveSofiePrivacyPolicyCfg: (data: SavePrivacyPolicyCfgReq) => {
        return request.post<SavePrivacyPolicyCfgRes>('/sofiePrivacyPolicy/saveSofiePrivacyPolicyCfg', data)
    },
}

export default sofiePrivacyPolicyApi
