import service from '@/utils/requests'
import type { ApiResponse } from '@/types'

// 通用GET请求
export function get<T>(url: string, params?: object): Promise<ApiResponse<T>> {
  return service.get(url, { params })
}

// 通用POST请求
export function post<T>(url: string, data?: object): Promise<ApiResponse<T>> {
  return service.post(url, data)
}
