/** 单文件分片上传时的并发分片数 */
export const APK_CHUNK_UPLOAD_CONCURRENCY = 4

export interface ApkChunkUploadSession {
  uploadId: string
  chunkSize: number
  totalChunk: number
}

export interface ApkChunkUploadHandlers {
  init: () => Promise<ApkChunkUploadSession>
  uploadChunk: (
      uploadId: string,
      chunkIndex: number,
      chunk: Blob,
      fileName: string,
      onChunkProgress?: (percent: number) => void,
  ) => Promise<void>
  complete: (uploadId: string) => Promise<void>
  abort: (uploadId: string) => Promise<void>
}

async function runWithConcurrency<T>(
    items: T[],
    concurrency: number,
    worker: (item: T) => Promise<void>,
): Promise<void> {
  if (items.length === 0) {
    return
  }
  const limit = Math.max(1, Math.min(concurrency, items.length))
  let next = 0
  const runners = Array.from({length: limit}, async () => {
    for (;;) {
      const index = next
      next += 1
      if (index >= items.length) {
        return
      }
      await worker(items[index])
    }
  })
  await Promise.all(runners)
}

/** 分片并行上传 APK（各渠道独立 session，可多路并行调用） */
export async function uploadApkFileInParallelChunks(
    file: File,
    handlers: ApkChunkUploadHandlers,
    onProgress?: (percent: number) => void,
    concurrency = APK_CHUNK_UPLOAD_CONCURRENCY,
): Promise<void> {
  const session = await handlers.init()
  if (!session.uploadId || session.chunkSize <= 0 || session.totalChunk <= 0) {
    throw new Error('invalid upload session')
  }
  const {uploadId, chunkSize, totalChunk} = session
  const chunkPercents = new Array<number>(totalChunk).fill(0)

  const reportProgress = () => {
    if (!onProgress) {
      return
    }
    let loaded = 0
    for (let i = 0; i < totalChunk; i += 1) {
      const start = i * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      const size = end - start
      loaded += (size * chunkPercents[i]) / 100
    }
    onProgress(Math.min(99, Math.round((loaded * 100) / file.size)))
  }

  const chunkIndexes = Array.from({length: totalChunk}, (_, i) => i)
  try {
    await runWithConcurrency(chunkIndexes, concurrency, async (chunkIndex) => {
      const start = chunkIndex * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      const chunk = file.slice(start, end)
      await handlers.uploadChunk(uploadId, chunkIndex, chunk, file.name, (p) => {
        chunkPercents[chunkIndex] = p
        reportProgress()
      })
      chunkPercents[chunkIndex] = 100
      reportProgress()
    })
    await handlers.complete(uploadId)
    onProgress?.(100)
  } catch (error) {
    try {
      await handlers.abort(uploadId)
    } catch {
      // ignore
    }
    throw error
  }
}
