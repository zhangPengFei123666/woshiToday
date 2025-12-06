<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { getExecutorList } from '@/api/executor'
import { getAllGroups } from '@/api/group'
import type { ExecutorNode, TaskGroup } from '@/api/types'

// 列表数据
const tableData = ref<ExecutorNode[]>([])
const groups = ref<TaskGroup[]>([])
const loading = ref(false)
const total = ref(0)

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  group_id: undefined as number | undefined,
  status: -1
})

// 状态选项
const statusOptions = [
  { label: '全部', value: -1 },
  { label: '在线', value: 1 },
  { label: '离线', value: 0 }
]

// 格式化时间
const formatTime = (time: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getExecutorList(queryParams)
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

onMounted(() => {
  loadGroups()
  loadData()
})
</script>

<template>
  <div class="executor-page">
    <div class="page-card">
      <div class="page-title">执行器管理</div>
      
      <!-- 搜索栏 -->
      <div class="table-toolbar">
        <div class="table-toolbar-left">
          <el-select v-model="queryParams.group_id" placeholder="选择任务组" clearable style="width: 160px">
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
          <el-select v-model="queryParams.status" placeholder="状态" style="width: 120px">
            <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="执行器ID" width="280" show-overflow-tooltip />
        <el-table-column prop="app_name" label="应用名称" width="150" />
        <el-table-column label="节点地址" width="160">
          <template #default="{ row }">{{ row.host }}:{{ row.port }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              <span class="status-dot" :class="row.status === 1 ? 'success' : 'danger'"></span>
              {{ row.status === 1 ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="负载" width="120">
          <template #default="{ row }">
            <el-progress :percentage="Math.min((row.current_load / row.max_concurrent) * 100, 100)" :stroke-width="6" />
            <span class="load-text">{{ row.current_load }}/{{ row.max_concurrent }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="100">
          <template #default="{ row }">
            <el-progress 
              :percentage="row.cpu_usage" 
              :stroke-width="6" 
              :color="row.cpu_usage > 80 ? '#f56c6c' : row.cpu_usage > 60 ? '#e6a23c' : '#67c23a'"
            />
          </template>
        </el-table-column>
        <el-table-column label="内存" width="100">
          <template #default="{ row }">
            <el-progress 
              :percentage="row.memory_usage" 
              :stroke-width="6" 
              :color="row.memory_usage > 80 ? '#f56c6c' : row.memory_usage > 60 ? '#e6a23c' : '#67c23a'"
            />
          </template>
        </el-table-column>
        <el-table-column prop="last_heartbeat" label="最后心跳" width="160">
          <template #default="{ row }">{{ formatTime(row.last_heartbeat) }}</template>
        </el-table-column>
        <el-table-column prop="registered_at" label="注册时间" width="160">
          <template #default="{ row }">{{ formatTime(row.registered_at) }}</template>
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
  </div>
</template>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 4px;
}

.status-dot.success {
  background-color: #67c23a;
  animation: pulse 2s infinite;
}

.status-dot.danger {
  background-color: #f56c6c;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.load-text {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}
</style>

