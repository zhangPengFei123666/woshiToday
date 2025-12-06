import { get } from '@/utils/request'
import type { ApiResponse, PageResult } from '@/utils/request'
import type { ExecutorNode, ExecutorListParams } from './types'

// 获取执行器列表
export function getExecutorList(params: ExecutorListParams): Promise<ApiResponse<PageResult<ExecutorNode>>> {
  return get('/executor', params)
}

// 获取执行器详情
export function getExecutorDetail(id: string): Promise<ApiResponse<ExecutorNode>> {
  return get(`/executor/${id}`)
}

// 获取在线执行器
export function getOnlineExecutors(groupId: number): Promise<ApiResponse<ExecutorNode[]>> {
  return get('/executor/online', { group_id: groupId })
}

