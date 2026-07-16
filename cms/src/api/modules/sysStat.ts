import {request} from '../request'
import type {ResourceMetricTrend, SysStat, UserStatTrend} from '@/types/api'

export const sysStatApi = {
    getSysStat: () => {
        return request.post<SysStat>('/sysStat/getSysStat', {})
    },
    getUserStatTrend: () => {
        return request.post<UserStatTrend>('/sysStat/getUserStatTrend', {})
    },
    getResourceMetricTrend: () => {
        return request.post<ResourceMetricTrend>('/sysStat/getResourceMetricTrend', {})
    },
}

export default sysStatApi
