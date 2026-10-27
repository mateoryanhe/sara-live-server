import {request} from '../request'

export interface AuditCodeCfg {
    id: string
    code: string
    createdAt: string
    updatedAt: string
}

export interface GetAuditCodeCfgRes {
    cfg: AuditCodeCfg | null
}

export interface GenerateAuditCodeCfgRes {
    success: boolean
    cfg: AuditCodeCfg | null
}

export const auditCodeApi = {
    getAuditCodeCfg: () => {
        return request.post<GetAuditCodeCfgRes>('/auditCode/getAuditCodeCfg', {})
    },
    generateAuditCodeCfg: (length = 8) => {
        return request.post<GenerateAuditCodeCfgRes>('/auditCode/generateAuditCodeCfg', {length})
    },
}

export default auditCodeApi
