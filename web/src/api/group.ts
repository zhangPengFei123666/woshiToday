import { get, post, put, del } from '@/utils/request'
import type { ApiResponse, PageResult } from '@/utils/request'
import type { TaskGroup, PageParams, CreateGroupRequest } from './types'

// 获取任务组列表
export function getGroupList(params: PageParams & { keyword?: string }): Promise<ApiResponse<PageResult<TaskGroup>>> {
  return get('/group', params)
}

// 获取所有任务组
export function getAllGroups(): Promise<ApiResponse<TaskGroup[]>> {
  return get('/group/all')
}

// 获取任务组详情
export function getGroupDetail(id: number): Promise<ApiResponse<TaskGroup>> {
  return get(`/group/${id}`)
}

// 创建任务组
export function createGroup(data: CreateGroupRequest): Promise<ApiResponse<TaskGroup>> {
  return post('/group', data)
}

// 更新任务组
export function updateGroup(id: number, data: CreateGroupRequest): Promise<ApiResponse<TaskGroup>> {
  return put(`/group/${id}`, data)
}

// 删除任务组
export function deleteGroup(id: number): Promise<ApiResponse<null>> {
  return del(`/group/${id}`)
}

