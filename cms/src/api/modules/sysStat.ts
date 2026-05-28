import {request} from '../request'
import type {SysStat} from '@/types/api'

export const sysStatApi = {
    getSysStat: () => {
        return request.post<SysStat>('/sysStat/getSysStat', {})
    },
}

export default sysStatApi
