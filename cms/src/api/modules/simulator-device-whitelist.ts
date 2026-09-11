import request from '@/api/request'
import type {PageResponse, SimulatorDeviceWhitelist, SimulatorDeviceWhitelistQuery} from '@/types/api'

export const simulatorDeviceWhitelistApi = {
  getList: (params: SimulatorDeviceWhitelistQuery) => {
    return request.post<PageResponse<SimulatorDeviceWhitelist>>('/simulatorDeviceWhitelist/simulatorDeviceWhitelistList', params)
  },
  create: (data: {deviceId: string}) => {
    return request.post<{id: string}>('/simulatorDeviceWhitelist/createSimulatorDeviceWhitelist', data)
  },
  update: (data: {id: string; deviceId: string}) => {
    return request.post<{success: boolean}>('/simulatorDeviceWhitelist/updateSimulatorDeviceWhitelist', data)
  },
  remove: (id: string) => {
    return request.post<{success: boolean}>('/simulatorDeviceWhitelist/deleteSimulatorDeviceWhitelist', {id})
  },
}

export default simulatorDeviceWhitelistApi
