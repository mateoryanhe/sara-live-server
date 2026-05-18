import {request} from '../request'

// 上传管理API
export const uploadApi = {
    // CMS后台上传图片或礼物动画资源,返回保存后的文件名
    uploadFile: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<{ fileName: string }>('/upload/uploadFile', formData, {
            headers: {'Content-Type': 'multipart/form-data'}
        })
    }
}

export default uploadApi
