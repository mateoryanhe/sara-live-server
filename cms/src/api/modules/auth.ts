import {request} from '../request'
import type {LoginReq, LoginRes} from '@/types/api'

const authApi = {
    // CMS登录
    cmsLogin: (data: LoginReq) => {
        return request.post<LoginRes>('/auth/cmsLogin', data)
    },

}

export default authApi
