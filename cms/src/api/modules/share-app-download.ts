import request from '../request'
import type {
    GetShareAppDownloadCfgRes,
    SaveShareAppDownloadCfgReq,
    SaveShareAppDownloadCfgRes,
} from '@/types/api'

export const shareAppDownloadApi = {
    getShareAppDownloadCfg: () => {
        return request.post<GetShareAppDownloadCfgRes>('/shareAppDownload/getShareAppDownloadCfg', {})
    },

    saveShareAppDownloadCfg: (data: SaveShareAppDownloadCfgReq) => {
        return request.post<SaveShareAppDownloadCfgRes>('/shareAppDownload/saveShareAppDownloadCfg', {
            id: data.id,
            shareText: data.shareText,
        })
    },
}

export default shareAppDownloadApi
