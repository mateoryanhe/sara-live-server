import {request} from '../request'
import type {GetAttrReportCfgRes, SaveAttrReportCfgReq, SaveAttrReportCfgRes} from '@/types/api'

export const attrReportCfgApi = {
    getAttrReportCfg: () => {
        return request.post<GetAttrReportCfgRes>('/attrReportCfg/getAttrReportCfg', {})
    },

    saveAttrReportCfg: (data: SaveAttrReportCfgReq) => {
        return request.post<SaveAttrReportCfgRes>('/attrReportCfg/saveAttrReportCfg', data)
    },
}

export default attrReportCfgApi
