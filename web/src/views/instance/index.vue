<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { getInstanceList, cancelInstance, retryInstance, getInstanceLogs } from '@/api/instance'
import type { TaskInstance, TaskLog } from '@/api/types'

// 列表数据
const tableData = ref<TaskInstance[]>([])
const loading = ref(false)
const total = ref(0)

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  task_id: undefined as number | undefined,
  status: -1,
  start_time: '',
  end_time: ''
})

// 日志对话框
const logDialogVisible = ref(false)
const logData = ref<TaskLog[]>([])
const logLoading = ref(false)
const currentInstance = ref<TaskInstance | null>(null)

// 状态选项
const statusOptions = [
  { label: '全部', value: -1 },
  { label: '待调度', value: 0 },
  { label: '调度中', value: 1 },
  { label: '执行中', value: 2 },
  { label: '成功', value: 3 },
  { label: '失败', value: 4 },
  { label: '已取消', value: 5 }
]

// 状态标签
const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'info', 1: 'warning', 2: 'primary', 3: 'success', 4: 'danger', 5: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待调度', 1: '调度中', 2: '执行中', 3: '成功', 4: '失败', 5: '已取消'
  }
  return map[status] || '未知'
}

// 格式化时间
const formatTime = (time: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

// 计算执行时长
const getDuration = (row: TaskInstance) => {
  if (!row.start_time || !row.end_time) return '-'
  const start = dayjs(row.start_time)
  const end = dayjs(row.end_time)
  const ms = end.diff(start)
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${(ms / 60000).toFixed(1)}min`
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getInstanceList(queryParams)
    tableData.value = res.data.list || []
    total.value = res.data.total
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  queryParams.task_id = undefined
  queryParams.status = -1
  queryParams.start_time = ''
  queryParams.end_time = ''
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

// 取消
const handleCancel = async (row: TaskInstance) => {
  try {
    await ElMessageBox.confirm('确定要取消该任务实例吗？', '提示', { type: 'warning' })
    await cancelInstance(row.id)
    ElMessage.success('取消成功')
    loadData()
  } catch (error) {}
}

// 重试
const handleRetry = async (row: TaskInstance) => {
  try {
    await ElMessageBox.confirm('确定要重试该任务实例吗？', '提示', { type: 'info' })
    await retryInstance(row.id)
    ElMessage.success('重试任务已创建')
    loadData()
  } catch (error) {}
}

// 查看日志
const handleViewLogs = async (row: TaskInstance) => {
  currentInstance.value = row
  logDialogVisible.value = true
  logLoading.value = true
  
  try {
    const res = await getInstanceLogs(row.id, 1, 100)
    logData.value = res.data.list || []
  } catch (error) {
    console.error('加载日志失败:', error)
  } finally {
    logLoading.value = false
  }
}

// 日志级别颜色
const getLogLevelColor = (level: string) => {
  const map: Record<string, string> = {
    DEBUG: '#909399',
    INFO: '#409eff',
    WARN: '#e6a23c',
    ERROR: '#f56c6c'
  }
  return map[level] || '#909399'
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="instance-page">
    <div class="page-card">
      <div class="page-title">执行记录</div>
      
      <!-- 搜索栏 -->
      <div class="table-toolbar">
        <div class="table-toolbar-left">
          <el-input v-model.number="queryParams.task_id" placeholder="任务ID" style="width: 120px" clearable />
          <el-select v-model="queryParams.status" placeholder="状态" style="width: 120px">
            <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
          <el-date-picker
            v-model="queryParams.start_time"
            type="datetime"
            placeholder="开始时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 180px"
          />
          <el-date-picker
            v-model="queryParams.end_time"
            type="datetime"
            placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 180px"
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="实例ID" width="80" />
        <el-table-column prop="task.name" label="任务名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="trigger_type" label="触发类型" width="90" />
        <el-table-column prop="trigger_time" label="触发时间" width="160">
          <template #default="{ row }">{{ formatTime(row.trigger_time) }}</template>
        </el-table-column>
        <el-table-column prop="executor_address" label="执行节点" width="140" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="执行时长" width="100">
          <template #default="{ row }">{{ getDuration(row) }}</template>
        </el-table-column>
        <el-table-column prop="result_msg" label="执行结果" show-overflow-tooltip />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewLogs(row)">日志</el-button>
            <el-button v-if="row.status === 0 || row.status === 1" type="warning" link @click="handleCancel(row)">取消</el-button>
            <el-button v-if="row.status === 4" type="success" link @click="handleRetry(row)">重试</el-button>
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

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="执行日志" width="800px" destroy-on-close>
      <div class="log-info" v-if="currentInstance">
        <el-descriptions :column="3" size="small" border>
          <el-descriptions-item label="实例ID">{{ currentInstance.id }}</el-descriptions-item>
          <el-descriptions-item label="任务名称">{{ currentInstance.task?.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentInstance.status)" size="small">{{ getStatusText(currentInstance.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="触发时间">{{ formatTime(currentInstance.trigger_time) }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">{{ formatTime(currentInstance.start_time) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间">{{ formatTime(currentInstance.end_time) }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <div class="log-container" v-loading="logLoading">
        <div v-if="logData.length === 0" class="log-empty">暂无日志</div>
        <div v-else class="log-list">
          <div v-for="log in logData" :key="log.id" class="log-item">
            <span class="log-time">{{ dayjs(log.log_time).format('HH:mm:ss.SSS') }}</span>
            <span class="log-level" :style="{ color: getLogLevelColor(log.log_level) }">[{{ log.log_level }}]</span>
            <span class="log-content">{{ log.log_content }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.log-info {
  margin-bottom: 16px;
}

.log-container {
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.log-empty {
  color: #909399;
  text-align: center;
  padding: 40px 0;
}

.log-list {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  line-height: 1.8;
}

.log-item {
  color: #d4d4d4;
}

.log-time {
  color: #6a9955;
  margin-right: 8px;
}

.log-level {
  margin-right: 8px;
  font-weight: 600;
}

.log-content {
  color: #d4d4d4;
}
</style>

