import {request} from '../request'
import type {CMSProcessDownLogItem, CMSProcessDownLogListQuery, PageResponse} from '@/types/api'

export const processDownLogApi = {
    listProcessDownLog: (data: CMSProcessDownLogListQuery) => {
        return request.post<PageResponse<CMSProcessDownLogItem>>('/processDownLog/listProcessDownLog', {
            host: data.host || '',
            startTime: data.startTime || 0,
            endTime: data.endTime || 0,
            pageIndex: data.pageIndex,
            pageSize: data.pageSize,
        })
    },
}
