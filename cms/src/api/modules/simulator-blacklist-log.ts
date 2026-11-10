import {request} from '../request'
import type {
  CMSSimulatorBlacklistLogItem,
  CMSSimulatorBlacklistLogListQuery,
  PageResponse,
} from '@/types/api'

export const simulatorBlacklistLogApi = {
  listSimulatorBlacklistLog: (data: CMSSimulatorBlacklistLogListQuery) => {
    return request.post<PageResponse<CMSSimulatorBlacklistLogItem>>('/simulatorBlacklistLog/listSimulatorBlacklistLog', {
      key: data.key || '',
      loginType: data.loginType || '',
      startTime: data.startTime || 0,
      endTime: data.endTime || 0,
      pageIndex: data.pageIndex,
      pageSize: data.pageSize,
    })
  },
}
