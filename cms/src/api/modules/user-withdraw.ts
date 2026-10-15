import request from '../request'

export interface UserWithdrawLogItem {
    id: string
    userId: string
    nickname?: string
    diamondAmount?: number
    goldAmount?: number
    withdrawUsd: number
    status: number
    transferOrderId?: string
    transferPlatformNo?: string
    transferCurrency?: string
    transferLocalAmount?: number
    transferFailMsg?: string
    rejectReason?: string
    createdAt?: string
    transferAt?: string
}

export interface ListUserWithdrawLogsReq {
    userId?: string
    status?: number
    page?: number
    pageSize?: number
}

export interface ListUserWithdrawLogsRes {
    total: number
    list: UserWithdrawLogItem[]
}

export interface BatchIdsReq {
    ids: string[]
    reason?: string
}

export interface BatchCountRes {
    successCount: number
    failCount: number
    message?: string
}

export const userWithdrawApi = {
    listWithdrawLogs: (data: ListUserWithdrawLogsReq) =>
        request.post<ListUserWithdrawLogsRes>('/userWithdrawCms/listWithdrawLogs', data),
    batchApproveWithdraw: (data: BatchIdsReq) =>
        request.post<BatchCountRes>('/userWithdrawCms/batchApproveWithdraw', data),
    batchRejectWithdraw: (data: BatchIdsReq) =>
        request.post<BatchCountRes>('/userWithdrawCms/batchRejectWithdraw', data),
    batchSubmitPayout: (data: BatchIdsReq) =>
        request.post<BatchCountRes>('/userWithdrawCms/batchSubmitPayout', data),
}

export default userWithdrawApi
