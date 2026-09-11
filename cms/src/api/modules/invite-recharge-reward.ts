import request from '../request'
import type {
  GetInviteRechargeRewardCfgRes,
  SaveInviteRechargeRewardCfgReq,
  SaveInviteRechargeRewardCfgRes,
} from '@/types/api'

export const inviteRechargeRewardApi = {
  getCfg: () => {
    return request.post<GetInviteRechargeRewardCfgRes>('/inviteRechargeReward/getInviteRechargeRewardCfg', {})
  },

  saveCfg: (data: SaveInviteRechargeRewardCfgReq) => {
    return request.post<SaveInviteRechargeRewardCfgRes>('/inviteRechargeReward/saveInviteRechargeRewardCfg', data)
  },
}

export default inviteRechargeRewardApi
