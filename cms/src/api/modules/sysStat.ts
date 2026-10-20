import {request} from '../request'
import type {
    GameStatTrend,
    ResourceMetricPoint,
    ResourceMetricTrend,
    ResourceMetricTrendQuery,
    SysStat,
    UserStatTrend,
    UserStatTrendQuery,
} from '@/types/api'

export const RESOURCE_METRIC_MAX_POINTS = 26000

const buildResourceMetricQuery = (params: ResourceMetricTrendQuery = {}) => ({
    startTime: params.startTime || '',
    endTime: params.endTime || '',
    limit: params.limit ?? RESOURCE_METRIC_MAX_POINTS,
})

export const sysStatApi = {
    getSysStat: () => {
        return request.post<SysStat>('/sysStat/getSysStat', {})
    },
    resetBasicDashboardStats: () => {
        return request.post<{ success: boolean }>('/sysStat/resetBasicDashboardStats', {})
    },
    getUserStatTrend: (params: UserStatTrendQuery = {}) => {
        return request.post<UserStatTrend>('/sysStat/getUserStatTrend', {
            granularity: params.granularity || '',
            startDate: params.startDate || '',
            endDate: params.endDate || '',
        })
    },
    getResourceMetricGameStatTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<GameStatTrend>(
            '/resourceMetric/getResourceMetricGameStatTrend',
            buildResourceMetricQuery(params),
        )
    },
    getResourceMetricMemoryTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<{ points: ResourceMetricPoint[] }>(
            '/resourceMetric/getResourceMetricMemoryTrend',
            buildResourceMetricQuery(params),
        )
    },
    getResourceMetricHeapTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<{ points: ResourceMetricPoint[] }>(
            '/resourceMetric/getResourceMetricHeapTrend',
            buildResourceMetricQuery(params),
        )
    },
    getResourceMetricRatioTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<{ points: ResourceMetricPoint[] }>(
            '/resourceMetric/getResourceMetricRatioTrend',
            buildResourceMetricQuery(params),
        )
    },
    getResourceMetricCpuTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<{ points: ResourceMetricPoint[] }>(
            '/resourceMetric/getResourceMetricCpuTrend',
            buildResourceMetricQuery(params),
        )
    },
    getResourceMetricOnlineTrend: (params: ResourceMetricTrendQuery = {}) => {
        return request.post<ResourceMetricTrend>(
            '/resourceMetric/getResourceMetricOnlineTrend',
            buildResourceMetricQuery(params),
        )
    },
}

export default sysStatApi
