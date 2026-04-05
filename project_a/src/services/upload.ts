import service from '@/utils/request'
import type { AxiosRequestConfig } from 'axios'

export interface UploadResponse {
  url: string
  id: string
  filename: string
  size: number
  width?: number
  height?: number
  created_at?: string
}

export interface UploadProgress {
  loaded: number
  total: number
  percentage: number
}

export interface UploadConfig extends AxiosRequestConfig {
  onProgress?: (progress: UploadProgress) => void
}

/**
 * 单个图片上传 - 适配 Go 后端
 */
export const uploadImage = async (file: File, config?: UploadConfig): Promise<UploadResponse> => {
  const formData = new FormData()
  formData.append('image', file)

  try {
    const response = await service.post<UploadResponse>('/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total && config?.onProgress) {
          config.onProgress({
            loaded: progressEvent.loaded,
            total: progressEvent.total,
            percentage: Math.round((progressEvent.loaded / progressEvent.total) * 100),
          })
        }
      },
      timeout: 1200000,
      ...config,
    })

    return response.data
  } catch (error) {
    console.error('图片上传失败:', error)
    throw error
  }
}

/**
 * 批量图片上传 - 适配 Go 后端
 */
export const uploadImages = async (
  files: File[],
  onProgress?: (fileIndex: number, progress: UploadProgress) => void,
): Promise<UploadResponse[]> => {
  const formData = new FormData()
  files.forEach((file) => {
    formData.append('images', file)
  })

  const uploadService = service.create({
    timeout: 60000,
  })

  try {
    const response = await uploadService.post<{ images: UploadResponse[] }>(
      '/upload/images',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total && onProgress) {
            onProgress(0, {
              loaded: progressEvent.loaded,
              total: progressEvent.total,
              percentage: Math.round((progressEvent.loaded / progressEvent.total) * 100),
            })
          }
        },
      },
    )

    return response.data.images
  } catch (error) {
    console.error('批量上传失败:', error)
    throw error
  }
}

/**
 * 分片上传初始化
 */
export const initChunkedUpload = async (
  filename: string,
  fileSize: number,
  fileType: string,
): Promise<{ upload_id: string }> => {
  const response = await service.post('/upload/init', {
    filename,
    file_size: fileSize,
    file_type: fileType,
  })
  return response.data
}

/**
 * 上传分片
 */
export const uploadChunk = async (
  uploadId: string,
  chunk: Blob,
  chunkIndex: number,
  totalChunks: number,
  onProgress?: (progress: UploadProgress) => void,
): Promise<void> => {
  const formData = new FormData()
  formData.append('chunk', chunk)
  formData.append('uploadId', uploadId)
  formData.append('chunkIndex', chunkIndex.toString())
  formData.append('totalChunks', totalChunks.toString())

  await service.post('/upload/chunk', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total && onProgress) {
        onProgress({
          loaded: progressEvent.loaded,
          total: progressEvent.total,
          percentage: Math.round((progressEvent.loaded / progressEvent.total) * 100),
        })
      }
    },
  })
}

/**
 * 合并分片
 */
export const mergeChunks = async (uploadId: string, filename: string): Promise<UploadResponse> => {
  const response = await service.post('/upload/merge', {
    uploadId,
    filename,
  })
  return response.data
}

/**
 * 分片上传完整流程
 */
export const uploadImageChunked = async (
  file: File,
  chunkSize: number = 2 * 1024 * 1024,
  onProgress?: (progress: UploadProgress) => void,
): Promise<UploadResponse> => {
  const fileSize = file.size
  const totalChunks = Math.ceil(fileSize / chunkSize)

  // 1. 初始化上传
  const { upload_id } = await initChunkedUpload(file.name, fileSize, file.type)

  // 2. 上传所有分片
  let uploadedBytes = 0
  for (let i = 0; i < totalChunks; i++) {
    const start = i * chunkSize
    const end = Math.min(start + chunkSize, fileSize)
    const chunk = file.slice(start, end)

    await uploadChunk(upload_id, chunk, i, totalChunks, (progress) => {
      if (onProgress) {
        onProgress({
          loaded: uploadedBytes + progress.loaded,
          total: fileSize,
          percentage: Math.round(((uploadedBytes + progress.loaded) / fileSize) * 100),
        })
      }
    })

    uploadedBytes += chunk.size
  }

  // 3. 合并分片
  return await mergeChunks(upload_id, file.name)
}
