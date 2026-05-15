// 环境配置
interface EnvConfig {
    BASE_API: string
    ENV: string
    TIMEOUT: number
    SUCCESS_CODE: number
    BASE_PATH: string  // 添加基础路径配置
}

const config: Record<string, EnvConfig> = {
    development: {
        BASE_API: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8898',
        ENV: 'development',
        TIMEOUT: 10000,
        SUCCESS_CODE: 200,
        BASE_PATH: ''  // 开发环境不需要额外的基础路径
    },
    production: {
        BASE_API: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:1000',
        ENV: 'production',
        TIMEOUT: 10000,
        SUCCESS_CODE: 200,
        BASE_PATH: '/res'  // 生产环境的基础路径
    },
    test: {
        BASE_API: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8898',
        ENV: 'test',
        TIMEOUT: 10000,
        SUCCESS_CODE: 200,
        BASE_PATH: ''  // 测试环境不需要额外的基础路径
    }
}

// 获取当前环境
const getCurrentEnv = (): string => {
    return import.meta.env.MODE || 'development'
}

const env = getCurrentEnv()

export default {
    ...config[env as keyof typeof config]
}