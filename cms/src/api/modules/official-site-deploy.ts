import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface StaticSiteDeployInfo {
    domain: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    lastUploadAt: string
}

export interface DeployStaticSiteZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetStaticSiteDeployInfoRes {
    info: StaticSiteDeployInfo | null
}

export interface SaveStaticSiteDeployCfgRes {
    success: boolean
}

export interface InitStaticSiteFileUploadRes {
    uploadId: string
    chunkSize: number
    totalChunks: number
}

export interface UploadStaticSiteFileChunkRes {
    chunkIndex: number
    fileSize: number
}

export interface CompleteStaticSiteFileUploadRes {
    fileType: 'zip' | 'apk'
    fileName: string
    fileSize: number
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
    downloadUrl: string
}

export interface StaticSiteDeployApi {
    getDeployInfo: () => Promise<GetStaticSiteDeployInfoRes>
    saveMapping?: (domain: string, deployPath: string) => Promise<SaveStaticSiteDeployCfgRes>
    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => Promise<DeployStaticSiteZipRes>
    initFileUpload?: (fileName: string, fileSize: number) => Promise<InitStaticSiteFileUploadRes>
    uploadFileChunk?: (
        uploadId: string,
        chunkIndex: number,
        chunk: Blob,
        fileName: string,
        onUploadProgress?: (percent: number) => void,
    ) => Promise<UploadStaticSiteFileChunkRes>
    completeFileUpload?: (uploadId: string) => Promise<CompleteStaticSiteFileUploadRes>
    abortFileUpload?: (uploadId: string) => Promise<void>
}

const deployZip = (
    url: string,
    file: File,
    onUploadProgress?: (percent: number) => void,
) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<DeployStaticSiteZipRes>(url, formData, {
        timeout: 0,
        onUploadProgress: (event: AxiosProgressEvent) => {
            if (!onUploadProgress || !event.total) {
                return
            }
            onUploadProgress(Math.min(100, Math.round((event.loaded * 100) / event.total)))
        },
    })
}

export const officialSiteDeployApi: StaticSiteDeployApi = {
    getDeployInfo: () => request.post<GetStaticSiteDeployInfoRes>(
        '/officialSiteDeploy/getOfficialSiteDeployInfo',
        {},
    ),
    saveMapping: (domain, deployPath) => request.post<SaveStaticSiteDeployCfgRes>(
        '/officialSiteDeploy/saveOfficialSiteDeployCfg',
        {domain, deployPath},
    ),
    deployZip: (file, onUploadProgress) => deployZip(
        '/officialSiteDeploy/deployZip',
        file,
        onUploadProgress,
    ),
    initFileUpload: (fileName, fileSize) => request.post<InitStaticSiteFileUploadRes>(
        '/officialSiteDeploy/initFileUpload',
        {fileName, fileSize},
    ),
    uploadFileChunk: (uploadId, chunkIndex, chunk, fileName, onUploadProgress) => {
        const formData = new FormData()
        formData.append('file', chunk, fileName)
        return request.post<UploadStaticSiteFileChunkRes>(
            '/officialSiteDeploy/uploadFileChunk',
            formData,
            {
                timeout: 120_000,
                headers: {
                    'X-Upload-Id': uploadId,
                    'X-Chunk-Index': String(chunkIndex),
                },
                onUploadProgress: (event: AxiosProgressEvent) => {
                    if (!onUploadProgress || !event.total) {
                        return
                    }
                    onUploadProgress(Math.min(100, Math.round((event.loaded * 100) / event.total)))
                },
            },
        )
    },
    completeFileUpload: (uploadId) => request.post<CompleteStaticSiteFileUploadRes>(
        '/officialSiteDeploy/completeFileUpload',
        {uploadId},
        {timeout: 0},
    ),
    abortFileUpload: (uploadId) => request.post<void>(
        '/officialSiteDeploy/abortFileUpload',
        {uploadId},
    ),
}

export const thirdPayOfficialSiteDeployApi: StaticSiteDeployApi = {
    getDeployInfo: () => request.post<GetStaticSiteDeployInfoRes>(
        '/thirdPayOfficialSiteDeploy/getThirdPayOfficialSiteDeployInfo',
        {},
    ),
    saveMapping: (domain, deployPath) => request.post<SaveStaticSiteDeployCfgRes>(
        '/thirdPayOfficialSiteDeploy/saveThirdPayOfficialSiteDeployCfg',
        {domain, deployPath},
    ),
    deployZip: (file, onUploadProgress) => deployZip(
        '/thirdPayOfficialSiteDeploy/deployZip',
        file,
        onUploadProgress,
    ),
    initFileUpload: (fileName, fileSize) => request.post<InitStaticSiteFileUploadRes>(
        '/thirdPayOfficialSiteDeploy/initFileUpload',
        {fileName, fileSize},
    ),
    uploadFileChunk: (uploadId, chunkIndex, chunk, fileName, onUploadProgress) => {
        const formData = new FormData()
        formData.append('file', chunk, fileName)
        return request.post<UploadStaticSiteFileChunkRes>(
            '/thirdPayOfficialSiteDeploy/uploadFileChunk',
            formData,
            {
                timeout: 120_000,
                headers: {
                    'X-Upload-Id': uploadId,
                    'X-Chunk-Index': String(chunkIndex),
                },
                onUploadProgress: (event: AxiosProgressEvent) => {
                    if (!onUploadProgress || !event.total) {
                        return
                    }
                    onUploadProgress(Math.min(100, Math.round((event.loaded * 100) / event.total)))
                },
            },
        )
    },
    completeFileUpload: (uploadId) => request.post<CompleteStaticSiteFileUploadRes>(
        '/thirdPayOfficialSiteDeploy/completeFileUpload',
        {uploadId},
        {timeout: 0},
    ),
    abortFileUpload: (uploadId) => request.post<void>(
        '/thirdPayOfficialSiteDeploy/abortFileUpload',
        {uploadId},
    ),
}
