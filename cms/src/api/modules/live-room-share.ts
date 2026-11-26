import request from '@/api/request'
import type {GetLiveRoomShareCfgRes, SaveLiveRoomShareCfgReq, SaveLiveRoomShareCfgRes} from '@/types/api'

export const liveRoomShareApi = {
    getLiveRoomShareCfg: () => request.post<GetLiveRoomShareCfgRes>('/liveRoomShare/getLiveRoomShareCfg', {}),
    saveLiveRoomShareCfg: (data: SaveLiveRoomShareCfgReq) =>
        request.post<SaveLiveRoomShareCfgRes>('/liveRoomShare/saveLiveRoomShareCfg', data),
}

export default liveRoomShareApi