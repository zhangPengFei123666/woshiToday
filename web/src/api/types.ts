// 用户相关类型
export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  last_login_time: string
  last_login_ip: string
  created_at: string
  roles: Role[]
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

// 任务组相关类型
export interface TaskGroup {
  id: number
  name: string
  description: string
  app_name: string
  status: number
  created_by: number
  created_at: string
  updated_at: string
}

// 任务相关类型
export interface Task {
  id: number
  group_id: number
  name: string
  description: string
  cron: string
  executor_type: string
  executor_handler: string
  executor_param: string
  route_strategy: string
  block_strategy: string
  shard_num: number
  retry_count: number
  retry_interval: number
  timeout: number
  alarm_email: string
  priority: number
  status: number
  version: number
  next_trigger_time: string
  last_trigger_time: string
  created_by: number
  created_at: string
  updated_at: string
  group?: TaskGroup
}

// 任务实例相关类型
export interface TaskInstance {
  id: number
  task_id: number
  group_id: number
  executor_id: string
  executor_address: string
  executor_handler: string
  executor_param: string
  shard_index: number
  shard_total: number
  trigger_type: string
  trigger_time: string
  schedule_time: string
  start_time: string
  end_time: string
  status: number
  result_code: number
  result_msg: string
  retry_count: number
  alarm_status: number
  created_at: string
  updated_at: string
  task?: Task
}

// 任务日志
export interface TaskLog {
  id: number
  instance_id: number
  task_id: number
  log_time: string
  log_level: string
  log_content: string
  created_at: string
}

// 执行器节点
export interface ExecutorNode {
  id: string
  group_id: number
  app_name: string
  host: string
  port: number
  weight: number
  max_concurrent: number
  current_load: number
  cpu_usage: number
  memory_usage: number
  status: number
  last_heartbeat: string
  registered_at: string
  updated_at: string
}

// 统计信息
export interface InstanceStatistics {
  total: number
  success: number
  failed: number
  running: number
  pending: number
  cancelled: number
  rate: number
}

// 分页请求参数
export interface PageParams {
  page: number
  page_size: number
}

// 任务列表请求参数
export interface TaskListParams extends PageParams {
  group_id?: number
  keyword?: string
  status?: number
}

// 实例列表请求参数
export interface InstanceListParams extends PageParams {
  task_id?: number
  status?: number
  start_time?: string
  end_time?: string
}

// 执行器列表请求参数
export interface ExecutorListParams extends PageParams {
  group_id?: number
  status?: number
}

// 创建任务请求
export interface CreateTaskRequest {
  group_id: number
  name: string
  description?: string
  cron: string
  executor_type: string
  executor_handler: string
  executor_param?: string
  route_strategy?: string
  block_strategy?: string
  shard_num?: number
  retry_count?: number
  retry_interval?: number
  timeout?: number
  alarm_email?: string
  priority?: number
  dependency_ids?: number[]
}

// 创建任务组请求
export interface CreateGroupRequest {
  name: string
  description?: string
  app_name: string
}

