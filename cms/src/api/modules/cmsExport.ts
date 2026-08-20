import {request} from '../request'
import type {
    CMSExportJobResult,
    CMSExportJobSubmitResult,
    CMSExportResult,
} from '@/types/api'

export const cmsExportApi = {
    submitJob: (data: {exportType: string; payload: object}) => {
        return request.post<CMSExportJobSubmitResult>('/cmsExport/submitJob', data)
    },
    getJob: (data: {jobId: string}) => {
        return request.post<CMSExportJobResult>('/cmsExport/getJob', data)
    },
    deleteExport: (data: {exportId: string}) => {
        return request.post<{success: boolean}>('/cmsExport/deleteExport', data)
    },
}

export default cmsExportApi
