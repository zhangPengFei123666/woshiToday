import { get, post, put, del } from '@/utils/request'
import type { ApiResponse, PageResult } from '@/utils/request'
import type { Task, TaskListParams, CreateTaskRequest } from './types'

// 获取任务列表
export function getTaskList(params: TaskListParams): Promise<ApiResponse<PageResult<Task>>> {
  return get('/task', params)
}

// 获取任务详情
export function getTaskDetail(id: number): Promise<ApiResponse<Task>> {
  return get(`/task/${id}`)
}

// 创建任务
export function createTask(data: CreateTaskRequest): Promise<ApiResponse<Task>> {
  return post('/task', data)
}

// 更新任务
export function updateTask(id: number, data: CreateTaskRequest): Promise<ApiResponse<Task>> {
  return put(`/task/${id}`, data)
}

// 删除任务
export function deleteTask(id: number): Promise<ApiResponse<null>> {
  return del(`/task/${id}`)
}

// 启动任务
export function startTask(id: number): Promise<ApiResponse<null>> {
  return post(`/task/${id}/start`)
}

// 停止任务
export function stopTask(id: number): Promise<ApiResponse<null>> {
  return post(`/task/${id}/stop`)
}

// 手动触发任务
export function triggerTask(id: number, param?: string): Promise<ApiResponse<null>> {
  return post(`/task/${id}/trigger`, { param })
}

// 获取下次触发时间
export function getNextTriggerTimes(cron: string, count?: number): Promise<ApiResponse<string[]>> {
  return get('/task/next-trigger-times', { cron, count: count || 5 })
}

