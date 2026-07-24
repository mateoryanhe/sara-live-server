import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

// 上传管理API
export const uploadApi = {
    // CMS后台上传图片或礼物动画资源,返回保存后的文件名
    uploadFile: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<{ fileName: string; fileUrl: string }>('/upload/uploadFile', formData, {
            timeout: 0, // 大文件上传不限制 axios 超时
            onUploadProgress: (event: AxiosProgressEvent) => {
                if (!onUploadProgress || !event.total) {
                    return
                }
                onUploadProgress(Math.min(100, Math.round((event.loaded * 100) / event.total)))
            },
        })
    }
}

export default uploadApi
