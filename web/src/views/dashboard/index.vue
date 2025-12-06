<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import { getStatistics, getRecentInstances } from '@/api/instance'
import type { InstanceStatistics, TaskInstance } from '@/api/types'

const statistics = ref<InstanceStatistics>({
  total: 0,
  success: 0,
  failed: 0,
  running: 0,
  pending: 0,
  cancelled: 0,
  rate: 0
})

const recentInstances = ref<TaskInstance[]>([])
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 状态标签样式
const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'info',
    1: 'warning',
    2: 'primary',
    3: 'success',
    4: 'danger',
    5: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待调度',
    1: '调度中',
    2: '执行中',
    3: '成功',
    4: '失败',
    5: '已取消'
  }
  return map[status] || '未知'
}

// 格式化时间
const formatTime = (time: string) => {
  return time ? dayjs(time).format('MM-DD HH:mm:ss') : '-'
}

// 加载数据
const loadData = async () => {
  try {
    const [statsRes, recentRes] = await Promise.all([
      getStatistics(),
      getRecentInstances(10)
    ])
    statistics.value = statsRes.data
    recentInstances.value = recentRes.data || []
    
    initChart()
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return
  
  if (chart) {
    chart.dispose()
  }
  
  chart = echarts.init(chartRef.value)
  
  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'center'
    },
    series: [
      {
        name: '任务状态',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['60%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: statistics.value.success, name: '成功', itemStyle: { color: '#67c23a' } },
          { value: statistics.value.failed, name: '失败', itemStyle: { color: '#f56c6c' } },
          { value: statistics.value.running, name: '执行中', itemStyle: { color: '#409eff' } },
          { value: statistics.value.pending, name: '待执行', itemStyle: { color: '#909399' } }
        ]
      }
    ]
  }
  
  chart.setOption(option)
}

// 窗口大小变化时重绘图表
const handleResize = () => {
  chart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>

<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon size="32"><List /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.total }}</div>
            <div class="stat-label">总执行次数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card success">
          <div class="stat-icon">
            <el-icon size="32"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.success }}</div>
            <div class="stat-label">执行成功</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card warning">
          <div class="stat-icon">
            <el-icon size="32"><CircleClose /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.failed }}</div>
            <div class="stat-label">执行失败</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card info">
          <div class="stat-icon">
            <el-icon size="32"><TrendCharts /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.rate.toFixed(1) }}%</div>
            <div class="stat-label">成功率</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 图表和最近执行 -->
    <el-row :gutter="20">
      <el-col :span="12">
        <div class="page-card">
          <div class="card-header">
            <span class="card-title">执行状态分布</span>
          </div>
          <div ref="chartRef" class="chart-container"></div>
        </div>
      </el-col>
      
      <el-col :span="12">
        <div class="page-card">
          <div class="card-header">
            <span class="card-title">最近执行记录</span>
            <el-button type="primary" link @click="$router.push('/instance')">
              查看更多
            </el-button>
          </div>
          <el-table :data="recentInstances" stripe size="small" max-height="320">
            <el-table-column prop="id" label="实例ID" width="80" />
            <el-table-column prop="task.name" label="任务名称" show-overflow-tooltip />
            <el-table-column prop="trigger_time" label="触发时间" width="140">
              <template #default="{ row }">
                {{ formatTime(row.trigger_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-row {
  margin-bottom: 20px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 24px;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 20px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.15);
}

.stat-card.success {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-card.warning {
  background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
}

.stat-card.info {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon {
  width: 64px;
  height: 64px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.85;
}

.page-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chart-container {
  height: 320px;
}
</style>

