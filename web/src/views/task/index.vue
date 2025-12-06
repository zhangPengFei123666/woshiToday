<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getTaskList, createTask, updateTask, deleteTask, startTask, stopTask, triggerTask } from '@/api/task'
import { getAllGroups } from '@/api/group'
import type { Task, TaskGroup, CreateTaskRequest } from '@/api/types'

// 列表数据
const tableData = ref<Task[]>([])
const groups = ref<TaskGroup[]>([])
const loading = ref(false)
const total = ref(0)

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  group_id: undefined as number | undefined,
  keyword: '',
  status: -1
})

// 表单相关
const dialogVisible = ref(false)
const dialogTitle = ref('新增任务')
const formRef = ref<FormInstance>()
const formLoading = ref(false)

const formData = reactive<CreateTaskRequest & { id?: number }>({
  group_id: 0,
  name: '',
  description: '',
  cron: '',
  executor_type: 'HTTP',
  executor_handler: '',
  executor_param: '',
  route_strategy: 'ROUND_ROBIN',
  block_strategy: 'SERIAL_EXECUTION',
  shard_num: 1,
  retry_count: 0,
  retry_interval: 0,
  timeout: 0,
  alarm_email: '',
  priority: 0
})

const rules: FormRules = {
  group_id: [{ required: true, message: '请选择任务组', trigger: 'change' }],
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }],
  executor_type: [{ required: true, message: '请选择执行器类型', trigger: 'change' }],
  executor_handler: [{ required: true, message: '请输入执行器Handler', trigger: 'blur' }]
}

// 执行器类型选项
const executorTypeOptions = [
  { label: 'HTTP', value: 'HTTP' },
  { label: 'GRPC', value: 'GRPC' },
  { label: 'SCRIPT', value: 'SCRIPT' }
]

// 路由策略选项
const routeStrategyOptions = [
  { label: '轮询', value: 'ROUND_ROBIN' },
  { label: '随机', value: 'RANDOM' },
  { label: '一致性哈希', value: 'CONSISTENT_HASH' },
  { label: '最不经常使用', value: 'LEAST_FREQUENTLY_USED' },
  { label: '最近最少使用', value: 'LEAST_RECENTLY_USED' },
  { label: '故障转移', value: 'FAILOVER' },
  { label: '分片广播', value: 'SHARDING_BROADCAST' }
]

