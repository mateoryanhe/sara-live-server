import {request} from '../request'
import type {PageResponse, SimulatorCpuKeyword, SimulatorCpuKeywordQuery} from '@/types/api'

export const simulatorCpuKeywordApi = {
  getList: (params: SimulatorCpuKeywordQuery) => {
    return request.post<PageResponse<SimulatorCpuKeyword>>('/simulatorCpuKeyword/simulatorCpuKeywordList', params)
  },

  create: (data: {keyword: string; remark: string}) => {
    return request.post<{id: string}>('/simulatorCpuKeyword/createSimulatorCpuKeyword', data)
  },

  update: (data: {id: string | number; keyword: string; remark: string}) => {
    return request.post<{success: boolean}>('/simulatorCpuKeyword/updateSimulatorCpuKeyword', data)
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/simulatorCpuKeyword/deleteSimulatorCpuKeyword', {id})
  },
}
