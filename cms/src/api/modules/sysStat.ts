import {request} from '../request'
import type {SysStat, UserStatTrend} from '@/types/api'

export const sysStatApi = {
    getSysStat: () => {
        return request.post<SysStat>('/sysStat/getSysStat', {})
    },
    getUserStatTrend: () => {
        return request.post<UserStatTrend>('/sysStat/getUserStatTrend', {})
    },
}

export default sysStatApi
