import { get, post } from '@/utils/request'
import type { ApiResponse, PageResult } from '@/utils/request'
import type { TaskInstance, TaskLog, InstanceListParams, InstanceStatistics } from './types'

// 获取任务实例列表
export function getInstanceList(params: InstanceListParams): Promise<ApiResponse<PageResult<TaskInstance>>> {
  return get('/instance', params)
}

// 获取任务实例详情
export function getInstanceDetail(id: number): Promise<ApiResponse<TaskInstance>> {
  return get(`/instance/${id}`)
}

// 取消任务实例
export function cancelInstance(id: number): Promise<ApiResponse<null>> {
  return post(`/instance/${id}/cancel`)
}

// 重试任务实例
export function retryInstance(id: number): Promise<ApiResponse<TaskInstance>> {
  return post(`/instance/${id}/retry`)
}

// 获取任务实例日志
export function getInstanceLogs(id: number, page: number, pageSize: number): Promise<ApiResponse<PageResult<TaskLog>>> {
  return get(`/instance/${id}/logs`, { page, page_size: pageSize })
}

// 获取统计信息
export function getStatistics(params?: { task_id?: number; start_time?: string; end_time?: string }): Promise<ApiResponse<InstanceStatistics>> {
  return get('/instance/statistics', params)
}

// 获取最近实例
export function getRecentInstances(limit?: number): Promise<ApiResponse<TaskInstance[]>> {
  return get('/instance/recent', { limit: limit || 10 })
}