// 阻塞策略选项
const blockStrategyOptions = [
  { label: '串行执行', value: 'SERIAL_EXECUTION' },
  { label: '丢弃后续', value: 'DISCARD_LATER' },
  { label: '覆盖之前', value: 'COVER_EARLY' }
]

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getTaskList(queryParams)
    tableData.value = res.data.list || []
    total.value = res.data.total
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载任务组
const loadGroups = async () => {
  try {
    const res = await getAllGroups()
    groups.value = res.data || []
  } catch (error) {
    console.error('加载任务组失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  queryParams.keyword = ''
  queryParams.group_id = undefined
  queryParams.status = -1
  queryParams.page = 1
  loadData()
}

// 分页
const handlePageChange = (page: number) => {
  queryParams.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  queryParams.page_size = size
  queryParams.page = 1
  loadData()
}

// 新增
const handleAdd = () => {
  dialogTitle.value = '新增任务'
  Object.assign(formData, {
    id: undefined,
    group_id: groups.value[0]?.id || 0,
    name: '',
    description: '',
    cron: '0 0 * * * ?',
    executor_type: 'HTTP',
    executor_handler: '',
    executor_param: '',
    route_strategy: 'ROUND_ROBIN',
    block_strategy: 'SERIAL_EXECUTION',
    shard_num: 1,
    retry_count: 0,
    retry_interval: 0,
    timeout: 0,
    alarm_email: '',
    priority: 0
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: Task) => {
  dialogTitle.value = '编辑任务'
  Object.assign(formData, {
    id: row.id,
    group_id: row.group_id,
    name: row.name,
    description: row.description,
    cron: row.cron,
    executor_type: row.executor_type,
    executor_handler: row.executor_handler,
    executor_param: row.executor_param,
    route_strategy: row.route_strategy,
    block_strategy: row.block_strategy,
    shard_num: row.shard_num,
    retry_count: row.retry_count,
    retry_interval: row.retry_interval,
    timeout: row.timeout,
    alarm_email: row.alarm_email,
    priority: row.priority
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: Task) => {
  try {
    await ElMessageBox.confirm(`确定要删除任务"${row.name}"吗？`, '提示', { type: 'warning' })
    await deleteTask(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {}
}

// 启动
const handleStart = async (row: Task) => {
  try {
    await startTask(row.id)
    ElMessage.success('启动成功')
    loadData()
  } catch (error) {}
}

// 停止
const handleStop = async (row: Task) => {
  try {
    await stopTask(row.id)
    ElMessage.success('停止成功')
    loadData()
  } catch (error) {}
}

// 手动触发
const handleTrigger = async (row: Task) => {
  try {
    await ElMessageBox.confirm(`确定要手动触发任务"${row.name}"吗？`, '提示', { type: 'info' })
    await triggerTask(row.id)
    ElMessage.success('触发成功')
  } catch (error) {}
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    formLoading.value = true
    
    const { id, ...data } = formData
    if (id) {
      await updateTask(id, data)
      ElMessage.success('更新成功')
    } else {
      await createTask(data)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadData()
  } catch (error) {
  } finally {
    formLoading.value = false
  }
}

onMounted(() => {
  loadGroups()
  loadData()
})
</script>

<template>
  <div class="task-page">
    <div class="page-card">
      <div class="page-title">任务管理</div>
      
      <!-- 搜索栏 -->
      <div class="table-toolbar">
        <div class="table-toolbar-left">
          <el-select v-model="queryParams.group_id" placeholder="选择任务组" clearable style="width: 160px">
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
          <el-input v-model="queryParams.keyword" placeholder="搜索任务名称" style="width: 200px" clearable />
          <el-select v-model="queryParams.status" placeholder="状态" style="width: 100px">
            <el-option label="全部" :value="-1" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <div class="table-toolbar-right">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增任务
          </el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="任务名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="group.name" label="任务组" width="120" />
        <el-table-column prop="cron" label="Cron表达式" width="140" />
        <el-table-column prop="executor_handler" label="Handler" width="140" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="next_trigger_time" label="下次触发" width="160" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 0" type="success" link @click="handleStart(row)">启动</el-button>
            <el-button v-else type="warning" link @click="handleStop(row)">停止</el-button>
            <el-button type="primary" link @click="handleTrigger(row)">触发</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="任务组" prop="group_id">
              <el-select v-model="formData.group_id" placeholder="请选择任务组" style="width: 100%">
                <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="任务名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入任务名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Cron表达式" prop="cron">
              <el-input v-model="formData.cron" placeholder="如: 0 0 * * * ?" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="执行器类型" prop="executor_type">
              <el-select v-model="formData.executor_type" style="width: 100%">
                <el-option v-for="opt in executorTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="Handler" prop="executor_handler">
          <el-input v-model="formData.executor_handler" placeholder="HTTP类型填写URL，如: http://localhost:8081/job/demo" />
        </el-form-item>
        <el-form-item label="执行参数" prop="executor_param">
          <el-input v-model="formData.executor_param" type="textarea" :rows="2" placeholder="JSON格式参数" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路由策略">
              <el-select v-model="formData.route_strategy" style="width: 100%">
                <el-option v-for="opt in routeStrategyOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="阻塞策略">
              <el-select v-model="formData.block_strategy" style="width: 100%">
                <el-option v-for="opt in blockStrategyOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="重试次数">
              <el-input-number v-model="formData.retry_count" :min="0" :max="10" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="超时时间">
              <el-input-number v-model="formData.timeout" :min="0" placeholder="秒" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="优先级">
              <el-input-number v-model="formData.priority" :min="0" :max="100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="任务描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

