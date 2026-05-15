import axios, {
    type AxiosInstance,
    type AxiosRequestConfig,
    type AxiosResponse,
    type InternalAxiosRequestConfig
} from 'axios'
import envConfig from '@/config/env'
import {ElMessage} from "element-plus";

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
        // 在发送请求之前做些什么，比如添加token和authId
        const token = localStorage.getItem('token')
        const authId = localStorage.getItem('authId')

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
            ElMessage.error(`出现错误，错误码：${res.code}`)
            return Promise.reject(new Error(res.code || 'Error'))
        } else {
            return res.data
        }
    },
    (error: any) => {
        // 对响应错误做点什么
        console.error('Response Error:', error)

        if (error.response?.status === 401) {
            // token过期，跳转到登录页
            localStorage.removeItem('token')
            localStorage.removeItem('authId')
            window.location.href = '/login'
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