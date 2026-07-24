import axios, {
    type AxiosInstance,
    type AxiosRequestConfig,
    type AxiosResponse,
    type InternalAxiosRequestConfig
} from 'axios'
import envConfig from '@/config/env'
import {ElMessage} from "element-plus"
import router from '@/router'
import {clearAuthSession, getAuthId, getToken} from '@/utils/auth'

const AUTH_ERROR_CODES = new Set([1, 2, 3])

function redirectToLogin(): void {
    clearAuthSession()
    const current = router.currentRoute.value
    if (current.name === 'Login') {
        return
    }
    const query = current.path !== '/' && current.path !== '/login'
        ? {redirect: current.fullPath}
        : undefined
    void router.replace({path: '/login', query})
}

function isAuthErrorCode(code: unknown): boolean {
    return typeof code === 'number' && AUTH_ERROR_CODES.has(code)
}

// 创建axios实例
const service: AxiosInstance = axios.create({
    baseURL: envConfig.BASE_API,
    timeout: envConfig.TIMEOUT,
    headers: {
        'Content-Type': 'application/json'
    }
})

// 请求拦截器
service.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // FormData 须由浏览器自动设置 multipart boundary,不能沿用 application/json
        if (config.data instanceof FormData) {
            if (config.headers) {
                delete config.headers['Content-Type']
            }
        }

        // 在发送请求之前做些什么，比如添加token和authId
        const token = getToken()
        const authId = getAuthId()

        if (token) {
            config.headers!['token'] = token
        }
        if (authId) {
            config.headers!['authId'] = authId
        }

        return config
    },
    (error: any) => {
        // 对请求错误做些什么
        console.error('Request Error:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    (response: AxiosResponse) => {
        // 对响应数据做点什么
        const res = response.data

        // 如果自定义code不是0，则判断为错误
        if (res.code !== 0) {
            if (isAuthErrorCode(res.code)) {
                redirectToLogin()
                return Promise.reject(new Error(String(res.code)))
            }
            ElMessage.error(res.message || `出现错误，错误码：${res.code}`)
            return Promise.reject(new Error(res.code || 'Error'))
        } else {
            return res.data
        }
    },
    (error: any) => {
        // 对响应错误做点什么
        console.error('Response Error:', error)

        if (error.response?.status === 401 || isAuthErrorCode(error.response?.data?.code)) {
            redirectToLogin()
        }

        return Promise.reject(error)
    }
)

// 通用请求方法
export const request = {
    get: <T>(url: string, config?: AxiosRequestConfig): Promise<T> => {
        return service.get(url, config)
    },
    post: <T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
        return service.post(url, data, config)
    },
    put: <T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
        return service.put(url, data, config)
    },
    delete: <T>(url: string, config?: AxiosRequestConfig): Promise<T> => {
        return service.delete(url, config)
    }
}

export default service