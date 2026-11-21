import request from '../request'
import type {
    GetAudienceOnlineAnchorGreetCfgRes,
    SaveAudienceOnlineAnchorGreetCfgReq,
    SaveAudienceOnlineAnchorGreetCfgRes,
} from '@/types/api'

export const audienceOnlineAnchorGreetApi = {
    getCfg: () => {
        return request.post<GetAudienceOnlineAnchorGreetCfgRes>('/audienceOnlineAnchorGreet/getAudienceOnlineAnchorGreetCfg', {})
    },

    saveCfg: (data: SaveAudienceOnlineAnchorGreetCfgReq) => {
        return request.post<SaveAudienceOnlineAnchorGreetCfgRes>('/audienceOnlineAnchorGreet/saveAudienceOnlineAnchorGreetCfg', data)
    },
}

export default audienceOnlineAnchorGreetApi