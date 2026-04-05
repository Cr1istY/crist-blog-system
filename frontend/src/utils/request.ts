import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { message } from '@/utils/naiveDiscrete'
import router from '@/router'

interface RefreshResponse {
  access_token: string
}

export interface AxiosRetryConfig extends AxiosRequestConfig {
  _retry?: boolean
  customErrorHandling?: boolean
}

interface FailedQueueCallback {
  resolve: (value?: string) => void
  reject: (error?: Error) => void
}

const service: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 5000,
  withCredentials: true,
})

service.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let isRefreshing = false
let failedQueue: FailedQueueCallback[] = []

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach((callback) => {
    if (error) {
      callback.reject(error)
    } else {
      callback.resolve(token || undefined)
    }
  })
  failedQueue = []
}

service.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as AxiosRetryConfig

    // 检查是否有自定义错误处理
    if (originalRequest.customErrorHandling) {
      return Promise.reject(error)
    }

    // 安全访问error.response.status
    const status = error.response?.status

    // 处理401未授权错误
    if (status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // 如果已经在刷新token，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then((token) => {
            if (originalRequest.headers && token) {
              originalRequest.headers.Authorization = `Bearer ${token}`
              return service(originalRequest)
            }
            return service(originalRequest)
          })
          .catch((err) => {
            return Promise.reject(err)
          })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        // 关键修改：在刷新token时不要使用带token的service实例，避免循环
        const refreshToken = localStorage.getItem('refresh_token')
        const refreshConfig: AxiosRequestConfig = {
          headers: refreshToken ? { Authorization: `Bearer ${refreshToken}` } : {},
        }

        const {
          data: { access_token },
        } = await axios.post<RefreshResponse>('/api/auth/refresh', {}, refreshConfig)

        console.log('refresh token success')
        localStorage.setItem('access_token', access_token)

        // 更新原始请求的token
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`
        }

        // 处理队列中的请求
        processQueue(null, access_token)

        return service(originalRequest)
      } catch (refreshError) {
        console.log('refresh 401') // 确保这行代码能执行

        // 清理认证信息
        // localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')

        // 通知用户
        message.error('登录已过期，请重新登录')

        // 跳转到登录页
        try {
          router.push({ name: 'login' })
        } catch (routerError) {
          console.error('路由跳转失败:', routerError)
          window.location.href = '/admin'
        }

        // 处理失败队列
        if (refreshError instanceof Error) {
          processQueue(refreshError, null)
        } else {
          processQueue(new Error(String(refreshError)), null)
        }

        isRefreshing = false
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // 其他错误处理
    let errorMessage = '请求失败'
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }

    message.error(errorMessage)
    return Promise.reject(error)
  },
)

export default service
