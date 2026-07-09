import request from '../request'
import type {PageResponse, VideoCallLogItem, VideoCallLogQuery} from '@/types/api'

export const videoCallLogApi = {
    getVideoCallLogList: (params: VideoCallLogQuery) => {
        return request.post<PageResponse<VideoCallLogItem>>('/call/cmsVideoCallLogList', params)
    },
}
