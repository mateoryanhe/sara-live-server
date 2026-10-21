import {request} from '../request'
import type {LoginReq, LoginRes} from '@/types/api'

export interface CmsSessionModulesRes {
    admin: boolean
    superAdmin: boolean
    modules: LoginRes['modules']
}

const authApi = {
    // CMS登录
    cmsLogin: (data: LoginReq) => {
        return request.post<LoginRes>('/auth/cmsLogin', data)
    },

    /** 刷新当前登录用户的模块权限（Layout 初始化） */
    cmsSessionModules: (data: Record<string, never> = {}) => {
        return request.post<CmsSessionModulesRes>('/auth/cmsSessionModules', data)
    },

}

export default authApi
