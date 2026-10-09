import {request} from '../request'
import type {GetMetaPixelCfgRes, SaveMetaPixelCfgReq, SaveMetaPixelCfgRes} from '@/types/api'

export const metaPixelApi = {
    getMetaPixelCfg: () => {
        return request.post<GetMetaPixelCfgRes>('/metaPixel/getMetaPixelCfg', {})
    },

    saveMetaPixelCfg: (data: SaveMetaPixelCfgReq) => {
        return request.post<SaveMetaPixelCfgRes>('/metaPixel/saveMetaPixelCfg', data)
    },
}

export default metaPixelApi
