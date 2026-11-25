import request from '../request'
import type {
    GetSofieShareAppDownloadCfgRes,
    SaveSofieShareAppDownloadCfgReq,
    SaveSofieShareAppDownloadCfgRes,
} from '@/types/api'

export const sofieShareAppDownloadApi = {
    getSofieShareAppDownloadCfg: () => {
        return request.post<GetSofieShareAppDownloadCfgRes>('/sofieShareAppDownload/getSofieShareAppDownloadCfg', {})
    },

    saveSofieShareAppDownloadCfg: (data: SaveSofieShareAppDownloadCfgReq) => {
        return request.post<SaveSofieShareAppDownloadCfgRes>('/sofieShareAppDownload/saveSofieShareAppDownloadCfg', data)
    },
}

export default sofieShareAppDownloadApi