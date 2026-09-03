import {request} from '../request'
import type {
    ResourceMetricPoint,
    ResourceMetricTrendQuery,
    SysStat,
    UserStatTrend,
} from '@/types/api'

export const RESOURCE_METRIC_MAX_POINTS = 10000

const buildResourceMetricQuery = (params: ResourceMetricTrendQuery = {}) => ({
    startTime: params.startTime || '',
    endTime: params.endTime || '',
    limit: params.limit ?? RESOURCE_METRIC_MAX_POINTS,
})

export const sysStatApi = {
    getSysStat: () => {
        return request.post<SysStat>('/sysStat/getSysStat', {})
    },
    getUserStatTrend: () => {
        return request.post<UserStatTrend>('/sysStat/getUserStatTrend', {})
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
        return request.post<{ points: ResourceMetricPoint[] }>(
            '/resourceMetric/getResourceMetricOnlineTrend',
            buildResourceMetricQuery(params),
        )
    },
}

export default sysStatApi
