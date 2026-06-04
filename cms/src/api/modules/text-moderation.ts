import request from '../request'
import type {GetTextModerationCfgRes, SaveTextModerationCfgReq, SaveTextModerationCfgRes} from '@/types/api'

export const textModerationApi = {
    getTextModerationCfg: () => {
        return request.post<GetTextModerationCfgRes>('/textModeration/getTextModerationCfg', {})
    },

    saveTextModerationCfg: (data: SaveTextModerationCfgReq) => {
        return request.post<SaveTextModerationCfgRes>('/textModeration/saveTextModerationCfg', data)
    },
}

export default textModerationApi
