import {request} from '../request'
import type {
  ApkPromotionChannelDailyStatListResponse,
  ApkPromotionChannelDailyStatQuery,
  ApkPromotionChannelListResponse,
  ApkPromotionChannelQuery,
} from '@/types/api'

export const apkPromotionChannelApi = {
  getList: (params: ApkPromotionChannelQuery) => {
    return request.post<ApkPromotionChannelListResponse>('/apkPromotionChannel/apkPromotionChannelList', params)
  },

  create: (data: { channelName: string; channelCode: string }) => {
    return request.post<{ id: string }>('/apkPromotionChannel/createApkPromotionChannel', data)
  },

  update: (data: { id: string; channelName: string; channelCode: string }) => {
    return request.post<{ success: boolean }>('/apkPromotionChannel/updateApkPromotionChannel', data)
  },

  delete: (id: string) => {
    return request.post<{ success: boolean }>('/apkPromotionChannel/deleteApkPromotionChannel', { id })
  },

  initApkUpload: (channelId: string, fileName: string, fileSize: number) => {
    return request.post<{ uploadId: string; chunkSize: number; totalChunk: number }>(
        '/apkPromotionChannel/initApkPromotionChannelApkUpload',
        { channelId, fileName, fileSize },
    )
  },

  uploadApkChunk: (
      uploadId: string,
      chunkIndex: number,
      chunk: Blob,
      fileName: string,
      onUploadProgress?: (percent: number) => void,
  ) => {
    const formData = new FormData()
    formData.append('file', chunk, fileName)
    return request.post<{ chunkIndex: number; fileSize: number }>(
        '/apkPromotionChannel/uploadApkFileChunk',
        formData,
        {
          headers: {
            'X-Upload-Id': uploadId,
            'X-Chunk-Index': String(chunkIndex),
          },
          onUploadProgress: (event) => {
            if (!onUploadProgress || !event.total) {
              return
            }
            onUploadProgress(Math.round((event.loaded * 100) / event.total))
          },
        },
    )
  },

  completeApkUpload: (uploadId: string, channelId: string) => {
    return request.post<{ apkDownloadUrl: string; apkFileName: string }>(
        '/apkPromotionChannel/completeApkPromotionChannelApkUpload',
        { uploadId, channelId },
    )
  },

  abortApkUpload: (uploadId: string) => {
    return request.post('/apkPromotionChannel/abortApkPromotionChannelApkUpload', { uploadId })
  },
}

export const apkPromotionChannelDailyStatApi = {
  getList: (params: ApkPromotionChannelDailyStatQuery) => {
    return request.post<ApkPromotionChannelDailyStatListResponse>(
        '/apkPromotionChannel/apkPromotionChannelDailyStatList',
        params,
    )
  },
}
