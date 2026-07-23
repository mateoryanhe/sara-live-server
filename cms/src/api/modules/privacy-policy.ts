import request from '../request'
import type {GetPrivacyPolicyCfgRes, SavePrivacyPolicyCfgReq, SavePrivacyPolicyCfgRes} from '@/types/api'

export const privacyPolicyApi = {
    getPrivacyPolicyCfg: () => {
        return request.post<GetPrivacyPolicyCfgRes>('/privacyPolicy/getPrivacyPolicyCfg', {})
    },

    savePrivacyPolicyCfg: (data: SavePrivacyPolicyCfgReq) => {
        return request.post<SavePrivacyPolicyCfgRes>('/privacyPolicy/savePrivacyPolicyCfg', data)
    },
}

export default privacyPolicyApi
