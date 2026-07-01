const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.apng']
const videoExts = ['.mp4', '.webm']

export type MediaPreviewType = 'image' | 'video' | 'file' | 'none'

export const getExt = (name: string): string => {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.slice(idx).toLowerCase() : ''
}

export const isImageFile = (fileName: string): boolean => imageExts.includes(getExt(fileName))

export const isVideoFile = (fileName: string): boolean => videoExts.includes(getExt(fileName))

const matchMediaByPath = (path: string, exts: string[]): boolean => exts.includes(getExt(path))

export const isImageUrl = (url: string): boolean => {
  if (!url) return false
  try {
    const pathname = new URL(url, window.location.origin).pathname
    return matchMediaByPath(pathname, imageExts)
  } catch {
    return matchMediaByPath(url, imageExts)
  }
}

export const isVideoUrl = (url: string): boolean => {
  if (!url) return false
  try {
    const pathname = new URL(url, window.location.origin).pathname
    return matchMediaByPath(pathname, videoExts)
  } catch {
    return matchMediaByPath(url, videoExts)
  }
}

export const resolveMediaPreviewType = (url: string, fileName = ''): MediaPreviewType => {
  if (url && isVideoUrl(url)) return 'video'
  if (url && isImageUrl(url)) return 'image'
  if (fileName && isVideoFile(fileName)) return 'video'
  if (fileName && isImageFile(fileName)) return 'image'
  if (fileName || url) return 'file'
  return 'none'
}

export const resolveFilePreviewType = (fileName: string): MediaPreviewType => {
  if (isVideoFile(fileName)) return 'video'
  if (isImageFile(fileName)) return 'image'
  return 'file'
}
