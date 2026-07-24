import request from '../request'
import type {
    ClearRandomNicknamesReq,
    ClearRandomNicknamesRes,
    GetRandomNicknameCfgRes,
    ImportRandomNicknamesReq,
    ImportRandomNicknamesRes,
} from '@/types/api'

export const randomNicknameApi = {
    getRandomNicknameCfg: () => {
        return request.post<GetRandomNicknameCfgRes>('/randomNickname/getRandomNicknameCfg', {})
    },

    importRandomNicknames: (data: ImportRandomNicknamesReq) => {
        return request.post<ImportRandomNicknamesRes>('/randomNickname/importRandomNicknames', data)
    },

    clearRandomNicknames: (data: ClearRandomNicknamesReq) => {
        return request.post<ClearRandomNicknamesRes>('/randomNickname/clearRandomNicknames', data)
    },
}

export default randomNicknameApi
