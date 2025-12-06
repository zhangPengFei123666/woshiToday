import { post, get } from '@/utils/request'
import type { ApiResponse } from '@/utils/request'
import type { LoginResult, UserInfo } from './types'

// 登录
export function login(username: string, password: string): Promise<ApiResponse<LoginResult>> {
  return post('/auth/login', { username, password })
}

// 退出登录
export function logout(): Promise<ApiResponse<null>> {
  return post('/auth/logout')
}

// 获取当前用户信息
export function getCurrentUser(): Promise<ApiResponse<UserInfo>> {
  return get('/user/current')
}

// 修改密码
export function changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse<null>> {
  return post('/user/password', { old_password: oldPassword, new_password: newPassword })
}

